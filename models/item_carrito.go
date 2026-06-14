package models

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type ControlItemCarrito interface {
	AgregarAlCarrito() error
	ActualizarCantidad() error
	QuitarDelCarrito() error
}

type ItemCarrito struct {
	ID             int
	CarritoID      int
	ProductoID     int
	Cantidad       int
	PrecioUnitario float64
}

func NuevoItemCarrito(carritoID, productoID, cantidad int, precioUnitario float64) (*ItemCarrito, error) {
	if cantidad <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a cero")
	}
	if precioUnitario <= 0 {
		return nil, errors.New("el precio unitario debe ser mayor a cero")
	}
	return &ItemCarrito{
		CarritoID:      carritoID,
		ProductoID:     productoID,
		Cantidad:       cantidad,
		PrecioUnitario: precioUnitario,
	}, nil
}

func (i *ItemCarrito) Subtotal() float64 {
	return float64(i.Cantidad) * i.PrecioUnitario
}

func (i *ItemCarrito) AgregarAlCarrito() error {
	query := `
		INSERT INTO item_carrito (carrito_id, producto_id, cantidad, precio_unitario)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return db.DB.QueryRow(
		query,
		i.CarritoID,
		i.ProductoID,
		i.Cantidad,
		i.PrecioUnitario,
	).Scan(&i.ID)
}

func (i *ItemCarrito) ActualizarCantidad() error {
	if i.Cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	query := `UPDATE item_carrito SET cantidad = $1 WHERE id = $2`
	resultado, err := db.DB.Exec(query, i.Cantidad, i.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("item no encontrado")
	}
	return nil
}

func (i *ItemCarrito) QuitarDelCarrito() error {
	resultado, err := db.DB.Exec(`DELETE FROM item_carrito WHERE id = $1`, i.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("item no encontrado")
	}
	return nil
}

func ListarItemsPorCarrito(carritoID int) ([]ItemCarrito, error) {
	query := `
		SELECT id, carrito_id, producto_id, cantidad, precio_unitario
		FROM item_carrito
		WHERE carrito_id = $1
	`
	filas, err := db.DB.Query(query, carritoID)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var items []ItemCarrito

	for filas.Next() {
		var item ItemCarrito
		err := filas.Scan(
			&item.ID,
			&item.CarritoID,
			&item.ProductoID,
			&item.Cantidad,
			&item.PrecioUnitario,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, filas.Err()
}
