package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// ProductoRepository gestiona todas las operaciones SQL sobre la tabla producto
type ProductoRepository struct{}

// NuevoProductoRepository construye una instancia del repositorio de productos
func NuevoProductoRepository() *ProductoRepository {
	return &ProductoRepository{}
}

// Insertar persiste un nuevo producto en la base de datos.
// El stock ya no se registra aquí, sino en producto_talla por cada talla.
func (r *ProductoRepository) Insertar(p *models.Producto) error {
	query := `
		INSERT INTO producto (categoria_id, material_id, nombre, descripcion, precio)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, fecha_ingreso
	`
	return db.DB.QueryRow(
		query,
		p.CategoriaID, p.MaterialID, p.Nombre, p.Descripcion, p.Precio,
	).Scan(&p.ID, &p.FechaIngreso)
}

// Actualizar modifica los datos de un producto existente
func (r *ProductoRepository) Actualizar(p *models.Producto) error {
	query := `
		UPDATE producto
		SET categoria_id = $1, material_id = $2, nombre = $3,
		    descripcion = $4, precio = $5
		WHERE id = $6
	`
	resultado, err := db.DB.Exec(
		query,
		p.CategoriaID, p.MaterialID, p.Nombre, p.Descripcion, p.Precio, p.ID,
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

// Eliminar borra un producto de la base de datos por su ID.
// Las filas de producto_talla se eliminan automáticamente por ON DELETE CASCADE.
func (r *ProductoRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(`DELETE FROM producto WHERE id = $1`, id)
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

// BuscarPorID recupera un producto según su identificador único
func (r *ProductoRepository) BuscarPorID(id int) (*models.Producto, bool) {
	query := `
		SELECT id, categoria_id, material_id, nombre, descripcion, precio, fecha_ingreso
		FROM producto
		WHERE id = $1
	`
	p := &models.Producto{}
	err := db.DB.QueryRow(query, id).Scan(
		&p.ID, &p.CategoriaID, &p.MaterialID, &p.Nombre,
		&p.Descripcion, &p.Precio, &p.FechaIngreso,
	)
	if err != nil {
		return nil, false
	}
	return p, true
}

// ListarTodos recupera todos los productos ordenados por identificador
func (r *ProductoRepository) ListarTodos() ([]models.Producto, error) {
	query := `
		SELECT id, categoria_id, material_id, nombre, descripcion, precio, fecha_ingreso
		FROM producto
		ORDER BY id ASC
	`
	filas, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var productos []models.Producto

	for filas.Next() {
		var p models.Producto
		if err := filas.Scan(
			&p.ID, &p.CategoriaID, &p.MaterialID, &p.Nombre,
			&p.Descripcion, &p.Precio, &p.FechaIngreso,
		); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	return productos, filas.Err()
}

// ListarPorCategoria recupera todos los productos de una categoría específica
func (r *ProductoRepository) ListarPorCategoria(categoriaID int) ([]models.Producto, error) {
	query := `
		SELECT id, categoria_id, material_id, nombre, descripcion, precio, fecha_ingreso
		FROM producto
		WHERE categoria_id = $1
		ORDER BY id ASC
	`
	filas, err := db.DB.Query(query, categoriaID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var productos []models.Producto

	for filas.Next() {
		var p models.Producto
		if err := filas.Scan(
			&p.ID, &p.CategoriaID, &p.MaterialID, &p.Nombre,
			&p.Descripcion, &p.Precio, &p.FechaIngreso,
		); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	return productos, filas.Err()
}

// ListarPorMaterial recupera todos los productos de un material específico
func (r *ProductoRepository) ListarPorMaterial(materialID int) ([]models.Producto, error) {
	query := `
		SELECT id, categoria_id, material_id, nombre, descripcion, precio, fecha_ingreso
		FROM producto
		WHERE material_id = $1
		ORDER BY id ASC
	`
	filas, err := db.DB.Query(query, materialID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var productos []models.Producto

	for filas.Next() {
		var p models.Producto
		if err := filas.Scan(
			&p.ID, &p.CategoriaID, &p.MaterialID, &p.Nombre,
			&p.Descripcion, &p.Precio, &p.FechaIngreso,
		); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	return productos, filas.Err()
}
