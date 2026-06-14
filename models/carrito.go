package models

import (
	"errors"
	"time"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type ControlCarrito interface {
	Abrir() error
	Cerrar() error
	Eliminar() error
}

type Carrito struct {
	ID            int
	UsuarioID     int
	Estado        string
	FechaApertura time.Time
}

func NuevoCarrito(usuarioID int) (*Carrito, error) {
	if usuarioID <= 0 {
		return nil, errors.New("usuario inválido")
	}
	return &Carrito{
		UsuarioID: usuarioID,
		Estado:    "activo",
	}, nil
}

// Función CambiarEstado encapsula las transiciones válidas del carrito.
func (c *Carrito) CambiarEstado(estado string) error {
	if estado != "activo" && estado != "cerrado" && estado != "abandonado" {
		return errors.New("estado inválido: debe ser 'activo', 'cerrado' o 'abandonado'")
	}
	c.Estado = estado
	return nil
}

func (c *Carrito) Abrir() error {
	query := `
		INSERT INTO carrito (usuario_id, estado)
		VALUES ($1, $2)
		RETURNING id, fecha_apertura
	`
	return db.DB.QueryRow(query, c.UsuarioID, c.Estado).Scan(&c.ID, &c.FechaApertura)
}

func (c *Carrito) Cerrar() error {
	query := `UPDATE carrito SET estado = $1 WHERE id = $2`
	resultado, err := db.DB.Exec(query, c.Estado, c.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("carrito no encontrado")
	}
	return nil
}

func (c *Carrito) Eliminar() error {
	resultado, err := db.DB.Exec(`DELETE FROM carrito WHERE id = $1`, c.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("carrito no encontrado")
	}
	return nil
}

func BuscarCarritoPorID(id int) (*Carrito, bool) {
	query := `SELECT id, usuario_id, estado, fecha_apertura FROM carrito WHERE id = $1`
	c := &Carrito{}
	err := db.DB.QueryRow(query, id).Scan(&c.ID, &c.UsuarioID, &c.Estado, &c.FechaApertura)
	if err != nil {
		return nil, false
	}
	return c, true
}

func BuscarCarritoActivoPorUsuario(usuarioID int) (*Carrito, bool) {
	query := `
		SELECT id, usuario_id, estado, fecha_apertura
		FROM carrito
		WHERE usuario_id = $1 AND estado = 'activo'
		LIMIT 1
	`
	c := &Carrito{}
	err := db.DB.QueryRow(query, usuarioID).Scan(&c.ID, &c.UsuarioID, &c.Estado, &c.FechaApertura)
	if err != nil {
		return nil, false
	}
	return c, true
}
