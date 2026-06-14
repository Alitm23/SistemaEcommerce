package models

// Registra lo que se compró y a qué precio

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
)

type ControlItemOrden interface {
	AgregarAOrden() error
	Eliminar() error
}

type ItemOrden struct {
	ID           int
	OrdenID      int
	ProductoID   int
	Cantidad     int
	PrecioCompra float64
}

func NuevoItemOrden(ordenID, productoID, cantidad int, precioCompra float64) (*ItemOrden, error) {
	if cantidad <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a cero")
	}
	if precioCompra <= 0 {
		return nil, errors.New("el precio de compra debe ser mayor a cero")
	}
	return &ItemOrden{
		OrdenID:      ordenID,
		ProductoID:   productoID,
		Cantidad:     cantidad,
		PrecioCompra: precioCompra,
	}, nil
}

func (i *ItemOrden) Subtotal() float64 {
	return float64(i.Cantidad) * i.PrecioCompra
}

func (i *ItemOrden) AgregarAOrden() error {
	query := `
		INSERT INTO item_orden (orden_id, producto_id, cantidad, precio_compra)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return db.DB.QueryRow(
		query,
		i.OrdenID,
		i.ProductoID,
		i.Cantidad,
		i.PrecioCompra,
	).Scan(&i.ID)
}

func (i *ItemOrden) Eliminar() error {
	resultado, err := db.DB.Exec(`DELETE FROM item_orden WHERE id = $1`, i.ID)
	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if filas == 0 {
		return errors.New("item de orden no encontrado")
	}
	return nil
}

func ListarItemsPorOrden(ordenID int) ([]ItemOrden, error) {
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

	var items []ItemOrden

	for filas.Next() {
		var item ItemOrden
		err := filas.Scan(
			&item.ID,
			&item.OrdenID,
			&item.ProductoID,
			&item.Cantidad,
			&item.PrecioCompra,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, filas.Err()
}
