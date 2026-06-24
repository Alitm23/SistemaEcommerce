package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// HomeHandler agrupa las funciones para las vistas generales del sitio
type HomeHandler struct{}

// NuevoHomeHandler es el constructor
func NuevoHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// MostrarInicio renderiza la landing page componiendo base.html e inicio.html
func (h *HomeHandler) MostrarInicio(w http.ResponseWriter, r *http.Request) {

	// 1. Declaramos los archivos que forman esta vista completa
	archivos := []string{
		"templates/base.html",
		"templates/inicio.html",
	}

	// 2. Parseamos los archivos
	tmpl, err := template.ParseFiles(archivos...)
	if err != nil {
		// Imprimimos el error en consola para depurar si nos equivocamos de ruta
		log.Println("Error al parsear templates:", err)
		http.Error(w, "Error interno del servidor al cargar la vista", http.StatusInternalServerError)
		return
	}

	// 3. Ejecutamos la plantilla principal llamada "base"
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println("Error al ejecutar template:", err)
		http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
	}
}
