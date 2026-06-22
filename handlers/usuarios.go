package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// UsuarioHandler agrupa los handlers HTTP del módulo usuario
type UsuarioHandler struct {
	service *services.UsuarioService
}

// NuevoUsuarioHandler construye el handler inyectando el service
func NuevoUsuarioHandler() *UsuarioHandler {
	return &UsuarioHandler{
		service: services.NuevoUsuarioService(),
	}
}

// MostrarRegistro renderiza el formulario de registro
func (h *UsuarioHandler) MostrarRegistro(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/usuario/registro.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// ProcesarRegistro recibe el formulario y registra un nuevo usuario
func (h *UsuarioHandler) ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	_, err := h.service.RegistrarUsuario(
		r.FormValue("nombre"),
		r.FormValue("apellido"),
		r.FormValue("email"),
		r.FormValue("password"),
		"cliente",
		r.FormValue("direccion"),
		r.FormValue("telefono"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// MostrarLogin renderiza el formulario de inicio de sesión
func (h *UsuarioHandler) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/usuario/login.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// ProcesarLogin valida las credenciales mediante el service
func (h *UsuarioHandler) ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	_, err := h.service.Autenticar(
		r.FormValue("email"),
		r.FormValue("password"),
	)
	if err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ListarUsuarios obtiene todos los usuarios y los renderiza
func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.service.ListarUsuarios()
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/usuario/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, usuarios)
}

// MostrarEdicion carga el usuario y renderiza el formulario de edición
func (h *UsuarioHandler) MostrarEdicion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	usuario, ok := h.service.BuscarPorID(id)
	if !ok {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/usuario/editar.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, usuario)
}

// ProcesarEdicion actualiza los datos del usuario con los valores del formulario
func (h *UsuarioHandler) ProcesarEdicion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Actualizar datos generales
	_, err = h.service.ActualizarUsuario(
		id,
		r.FormValue("nombre"),
		r.FormValue("apellido"),
		r.FormValue("email"),
		r.FormValue("direccion"),
		r.FormValue("telefono"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cambiar contraseña solo si se envió un valor
	if pw := r.FormValue("password"); pw != "" {
		if err := h.service.CambiarPassword(id, pw); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Cambiar rol solo si se envió un valor
	if rol := r.FormValue("rol"); rol != "" {
		if err := h.service.CambiarRol(id, rol); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

// EliminarUsuario elimina un usuario por su ID
func (h *UsuarioHandler) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.EliminarUsuario(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
