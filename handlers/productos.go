package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/gorilla/mux"
)

type catalogoViewData struct {
	Productos        []models.Producto
	CategoriaActiva  int
	Busqueda         string
	Material         string
	PrecioMin        string
	PrecioMax        string
	TotalResultados  int
}

type detalleProductoViewData struct {
	Producto *models.Producto
	Tallas   []models.ProductoTalla
}

// ListarProductos maneja la solicitud para mostrar todos los productos.
func ListarProductos(w http.ResponseWriter, r *http.Request) {
	categoriaID, _ := strconv.Atoi(r.URL.Query().Get("categoria"))
	precioMin, _ := strconv.ParseFloat(r.URL.Query().Get("precio_min"), 64)
	precioMax, _ := strconv.ParseFloat(r.URL.Query().Get("precio_max"), 64)

	filtros := models.FiltrosProducto{
		Query:       r.URL.Query().Get("q"),
		Material:    r.URL.Query().Get("material"),
		CategoriaID: categoriaID,
		PrecioMin:   precioMin,
		PrecioMax:   precioMax,
	}

	productos, err := models.ListarProductosFiltrados(filtros)
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	data := catalogoViewData{
		Productos:       productos,
		CategoriaActiva: categoriaID,
		Busqueda:        filtros.Query,
		Material:        filtros.Material,
		PrecioMin:       r.URL.Query().Get("precio_min"),
		PrecioMax:       r.URL.Query().Get("precio_max"),
		TotalResultados: len(productos),
	}
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/catalogo.html"}, "base", data)
}

// DetalleProducto maneja la solicitud para mostrar los detalles de un producto especifico.
func DetalleProducto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	producto, err := models.ObtenerProductoPorID(id)
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	tallas, err := models.ListarTallasPorProducto(id)
	if err != nil {
		http.Error(w, "Error al obtener tallas", http.StatusInternalServerError)
		return
	}

	data := detalleProductoViewData{
		Producto: producto,
		Tallas:   tallas,
	}
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/detalle.html"}, "base", data)
}
