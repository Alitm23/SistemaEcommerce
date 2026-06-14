package models

import (
	"errors"
	"time"

	"github.com/Alitm23/SistemaEcommerce/db"
)

// interfaz con métodos para genera ordenes, actualizar estado de la orden y cancelar orden
type ControlOrden interface {
	GenerarOrden() error
	ActualizarEstado() error
	CancelarOrden() error
}

type Orden struct {
	ID         int
	UsuarioID  int
	Total      float64
	Estado     string
	FechaOrden time.Time
}

// reigstra una nueva orden
func NuevaOrden(usuarioID int, total float64) (*Orden, error) {
	if total <= 0 {
		return nil, errors.New("el total debe ser mayor a cero") //control de error
	}
	return &Orden{
		UsuarioID: usuarioID,
		Total:     total,
		Estado:    "pendiente",
	}, nil
}

func (o *Orden) CambiarEstado(estado string) error {
	estados := map[string]bool{
		"pendiente":  true,
		"procesando": true,
		"enviada":    true,
		"entregada":  true,
		"cancelada":  true,
	}
	if !estados[estado] {
		return errors.New("estado de orden inválido")
	}

	// no modificar el estado de una orden ya entregada
	if o.Estado == "entregada" {
		return errors.New("una orden entregada no puede cambiar de estado")
	}
	o.Estado = estado
	return nil
}

func (o *Orden) GenerarOrden() error {
	query := `
		INSERT INTO orden (usuario_id, total, estado)
		VALUES ($1, $2, $3)
		RETURNING id, fecha_orden
	`
	return db.DB.QueryRow(query, o.UsuarioID, o.Total, o.Estado).Scan(&o.ID, &o.FechaOrden)
}

func (o *Orden) ActualizarEstado() error {
	query := `UPDATE orden SET estado = $1 WHERE id = $2`
	resultado, err := db.DB.Exec(query, o.Estado, o.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("orden no encontrada")
	}
	return nil
}

func (o *Orden) CancelarOrden() error {
	if err := o.CambiarEstado("cancelada"); err != nil {
		return err
	}
	return o.ActualizarEstado()
}

func BuscarOrdenPorID(id int) (*Orden, bool) {
	query := `
		SELECT id, usuario_id, total, estado, fecha_orden
		FROM orden WHERE id = $1
	`
	o := &Orden{}
	err := db.DB.QueryRow(query, id).Scan(&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.FechaOrden)
	if err != nil {
		return nil, false
	}
	return o, true
}

func ListarOrdenesPorUsuario(usuarioID int) ([]Orden, error) {
	query := `
		SELECT id, usuario_id, total, estado, fecha_orden
		FROM orden
		WHERE usuario_id = $1
		ORDER BY fecha_orden DESC
	`
	filas, err := db.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var ordenes []Orden

	for filas.Next() {
		var o Orden
		err := filas.Scan(&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.FechaOrden)
		if err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}

	return ordenes, filas.Err()
}
