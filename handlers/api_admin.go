package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/gorilla/mux"
)

// ApiAdminCrearProducto inserta una nueva joya en la base de datos
func ApiAdminCrearProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nuevoProducto models.Producto //variable nuevoProducto de tipo Producto para recibir los datos del frontend
	if err := json.NewDecoder(r.Body).Decode(&nuevoProducto); err != nil {
		w.WriteHeader(http.StatusBadRequest) // Si hay un error al decodificar el JSON, respondemos con un error 400
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos del producto inválidos"})
		return
	}

	nuevoID, err := models.CrearProducto(nuevoProducto)
	if err != nil {
		w.WriteHeader(statusProductoError(err))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Producto creado exitosamente",
		"id":      nuevoID,
	})
}

// ApiAdminActualizarProducto modifica los datos de una joya existente
func ApiAdminActualizarProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID invalido"})
		return
	}

	var productoActualizado models.Producto
	if err := json.NewDecoder(r.Body).Decode(&productoActualizado); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos inválidos"})
		return
	}
	productoActualizado.ID = id

	// AHORA RECIBIMOS EL INT Y EL ERROR
	filas, err := models.ActualizarProducto(productoActualizado)

	if err != nil {
		// Si el error es "producto no encontrado", enviamos un 404
		if err.Error() == "producto no encontrado" {
			w.WriteHeader(http.StatusNotFound)
		} else if esErrorValidacionProducto(err) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Producto actualizado correctamente",
		"filas":   filas, // Aprovechamos la información que el modelo nos dio
	})
}

func statusProductoError(err error) int {
	if esErrorValidacionProducto(err) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func esErrorValidacionProducto(err error) bool {
	switch err.Error() {
	case "el precio no puede ser negativo",
		"el nombre del producto es obligatorio",
		"el material del producto es obligatorio":
		return true
	default:
		return false
	}
}

// ApiAdminEliminarProducto borra una joya del catálogo
func ApiAdminEliminarProducto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}
	filas, err := models.EliminarProducto(id)
	if err != nil {
		if err.Error() == "producto no encontrado" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Producto eliminado exitosamente",
		"filas":   filas,
	})
}
