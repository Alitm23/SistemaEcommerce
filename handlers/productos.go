package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/services"
	"github.com/gorilla/mux"
)

// ProductoHandler expone los endpoints HTTP relacionados con la gestión de productos
type ProductoHandler struct {
	servicio *services.ProductoService
}

// NuevoProductoHandler construye el handler inyectando el servicio correspondiente
func NuevoProductoHandler() *ProductoHandler {
	return &ProductoHandler{
		servicio: services.NuevoProductoService(),
	}
}

// Crear registra un nuevo producto en el catálogo
func (h *ProductoHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		CategoriaID int     `json:"categoria_id"`
		MaterialID  int     `json:"material_id"`
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Precio      float64 `json:"precio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	producto, err := h.servicio.CrearProducto(
		datos.CategoriaID, datos.MaterialID,
		datos.Nombre, datos.Descripcion, datos.Precio,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(producto)
}

// ObtenerPorID recupera un producto según su identificador único
func (h *ProductoHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	producto, ok := h.servicio.BuscarPorID(id)
	if !ok {
		http.Error(w, "producto no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// Listar recupera todos los productos del catálogo
func (h *ProductoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	productos, err := h.servicio.ListarProductos()
	if err != nil {
		http.Error(w, "error al obtener productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

// ListarPorCategoria recupera todos los productos de una categoría específica
func (h *ProductoHandler) ListarPorCategoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoriaID, err := strconv.Atoi(vars["categoriaId"])
	if err != nil {
		http.Error(w, "identificador de categoría inválido", http.StatusBadRequest)
		return
	}

	productos, err := h.servicio.ListarPorCategoria(categoriaID)
	if err != nil {
		http.Error(w, "error al obtener productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

// ListarPorMaterial recupera todos los productos de un material específico
func (h *ProductoHandler) ListarPorMaterial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	materialID, err := strconv.Atoi(vars["materialId"])
	if err != nil {
		http.Error(w, "identificador de material inválido", http.StatusBadRequest)
		return
	}

	productos, err := h.servicio.ListarPorMaterial(materialID)
	if err != nil {
		http.Error(w, "error al obtener productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

// Actualizar modifica los datos de un producto existente
func (h *ProductoHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		CategoriaID int     `json:"categoria_id"`
		MaterialID  int     `json:"material_id"`
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Precio      float64 `json:"precio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	producto, err := h.servicio.ActualizarProducto(
		id, datos.CategoriaID, datos.MaterialID,
		datos.Nombre, datos.Descripcion, datos.Precio,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(producto)
}

// Eliminar borra un producto del catálogo por su identificador.
// Las tallas asociadas se eliminan automáticamente por ON DELETE CASCADE.
func (h *ProductoHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.EliminarProducto(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AgregarTalla agrega una nueva talla con stock inicial a un producto existente
func (h *ProductoHandler) AgregarTalla(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador de producto inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Talla string `json:"talla"`
		Stock int    `json:"stock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	talla, err := h.servicio.AgregarTalla(productoID, datos.Talla, datos.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(talla)
}

// ObtenerTallas recupera todas las tallas disponibles para un producto específico
func (h *ProductoHandler) ObtenerTallas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "identificador de producto inválido", http.StatusBadRequest)
		return
	}

	tallas, err := h.servicio.ObtenerTallas(productoID)
	if err != nil {
		http.Error(w, "error al obtener tallas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tallas)
}

// ActualizarStock modifica el stock de una talla usando un delta positivo o negativo
func (h *ProductoHandler) ActualizarStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tallaID, err := strconv.Atoi(vars["tallaId"])
	if err != nil {
		http.Error(w, "identificador de talla inválido", http.StatusBadRequest)
		return
	}

	var datos struct {
		Delta int `json:"delta"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "cuerpo de la petición inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.ActualizarStockTalla(tallaID, datos.Delta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// EliminarTalla borra una talla específica de un producto por su identificador
func (h *ProductoHandler) EliminarTalla(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tallaID, err := strconv.Atoi(vars["tallaId"])
	if err != nil {
		http.Error(w, "identificador de talla inválido", http.StatusBadRequest)
		return
	}

	if err := h.servicio.EliminarTalla(tallaID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
