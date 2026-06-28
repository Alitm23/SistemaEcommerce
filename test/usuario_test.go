package test

import (
	"testing"

	"github.com/Alitm23/SistemaEcommerce/models"
)

// TestRegistroUsuario verifica que el registro de un usuario falle cuando los campos requeridos están vacíos
func TestRegistroUsuario(t *testing.T) {
	t.Run("Registro con campos vacíos", func(t *testing.T) { //t.run permite definir subtests dentro de un test principal.
		u := models.Usuario{Nombre: "", Email: ""} // Crea un usuario con campos vacíos
		err := models.RegistrarUsuario(u)          // Llama a la función RegistrarUsuario con el usuario vacío
		if err == nil {
			t.Errorf("Se esperaba error al registrar usuario vacío, pero no hubo")
		}
	})
}
