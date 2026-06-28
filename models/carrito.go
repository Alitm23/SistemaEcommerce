package models

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type Carrito struct {
	ID        int    `json:"id"`
	UsuarioID int    `json:"usuario_id"`
	Estado    string `json:"estado"`
}

type ItemCarrito struct {
	ID              int     `json:"id"`
	CarritoID       int     `json:"carrito_id"`
	ProductoTallaID int     `json:"producto_talla_id"`
	Cantidad        int     `json:"cantidad"`
	PrecioUnitario  float64 `json:"precio_unitario"`

	ProductoNombre string `json:"producto_nombre"`
	ProductoURL    string `json:"producto_url"`
	TallaNombre    string `json:"talla_nombre"`
}

// CanalAnaliticaCarrito registra qué joyas son las más añadidas que no ocurren al mismo tiempo que la inserción, para no bloquear la operación de agregar al carrito
var CanalAnalisisCarrito = make(chan int, 100)

// GestorAnalitica procesa las métricas de carritos en segundo plano
func GestorAnalisis() {
	go func() {
		for productoTallaID := range CanalAnalisisCarrito {
			log.Printf("Métrica asíncrona: Producto/Talla ID %d añadido a un carrito", productoTallaID)

		}
	}()
}

// ObtenerCarritoActivo busca o crea un carrito para el usuario
func ObtenerCarritoActivo(usuarioID int) (int, error) {
	var carritoID int
	err := db.DB.QueryRow("SELECT id FROM carrito WHERE usuario_id = $1 AND estado = 'activo' LIMIT 1", usuarioID).Scan(&carritoID)

	if err == sql.ErrNoRows {
		// Si no tiene, se lo crea
		err = db.DB.QueryRow("INSERT INTO carrito (usuario_id, estado) VALUES ($1, 'activo') RETURNING id", usuarioID).Scan(&carritoID)
	}
	if err != nil {
		return 0, errors.New("error al obtener o crear el carrito")
	}
	return carritoID, nil
}

// ObtenerItemsCarrito recupera los items con los nombres gracias al JOIN (Servicio Web 4)
func ObtenerItemsCarrito(carritoID int) ([]ItemCarrito, error) {
	query := `
		SELECT i.id, i.carrito_id, i.producto_talla_id, i.cantidad, i.precio_unitario,
			p.nombre, p.imagen_url, t.talla
		FROM item_carrito i
		INNER JOIN producto_talla t ON i.producto_talla_id = t.id
		INNER JOIN producto p ON t.producto_id = p.id
		WHERE i.carrito_id = $1
	`
	filas, err := db.DB.Query(query, carritoID)
	if err != nil {
		return nil, errors.New("error al obtener los items del carrito")
	}
	defer filas.Close()

	var items []ItemCarrito
	for filas.Next() {
		var item ItemCarrito
		if err := filas.Scan(&item.ID, &item.CarritoID, &item.ProductoTallaID, &item.Cantidad, &item.PrecioUnitario, &item.ProductoNombre, &item.ProductoURL, &item.TallaNombre); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// AgregarItem inserta una joya al carrito y envía el dato al canal de analítica (Servicio Web 6)
func AgregarItem(item ItemCarrito) error {
	if item.CarritoID <= 0 {
		return errors.New("carrito invalido")
	}
	if item.ProductoTallaID <= 0 {
		return errors.New("producto_talla_id invalido")
	}
	if item.Cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("SELECT pg_advisory_xact_lock($1, $2)", item.CarritoID, item.ProductoTallaID); err != nil {
		return err
	}

	var precioUnitario float64
	var stockDisponible int
	err = tx.QueryRow(`
		SELECT p.precio, pt.stock
		FROM producto_talla pt
		INNER JOIN producto p ON pt.producto_id = p.id
		WHERE pt.id = $1`,
		item.ProductoTallaID,
	).Scan(&precioUnitario, &stockDisponible)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("producto o talla no encontrado")
		}
		return err
	}

	var itemID int
	var cantidadActual int
	err = tx.QueryRow(`
		SELECT id, cantidad
		FROM item_carrito
		WHERE carrito_id = $1 AND producto_talla_id = $2
		FOR UPDATE`,
		item.CarritoID, item.ProductoTallaID,
	).Scan(&itemID, &cantidadActual)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	nuevaCantidad := item.Cantidad
	if err == nil {
		nuevaCantidad += cantidadActual
	}
	if nuevaCantidad > stockDisponible {
		return errors.New("stock insuficiente")
	}

	if itemID > 0 {
		_, err = tx.Exec(`
			UPDATE item_carrito
			SET cantidad = $1, precio_unitario = $2
			WHERE id = $3`,
			nuevaCantidad, precioUnitario, itemID,
		)
	} else {
		_, err = tx.Exec(`
			INSERT INTO item_carrito (carrito_id, producto_talla_id, cantidad, precio_unitario)
			VALUES ($1, $2, $3, $4)`,
			item.CarritoID, item.ProductoTallaID, item.Cantidad, precioUnitario,
		)
	}
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	CanalAnalisisCarrito <- item.ProductoTallaID
	return nil
}

// ActualizarCantidadItem cambia la cantidad de un item en el carrito (Servicio Web 6)
func ActualizarCantidadItem(itemID int, nuevaCantidad int) error {
	_, err := db.DB.Exec("UPDATE item_carrito SET cantidad = $1 WHERE id = $2", nuevaCantidad, itemID)
	return err
}

// EliminarItem borra un item del carrito (Servicio Web 6)
func EliminarItem(itemID int) error {
	_, err := db.DB.Exec("DELETE FROM item_carrito WHERE id = $1", itemID)
	return err
}
