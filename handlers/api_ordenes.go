package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/gorilla/mux"
)

// ApiListarOrdenes devuelve el historial de compras del usuario autenticado (Servicio Web 7)
func ApiListarOrdenes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usuarioID, ok := r.Context().Value("usuario_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado"})
		return
	}

	ordenes, err := models.ListarOrdenesPorUsuario(usuarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener el historial de órdenes"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ordenes)
}

// ApiDetalleOrden muestra qué items contiene una orden específica
func ApiDetalleOrden(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ordenID, _ := strconv.Atoi(vars["id"])

	detalles, err := models.ObtenerDetalleOrden(ordenID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Detalle no encontrado"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(detalles)
}

// ApiProcesarPago gestiona el pago y activa la notificación asíncrona (Servicio Web 7)
func ApiProcesarPago(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usuarioID, _ := r.Context().Value("usuario_id").(int)

	var p struct { // struct encarhados de recibir datos del frontend
		Total      float64 `json:"total"`
		MetodoPago string  `json:"metodo_pago"`
	}
	json.NewDecoder(r.Body).Decode(&p) // decodifica el JSON recibido del frontend y lo almacena en la variable p

	nuevaOrden := models.Orden{
		UsuarioID:  usuarioID,
		Total:      p.Total,
		Estado:     "pendiente",
		MetodoPago: p.MetodoPago,
	}

	id, err := models.CrearOrden(nuevaOrden)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear orden"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"orden_id": id, "status": 201})
}
