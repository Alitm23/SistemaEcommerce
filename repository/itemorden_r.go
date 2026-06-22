package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// ItemOrdenRepository gestiona las operaciones SQL sobre la tabla item_orden
type ItemOrdenRepository struct{}

// NuevoItemOrdenRepository construye una instancia del repositorio
func NuevoItemOrdenRepository() *ItemOrdenRepository {
	return &ItemOrdenRepository{}
}

// Insertar persiste un nuevo ítem dentro de una orden
func (r *ItemOrdenRepository) Insertar(i *models.ItemOrden) error {
	query := `
		INSERT INTO item_orden (orden_id, producto_id, cantidad, precio_compra)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return db.DB.QueryRow(
		query,
		i.OrdenID, i.ProductoID, i.Cantidad, i.PrecioCompra,
	).Scan(&i.ID)
}

// Eliminar borra un ítem de orden por su ID
func (r *ItemOrdenRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(
		`DELETE FROM item_orden WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("ítem de orden no encontrado")
	}
	return nil
}

// ListarPorOrden recupera todos los ítems asociados a una orden específica
func (r *ItemOrdenRepository) ListarPorOrden(ordenID int) ([]models.ItemOrden, error) {
	query := `
		SELECT id, orden_id, producto_id, cantidad, precio_compra
		FROM item_orden
		WHERE orden_id = $1
	`
	filas, err := db.DB.Query(query, ordenID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var items []models.ItemOrden

	for filas.Next() {
		var item models.ItemOrden
		err := filas.Scan(
			&item.ID, &item.OrdenID, &item.ProductoID,
			&item.Cantidad, &item.PrecioCompra,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, filas.Err()
}
