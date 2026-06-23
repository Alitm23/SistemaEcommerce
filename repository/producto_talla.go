package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// ProductoTallaRepository gestiona todas las operaciones SQL sobre la tabla producto_talla
type ProductoTallaRepository struct{}

// NuevoProductoTallaRepository construye una instancia del repositorio de tallas
func NuevoProductoTallaRepository() *ProductoTallaRepository {
	return &ProductoTallaRepository{}
}

// Insertar persiste una nueva talla con su stock inicial para un producto
func (r *ProductoTallaRepository) Insertar(pt *models.ProductoTalla) error {
	query := `
		INSERT INTO producto_talla (producto_id, talla, stock)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return db.DB.QueryRow(query, pt.ProductoID, pt.Talla, pt.Stock).Scan(&pt.ID)
}

// ActualizarStock modifica el stock de una talla específica
func (r *ProductoTallaRepository) ActualizarStock(pt *models.ProductoTalla) error {
	resultado, err := db.DB.Exec(
		`UPDATE producto_talla SET stock = $1 WHERE id = $2`,
		pt.Stock, pt.ID,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("talla no encontrada")
	}
	return nil
}

// Eliminar borra una talla específica de un producto por su ID
func (r *ProductoTallaRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(`DELETE FROM producto_talla WHERE id = $1`, id)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("talla no encontrada")
	}
	return nil
}

// BuscarPorID recupera una talla según su identificador único
func (r *ProductoTallaRepository) BuscarPorID(id int) (*models.ProductoTalla, bool) {
	query := `SELECT id, producto_id, talla, stock FROM producto_talla WHERE id = $1`
	pt := &models.ProductoTalla{}
	err := db.DB.QueryRow(query, id).Scan(&pt.ID, &pt.ProductoID, &pt.Talla, &pt.Stock)
	if err != nil {
		return nil, false
	}
	return pt, true
}

// ListarPorProducto recupera todas las tallas disponibles para un producto específico
func (r *ProductoTallaRepository) ListarPorProducto(productoID int) ([]models.ProductoTalla, error) {
	query := `
		SELECT id, producto_id, talla, stock
		FROM producto_talla
		WHERE producto_id = $1
		ORDER BY talla ASC
	`
	filas, err := db.DB.Query(query, productoID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var tallas []models.ProductoTalla

	for filas.Next() {
		var pt models.ProductoTalla
		if err := filas.Scan(&pt.ID, &pt.ProductoID, &pt.Talla, &pt.Stock); err != nil {
			return nil, err
		}
		tallas = append(tallas, pt)
	}

	return tallas, filas.Err()
}
