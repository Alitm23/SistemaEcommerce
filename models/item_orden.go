package models

// Registra lo que se compró y a qué precio
type ControlItemOrden interface {
	AgregarAOrden() error
	Eliminar() error
}

type ItemOrden struct {
	ID              int
	OrdenID         int
	ProductoTallaID int
	Cantidad        int
	PrecioCompra    float64
}
