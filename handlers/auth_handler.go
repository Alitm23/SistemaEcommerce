package handlers

import (
	"html/template"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/Alitm23/SistemaEcommerce/utils"
)

type AuthHandler struct {
	usuarioService *services.UsuarioService
}

func NuevoAuthHandler(us *services.UsuarioService) *AuthHandler {
	return &AuthHandler{
		usuarioService: us,
	}
}

// MostrarLogin renderiza el formulario dual (Login / Registro)
func (h *AuthHandler) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	archivos := []string{"templates/base.html", "templates/login.html"}

	tmpl, err := template.ParseFiles(archivos...)
	if err != nil {
		http.Error(w, "Error al cargar la vista", http.StatusInternalServerError)
		return
	}

	// Enviamos un mapa vacío para que la plantilla no falle al buscar errores
	tmpl.ExecuteTemplate(w, "base", map[string]string{})
}

// ProcesarLogin recibe los datos del formulario y crea la sesión
func (h *AuthHandler) ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	usuario, err := h.usuarioService.Autenticar(email, password)
	if err != nil {
		archivos := []string{
			"templates/base.html",
			"templates/login.html",
		}
		tmpl, _ := template.ParseFiles(archivos...)

		datosPantalla := map[string]string{
			"Error": "Credenciales inválidas. Verifica tu correo y contraseña",
			"Email": email, // Devolvemos el email para que no tenga que volver a escribirlo
		}

		tmpl.ExecuteTemplate(w, "base", datosPantalla)
		return
	}

	// Si todo sale bien, creamos la sesión
	session, _ := utils.Store.Get(r, "sesion-ecommerce")
	session.Values["usuario_id"] = usuario.ID
	session.Values["usuario_nombre"] = usuario.Nombre
	session.Save(r, w)

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

// ProcesarRegistro inserta un nuevo usuario y lo redirige al login
func (h *AuthHandler) ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	nombre := r.FormValue("nombre")
	apellido := r.FormValue("apellido")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Para simplificar, dejamos dirección y teléfono en blanco por ahora
	_, err := h.usuarioService.RegistrarUsuario(nombre, apellido, email, password, "cliente", "", "")
	if err != nil {
		http.Error(w, "Error al registrar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirigimos al login para que inicie sesión con su nueva cuenta
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// CerrarSesion destruye la cookie
func (h *AuthHandler) CerrarSesion(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "sesion-ecommerce")

	// Revocamos la sesión
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
