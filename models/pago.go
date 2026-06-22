package models

import "time"

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
