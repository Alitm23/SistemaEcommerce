package handlers

import "net/http"

// MostrarInicioGetHandler renderiza la pagina de inicio.
func MostrarInicioGetHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/inicio.html"}, "base", nil)
}
