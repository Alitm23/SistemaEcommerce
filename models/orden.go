package models

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type Orden struct {
	ID            int       `json:"id"`
	UsuarioID     int       `json:"usuario_id"`
	Total         float64   `json:"total"`
	Estado        string    `json:"estado"`
	MetodoPago    string    `json:"metodo_pago"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

type ProductoVendido struct {
	ID       int
	Nombre   string
	Cantidad int
}

// CanalNotificacionOrden gestiona el envío asíncrono de correos/alertas
var CanalNotificacionOrden = make(chan int, 100)

// GestorNotificaciones simula el envío de correos en segundo plano
func GestorNotificaciones() {
	go func() {
		for ordenID := range CanalNotificacionOrden {
			// Simulación de envío de correo
			log.Printf("Procesando notificación... Correo de confirmación enviado para la Orden #%d", ordenID)
			time.Sleep(2 * time.Second) // Simula la latencia de red
		}
	}()
}

// CrearOrden inserta la orden y despacha la notificación concurrente
func CrearOrden(o Orden) (int, error) {
	var nuevoID int
	query := `
		INSERT INTO orden (usuario_id, total, estado, metodo_pago, fecha_creacion) 
		VALUES ($1, $2, $3, $4, NOW()) RETURNING id`

	err := db.DB.QueryRow(query, o.UsuarioID, o.Total, o.Estado, o.MetodoPago).Scan(&nuevoID)

	if err == nil {
		//envia la notificación de manera asincrona al canal
		CanalNotificacionOrden <- nuevoID
		return nuevoID, nil
	}
	return 0, err
}

// Crear Orden
func CrearOrdenDesdeCarrito(o Orden, carritoID int) (int, error) {
	if o.Estado != "pagado" && o.Estado != "pendiente" && o.Estado != "cancelado" {
		o.Estado = "pendiente"
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var nuevoID int
	err = tx.QueryRow(`
		INSERT INTO orden (usuario_id, total, estado, metodo_pago, fecha_creacion)
		VALUES ($1, $2, $3, $4, NOW()) RETURNING id`,
		o.UsuarioID, o.Total, o.Estado, o.MetodoPago,
	).Scan(&nuevoID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`
		INSERT INTO detalle_orden (orden_id, producto_id, cantidad, precio)
		SELECT $1, pt.producto_id, ic.cantidad, ic.precio_unitario
		FROM item_carrito ic
		INNER JOIN producto_talla pt ON ic.producto_talla_id = pt.id
		WHERE ic.carrito_id = $2`,
		nuevoID, carritoID,
	)
	if err != nil {
		return 0, err
	}

	if o.Estado == "pagado" {
		resultado, err := tx.Exec(`
			UPDATE producto_talla pt
			SET stock = pt.stock - ic.cantidad
			FROM item_carrito ic
			WHERE ic.carrito_id = $1
				AND ic.producto_talla_id = pt.id
				AND pt.stock >= ic.cantidad`,
			carritoID,
		)
		if err != nil {
			return 0, err
		}

		var totalItems int64
		err = tx.QueryRow("SELECT COUNT(*) FROM item_carrito WHERE carrito_id = $1", carritoID).Scan(&totalItems)
		if err != nil {
			return 0, err
		}
		filasActualizadas, err := resultado.RowsAffected()
		if err != nil {
			return 0, err
		}
		if filasActualizadas != totalItems {
			return 0, errors.New("stock insuficiente para completar la orden")
		}
	}

	_, err = tx.Exec("DELETE FROM item_carrito WHERE carrito_id = $1", carritoID)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec("UPDATE carrito SET estado = 'comprado' WHERE id = $1", carritoID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	CanalNotificacionOrden <- nuevoID
	return nuevoID, nil
}

// ListarOrdenesPorUsuario devuelve el historial de compras de un cliente (Servicio Web 7)
func ListarOrdenesPorUsuario(usuarioID int) ([]Orden, error) {
	filas, err := db.DB.Query("SELECT id, usuario_id, total, estado, metodo_pago, fecha_creacion FROM orden WHERE usuario_id = $1", usuarioID)
	if err != nil {
		return nil, errors.New("error al listar las órdenes")
	}
	defer filas.Close()

	var ordenes []Orden
	for filas.Next() {
		var o Orden
		if err := filas.Scan(&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.MetodoPago, &o.FechaCreacion); err != nil {
			return nil, errors.New("error al escanear los datos de la orden")
		}
		ordenes = append(ordenes, o)
	}
	return ordenes, nil
}

func UsuarioTieneOrdenCancelada(usuarioID int) (bool, error) {
	var existe bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM orden WHERE usuario_id = $1 AND estado = 'cancelado')", usuarioID).Scan(&existe)
	if err != nil {
		return false, errors.New("error al validar orden cancelada")
	}
	return existe, nil
}

// ListarOrdenes devuelve todas las ordenes para el panel administrativo.
func ListarOrdenes() ([]Orden, error) {
	filas, err := db.DB.Query("SELECT id, usuario_id, total, estado, metodo_pago, fecha_creacion FROM orden ORDER BY fecha_creacion DESC")
	if err != nil {
		return nil, errors.New("error al listar las ordenes")
	}
	defer filas.Close()

	var ordenes []Orden
	for filas.Next() {
		var o Orden
		if err := filas.Scan(&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.MetodoPago, &o.FechaCreacion); err != nil {
			return nil, errors.New("error al escanear los datos de la orden")
		}
		ordenes = append(ordenes, o)
	}
	if err := filas.Err(); err != nil {
		return nil, errors.New("error al recorrer las ordenes")
	}
	return ordenes, nil
}

func ListarProductosVendidos(ascendente bool, limite int) ([]ProductoVendido, error) {
	if limite <= 0 {
		limite = 5
	}
	direccion := "DESC"
	if ascendente {
		direccion = "ASC"
	}

	query := `
		SELECT p.id, p.nombre, SUM(d.cantidad)::int AS vendidos
		FROM detalle_orden d
		INNER JOIN orden o ON o.id = d.orden_id AND o.estado = 'pagado'
		INNER JOIN producto p ON p.id = d.producto_id
		GROUP BY p.id, p.nombre
		ORDER BY vendidos ` + direccion + `, p.nombre ASC
		LIMIT $1`

	filas, err := db.DB.Query(query, limite)
	if err != nil {
		if strings.Contains(err.Error(), "detalle_orden") {
			return []ProductoVendido{}, nil
		}
		return nil, errors.New("error al listar productos vendidos")
	}
	defer filas.Close()

	var productos []ProductoVendido
	for filas.Next() {
		var p ProductoVendido
		if err := filas.Scan(&p.ID, &p.Nombre, &p.Cantidad); err != nil {
			return nil, errors.New("error al leer productos vendidos")
		}
		productos = append(productos, p)
	}
	if err := filas.Err(); err != nil {
		return nil, errors.New("error al recorrer productos vendidos")
	}
	return productos, nil
}

// ActualizarEstadoOrden permite al admin cambiar la orden a "Enviada" o "Cancelada"
func ActualizarEstadoOrden(ordenID int, nuevoEstado string) error {
	_, err := db.DB.Exec("UPDATE orden SET estado = $1 WHERE id = $2", nuevoEstado, ordenID)
	if err != nil {
		return errors.New("error al actualizar el estado de la orden")
	}
	return nil
}

func CancelarOrdenPendiente(ordenID int) error {
	resultado, err := db.DB.Exec("UPDATE orden SET estado = 'cancelado' WHERE id = $1 AND estado = 'pendiente'", ordenID)
	if err != nil {
		return errors.New("error al cancelar la orden")
	}
	filas, err := resultado.RowsAffected()
	if err != nil {
		return errors.New("error al validar la cancelacion")
	}
	if filas == 0 {
		return errors.New("solo se pueden cancelar ordenes pendientes")
	}
	return nil
}

// Funcion para obtener el detalle de una orden específica, incluyendo los productos asociados
func ObtenerDetalleOrden(ordenID int) ([]ItemCarrito, error) {
	query := `
		SELECT 
			i.id, i.producto_id, i.cantidad, i.precio,
			p.nombre, p.imagen_url
		FROM detalle_orden i
		INNER JOIN producto p ON i.producto_id = p.id
		WHERE i.orden_id = $1
	`
	filas, err := db.DB.Query(query, ordenID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var detalles []ItemCarrito
	for filas.Next() {
		var item ItemCarrito

		if err := filas.Scan(&item.ID, &item.ProductoTallaID, &item.Cantidad, &item.PrecioUnitario, &item.ProductoNombre, &item.ProductoURL); err != nil {
			return nil, err
		}
		detalles = append(detalles, item)
	}
	return detalles, nil
}
