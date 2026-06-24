package models

import "time"

type ControlProducto interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

type Producto struct {
	ID                int
	CategoriaID       int // FK hacia Categoria (Anillos, Collares, Aretes)
	Nombre            string
	Descripcion       string
	Material          string
	Precio            float64
	FechaIngreso      time.Time
	ImagenURL         string
	TallasDisponibles []ProductoTalla
}

// ControlProductoTalla define las operaciones de gestión de tallas y stock
type ControlProductoTalla interface {
	Registrar() error
	ActualizarStock() error
	Eliminar() error
}

// ProductoTalla registra cada variante de talla disponible para un producto y almacena el stock real de esa combinación producto-talla.
type ProductoTalla struct {
	ID         int
	ProductoID int
	Talla      string
	Stock      int
}
