package handlers

import (
	"html/template"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/models"
)

// Home renderiza la página principal del e-commerce: muestra el catálogo de
// productos con stock disponible junto con las categorías existentes.
func Home(w http.ResponseWriter, r *http.Request) {
	// Solo la ruta raíz se atiende como home; cualquier otra es 404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	productos, err := models.ListarProductos()
	if err != nil {
		http.Error(w, "Error al obtener los productos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Mostrar únicamente los productos con stock disponible (mayor a cero)
	disponibles := make([]models.Producto, 0, len(productos))
	for _, p := range productos {
		if p.Stock > 0 {
			disponibles = append(disponibles, p)
		}
	}

	categorias, err := models.ListarCategorias()
	if err != nil {
		http.Error(w, "Error al obtener las categorías: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Productos":  disponibles,
		"Categorias": categorias,
	})
}
