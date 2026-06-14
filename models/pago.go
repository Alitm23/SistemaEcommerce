package models

import (
	"errors"
	"time"

	"github.com/Alitm23/SistemaEcommerce/db"
)

// define las operaciones básicas para la gestión de pagos
type ControlPago interface {
	RegistrarPago() error
	ActualizarEstado() error
	CancelarPago() error
}

// Pago asociado a una orden
type Pago struct {
	ID        int
	OrdenID   int
	Monto     float64
	Estado    string
	FechaPago time.Time
}

// valida los datos recibidos y crea una instancia de pago
func NuevoPago(ordenID int, monto float64) (*Pago, error) {
	if monto <= 0 {
		return nil, errors.New("el monto debe ser mayor a cero")
	}
	return &Pago{
		OrdenID: ordenID,
		Monto:   monto,
		Estado:  "pendiente",
	}, nil
}

// controla la actualización del estado del pago aplicando las reglas definidas por el sistema
func (p *Pago) CambiarEstado(estado string) error {
	estados := map[string]bool{
		"pendiente":  true,
		"completado": true,
		"fallido":    true,
	}
	if !estados[estado] {
		return errors.New("estado inválido: debe ser 'pendiente', 'completado' o 'fallido'")
	}

	// Evitar que un pago completado vuelva al estado pendiente
	if p.Estado == "completado" && estado == "pendiente" {
		return errors.New("un pago completado no puede volver a pendiente")
	}
	p.Estado = estado
	return nil
}

func (p *Pago) RegistrarPago() error {
	query := `
		INSERT INTO pago (orden_id, monto, estado)
		VALUES ($1, $2, $3)
		RETURNING id, fecha_pago
	`
	return db.DB.QueryRow(query, p.OrdenID, p.Monto, p.Estado).Scan(&p.ID, &p.FechaPago)
}

func (p *Pago) ActualizarEstado() error {
	//actualiza el estado del pago utilizando su identificador mediante consulta sql
	query := `UPDATE pago SET estado = $1 WHERE id = $2`
	resultado, err := db.DB.Exec(query, p.Estado, p.ID)
	if err != nil {
		return err
	}

	// Verifica si la operación afectó algún registro
	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("pago no encontrado")
	}
	return nil
}

func (p *Pago) AnularPago() error {
	//// Un pago completado no puede ser anulado.
	if p.Estado == "completado" {
		return errors.New("no se puede cancelar si el pago completado")
	}
	if err := p.CambiarEstado("fallido"); err != nil {
		return err
	}
	return p.ActualizarEstado()
}

func BuscarPagoPorOrden(ordenID int) (*Pago, bool) {
	//obtiene los datos del pago correspondiente al id de la orden
	query := `
		SELECT id, orden_id, monto, estado, fecha_pago
		FROM pago WHERE orden_id = $1
	`
	//instancia vacia para datos recuperados de la consulta
	p := &Pago{}
	err := db.DB.QueryRow(query, ordenID).Scan(&p.ID, &p.OrdenID, &p.Monto, &p.Estado, &p.FechaPago)
	if err != nil {
		return nil, false
	}
	return p, true
}

func ListarPagos() ([]Pago, error) {
	query := `
		SELECT id, orden_id, monto, estado, fecha_pago
		FROM pago ORDER BY id ASC
	`
	filas, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var pagos []Pago

	for filas.Next() {
		var p Pago
		err := filas.Scan(&p.ID, &p.OrdenID, &p.Monto, &p.Estado, &p.FechaPago)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, p)
	}

	return pagos, filas.Err()
}
