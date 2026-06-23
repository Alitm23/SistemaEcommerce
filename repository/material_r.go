package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// MaterialRepository gestiona todas las operaciones SQL sobre la tabla material
type MaterialRepository struct{}

// NuevoMaterialRepository construye una instancia del repositorio de materiales
func NuevoMaterialRepository() *MaterialRepository {
	return &MaterialRepository{}
}

// Insertar persiste un nuevo material en la base de datos
func (r *MaterialRepository) Insertar(m *models.Material) error {
	query := `INSERT INTO material (nombre) VALUES ($1) RETURNING id`
	return db.DB.QueryRow(query, m.Nombre).Scan(&m.ID)
}

// Actualizar modifica el nombre de un material existente
func (r *MaterialRepository) Actualizar(m *models.Material) error {
	resultado, err := db.DB.Exec(
		`UPDATE material SET nombre = $1 WHERE id = $2`,
		m.Nombre, m.ID,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("material no encontrado")
	}
	return nil
}

// Eliminar borra un material de la base de datos por su ID
func (r *MaterialRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(`DELETE FROM material WHERE id = $1`, id)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("material no encontrado")
	}
	return nil
}

// BuscarPorID recupera un material según su identificador único
func (r *MaterialRepository) BuscarPorID(id int) (*models.Material, bool) {
	query := `SELECT id, nombre FROM material WHERE id = $1`
	m := &models.Material{}
	err := db.DB.QueryRow(query, id).Scan(&m.ID, &m.Nombre)
	if err != nil {
		return nil, false
	}
	return m, true
}

// ListarTodos recupera todos los materiales ordenados por identificador
func (r *MaterialRepository) ListarTodos() ([]models.Material, error) {
	filas, err := db.DB.Query(`SELECT id, nombre FROM material ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var materiales []models.Material

	for filas.Next() {
		var m models.Material
		if err := filas.Scan(&m.ID, &m.Nombre); err != nil {
			return nil, err
		}
		materiales = append(materiales, m)
	}

	return materiales, filas.Err()
}
