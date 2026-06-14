package models

import (
	"errors"
	"time"

	"github.com/Alitm23/SistemaEcommerce/db"
)

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

// Valida los datos recibidos y crea una instancia de Producto
func NuevoProducto(categoriaID int, nombre, descripcion string, precio float64, stock int) (*Producto, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}

	if stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	return &Producto{
		CategoriaID: categoriaID,
		Nombre:      nombre,
		Descripcion: descripcion,
		Precio:      precio,
		Stock:       stock,
	}, nil
}

// Controlar la modificación del stock y evita valores negativos
func (p *Producto) ActualizarStock(cantidad int) error {
	if p.Stock+cantidad < 0 {
		return errors.New("stock insuficiente")
	}

	p.Stock += cantidad
	return nil
}

// Función Registrar almacena un nuevo producto en la base de datos.
func (p *Producto) Registrar() error {
	query := `
		INSERT INTO producto (categoria_id, nombre, descripcion, precio, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, fecha_ingreso
	`

	return db.DB.QueryRow(
		query,
		p.CategoriaID,
		p.Nombre,
		p.Descripcion,
		p.Precio,
		p.Stock,
	).Scan(&p.ID, &p.FechaIngreso)
}

// Actualizar, modifica la información de un producto existente.
func (p *Producto) Actualizar() error {
	query := `
		UPDATE producto
		SET categoria_id = $1,
		    nombre       = $2,
			descripcion  = $3,
		    precio       = $4,
		    stock        = $5
		WHERE id = $6
	`

	resultado, err := db.DB.Exec(
		query,
		p.CategoriaID,
		p.Nombre,
		p.Descripcion,
		p.Precio,
		p.Stock,
		p.ID,
	)

	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if filas == 0 {
		return errors.New("producto no encontrado")
	}

	return nil
}

// Eliminar un producto de la base de datos
func (p *Producto) Eliminar() error {
	resultado, err := db.DB.Exec(
		`DELETE FROM producto WHERE id = $1`,
		p.ID,
	)

	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if filas == 0 {
		return errors.New("producto no encontrado")
	}

	return nil
}

// Funciones de consulta y búsqueda de productos

func BuscarProductoPorID(id int) (*Producto, bool) {
	query := `
		SELECT id, categoria_id, nombre, descripcion, precio, stock, fecha_ingreso
		FROM producto
		WHERE id = $1
	`

	p := &Producto{}

	err := db.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.CategoriaID,
		&p.Nombre,
		&p.Descripcion,
		&p.Precio,
		&p.Stock,
		&p.FechaIngreso,
	)

	if err != nil {
		return nil, false
	}

	return p, true
}

// Listar todos loa productos registrados
func ListarProductos() ([]Producto, error) {
	query := `
		SELECT id, categoria_id, nombre, descripcion, precio, stock, fecha_ingreso
		FROM producto
		ORDER BY id ASC
	`

	filas, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer filas.Close()

	var productos []Producto

	for filas.Next() {
		var p Producto

		err := filas.Scan(
			&p.ID,
			&p.CategoriaID,
			&p.Nombre,
			&p.Descripcion,
			&p.Precio,
			&p.Stock,
			&p.FechaIngreso,
		)

		if err != nil {
			return nil, err
		}

		productos = append(productos, p)
	}

	return productos, filas.Err()
}
