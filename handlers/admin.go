package handlers

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Alitm23/SistemaEcommerce/models"
)

type adminDashboardData struct {
	Active         string
	TotalVentas    float64
	TotalProductos int
	TotalOrdenes   int
	TotalUsuarios  int
	MasVendidos    []models.ProductoVendido
	MenosVendidos  []models.ProductoVendido
}

type adminProductosData struct {
	Active     string
	Productos  []models.Producto
	Categorias []models.Categoria
}

type adminOrdenesData struct {
	Active   string
	Ordenes  []models.Orden
	Detalles map[int][]models.ItemCarrito
}

type adminUsuariosData struct {
	Active   string
	Usuarios []models.Usuario
}

func renderAdminTemplate(w http.ResponseWriter, archivo string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/admin_base.html", archivo)
	if err != nil {
		http.Error(w, "Error al cargar el panel administrativo", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "admin_base", data); err != nil {
		http.Error(w, "Error al renderizar el panel administrativo", http.StatusInternalServerError)
	}
}

func RedireccionarDashboardAdmin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// DashboardAdmin muestra el resumen de la tienda
func DashboardAdmin(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	productos, _ := models.ListarProductos()
	ordenes, _ := models.ListarOrdenes()
	usuarios, _ := models.ListarUsuarios()

	var totalVentas float64
	for _, orden := range ordenes {
		if orden.Estado == "pagado" {
			totalVentas += orden.Total
		}
	}
	masVendidos, _ := models.ListarProductosVendidos(false, 5)
	menosVendidos, _ := models.ListarProductosVendidos(true, 5)

	data := adminDashboardData{
		Active:         "dashboard",
		TotalVentas:    totalVentas,
		TotalProductos: len(productos),
		TotalOrdenes:   len(ordenes),
		TotalUsuarios:  len(usuarios),
		MasVendidos:    masVendidos,
		MenosVendidos:  menosVendidos,
	}
	renderAdminTemplate(w, "templates/admin_dashboard.html", data)
}

// GestionInventario muestra el listado para editar/borrar
func GestionInventario(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	productos, err := models.ListarProductos()
	if err != nil {
		http.Error(w, "Error al cargar productos", http.StatusInternalServerError)
		return
	}
	for i := range productos {
		tallas, err := models.ListarTallasPorProducto(productos[i].ID)
		if err != nil {
			http.Error(w, "Error al cargar stock", http.StatusInternalServerError)
			return
		}
		productos[i].Tallas = tallas
	}
	categorias, err := models.ListarCategorias()
	if err != nil {
		http.Error(w, "Error al cargar categorias", http.StatusInternalServerError)
		return
	}
	renderAdminTemplate(w, "templates/admin_productos.html", adminProductosData{
		Active:     "productos",
		Productos:  productos,
		Categorias: categorias,
	})
}

// CrearProducto es el POST para el formulario del admin
func CrearProducto(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Error al leer el formulario", http.StatusBadRequest)
			return
		}
		precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
		catID, _ := strconv.Atoi(r.FormValue("categoria_id"))
		stock, _ := strconv.Atoi(r.FormValue("stock"))
		imagenURL, err := guardarImagenProducto(r, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nuevo := models.Producto{
			Nombre:      r.FormValue("nombre"),
			Material:    r.FormValue("material"),
			Descripcion: r.FormValue("descripcion"),
			Precio:      precio,
			CategoriaID: catID,
			ImagenURL:   imagenURL,
		}

		nuevoID, err := models.CrearProducto(nuevo)
		if err != nil {
			http.Error(w, "Error al crear", http.StatusInternalServerError)
			return
		}
		if err := models.CrearTallaProducto(nuevoID, r.FormValue("talla"), stock); err != nil {
			http.Error(w, "Error al crear stock", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
	}
}

func EliminarProductoAdmin(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID de producto invalido", http.StatusBadRequest)
		return
	}
	if _, err := models.EliminarProducto(id); err != nil {
		http.Error(w, "Error al eliminar producto", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
}

func CrearCategoriaAdmin(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	categoria := models.Categoria{
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		Descripcion: strings.TrimSpace(r.FormValue("descripcion")),
	}
	if categoria.Nombre == "" {
		http.Error(w, "El nombre de la categoria es obligatorio", http.StatusBadRequest)
		return
	}
	if _, err := models.CrearCategoria(categoria); err != nil {
		http.Error(w, "Error al crear categoria", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
}

func EditarProductoAdmin(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID de producto invalido", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al leer el formulario", http.StatusBadRequest)
		return
	}
	precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
	catID, _ := strconv.Atoi(r.FormValue("categoria_id"))
	imagenURL, err := guardarImagenProducto(r, r.FormValue("imagen_actual"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	producto := models.Producto{
		ID:          id,
		CategoriaID: catID,
		Nombre:      r.FormValue("nombre"),
		Material:    r.FormValue("material"),
		Descripcion: r.FormValue("descripcion"),
		Precio:      precio,
		ImagenURL:   imagenURL,
	}
	if _, err := models.ActualizarProducto(producto); err != nil {
		http.Error(w, "Error al actualizar producto", http.StatusInternalServerError)
		return
	}
	for clave, valores := range r.MultipartForm.Value {
		if !strings.HasPrefix(clave, "stock_") || len(valores) == 0 {
			continue
		}
		tallaID, err := strconv.Atoi(strings.TrimPrefix(clave, "stock_"))
		if err != nil {
			continue
		}
		stock, err := strconv.Atoi(valores[0])
		if err != nil {
			continue
		}
		if err := models.ActualizarStockTalla(tallaID, stock); err != nil {
			http.Error(w, "Error al actualizar stock", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
}

func GestionOrdenes(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	ordenes, err := models.ListarOrdenes()
	if err != nil {
		http.Error(w, "Error al cargar ordenes", http.StatusInternalServerError)
		return
	}
	detalles := make(map[int][]models.ItemCarrito)
	for _, orden := range ordenes {
		items, err := models.ObtenerDetalleOrden(orden.ID)
		if err == nil {
			detalles[orden.ID] = items
		}
	}
	renderAdminTemplate(w, "templates/admin_ordenes.html", adminOrdenesData{
		Active:   "ordenes",
		Ordenes:  ordenes,
		Detalles: detalles,
	})
}

func GestionUsuarios(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	usuarios, err := models.ListarUsuarios()
	if err != nil {
		http.Error(w, "Error al cargar usuarios", http.StatusInternalServerError)
		return
	}
	renderAdminTemplate(w, "templates/admin_usuarios.html", adminUsuariosData{
		Active:   "usuarios",
		Usuarios: usuarios,
	})
}

func CambiarRolUsuarioAdmin(w http.ResponseWriter, r *http.Request) {
	if !redirigirSiNoAdmin(w, r) {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
		return
	}
	nuevoRol := r.FormValue("rol")
	if nuevoRol != "admin" && nuevoRol != "cliente" {
		http.Error(w, "Rol invalido", http.StatusBadRequest)
		return
	}
	if err := models.ActualizarRol(id, nuevoRol); err != nil {
		http.Error(w, "Error al cambiar rol", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/usuarios", http.StatusSeeOther)
}

func guardarImagenProducto(r *http.Request, actual string) (string, error) {
	archivo, encabezado, err := r.FormFile("imagen")
	if err == http.ErrMissingFile {
		return actual, nil
	}
	if err != nil {
		return "", err
	}
	defer archivo.Close()

	extension := strings.ToLower(filepath.Ext(encabezado.Filename))
	switch extension {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
	default:
		return "", errors.New("solo se permiten imagenes jpg, png, webp o gif")
	}

	carpeta := filepath.Join("static", "img", "uploads")
	if err := os.MkdirAll(carpeta, 0755); err != nil {
		return "", err
	}
	nombre := "producto_" + strconv.FormatInt(time.Now().UnixNano(), 10) + extension
	rutaDisco := filepath.Join(carpeta, nombre)
	destino, err := os.Create(rutaDisco)
	if err != nil {
		return "", err
	}
	defer destino.Close()

	if _, err := io.Copy(destino, archivo); err != nil {
		return "", err
	}
	return "/static/img/uploads/" + nombre, nil
}
