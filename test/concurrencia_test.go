package test

import (
	"errors"
	"sync"
	"testing"
)

type ordenConcurrente struct {
	ID     int
	Estado string
}

type tiendaConcurrente struct {
	mu        sync.Mutex
	stock     int
	siguiente int
	ordenes   []ordenConcurrente
}

func (t *tiendaConcurrente) crearOrden(cantidad int) (ordenConcurrente, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if cantidad <= 0 {
		return ordenConcurrente{}, errors.New("cantidad invalida")
	}
	if t.stock < cantidad {
		return ordenConcurrente{}, errors.New("stock insuficiente")
	}

	t.stock -= cantidad
	t.siguiente++
	orden := ordenConcurrente{
		ID:     t.siguiente,
		Estado: "pagado",
	}
	t.ordenes = append(t.ordenes, orden)
	return orden, nil
}

func (t *tiendaConcurrente) procesarOrden(id int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, orden := range t.ordenes {
		if orden.ID == id {
			if orden.Estado == "cancelado" {
				return errors.New("no se pueden procesar ordenes canceladas")
			}
			return nil
		}
	}
	return errors.New("orden no encontrada")
}

func TestConcurrenciaStockYOrdenes(t *testing.T) {
	tienda := tiendaConcurrente{stock: 10}
	const usuarios = 25

	var wg sync.WaitGroup
	errores := make(chan error, usuarios)

	for i := 0; i < usuarios; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := tienda.crearOrden(1)
			errores <- err
		}()
	}

	wg.Wait()
	close(errores)

	var ordenesCreadas int
	var stockInsuficiente int
	for err := range errores {
		if err == nil {
			ordenesCreadas++
			continue
		}
		if err.Error() == "stock insuficiente" {
			stockInsuficiente++
		}
	}

	if tienda.stock < 0 {
		t.Fatalf("el stock no debe quedar negativo, stock final: %d", tienda.stock)
	}
	if tienda.stock != 0 {
		t.Fatalf("se esperaba stock final 0, se obtuvo %d", tienda.stock)
	}
	if ordenesCreadas != 10 {
		t.Fatalf("se esperaban 10 ordenes creadas, se crearon %d", ordenesCreadas)
	}
	if stockInsuficiente != usuarios-10 {
		t.Fatalf("se esperaban %d rechazos por stock insuficiente, se obtuvieron %d", usuarios-10, stockInsuficiente)
	}
	if len(tienda.ordenes) != ordenesCreadas {
		t.Fatalf("el total de ordenes guardadas no coincide: %d", len(tienda.ordenes))
	}
}

func TestConcurrenciaNoProcesaOrdenCancelada(t *testing.T) {
	tienda := tiendaConcurrente{
		stock: 5,
		ordenes: []ordenConcurrente{
			{ID: 1, Estado: "cancelado"},
		},
	}
	const intentos = 12

	var wg sync.WaitGroup
	errores := make(chan error, intentos)

	for i := 0; i < intentos; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errores <- tienda.procesarOrden(1)
		}()
	}

	wg.Wait()
	close(errores)

	for err := range errores {
		if err == nil {
			t.Fatal("no se debe permitir procesar una orden cancelada")
		}
		if err.Error() != "no se pueden procesar ordenes canceladas" {
			t.Fatalf("error inesperado: %v", err)
		}
	}
	if tienda.ordenes[0].Estado != "cancelado" {
		t.Fatalf("la orden cancelada no debe cambiar de estado, estado final: %s", tienda.ordenes[0].Estado)
	}
}
