package models

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type ControlCategoria interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

type Categoria struct {
	ID          int
	Nombre      string
	Descripcion string
}

func NuevaCategoria(nombre, descripcion string) (*Categoria, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	return &Categoria{
		Nombre:      nombre,
		Descripcion: descripcion,
	}, nil
}

func (c *Categoria) Registrar() error {
	query := `
		INSERT INTO categoria (nombre, descripcion)
		VALUES ($1, $2)
		RETURNING id
	`
	return db.DB.QueryRow(query, c.Nombre, c.Descripcion).Scan(&c.ID)
}

func (c *Categoria) Actualizar() error {
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

func (c *Categoria) Eliminar() error {
	resultado, err := db.DB.Exec(`DELETE FROM categoria WHERE id = $1`, c.ID)
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

func BuscarCategoriaPorID(id int) (*Categoria, bool) {
	query := `SELECT id, nombre, descripcion FROM categoria WHERE id = $1`
	c := &Categoria{}
	err := db.DB.QueryRow(query, id).Scan(&c.ID, &c.Nombre, &c.Descripcion)
	if err != nil {
		return nil, false
	}
	return c, true
}

func ListarCategorias() ([]Categoria, error) {
	filas, err := db.DB.Query(`SELECT id, nombre, descripcion FROM categoria ORDER BY id ASC`)
	if err != nil {
		return nil, err
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

	return categorias, filas.Err()
}
