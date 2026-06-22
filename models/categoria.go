package models

type ControlCategoria interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

type Categoria struct {
	ID     int
	Nombre string
}
