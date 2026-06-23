package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// CategoriaRepository gestiona todas las operaciones SQL sobre la tabla categoria
type CategoriaRepository struct{}

// NuevoCategoriaRepository construye una instancia del repositorio de categorías
func NuevoCategoriaRepository() *CategoriaRepository {
	return &CategoriaRepository{}
}

// Insertar persiste una nueva categoría en la base de datos
func (r *CategoriaRepository) Insertar(c *models.Categoria) error {
	query := `
		INSERT INTO categoria (nombre, descripcion)
		VALUES ($1, $2)
		RETURNING id
	`
	return db.DB.QueryRow(query, c.Nombre, c.Descripcion).Scan(&c.ID)
}

// Actualizar modifica el nombre y descripción de una categoría existente
func (r *CategoriaRepository) Actualizar(c *models.Categoria) error {
	query := `
		UPDATE categoria
		SET nombre = $1, descripcion = $2
		WHERE id = $3
	`
	resultado, err := db.DB.Exec(query, c.Nombre, c.Descripcion, c.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("categoría no encontrada")
	}
	return nil
}

// Eliminar borra una categoría de la base de datos por su ID
func (r *CategoriaRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(`DELETE FROM categoria WHERE id = $1`, id)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("categoría no encontrada")
	}
	return nil
}

// BuscarPorID recupera una categoría según su identificador único
func (r *CategoriaRepository) BuscarPorID(id int) (*models.Categoria, bool) {
	query := `SELECT id, nombre, descripcion FROM categoria WHERE id = $1`
	c := &models.Categoria{}
	err := db.DB.QueryRow(query, id).Scan(&c.ID, &c.Nombre, &c.Descripcion)
	if err != nil {
		return nil, false
	}
	return c, true
}

// ListarTodas recupera todas las categorías ordenadas por identificador
func (r *CategoriaRepository) ListarTodas() ([]models.Categoria, error) {
	filas, err := db.DB.Query(
		`SELECT id, nombre, descripcion FROM categoria ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var categorias []models.Categoria

	for filas.Next() {
		var c models.Categoria
		if err := filas.Scan(&c.ID, &c.Nombre, &c.Descripcion); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}

	return categorias, filas.Err()
}
