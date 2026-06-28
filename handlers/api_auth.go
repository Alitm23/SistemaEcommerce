package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/utils"
)

// Credenciales recibe los datos del login.
type Credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ApiRegistro crea un nuevo cliente y devuelve un JWT.
func ApiRegistro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nuevoUsuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&nuevoUsuario); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos invalidos"})
		return
	}

	if err := models.RegistrarUsuario(nuevoUsuario); err != nil {
		status := http.StatusInternalServerError
		if esErrorValidacion(err) {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	usuarioRegistrado, err := models.AutenticarUsuario(nuevoUsuario.Email, nuevoUsuario.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al generar sesion"})
		return
	}

	token, err := utils.GenerarToken(usuarioRegistrado.ID, usuarioRegistrado.Rol)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al generar token"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ApiLogin verifica credenciales y devuelve el JWT.
func ApiLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds Credenciales
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON invalido"})
		return
	}
	if strings.TrimSpace(creds.Email) == "" || strings.TrimSpace(creds.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email y contrasena son obligatorios"})
		return
	}

	usuario, err := models.AutenticarUsuario(creds.Email, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Correo o contrasena incorrectos"})
		return
	}

	token, err := utils.GenerarToken(usuario.ID, usuario.Rol)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al generar el token de acceso"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func esErrorValidacion(err error) bool {
	mensaje := err.Error()
	return strings.Contains(mensaje, "obligatorio") || strings.Contains(mensaje, "obligatoria")
}
