package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/services"
)

type ProductoHandler struct {
	productoService *services.ProductoService
}

func NuevoProductoHandler(ps *services.ProductoService) *ProductoHandler {
	return &ProductoHandler{
		productoService: ps,
	}
}

// ListarProductos extrae el catálogo filtrado o completo y renderiza la vista
func (h *ProductoHandler) ListarProductos(w http.ResponseWriter, r *http.Request) {
	// 1. Leemos si hay un filtro en la URL (ej: /productos?categoria=2)
	categoriaStr := r.URL.Query().Get("categoria")

	var productos []models.Producto
	var err error
	categoriaActiva := 0 // 0 significará "Todos"

	// 2. Decidimos qué pedirle a la Base de Datos
	if categoriaStr != "" {
		categoriaActiva, _ = strconv.Atoi(categoriaStr)
		productos, err = h.productoService.ListarPorCategoria(categoriaActiva)
	} else {
		productos, err = h.productoService.ListarProductos()
	}

	if err != nil {
		log.Println("Error al cargar productos:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// 3. Empaquetamos las joyas y la categoría activa
	datosPantalla := map[string]interface{}{
		"Productos":       productos,
		"CategoriaActiva": categoriaActiva,
	}

	// 4. Renderizamos
	archivos := []string{
		"templates/base.html",
		"templates/catalogo.html",
	}

	tmpl, err := template.ParseFiles(archivos...)
	if err != nil {
		log.Println("Error al parsear templates:", err)
		http.Error(w, "Error al procesar la vista", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "base", datosPantalla)
}
