package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/handlers"
	"github.com/Alitm23/SistemaEcommerce/services"

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
		log.Fatal("Error crítico al conectar a la BD: ", errDB)
	}
	// Cierra la conexión cuando el servidor se detenga
	defer db.DB.Close()

	r := mux.NewRouter()

	// Instanciar Servicios (Capa de Negocio)
	usuarioService := services.NuevoUsuarioService()
	productoService := services.NuevoProductoService()
	// Instanciar Controladores (Capa HTTP)
	productoHandler := handlers.NuevoProductoHandler(productoService)
	homeHandler := handlers.NuevoHomeHandler()
	authHandler := handlers.NuevoAuthHandler(usuarioService)

	//RUTAS
	// Rutas de Navegación Pública
	r.HandleFunc("/productos", productoHandler.ListarProductos).Methods("GET")
	r.HandleFunc("/", homeHandler.MostrarInicio).Methods("GET")

	// Rutas de Autenticación (Login / Registro / Logout)
	r.HandleFunc("/login", authHandler.MostrarLogin).Methods("GET")
	r.HandleFunc("/login", authHandler.ProcesarLogin).Methods("POST")
	r.HandleFunc("/registro", authHandler.ProcesarRegistro).Methods("POST")
	r.HandleFunc("/logout", authHandler.CerrarSesion).Methods("GET")

	//Archivos estátivos y servidor

	// Servir la carpeta de archivos estáticos (CSS, JS, imágenes)
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Leer el puerto desde el .env
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	log.Printf("Servidor corriendo en http://localhost:%s", puerto)

	// Levantar el servidor HTTP
	err = http.ListenAndServe(":"+puerto, r)
	if err != nil {
		log.Fatal("Error crítico al iniciar el servidor: ", err)
	}
}
