package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/handlers"
	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar el archivo .env al inicio de la aplicación
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables del sistema")
	}

	// Conectar a PostgreSQL
	_, errDB := db.Connect()
	if errDB != nil {
		log.Fatal("Error al conectar a la BD: ", errDB)
	}
	// Cierra la conexión cuando el servidor se detenga
	defer db.DB.Close()

	models.GestiondeStock()
	models.GestorNotificaciones()
	models.GestorAnalisis()

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.MostrarInicioGetHandler).Methods("GET")
	r.HandleFunc("/productos", handlers.ListarProductos).Methods("GET")
	r.HandleFunc("/productos/{id:[0-9]+}", handlers.DetalleProducto).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("GET", "POST")
	r.HandleFunc("/registro", handlers.Registro).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")
	r.HandleFunc("/carrito", handlers.VerCarrito).Methods("GET")
	r.HandleFunc("/carrito/agregar", handlers.AgregarItems).Methods("POST")
	r.HandleFunc("/carrito/quitar/{id:[0-9]+}", handlers.QuitarItem).Methods("POST")
	r.HandleFunc("/checkout", handlers.MostrarCheckout).Methods("GET")
	r.HandleFunc("/checkout/procesar", handlers.ProcesarCheckout).Methods("POST")

	// --- 2. RUTAS PÚBLICAS (API JSON) ---
	r.HandleFunc("/api/auth/login", handlers.ApiLogin).Methods("POST")
	r.HandleFunc("/api/auth/registro", handlers.ApiRegistro).Methods("POST")
	r.HandleFunc("/api/productos", handlers.ApiListarProductos).Methods("GET")
	r.HandleFunc("/api/productos/buscar", handlers.ApiBuscarProductos).Methods("GET")

	// --- 3. RUTAS PRIVADAS (API PROTEGIDA CON JWT) ---
	api := r.PathPrefix("/api").Subrouter()
	api.Use(utils.MiddlewareJWT)

	api.HandleFunc("/carrito", handlers.ApiVerCarrito).Methods("GET")
	api.HandleFunc("/carrito/agregar", handlers.ApiAgregarItem).Methods("POST")
	api.HandleFunc("/ordenes", handlers.ApiListarOrdenes).Methods("GET")
	api.HandleFunc("/ordenes/pagar", handlers.ApiProcesarPago).Methods("POST")

	// --- 4. RUTAS ADMIN (PROTEGIDAS) ---
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(utils.SoloAdminMiddleware)

	admin.HandleFunc("/productos", handlers.ApiAdminCrearProducto).Methods("POST")
	admin.HandleFunc("/productos/{id:[0-9]+}", handlers.ApiAdminActualizarProducto).Methods("PUT")
	admin.HandleFunc("/productos/{id:[0-9]+}", handlers.ApiAdminEliminarProducto).Methods("DELETE")

	// --- 5. PANEL ADMINISTRATIVO (VISTAS HTML) ---
	adminWeb := r.PathPrefix("/admin").Subrouter()
	adminWeb.HandleFunc("", handlers.RedireccionarDashboardAdmin).Methods("GET")
	adminWeb.HandleFunc("/", handlers.RedireccionarDashboardAdmin).Methods("GET")
	adminWeb.HandleFunc("/dashboard", handlers.DashboardAdmin).Methods("GET")
	adminWeb.HandleFunc("/productos", handlers.GestionInventario).Methods("GET")
	adminWeb.HandleFunc("/productos/nuevo", handlers.CrearProducto).Methods("POST")
	adminWeb.HandleFunc("/categorias/nueva", handlers.CrearCategoriaAdmin).Methods("POST")
	adminWeb.HandleFunc("/productos/editar", handlers.EditarProductoAdmin).Methods("POST")
	adminWeb.HandleFunc("/productos/eliminar", handlers.EliminarProductoAdmin).Methods("GET", "POST")
	adminWeb.HandleFunc("/ordenes", handlers.GestionOrdenes).Methods("GET")
	adminWeb.HandleFunc("/ordenes/cancelar", handlers.CancelarOrdenAdmin).Methods("POST")
	adminWeb.HandleFunc("/usuarios", handlers.GestionUsuarios).Methods("GET")
	adminWeb.HandleFunc("/usuarios/rol", handlers.CambiarRolUsuarioAdmin).Methods("POST")

	// Archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	log.Printf("Servidor corriendo en http://localhost:%s", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, r))
}
