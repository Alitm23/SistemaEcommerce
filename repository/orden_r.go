package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// OrdenRepository gestiona todas las operaciones SQL sobre la tabla orden
type OrdenRepository struct{}

// NuevoOrdenRepository construye una instancia del repositorio de órdenes
func NuevoOrdenRepository() *OrdenRepository {
	return &OrdenRepository{}
}

// Insertar persiste una nueva orden en la base de datos
func (r *OrdenRepository) Insertar(o *models.Orden) error {
	query := `
		INSERT INTO orden (usuario_id, total, estado)
		VALUES ($1, $2, $3)
		RETURNING id, fecha_orden
	`
	return db.DB.QueryRow(
		query, o.UsuarioID, o.Total, o.Estado,
	).Scan(&o.ID, &o.FechaOrden)
}

// ActualizarEstado modifica el estado de una orden existente
func (r *OrdenRepository) ActualizarEstado(o *models.Orden) error {
	resultado, err := db.DB.Exec(
		`UPDATE orden SET estado = $1 WHERE id = $2`,
		o.Estado, o.ID,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("orden no encontrada")
	}
	return nil
}

// BuscarPorID recupera una orden según su identificador único
func (r *OrdenRepository) BuscarPorID(id int) (*models.Orden, bool) {
	query := `
		SELECT id, usuario_id, total, estado, fecha_orden
		FROM orden WHERE id = $1
	`
	o := &models.Orden{}
	err := db.DB.QueryRow(query, id).Scan(
		&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.FechaOrden,
	)
	if err != nil {
		return nil, false
	}
	return o, true
}

// ListarPorUsuario recupera todas las órdenes de un usuario específico
func (r *OrdenRepository) ListarPorUsuario(usuarioID int) ([]models.Orden, error) {
	query := `
		SELECT id, usuario_id, total, estado, fecha_orden
		FROM orden
		WHERE usuario_id = $1
		ORDER BY fecha_orden DESC
	`
	filas, err := db.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var ordenes []models.Orden

	for filas.Next() {
		var o models.Orden
		err := filas.Scan(
			&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.FechaOrden,
		)
		if err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}

	return ordenes, filas.Err()
}

// ListarTodas recupera todas las órdenes del sistema ordenadas por fecha
func (r *OrdenRepository) ListarTodas() ([]models.Orden, error) {
	query := `
		SELECT id, usuario_id, total, estado, fecha_orden
		FROM orden ORDER BY fecha_orden DESC
	`
	filas, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var ordenes []models.Orden

	for filas.Next() {
		var o models.Orden
		err := filas.Scan(
			&o.ID, &o.UsuarioID, &o.Total, &o.Estado, &o.FechaOrden,
		)
		if err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}

	return ordenes, filas.Err()
}
