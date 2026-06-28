package models

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Alitm23/SistemaEcommerce/db"
)

// Producto representa una joya en la tienda.
type Producto struct {
	ID          int     `json:"id"`
	CategoriaID int     `json:"categoria_id"`
	Nombre      string  `json:"nombre"`
	Material    string  `json:"material"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	ImagenURL   string  `json:"imagen_url"`
	Tallas      []ProductoTalla
}

type FiltrosProducto struct {
	Query       string
	Nombre      string
	Material    string
	CategoriaID int
	PrecioMin   float64
	PrecioMax   float64
}

type ProductoTalla struct {
	ID         int    `json:"id"`
	ProductoID int    `json:"producto_id"`
	Talla      string `json:"talla"`
	Stock      int    `json:"stock"`
}

// PeticionStock es la estructura que viaja por el canal.
type PeticionStock struct {
	ProductoID int
	Cantidad   int
}

// CanalStock es el canal global para gestionar el inventario sin bloquear peticiones.
var CanalStock = make(chan PeticionStock, 100)

// GestiondeStock se ejecuta en segundo plano escuchando el canal.
func GestiondeStock() {
	go func() {
		for peticion := range CanalStock {
			_, err := db.DB.Exec(`
				UPDATE producto_talla
				SET stock = stock - $1
				WHERE producto_id = $2 AND stock >= $1`,
				peticion.Cantidad, peticion.ProductoID,
			)

			if err != nil {
				log.Printf("Error critico en GestiondeStock para ID %d: %v", peticion.ProductoID, err)
			}
		}
	}()
}

// ListarProductos devuelve todo el catalogo de joyas.
func ListarProductos() ([]Producto, error) {
	return ListarProductosFiltrados(FiltrosProducto{})
}

// ListarProductosFiltrados devuelve el catalogo aplicando filtros opcionales.
func ListarProductosFiltrados(f FiltrosProducto) ([]Producto, error) {
	query := "SELECT id, categoria_id, nombre, material, descripcion, precio, imagen_url FROM producto"
	var condiciones []string
	var args []interface{}

	agregarFiltroTexto := func(campo, valor string) {
		valor = strings.TrimSpace(valor)
		if valor == "" {
			return
		}
		args = append(args, "%"+strings.ToLower(valor)+"%")
		condiciones = append(condiciones, "LOWER("+campo+") LIKE $"+strconv.Itoa(len(args)))
	}

	if strings.TrimSpace(f.Query) != "" {
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(f.Query))+"%")
		placeholder := "$" + strconv.Itoa(len(args))
		condiciones = append(condiciones, "(LOWER(nombre) LIKE "+placeholder+" OR LOWER(material) LIKE "+placeholder+")")
	}
	agregarFiltroTexto("nombre", f.Nombre)
	agregarFiltroTexto("material", f.Material)

	if f.CategoriaID > 0 {
		args = append(args, f.CategoriaID)
		condiciones = append(condiciones, "categoria_id = $"+strconv.Itoa(len(args)))
	}
	if f.PrecioMin > 0 {
		args = append(args, f.PrecioMin)
		condiciones = append(condiciones, "precio >= $"+strconv.Itoa(len(args)))
	}
	if f.PrecioMax > 0 {
		args = append(args, f.PrecioMax)
		condiciones = append(condiciones, "precio <= $"+strconv.Itoa(len(args)))
	}

	if len(condiciones) > 0 {
		query += " WHERE " + strings.Join(condiciones, " AND ")
	}
	query += " ORDER BY id"

	filas, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, errors.New("error al listar productos")
	}
	defer filas.Close()

	var productos []Producto
	for filas.Next() {
		var p Producto
		if err := filas.Scan(&p.ID, &p.CategoriaID, &p.Nombre, &p.Material, &p.Descripcion, &p.Precio, &p.ImagenURL); err != nil {
			return nil, errors.New("error al leer los datos del producto")
		}
		productos = append(productos, p)
	}
	if err := filas.Err(); err != nil {
		return nil, errors.New("error al recorrer la lista de productos")
	}
	return productos, nil
}

// CrearProducto inserta una nueva joya en la base de datos.
func CrearProducto(p Producto) (int, error) {
	if err := validarProducto(p); err != nil {
		return 0, err
	}

	var nuevoID int
	err := db.DB.QueryRow(`
		INSERT INTO producto (categoria_id, nombre, material, descripcion, precio, imagen_url)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		p.CategoriaID, p.Nombre, p.Material, p.Descripcion, p.Precio, p.ImagenURL,
	).Scan(&nuevoID)

	return nuevoID, err
}

func CrearTallaProducto(productoID int, talla string, stock int) error {
	if strings.TrimSpace(talla) == "" {
		talla = "Unica"
	}
	if stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}
	_, err := db.DB.Exec(
		"INSERT INTO producto_talla (producto_id, talla, stock) VALUES ($1, $2, $3)",
		productoID, talla, stock,
	)
	return err
}

// ObtenerProductoPorID devuelve los detalles de una joya especifica.
func ObtenerProductoPorID(id int) (*Producto, error) {
	var p Producto
	err := db.DB.QueryRow(
		"SELECT id, categoria_id, nombre, material, descripcion, precio, imagen_url FROM producto WHERE id = $1",
		id,
	).Scan(&p.ID, &p.CategoriaID, &p.Nombre, &p.Material, &p.Descripcion, &p.Precio, &p.ImagenURL)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("producto no encontrado")
		}
		return nil, errors.New("error al obtener el producto")
	}
	return &p, nil
}

// ActualizarProducto actualiza los detalles de una joya existente en la base de datos.
func ActualizarProducto(p Producto) (int, error) {
	if err := validarProducto(p); err != nil {
		return 0, err
	}

	resultado, err := db.DB.Exec(`
		UPDATE producto
		SET categoria_id = $1, nombre = $2, material = $3, descripcion = $4, precio = $5, imagen_url = $6
		WHERE id = $7`,
		p.CategoriaID, p.Nombre, p.Material, p.Descripcion, p.Precio, p.ImagenURL, p.ID,
	)

	if err != nil {
		return 0, errors.New("error al actualizar el producto")
	}
	filasAfectadas, err := resultado.RowsAffected()
	if err != nil {
		return 0, errors.New("error al verificar la actualizacion del producto")
	}
	if filasAfectadas == 0 {
		return 0, errors.New("producto no encontrado")
	}
	return int(filasAfectadas), nil
}

// EliminarProducto elimina una joya de la base de datos.
func EliminarProducto(id int) (int, error) {
	resultado, err := db.DB.Exec("DELETE FROM producto WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return 0, err
	}

	if filas == 0 {
		return 0, errors.New("producto no encontrado")
	}

	return int(filas), nil
}

func ListarTallasPorProducto(productoID int) ([]ProductoTalla, error) {
	filas, err := db.DB.Query(
		"SELECT id, producto_id, talla, stock FROM producto_talla WHERE producto_id = $1 ORDER BY talla",
		productoID,
	)
	if err != nil {
		return nil, errors.New("error al listar tallas")
	}
	defer filas.Close()

	var tallas []ProductoTalla
	for filas.Next() {
		var t ProductoTalla
		if err := filas.Scan(&t.ID, &t.ProductoID, &t.Talla, &t.Stock); err != nil {
			return nil, errors.New("error al leer tallas")
		}
		tallas = append(tallas, t)
	}
	if err := filas.Err(); err != nil {
		return nil, errors.New("error al recorrer tallas")
	}
	return tallas, nil
}

func ActualizarStockTalla(tallaID int, stock int) error {
	if stock < 0 {
		return errors.New("el stock no puede ser negativo")
	}
	resultado, err := db.DB.Exec("UPDATE producto_talla SET stock = $1 WHERE id = $2", stock, tallaID)
	if err != nil {
		return errors.New("error al actualizar stock")
	}
	filas, err := resultado.RowsAffected()
	if err != nil {
		return errors.New("error al verificar stock")
	}
	if filas == 0 {
		return errors.New("talla no encontrada")
	}
	return nil
}

func validarProducto(p Producto) error {
	if p.Precio < 0 {
		return errors.New("el precio no puede ser negativo")
	}
	if strings.TrimSpace(p.Nombre) == "" {
		return errors.New("el nombre del producto es obligatorio")
	}
	if strings.TrimSpace(p.Material) == "" {
		return errors.New("el material del producto es obligatorio")
	}
	return nil
}
