package models

import "time"

// Interfaz que controla las operaciones del producto implementando métodos de registras, actualizar y eliminar
type ControlProducto interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

// Producto representa la información de un producto dentro del sistema
type Producto struct {
	ID           int
	CategoriaID  int
	Nombre       string
	Descripcion  string
	Precio       float64
	Stock        int
	FechaIngreso time.Time
}
