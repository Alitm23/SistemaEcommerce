package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Lista los ítems pertenecientes a una orden específica
func ListarItemsOrden(w http.ResponseWriter, r *http.Request) {
	ordenID, err := strconv.Atoi(mux.Vars(r)["ordenID"])
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	items, err := models.ListarItemsPorOrden(ordenID)
	if err != nil {
		http.Error(w, "Error al obtener los ítems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/item_orden/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, items)
}

// Agrega un ítem a una orden registrando el producto, cantidad y precio de compra
func AgregarItemOrden(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	ordenID, err := strconv.Atoi(r.FormValue("orden_id"))
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	productoID, err := strconv.Atoi(r.FormValue("producto_id"))
	if err != nil {
		http.Error(w, "Producto inválido", http.StatusBadRequest)
		return
	}

	cantidad, err := strconv.Atoi(r.FormValue("cantidad"))
	if err != nil {
		http.Error(w, "Cantidad inválida", http.StatusBadRequest)
		return
	}

	precioCompra, err := strconv.ParseFloat(r.FormValue("precio_compra"), 64)
	if err != nil {
		http.Error(w, "Precio de compra inválido", http.StatusBadRequest)
		return
	}

	item, err := models.NuevoItemOrden(ordenID, productoID, cantidad, precioCompra)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := item.AgregarAOrden(); err != nil {
		http.Error(w, "Error al agregar el ítem: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ordenes/items/"+strconv.Itoa(ordenID), http.StatusSeeOther)
}

// Elimina un ítem de una orden por su ID
func EliminarItemOrden(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	ordenID, err := strconv.Atoi(r.FormValue("orden_id"))
	if err != nil {
		http.Error(w, "Orden inválida", http.StatusBadRequest)
		return
	}

	// Solo necesitamos el ID para la operación de eliminación
	item := &models.ItemOrden{ID: id}
	if err := item.Eliminar(); err != nil {
		http.Error(w, "Error al eliminar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ordenes/items/"+strconv.Itoa(ordenID), http.StatusSeeOther)
}
