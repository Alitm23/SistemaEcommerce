package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// UsuarioHandler expone los endpoints HTTP relacionados con la gestión de usuarios
type UsuarioHandler struct {
	servicio *services.UsuarioService
}

// NuevoUsuarioHandler construye el handler inyectando el servicio correspondiente
func NuevoUsuarioHandler() *UsuarioHandler {
	return &UsuarioHandler{
		servicio: services.NuevoUsuarioService(),
	}
}

// Registrar crea un nuevo usuario con los datos enviados en el cuerpo de la petición
func (h *UsuarioHandler) Registrar(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		Nombre    string `json:"nombre"`
		Apellido  string `json:"apellido"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Rol       string `json:"rol"`
		Direccion string `json:"direccion"`
		Telefono  string `json:"telefono"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.RegistrarUsuario(
		datos.Nombre, datos.Apellido, datos.Email,
		datos.Password, datos.Rol, datos.Direccion, datos.Telefono,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// Autenticar verifica las credenciales y retorna el usuario si son válidas
func (h *UsuarioHandler) Autenticar(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.Autenticar(datos.Email, datos.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// ObtenerPorID recupera un usuario según su identificador único
func (h *UsuarioHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	usuario, ok := h.servicio.BuscarPorID(id)
	if !ok {
		http.Error(w, "usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// Listar recupera todos los usuarios registrados en el sistema
func (h *UsuarioHandler) Listar(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.servicio.ListarUsuarios()
	if err != nil {
		http.Error(w, "error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

// Actualizar modifica los datos de un usuario existente
func (h *UsuarioHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Nombre    string `json:"nombre"`
		Apellido  string `json:"apellido"`
		Email     string `json:"email"`
		Direccion string `json:"direccion"`
		Telefono  string `json:"telefono"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	usuario, err := h.servicio.ActualizarUsuario(
		id, datos.Nombre, datos.Apellido,
		datos.Email, datos.Direccion, datos.Telefono,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// CambiarPassword actualiza la contraseña de un usuario existente
func (h *UsuarioHandler) CambiarPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		NuevaPassword string `json:"nueva_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.CambiarPassword(id, datos.NuevaPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CambiarRol actualiza el rol de un usuario existente
func (h *UsuarioHandler) CambiarRol(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Rol string `json:"rol"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.CambiarRol(id, datos.Rol); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Eliminar borra un usuario del sistema por su identificador
func (h *UsuarioHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.EliminarUsuario(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
