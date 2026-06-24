package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Store de sesiones (debe inicializarse con una clave secreta en main)
var Store *sessions.CookieStore

// UsuarioService instancia global (se puede inyectar)
var usuarioService = services.NuevoUsuarioService()

// MostrarLoginRegistro renderiza la vista con los formularios duales
func MostrarLoginRegistro(w http.ResponseWriter, r *http.Request) {
	// Verificar si ya hay sesión activa
	session, _ := Store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Renderizar template login.html
	tmpl := parseTemplate("login.html")
	err := tmpl.Execute(w, nil) // Puedes pasar datos como errores si usas flash messages
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ProcesarRegistro maneja el POST para crear un nuevo usuario
func ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener datos del formulario
	nombre := r.FormValue("nombre")
	apellido := r.FormValue("apellido")
	email := r.FormValue("email")
	password := r.FormValue("password")
	rol := r.FormValue("rol") // opcional, por defecto "cliente"
	direccion := r.FormValue("direccion")
	telefono := r.FormValue("telefono")

	// Llamar al servicio para registrar
	usuario, err := usuarioService.RegistrarUsuario(nombre, apellido, email, password, rol, direccion, telefono)
	if err != nil {
		// Redirigir con mensaje de error (usando flash o query param)
		http.Redirect(w, r, "/login?error="+err.Error(), http.StatusSeeOther)
		return
	}

	// Opcional: iniciar sesión automáticamente después del registro
	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["user_id"] = usuario.ID
	session.Values["user_name"] = usuario.Nombre
	session.Values["user_rol"] = usuario.Rol
	session.Save(r, w)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// ProcesarLogin valida credenciales e inicia sesión
func ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	usuario, err := usuarioService.Autenticar(email, password)
	if err != nil {
		http.Redirect(w, r, "/login?error=Credenciales inválidas", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["user_id"] = usuario.ID
	session.Values["user_name"] = usuario.Nombre
	session.Values["user_rol"] = usuario.Rol
	session.Save(r, w)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// CerrarSesion destruye la sesión actual
func CerrarSesion(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // Elimina la cookie
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// --- Handlers adicionales para gestión de usuarios (CRUD) ---

// ListarUsuarios devuelve todos los usuarios
func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := usuarioService.ListarUsuarios()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

// ObtenerUsuario devuelve un usuario por ID
func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	usuario, ok := usuarioService.BuscarPorID(id)
	if !ok {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// ActualizarUsuario maneja la actualización de datos (sin password)
func ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	nombre := r.FormValue("nombre")
	apellido := r.FormValue("apellido")
	email := r.FormValue("email")
	direccion := r.FormValue("direccion")
	telefono := r.FormValue("telefono")

	usuario, err := usuarioService.ActualizarUsuario(id, nombre, apellido, email, direccion, telefono)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// EliminarUsuario borra un usuario
func EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	err = usuarioService.EliminarUsuario(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CambiarRol actualiza el rol de un usuario
func CambiarRol(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	rol := r.FormValue("rol")
	err = usuarioService.CambiarRol(id, rol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CambiarPassword permite al usuario cambiar su contraseña
func CambiarPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	nuevaPassword := r.FormValue("nueva_password")
	err = usuarioService.CambiarPassword(id, nuevaPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// helper para parsear templates (puedes mejorarlo con caché)
func parseTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles("templates/" + name))
}
