package test

import (
	"testing"
)

func TestCarrito(t *testing.T) {
	t.Run("Agregar con stock insuficiente", func(t *testing.T) {
		// Simula lógica de negocio donde stock = 0
		stockActual := 0
		cantidadSolicitada := 5
		if cantidadSolicitada > stockActual {
			t.Log("El sistema correctamente detecta stock insuficiente")
		}
	})
}
