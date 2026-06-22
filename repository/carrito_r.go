package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// CarritoRepository gestiona todas las operaciones SQL sobre la tabla carrito
type CarritoRepository struct{}

// NuevoCarritoRepository construye una instancia del repositorio de carritos
func NuevoCarritoRepository() *CarritoRepository {
	return &CarritoRepository{}
}

// Insertar persiste un nuevo carrito en la base de datos
func (r *CarritoRepository) Insertar(c *models.Carrito) error {
	query := `
		INSERT INTO carrito (usuario_id, estado)
		VALUES ($1, $2)
		RETURNING id, fecha_apertura
	`
	return db.DB.QueryRow(query, c.UsuarioID, c.Estado).Scan(&c.ID, &c.FechaApertura)
}

// ActualizarEstado modifica el estado de un carrito existente
func (r *CarritoRepository) ActualizarEstado(c *models.Carrito) error {
	resultado, err := db.DB.Exec(
		`UPDATE carrito SET estado = $1 WHERE id = $2`,
		c.Estado, c.ID,
	)
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

// Eliminar borra un carrito de la base de datos por su ID
func (r *CarritoRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(`DELETE FROM carrito WHERE id = $1`, id)
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

// BuscarPorID recupera un carrito según su identificador único
func (r *CarritoRepository) BuscarPorID(id int) (*models.Carrito, bool) {
	query := `
		SELECT id, usuario_id, estado, fecha_apertura
		FROM carrito WHERE id = $1
	`
	c := &models.Carrito{}
	err := db.DB.QueryRow(query, id).Scan(
		&c.ID, &c.UsuarioID, &c.Estado, &c.FechaApertura,
	)
	if err != nil {
		return nil, false
	}
	return c, true
}

// BuscarActivoPorUsuario recupera el carrito con estado 'activo' de un usuario
func (r *CarritoRepository) BuscarActivoPorUsuario(usuarioID int) (*models.Carrito, bool) {
	query := `
		SELECT id, usuario_id, estado, fecha_apertura
		FROM carrito
		WHERE usuario_id = $1 AND estado = 'activo'
		LIMIT 1
	`
	c := &models.Carrito{}
	err := db.DB.QueryRow(query, usuarioID).Scan(
		&c.ID, &c.UsuarioID, &c.Estado, &c.FechaApertura,
	)
	if err != nil {
		return nil, false
	}
	return c, true
}
