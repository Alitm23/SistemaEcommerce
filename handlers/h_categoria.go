package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Obtiene todas las categorías y las renderiza en el template
func ListarCategorias(w http.ResponseWriter, r *http.Request) {
	categorias, err := models.ListarCategorias()
	if err != nil {
		http.Error(w, "Error al obtener las categorías: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/categoria/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, categorias)
}

// Renderiza el formulario vacío para crear una nueva categoría
func MostrarFormCategoria(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/categoria/form.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Procesa el formulario y registra una nueva categoría en la BD
func CrearCategoria(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Valida que el nombre no esté vacío
	categoria, err := models.NuevaCategoria(
		r.FormValue("nombre"))
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := categoria.Registrar(); err != nil {
		http.Error(w, "Error al crear la categoría: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categorias", http.StatusSeeOther)
}

// Cargar los datos actuales y muestra el formulario de edición.
func MostrarEditarCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// retorna (categoria, bool)
	categoria, ok := models.BuscarCategoriaPorID(id)
	if !ok {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/categoria/form.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	// El template detecta si catergori.ID > 0 para mostrar editar o crear
	tmpl.Execute(w, categoria)
}

// Función para guardar los cambios del formulario de edición en la BD
func ActualizarCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Verificar que la categoría existe antes de modificarla
	categoria, ok := models.BuscarCategoriaPorID(id)
	if !ok {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Sobreescribir los campos con los nuevos valores del formulario

	categoria.Nombre = r.FormValue("nombre")
	if err := categoria.Actualizar(); err != nil {
		http.Error(w, "Error al actualizar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categorias", http.StatusSeeOther)
}

// elimina una categoría por su ID
func EliminarCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Solo necesitamos el ID para la operación de eliminación
	categoria := &models.Categoria{ID: id}
	if err := categoria.Eliminar(); err != nil {
		http.Error(w, "Error al eliminar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categorias", http.StatusSeeOther)
}
