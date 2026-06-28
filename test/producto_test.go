package test

import (
	"testing"

	"github.com/Alitm23/SistemaEcommerce/models"
)

// TestValidacionesProducto verifica que la creación de un producto falle cuando el precio es negativo
func TestValidacionesProducto(t *testing.T) {
	t.Run("Crear con precio negativo", func(t *testing.T) {
		// Datos erroneos para el producto
		p := models.Producto{Nombre: "Joya", Precio: -10.0}

		// Ejecutar la función
		_, err := models.CrearProducto(p)

		// Verificación de resultados
		if err == nil {
			t.Errorf("Se esperaba un error por precio negativo, pero la función no devolvió nada")
		} else if err.Error() != "el precio no puede ser negativo" {
			t.Errorf("El error recibido no es el esperado: %v", err)
		} else {
			t.Log("Test exitoso: El sistema bloqueó el precio negativo correctamente")
		}
	})
	t.Run("Crear sin nombre", func(t *testing.T) {
		p := models.Producto{Nombre: "", Precio: 10.0}

		_, err := models.CrearProducto(p)

		if err == nil {
			t.Errorf("Se esperaba un error por nombre vacío, pero la función no devolvió nada")
		} else if err.Error() != "el nombre del producto es obligatorio" {
			t.Errorf("El error recibido no es el esperado: %v", err)
		} else {
			t.Log("Test exitoso: El sistema bloqueó el nombre vacío correctamente")
		}
	})
}
