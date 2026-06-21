package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Obtiene todos los pagos y los renderiza en el template
func ListarPagos(w http.ResponseWriter, r *http.Request) {
	pagos, err := models.ListarPagos()
	if err != nil {
		http.Error(w, "Error al obtener los pagos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/pago/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, pagos)
}

// Muestra el pago asociado a una orden
func VerPagoPorOrden(w http.ResponseWriter, r *http.Request) {
	ordenID, err := strconv.Atoi(mux.Vars(r)["ordenID"])
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	pago, ok := models.BuscarPagoPorOrden(ordenID)
	if !ok {
		http.Error(w, "Pago no encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/pago/detalle.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, pago)
}

// Registra un pago simulado para una orden. El pago se aprueba automáticamente
// (estado "completado") si la orden existe; de lo contrario queda pendiente.
func CrearPago(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	ordenID, err := strconv.Atoi(r.FormValue("orden_id"))
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	monto, err := strconv.ParseFloat(r.FormValue("monto"), 64)
	if err != nil {
		http.Error(w, "Monto inválido", http.StatusBadRequest)
		return
	}

	pago, err := models.NuevoPago(ordenID, monto)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := pago.RegistrarPago(); err != nil {
		http.Error(w, "Error al registrar el pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pagos", http.StatusSeeOther)
}

// Actualiza el estado de un pago aplicando las reglas del modelo
func ActualizarEstadoPago(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	ordenID, err := strconv.Atoi(mux.Vars(r)["ordenID"])
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	pago, ok := models.BuscarPagoPorOrden(ordenID)
	if !ok {
		http.Error(w, "Pago no encontrado", http.StatusNotFound)
		return
	}

	// CambiarEstado valida la transición antes de persistirla
	if err := pago.CambiarEstado(r.FormValue("estado")); err != nil {
		http.Error(w, "Estado inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := pago.ActualizarEstado(); err != nil {
		http.Error(w, "Error al actualizar el estado: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pagos", http.StatusSeeOther)
}

// Anula un pago marcándolo como fallido
func AnularPago(w http.ResponseWriter, r *http.Request) {
	ordenID, err := strconv.Atoi(mux.Vars(r)["ordenID"])
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	pago, ok := models.BuscarPagoPorOrden(ordenID)
	if !ok {
		http.Error(w, "Pago no encontrado", http.StatusNotFound)
		return
	}

	if err := pago.AnularPago(); err != nil {
		http.Error(w, "Error al anular el pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pagos", http.StatusSeeOther)
}
