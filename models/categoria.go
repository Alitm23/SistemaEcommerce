package models

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type Categoria struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// ListarCategorias obtiene todas las categorías para los filtros del catálogo
func ListarCategorias() ([]Categoria, error) {
	filas, err := db.DB.Query("SELECT id, nombre, descripcion FROM categoria")
	if err != nil {
		return nil, errors.New("error al listar categorías")
	}
	defer filas.Close()

	var categorias []Categoria
	for filas.Next() {
		var c Categoria
		if err := filas.Scan(&c.ID, &c.Nombre, &c.Descripcion); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}
	return categorias, nil
}

// CrearCategoria inserta una nueva clasificación de joyas
func CrearCategoria(c Categoria) (int, error) {
	var nuevoID int
	err := db.DB.QueryRow(
		"INSERT INTO categoria (nombre, descripcion) VALUES ($1, $2) RETURNING id",
		c.Nombre, c.Descripcion,
	).Scan(&nuevoID)

	if err != nil {
		return 0, errors.New("error al crear la categoría")
	}
	return nuevoID, nil
}

// ObtenerCategoriaPorID recupera los datos de una categoría específica
func ObtenerCategoriaPorID(id int) (*Categoria, error) {
	var c Categoria
	err := db.DB.QueryRow(
		"SELECT id, nombre, descripcion FROM categoria WHERE id = $1",
		id,
	).Scan(&c.ID, &c.Nombre, &c.Descripcion)

	if err != nil {
		return nil, errors.New("error al obtener la categoría")
	}
	return &c, nil
}

// ActualizarCategoria modifica el nombre o descripción de la categoría
func ActualizarCategoria(c Categoria) error {
	_, err := db.DB.Exec(
		"UPDATE categoria SET nombre = $1, descripcion = $2 WHERE id = $3",
		c.Nombre, c.Descripcion, c.ID,
	)
	if err != nil {
		return errors.New("error al actualizar la categoría")
	}
	return nil
}

// EliminarCategoria borra la categoría (fallará si hay productos usándola debido a llaves foráneas)
func EliminarCategoria(id int) error {
	_, err := db.DB.Exec("DELETE FROM categoria WHERE id = $1", id)
	if err != nil {
		return errors.New("error al eliminar la categoría")
	}
	return nil
}
