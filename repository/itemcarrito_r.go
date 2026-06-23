package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// ItemCarritoRepository gestiona las operaciones SQL sobre la tabla item_carrito
type ItemCarritoRepository struct{}

// NuevoItemCarritoRepository construye una instancia del repositorio
func NuevoItemCarritoRepository() *ItemCarritoRepository {
	return &ItemCarritoRepository{}
}

// Insertar persiste un nuevo ítem en el carrito.
// Referencia producto_talla_id para registrar la talla exacta seleccionada por el usuario.
func (r *ItemCarritoRepository) Insertar(i *models.ItemCarrito) error {
	query := `
		INSERT INTO item_carrito (carrito_id, producto_talla_id, cantidad, precio_unitario)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return db.DB.QueryRow(
		query,
		i.CarritoID, i.ProductoTallaID, i.Cantidad, i.PrecioUnitario,
	).Scan(&i.ID)
}

// ActualizarCantidad modifica la cantidad de un ítem existente en el carrito
func (r *ItemCarritoRepository) ActualizarCantidad(i *models.ItemCarrito) error {
	resultado, err := db.DB.Exec(
		`UPDATE item_carrito SET cantidad = $1 WHERE id = $2`,
		i.Cantidad, i.ID,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("ítem no encontrado")
	}
	return nil
}

// Eliminar borra un ítem del carrito por su ID
func (r *ItemCarritoRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(
		`DELETE FROM item_carrito WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("ítem no encontrado")
	}
	return nil
}

// ListarPorCarrito recupera todos los ítems asociados a un carrito específico
func (r *ItemCarritoRepository) ListarPorCarrito(carritoID int) ([]models.ItemCarrito, error) {
	query := `
		SELECT id, carrito_id, producto_talla_id, cantidad, precio_unitario
		FROM item_carrito
		WHERE carrito_id = $1
	`
	filas, err := db.DB.Query(query, carritoID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var items []models.ItemCarrito

	for filas.Next() {
		var item models.ItemCarrito
		if err := filas.Scan(
			&item.ID, &item.CarritoID, &item.ProductoTallaID,
			&item.Cantidad, &item.PrecioUnitario,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, filas.Err()
}
