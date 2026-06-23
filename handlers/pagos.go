package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// PagoHandler expone los endpoints HTTP relacionados con la gestión de pagos
type PagoHandler struct {
	servicio *services.PagoService
}

// NuevoPagoHandler construye el handler inyectando el servicio correspondiente
func NuevoPagoHandler() *PagoHandler {
	return &PagoHandler{
		servicio: services.NuevoPagoService(),
	}
}

// Registrar crea un nuevo pago en estado pendiente para la orden indicada
func (h *PagoHandler) Registrar(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		OrdenID int     `json:"orden_id"`
		Monto   float64 `json:"monto"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	pago, err := h.servicio.RegistrarPago(datos.OrdenID, datos.Monto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pago)
}

// ActualizarEstado aplica las reglas de transición de estado del pago
func (h *PagoHandler) ActualizarEstado(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ordenID, err := strconv.Atoi(vars["ordenId"])
	if err != nil {
		http.Error(w, "identificador de orden inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Estado string `json:"estado"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.ActualizarEstado(ordenID, datos.Estado); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Anular marca el pago como fallido si aún no fue completado
func (h *PagoHandler) Anular(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ordenID, err := strconv.Atoi(vars["ordenId"])
	if err != nil {
		http.Error(w, "identificador de orden inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.AnularPago(ordenID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Listar recupera todos los pagos registrados en el sistema
func (h *PagoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	pagos, err := h.servicio.ListarPagos()
	if err != nil {
		http.Error(w, "error al obtener pagos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pagos)
}
