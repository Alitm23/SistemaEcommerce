package models

import "time"

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
