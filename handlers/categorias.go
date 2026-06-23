package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// CategoriaHandler expone los endpoints HTTP relacionados con la gestión de categorías
type CategoriaHandler struct {
	servicio *services.CategoriaService
}

// NuevoCategoriaHandler construye el handler inyectando el servicio correspondiente
func NuevoCategoriaHandler() *CategoriaHandler {
	return &CategoriaHandler{
		servicio: services.NuevoCategoriaService(),
	}
}

// Crear registra una nueva categoría en el sistema
func (h *CategoriaHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	categoria, err := h.servicio.CrearCategoria(datos.Nombre, datos.Descripcion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoria)
}

// ObtenerPorID recupera una categoría según su identificador único
func (h *CategoriaHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	categoria, ok := h.servicio.BuscarPorID(id)
	if !ok {
		http.Error(w, "categoría no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

// Listar recupera todas las categorías registradas en el sistema
func (h *CategoriaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	categorias, err := h.servicio.ListarCategorias()
	if err != nil {
		http.Error(w, "error al obtener categorías", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

// Actualizar modifica el nombre y descripción de una categoría existente
func (h *CategoriaHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	categoria, err := h.servicio.ActualizarCategoria(id, datos.Nombre, datos.Descripcion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

// Eliminar borra una categoría del sistema por su identificador
func (h *CategoriaHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.EliminarCategoria(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
