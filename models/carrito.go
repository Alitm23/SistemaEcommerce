package models

import "time"

type ControlCarrito interface {
	Abrir() error
	Cerrar() error
	Eliminar() error
}

type Carrito struct {
	ID            int
	UsuarioID     int
	Estado        string
	FechaApertura time.Time
}
