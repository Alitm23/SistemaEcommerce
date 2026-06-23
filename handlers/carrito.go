package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// CarritoHandler expone los endpoints HTTP relacionados con la gestión del carrito de compras
type CarritoHandler struct {
	servicio *services.CarritoService
}

// NuevoCarritoHandler construye el handler inyectando el servicio correspondiente
func NuevoCarritoHandler() *CarritoHandler {
	return &CarritoHandler{
		servicio: services.NuevoCarritoService(),
	}
}

// Abrir crea un nuevo carrito activo para el usuario indicado
func (h *CarritoHandler) Abrir(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		UsuarioID int `json:"usuario_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	carrito, err := h.servicio.AbrirCarrito(datos.UsuarioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(carrito)
}

// Cerrar cambia el estado del carrito a cerrado
func (h *CarritoHandler) Cerrar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.CerrarCarrito(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AgregarItem agrega una talla de producto al carrito con su cantidad y precio
func (h *CarritoHandler) AgregarItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carritoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador de carrito inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		ProductoTallaID int     `json:"producto_talla_id"`
		Cantidad        int     `json:"cantidad"`
		PrecioUnitario  float64 `json:"precio_unitario"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	item, err := h.servicio.AgregarItem(
		carritoID, datos.ProductoTallaID,
		datos.Cantidad, datos.PrecioUnitario,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// ObtenerItems recupera todos los ítems del carrito junto con el total acumulado
func (h *CarritoHandler) ObtenerItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carritoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador de carrito inválido", http.StatusBadRequest)
		return
	}

	items, total, err := h.servicio.ObtenerItems(carritoID)
	if err != nil {
		http.Error(w, "error al obtener ítems del carrito", http.StatusInternalServerError)
		return
	}

	respuesta := struct {
		Items []interface{} `json:"items"`
		Total float64       `json:"total"`
	}{
		Total: total,
	}

	// Convertir ítems a interfaz para incluirlos en la respuesta
	for _, item := range items {
		respuesta.Items = append(respuesta.Items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

// ActualizarCantidadItem modifica la cantidad de un ítem existente en el carrito
func (h *CarritoHandler) ActualizarCantidadItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["itemId"])
	if err != nil {
		http.Error(w, "identificador de ítem inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Cantidad int `json:"cantidad"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.ActualizarCantidadItem(itemID, datos.Cantidad); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// QuitarItem elimina un ítem del carrito por su identificador
func (h *CarritoHandler) QuitarItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["itemId"])
	if err != nil {
		http.Error(w, "identificador de ítem inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.QuitarItem(itemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
