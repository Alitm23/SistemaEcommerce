package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Obtiene todos los productos y los renderiza en el template
func ListarProductos(w http.ResponseWriter, r *http.Request) {
	productos, err := models.ListarProductos()
	if err != nil {
		http.Error(w, "Error al obtener los productos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/producto/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, productos)
}

// Carga un producto por su ID y renderiza la vista de detalle
func VerProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	producto, ok := models.BuscarProductoPorID(id)
	if !ok {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/producto/detalle.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, producto)
}

// Renderiza el formulario vacío para crear un nuevo producto
func MostrarFormProducto(w http.ResponseWriter, r *http.Request) {
	// Cargar las categorías para el selector del formulario
	categorias, err := models.ListarCategorias()
	if err != nil {
		http.Error(w, "Error al obtener las categorías: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/producto/form.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	// El template recibe las categorías disponibles; el producto aún no existe
	tmpl.Execute(w, map[string]interface{}{
		"Categorias": categorias,
		"Producto":   nil,
	})
}

// Procesa el formulario y registra un nuevo producto en la BD
func CrearProducto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	categoriaID, err := strconv.Atoi(r.FormValue("categoria_id"))
	if err != nil {
		http.Error(w, "Categoría inválida", http.StatusBadRequest)
		return
	}

	precio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "Precio inválido", http.StatusBadRequest)
		return
	}

	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		http.Error(w, "Stock inválido", http.StatusBadRequest)
		return
	}

	producto, err := models.NuevoProducto(
		categoriaID,
		r.FormValue("nombre"),
		r.FormValue("descripcion"),
		precio,
		stock,
	)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := producto.Registrar(); err != nil {
		http.Error(w, "Error al crear el producto: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

// Cargar los datos actuales y muestra el formulario de edición.
func MostrarEditarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	producto, ok := models.BuscarProductoPorID(id)
	if !ok {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	// Cargar las categorías para el selector del formulario
	categorias, err := models.ListarCategorias()
	if err != nil {
		http.Error(w, "Error al obtener las categorías: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/producto/form.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Categorias": categorias,
		"Producto":   producto,
	})
}

// Función para guardar los cambios del formulario de edición en la BD
func ActualizarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Verificar que el producto existe antes de modificarlo
	producto, ok := models.BuscarProductoPorID(id)
	if !ok {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	categoriaID, err := strconv.Atoi(r.FormValue("categoria_id"))
	if err != nil {
		http.Error(w, "Categoría inválida", http.StatusBadRequest)
		return
	}

	precio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "Precio inválido", http.StatusBadRequest)
		return
	}

	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		http.Error(w, "Stock inválido", http.StatusBadRequest)
		return
	}

	// Sobreescribir los campos con los nuevos valores del formulario
	producto.CategoriaID = categoriaID
	producto.Nombre = r.FormValue("nombre")
	producto.Descripcion = r.FormValue("descripcion")
	producto.Precio = precio
	producto.Stock = stock

	if err := producto.Actualizar(); err != nil {
		http.Error(w, "Error al actualizar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

// elimina un producto por su ID
func EliminarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Solo necesitamos el ID para la operación de eliminación
	producto := &models.Producto{ID: id}
	if err := producto.Eliminar(); err != nil {
		http.Error(w, "Error al eliminar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}
