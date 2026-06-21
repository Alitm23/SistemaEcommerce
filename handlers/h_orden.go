package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Lista las órdenes de un usuario específico
func ListarOrdenes(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := strconv.Atoi(mux.Vars(r)["usuarioID"])
	if err != nil {
		http.Error(w, "Usuario inválido", http.StatusBadRequest)
		return
	}

	ordenes, err := models.ListarOrdenesPorUsuario(usuarioID)
	if err != nil {
		http.Error(w, "Error al obtener las órdenes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/orden/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, ordenes)
}

// Muestra el detalle de una orden junto con sus ítems
func VerOrden(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	orden, ok := models.BuscarOrdenPorID(id)
	if !ok {
		http.Error(w, "Orden no encontrada", http.StatusNotFound)
		return
	}

	items, err := models.ListarItemsPorOrden(orden.ID)
	if err != nil {
		http.Error(w, "Error al obtener los ítems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/orden/detalle.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Orden": orden,
		"Items": items,
	})
}

// Genera una nueva orden para un usuario a partir del total enviado
func CrearOrden(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	usuarioID, err := strconv.Atoi(r.FormValue("usuario_id"))
	if err != nil {
		http.Error(w, "Usuario inválido", http.StatusBadRequest)
		return
	}

	total, err := strconv.ParseFloat(r.FormValue("total"), 64)
	if err != nil {
		http.Error(w, "Total inválido", http.StatusBadRequest)
		return
	}

	orden, err := models.NuevaOrden(usuarioID, total)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := orden.GenerarOrden(); err != nil {
		http.Error(w, "Error al generar la orden: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ordenes/usuario/"+strconv.Itoa(usuarioID), http.StatusSeeOther)
}

// Actualiza el estado de una orden aplicando las transiciones válidas
func ActualizarEstadoOrden(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	orden, ok := models.BuscarOrdenPorID(id)
	if !ok {
		http.Error(w, "Orden no encontrada", http.StatusNotFound)
		return
	}

	// CambiarEstado valida la transición antes de persistirla
	if err := orden.CambiarEstado(r.FormValue("estado")); err != nil {
		http.Error(w, "Estado inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := orden.ActualizarEstado(); err != nil {
		http.Error(w, "Error al actualizar el estado: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ordenes/"+strconv.Itoa(orden.ID), http.StatusSeeOther)
}

// Cancela una orden cambiando su estado a "cancelada"
func CancelarOrden(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	orden, ok := models.BuscarOrdenPorID(id)
	if !ok {
		http.Error(w, "Orden no encontrada", http.StatusNotFound)
		return
	}

	if err := orden.CancelarOrden(); err != nil {
		http.Error(w, "Error al cancelar la orden: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ordenes/"+strconv.Itoa(orden.ID), http.StatusSeeOther)
}
