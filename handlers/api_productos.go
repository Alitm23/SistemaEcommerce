package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/gorilla/mux"
)

// ApiListarProductos devuelve el catalogo en JSON con filtros opcionales.
func ApiListarProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filtros, err := filtrosProductoDesdeQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	productos, err := models.ListarProductosFiltrados(filtros)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al cargar el catalogo de joyas"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productos)
}

// ApiDetalleProducto devuelve la informacion de una sola joya.
func ApiDetalleProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID de producto invalido"})
		return
	}

	producto, err := models.ObtenerProductoPorID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Joya no encontrada"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(producto)
}

// ApiBuscarProductos busca joyas por nombre o material.
func ApiBuscarProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "El parametro de busqueda esta vacio"})
		return
	}

	resultados, err := models.ListarProductosFiltrados(models.FiltrosProducto{Query: query})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al realizar la busqueda"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultados)
}

func filtrosProductoDesdeQuery(r *http.Request) (models.FiltrosProducto, error) {
	q := r.URL.Query()
	filtros := models.FiltrosProducto{
		Query:    q.Get("q"),
		Nombre:   q.Get("nombre"),
		Material: q.Get("material"),
	}

	if categoriaID := strings.TrimSpace(q.Get("categoria_id")); categoriaID != "" {
		id, err := strconv.Atoi(categoriaID)
		if err != nil || id <= 0 {
			return filtros, errParametroInvalido("categoria_id")
		}
		filtros.CategoriaID = id
	}

	if precioMin := strings.TrimSpace(q.Get("precio_min")); precioMin != "" {
		precio, err := strconv.ParseFloat(precioMin, 64)
		if err != nil || precio < 0 {
			return filtros, errParametroInvalido("precio_min")
		}
		filtros.PrecioMin = precio
	}

	if precioMax := strings.TrimSpace(q.Get("precio_max")); precioMax != "" {
		precio, err := strconv.ParseFloat(precioMax, 64)
		if err != nil || precio < 0 {
			return filtros, errParametroInvalido("precio_max")
		}
		filtros.PrecioMax = precio
	}

	if filtros.PrecioMax > 0 && filtros.PrecioMin > filtros.PrecioMax {
		return filtros, errParametroInvalido("precio_min")
	}

	return filtros, nil
}

func errParametroInvalido(nombre string) error {
	return &parametroInvalido{nombre: nombre}
}

type parametroInvalido struct {
	nombre string
}

func (e *parametroInvalido) Error() string {
	return "parametro invalido: " + e.nombre
}
