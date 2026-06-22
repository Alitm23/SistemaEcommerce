package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// PagoRepository gestiona todas las operaciones SQL sobre la tabla pago
type PagoRepository struct{}

// NuevoPagoRepository construye una instancia del repositorio de pagos
func NuevoPagoRepository() *PagoRepository {
	return &PagoRepository{}
}

// Insertar persiste un nuevo pago en la base de datos
func (r *PagoRepository) Insertar(p *models.Pago) error {
	query := `
		INSERT INTO pago (orden_id, monto, estado)
		VALUES ($1, $2, $3)
		RETURNING id, fecha_pago
	`
	return db.DB.QueryRow(
		query, p.OrdenID, p.Monto, p.Estado,
	).Scan(&p.ID, &p.FechaPago)
}

// ActualizarEstado modifica el estado de un pago existente
func (r *PagoRepository) ActualizarEstado(p *models.Pago) error {
	resultado, err := db.DB.Exec(
		`UPDATE pago SET estado = $1 WHERE id = $2`,
		p.Estado, p.ID,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("pago no encontrado")
	}
	return nil
}

// BuscarPorOrden recupera el pago asociado a una orden específica
func (r *PagoRepository) BuscarPorOrden(ordenID int) (*models.Pago, bool) {
	query := `
		SELECT id, orden_id, monto, estado, fecha_pago
		FROM pago WHERE orden_id = $1
	`
	p := &models.Pago{}
	err := db.DB.QueryRow(query, ordenID).Scan(
		&p.ID, &p.OrdenID, &p.Monto, &p.Estado, &p.FechaPago,
	)
	if err != nil {
		return nil, false
	}
	return p, true
}

// ListarTodos recupera todos los pagos ordenados por identificador
func (r *PagoRepository) ListarTodos() ([]models.Pago, error) {
	filas, err := db.DB.Query(
		`SELECT id, orden_id, monto, estado, fecha_pago FROM pago ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var pagos []models.Pago

	for filas.Next() {
		var p models.Pago
		err := filas.Scan(
			&p.ID, &p.OrdenID, &p.Monto, &p.Estado, &p.FechaPago,
		)
		if err != nil {
			return nil, err
		}
		pagos = append(pagos, p)
	}

	return pagos, filas.Err()
}
