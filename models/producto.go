package models

import "time"

type ControlProducto interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

type Producto struct {
	ID           int
	CategoriaID  int // FK hacia Categoria (Anillos, Collares, Aretes)
	MaterialID   int // FK hacia Material (Plata, Baño en oro)
	Nombre       string
	Descripcion  string
	Precio       float64
	FechaIngreso time.Time
}

// ControlProductoTalla define las operaciones de gestión de tallas y stock
type ControlProductoTalla interface {
	Registrar() error
	ActualizarStock() error
	Eliminar() error
}

// ProductoTalla registra cada variante de talla disponible para un producto
// y almacena el stock real de esa combinación producto-talla.
type ProductoTalla struct {
	ID         int
	ProductoID int
	Talla      string
	Stock      int
}

// ControlMaterial define las operaciones básicas de gestión de materiales
type ControlMaterial interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

// Material representa el tipo de material de una joya (oro, plata, especiales, etc.)
type Material struct {
	ID     int
	Nombre string
}
