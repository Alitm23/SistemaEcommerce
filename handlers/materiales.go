package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// MaterialHandler expone los endpoints HTTP relacionados con la gestión de materiales
type MaterialHandler struct {
	servicio *services.MaterialService
}

// NuevoMaterialHandler construye el handler inyectando el servicio correspondiente
func NuevoMaterialHandler() *MaterialHandler {
	return &MaterialHandler{
		servicio: services.NuevoMaterialService(),
	}
}

// Crear registra un nuevo material en el sistema
func (h *MaterialHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	material, err := h.servicio.CrearMaterial(datos.Nombre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(material)
}

// ObtenerPorID recupera un material según su identificador único
func (h *MaterialHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	material, ok := h.servicio.BuscarPorID(id)
	if !ok {
		http.Error(w, "material no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(material)
}

// Listar recupera todos los materiales registrados en el sistema
func (h *MaterialHandler) Listar(w http.ResponseWriter, r *http.Request) {
	materiales, err := h.servicio.ListarMateriales()
	if err != nil {
		http.Error(w, "error al obtener materiales", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materiales)
}

// Actualizar modifica el nombre de un material existente
func (h *MaterialHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	material, err := h.servicio.ActualizarMaterial(id, datos.Nombre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(material)
}

// Eliminar borra un material del sistema por su identificador
func (h *MaterialHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.EliminarMaterial(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
