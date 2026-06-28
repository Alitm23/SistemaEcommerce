package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/gorilla/mux"
)

type carritoViewData struct {
	Items []models.ItemCarrito
	Total float64
}

// VerCarrito muestra todos los productos del carrito actual.
func VerCarrito(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := redirigirSiNoAutenticado(w, r)
	if !ok {
		return
	}

	carritoID, err := models.ObtenerCarritoActivo(usuarioID)
	if err != nil {
		http.Error(w, "Error al cargar carrito", http.StatusInternalServerError)
		return
	}
	items, err := models.ObtenerItemsCarrito(carritoID)
	if err != nil {
		http.Error(w, "Error al cargar carrito", http.StatusInternalServerError)
		return
	}

	data := carritoViewData{Items: items, Total: totalCarrito(items)}
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/carrito.html"}, "base", data)
}

// AgregarItems maneja el POST del formulario.
func AgregarItems(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := redirigirSiNoAutenticado(w, r)
	if !ok {
		return
	}

	productoTallaID, _ := strconv.Atoi(r.FormValue("producto_talla_id"))
	cantidad, _ := strconv.Atoi(r.FormValue("cantidad"))
	if cantidad == 0 {
		cantidad = 1
	}

	carritoID, err := models.ObtenerCarritoActivo(usuarioID)
	if err != nil {
		http.Error(w, "No se pudo obtener el carrito", http.StatusInternalServerError)
		return
	}

	item := models.ItemCarrito{
		CarritoID:       carritoID,
		ProductoTallaID: productoTallaID,
		Cantidad:        cantidad,
	}

	if err := models.AgregarItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/carrito", http.StatusSeeOther)
}

// QuitarItem elimina un item especifico del carrito.
func QuitarItem(w http.ResponseWriter, r *http.Request) {
	if _, ok := redirigirSiNoAutenticado(w, r); !ok {
		return
	}

	vars := mux.Vars(r)
	itemID, _ := strconv.Atoi(vars["id"])

	if err := models.EliminarItem(itemID); err != nil {
		http.Error(w, "Error al eliminar item", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/carrito", http.StatusSeeOther)
}

// MostrarCheckout renderiza el formulario de pago.
func MostrarCheckout(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := redirigirSiNoAutenticado(w, r)
	if !ok {
		return
	}

	carritoID, err := models.ObtenerCarritoActivo(usuarioID)
	if err != nil {
		http.Error(w, "Error al cargar checkout", http.StatusInternalServerError)
		return
	}
	items, err := models.ObtenerItemsCarrito(carritoID)
	if err != nil {
		http.Error(w, "Error al cargar checkout", http.StatusInternalServerError)
		return
	}

	data := carritoViewData{Items: items, Total: totalCarrito(items)}
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/checkout.html"}, "base", data)
}

func ProcesarCheckout(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := redirigirSiNoAutenticado(w, r)
	if !ok {
		return
	}

	carritoID, err := models.ObtenerCarritoActivo(usuarioID)
	if err != nil {
		http.Error(w, "Error al procesar pago", http.StatusInternalServerError)
		return
	}
	items, err := models.ObtenerItemsCarrito(carritoID)
	if err != nil {
		http.Error(w, "Error al procesar pago", http.StatusInternalServerError)
		return
	}
	if len(items) == 0 {
		renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/mensaje.html"}, "base", map[string]string{
			"Titulo":  "Carrito vacio",
			"Mensaje": "No se registro ninguna orden porque tu carrito esta vacio.",
		})
		return
	}

	orden := models.Orden{
		UsuarioID:  usuarioID,
		Total:      totalCarrito(items),
		Estado:     r.FormValue("estado_pago"),
		MetodoPago: r.FormValue("metodo_pago"),
	}
	if orden.Estado != "pagado" && orden.Estado != "pendiente" && orden.Estado != "cancelado" {
		orden.Estado = "pagado"
	}
	if orden.MetodoPago == "" {
		orden.MetodoPago = "tarjeta"
	}
	if orden.Estado == "cancelado" {
		tieneCancelada, err := models.UsuarioTieneOrdenCancelada(usuarioID)
		if err != nil {
			http.Error(w, "Error al validar la orden", http.StatusInternalServerError)
			return
		}
		if tieneCancelada {
			renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/mensaje.html"}, "base", map[string]string{
				"Titulo":  "Orden no registrada",
				"Mensaje": "Este usuario ya tiene una orden cancelada. No se registro otra cancelacion.",
			})
			return
		}
	}

	if _, err := models.CrearOrdenDesdeCarrito(orden, carritoID); err != nil {
		http.Error(w, "Error al crear orden", http.StatusInternalServerError)
		return
	}
	titulo := "Orden registrada"
	mensaje := "Tu orden fue pagada correctamente."
	switch orden.Estado {
	case "cancelado":
		titulo = "Orden cancelada"
		mensaje = "Tu orden fue registrada como cancelada."
	case "pendiente":
		mensaje = "Tu orden quedo pendiente de pago."
	}
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/mensaje.html"}, "base", map[string]string{
		"Titulo":  titulo,
		"Mensaje": mensaje,
	})
}

func totalCarrito(items []models.ItemCarrito) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Cantidad) * item.PrecioUnitario
	}
	return total
}
