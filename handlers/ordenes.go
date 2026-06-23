package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// OrdenHandler expone los endpoints HTTP relacionados con la gestión de órdenes de compra
type OrdenHandler struct {
	servicio *services.OrdenService
}

// NuevoOrdenHandler construye el handler inyectando el servicio correspondiente
func NuevoOrdenHandler() *OrdenHandler {
	return &OrdenHandler{
		servicio: services.NuevoOrdenService(),
	}
}

// Generar crea una nueva orden en estado pendiente para el usuario indicado
func (h *OrdenHandler) Generar(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		UsuarioID int     `json:"usuario_id"`
		Total     float64 `json:"total"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	orden, err := h.servicio.GenerarOrden(datos.UsuarioID, datos.Total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orden)
}

// ObtenerPorID recupera una orden según su identificador único
func (h *OrdenHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	orden, ok := h.servicio.BuscarPorID(id)
	if !ok {
		http.Error(w, "orden no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orden)
}

// ListarPorUsuario recupera todas las órdenes de un usuario específico
func (h *OrdenHandler) ListarPorUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usuarioID, err := strconv.Atoi(vars["usuarioId"])
	if err != nil {
		http.Error(w, "identificador de usuario inválido", http.StatusBadRequest)
		return
	}

	ordenes, err := h.servicio.ListarPorUsuario(usuarioID)
	if err != nil {
		http.Error(w, "error al obtener órdenes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

// ListarTodas recupera todas las órdenes del sistema
func (h *OrdenHandler) ListarTodas(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.servicio.ListarTodas()
	if err != nil {
		http.Error(w, "error al obtener órdenes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

// ActualizarEstado aplica la máquina de estados de la orden
func (h *OrdenHandler) ActualizarEstado(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Estado string `json:"estado"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.ActualizarEstado(id, datos.Estado); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Cancelar establece la orden en estado cancelado
func (h *OrdenHandler) Cancelar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.CancelarOrden(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AgregarItem agrega una talla de producto a la orden y descuenta el stock correspondiente
func (h *OrdenHandler) AgregarItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ordenID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador de orden inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		ProductoTallaID int     `json:"producto_talla_id"`
		Cantidad        int     `json:"cantidad"`
		PrecioCompra    float64 `json:"precio_compra"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	item, err := h.servicio.AgregarItem(
		ordenID, datos.ProductoTallaID,
		datos.Cantidad, datos.PrecioCompra,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// ObtenerItems recupera todos los ítems de una orden específica
func (h *OrdenHandler) ObtenerItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ordenID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador de orden inválido", http.StatusBadRequest)
		return
	}

	items, err := h.servicio.ObtenerItems(ordenID)
	if err != nil {
		http.Error(w, "error al obtener ítems de la orden", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
