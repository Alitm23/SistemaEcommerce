package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Lista los ítems de un carrito específico
func ListarItemsCarrito(w http.ResponseWriter, r *http.Request) {
	carritoID, err := strconv.Atoi(mux.Vars(r)["carritoID"])
	if err != nil {
		http.Error(w, "Carrito inválido", http.StatusBadRequest)
		return
	}

	items, err := models.ListarItemsPorCarrito(carritoID)
	if err != nil {
		http.Error(w, "Error al obtener los ítems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/item_carrito/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, items)
}

// Agrega un producto al carrito. Si el producto ya existe en el carrito,
// incrementa la cantidad del ítem en lugar de duplicarlo.
func AgregarItemCarrito(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	carritoID, err := strconv.Atoi(r.FormValue("carrito_id"))
	if err != nil {
		http.Error(w, "Carrito inválido", http.StatusBadRequest)
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

	// Obtener el producto para tomar su precio unitario actual
	producto, ok := models.BuscarProductoPorID(productoID)
	if !ok {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	// Si el producto ya está en el carrito, incrementar su cantidad
	items, err := models.ListarItemsPorCarrito(carritoID)
	if err != nil {
		http.Error(w, "Error al obtener los ítems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, existente := range items {
		if existente.ProductoID == productoID {
			existente.Cantidad += cantidad
			if err := existente.ActualizarCantidad(); err != nil {
				http.Error(w, "Error al actualizar la cantidad: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/carritos/items/"+strconv.Itoa(carritoID), http.StatusSeeOther)
			return
		}
	}

	item, err := models.NuevoItemCarrito(carritoID, productoID, cantidad, producto.Precio)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := item.AgregarAlCarrito(); err != nil {
		http.Error(w, "Error al agregar el ítem: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/carritos/items/"+strconv.Itoa(carritoID), http.StatusSeeOther)
}

// Actualiza la cantidad de un ítem del carrito
func ActualizarItemCarrito(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	cantidad, err := strconv.Atoi(r.FormValue("cantidad"))
	if err != nil {
		http.Error(w, "Cantidad inválida", http.StatusBadRequest)
		return
	}

	carritoID, err := strconv.Atoi(r.FormValue("carrito_id"))
	if err != nil {
		http.Error(w, "Carrito inválido", http.StatusBadRequest)
		return
	}

	item := &models.ItemCarrito{ID: id, Cantidad: cantidad}
	if err := item.ActualizarCantidad(); err != nil {
		http.Error(w, "Error al actualizar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/carritos/items/"+strconv.Itoa(carritoID), http.StatusSeeOther)
}

// Quita un ítem del carrito por su ID
func EliminarItemCarrito(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	carritoID, err := strconv.Atoi(r.FormValue("carrito_id"))
	if err != nil {
		http.Error(w, "Carrito inválido", http.StatusBadRequest)
		return
	}

	// Solo necesitamos el ID para la operación de eliminación
	item := &models.ItemCarrito{ID: id}
	if err := item.QuitarDelCarrito(); err != nil {
		http.Error(w, "Error al eliminar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/carritos/items/"+strconv.Itoa(carritoID), http.StatusSeeOther)
}
