package handlers

import (
	"html/template"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/utils"
)

func renderTemplate(w http.ResponseWriter, status int, files []string, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Error al cargar la pagina", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "Error al renderizar la pagina", http.StatusInternalServerError)
	}
}

func usuarioIDDesdeCookie(r *http.Request) (int, bool) {
	reclamos, ok := reclamosDesdeCookie(r)
	if !ok {
		return 0, false
	}
	return reclamos.UsuarioID, true
}

func reclamosDesdeCookie(r *http.Request) (*utils.ReclamosEcommerce, bool) {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		return nil, false
	}
	reclamos, err := utils.ValidarToken(cookie.Value)
	if err != nil {
		return nil, false
	}
	return reclamos, true
}

func redirigirSiNoAutenticado(w http.ResponseWriter, r *http.Request) (int, bool) {
	usuarioID, ok := usuarioIDDesdeCookie(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0, false
	}
	return usuarioID, true
}

func redirigirSiNoAdmin(w http.ResponseWriter, r *http.Request) bool {
	if rol, ok := r.Context().Value("rol").(string); ok {
		if rol == "admin" {
			return true
		}
		http.Error(w, "Acceso denegado: requiere rol administrador", http.StatusForbidden)
		return false
	}

	reclamos, ok := reclamosDesdeCookie(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	if reclamos.Rol != "admin" {
		http.Error(w, "Acceso denegado: requiere rol administrador", http.StatusForbidden)
		return false
	}
	return true
}
