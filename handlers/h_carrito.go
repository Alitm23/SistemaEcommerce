package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Muestra el carrito de un usuario junto con sus ítems y el total calculado
func VerCarrito(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := strconv.Atoi(mux.Vars(r)["usuarioID"])
	if err != nil {
		http.Error(w, "Usuario inválido", http.StatusBadRequest)
		return
	}

	// Buscar el carrito activo del usuario — retorna (carrito, bool)
	carrito, ok := models.BuscarCarritoActivoPorUsuario(usuarioID)
	if !ok {
		http.Error(w, "El usuario no tiene un carrito activo", http.StatusNotFound)
		return
	}

	items, err := models.ListarItemsPorCarrito(carrito.ID)
	if err != nil {
		http.Error(w, "Error al obtener los ítems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calcular el total sumando el subtotal de cada ítem
	var total float64
	for _, item := range items {
		total += item.Subtotal()
	}

	tmpl, err := template.ParseFiles("templates/carrito/detalle.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Carrito": carrito,
		"Items":   items,
		"Total":   total,
	})
}

// Crea (abre) un nuevo carrito para el usuario indicado
func CrearCarrito(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	usuarioID, err := strconv.Atoi(r.FormValue("usuario_id"))
	if err != nil {
		http.Error(w, "Usuario inválido", http.StatusBadRequest)
		return
	}

	carrito, err := models.NuevoCarrito(usuarioID)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := carrito.Abrir(); err != nil {
		http.Error(w, "Error al crear el carrito: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/carritos/"+strconv.Itoa(usuarioID), http.StatusSeeOther)
}

// Cierra el carrito (checkout) cambiando su estado a "cerrado"
func CerrarCarrito(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	carrito, ok := models.BuscarCarritoPorID(id)
	if !ok {
		http.Error(w, "Carrito no encontrado", http.StatusNotFound)
		return
	}

	// CambiarEstado valida la transición antes de persistirla
	if err := carrito.CambiarEstado("cerrado"); err != nil {
		http.Error(w, "Estado inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := carrito.Cerrar(); err != nil {
		http.Error(w, "Error al cerrar el carrito: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/carritos/"+strconv.Itoa(carrito.UsuarioID), http.StatusSeeOther)
}

// Elimina un carrito por su ID
func EliminarCarrito(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Solo necesitamos el ID para la operación de eliminación
	carrito := &models.Carrito{ID: id}
	if err := carrito.Eliminar(); err != nil {
		http.Error(w, "Error al eliminar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
