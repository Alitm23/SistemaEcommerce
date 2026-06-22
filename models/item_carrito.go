package models

type ControlItemCarrito interface {
	AgregarAlCarrito() error
	ActualizarCantidad() error
	QuitarDelCarrito() error
}

type ItemCarrito struct {
	ID             int
	CarritoID      int
	ProductoID     int
	Cantidad       int
	PrecioUnitario float64
}
