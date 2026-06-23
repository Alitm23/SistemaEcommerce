package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Verificar si se ejecuta el ping de la conexion
	DB, err := db.Connect()
	if err != nil {
		// Fatalf nos permite inyectar el error dentro del texto
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	fmt.Println("Conectado exitosamente a la base de datos")
	defer DB.Close()

	// Instanciar todos los handlers del sistema
	usuarioHandler := handlers.NuevoUsuarioHandler()
	productoHandler := handlers.NuevoProductoHandler()
	categoriaHandler := handlers.NuevoCategoriaHandler()
	materialHandler := handlers.NuevoMaterialHandler()
	carritoHandler := handlers.NuevoCarritoHandler()
	ordenHandler := handlers.NuevoOrdenHandler()
	pagoHandler := handlers.NuevoPagoHandler()

	r := mux.NewRouter()

	// ==========================================
	// Rutas de usuarios
	// ==========================================
	r.HandleFunc("/usuarios/registro", usuarioHandler.Registrar).Methods("POST")
	r.HandleFunc("/usuarios/login", usuarioHandler.Autenticar).Methods("POST")
	r.HandleFunc("/usuarios", usuarioHandler.Listar).Methods("GET")
	r.HandleFunc("/usuarios/{id}", usuarioHandler.ObtenerPorID).Methods("GET")
	r.HandleFunc("/usuarios/{id}", usuarioHandler.Actualizar).Methods("PUT")
	r.HandleFunc("/usuarios/{id}/password", usuarioHandler.CambiarPassword).Methods("PATCH")
	r.HandleFunc("/usuarios/{id}/rol", usuarioHandler.CambiarRol).Methods("PATCH")
	r.HandleFunc("/usuarios/{id}", usuarioHandler.Eliminar).Methods("DELETE")

	// ==========================================
	// Rutas de categorías
	// ==========================================
	r.HandleFunc("/categorias", categoriaHandler.Crear).Methods("POST")
	r.HandleFunc("/categorias", categoriaHandler.Listar).Methods("GET")
	r.HandleFunc("/categorias/{id}", categoriaHandler.ObtenerPorID).Methods("GET")
	r.HandleFunc("/categorias/{id}", categoriaHandler.Actualizar).Methods("PUT")
	r.HandleFunc("/categorias/{id}", categoriaHandler.Eliminar).Methods("DELETE")

	// ==========================================
	// Rutas de materiales
	// ==========================================
	r.HandleFunc("/materiales", materialHandler.Crear).Methods("POST")
	r.HandleFunc("/materiales", materialHandler.Listar).Methods("GET")
	r.HandleFunc("/materiales/{id}", materialHandler.ObtenerPorID).Methods("GET")
	r.HandleFunc("/materiales/{id}", materialHandler.Actualizar).Methods("PUT")
	r.HandleFunc("/materiales/{id}", materialHandler.Eliminar).Methods("DELETE")

	// ==========================================
	// Rutas de productos y tallas
	// ==========================================
	r.HandleFunc("/productos", productoHandler.Crear).Methods("POST")
	r.HandleFunc("/productos", productoHandler.Listar).Methods("GET")
	r.HandleFunc("/productos/{id}", productoHandler.ObtenerPorID).Methods("GET")
	r.HandleFunc("/productos/{id}", productoHandler.Actualizar).Methods("PUT")
	r.HandleFunc("/productos/{id}", productoHandler.Eliminar).Methods("DELETE")
	r.HandleFunc("/productos/categoria/{categoriaId}", productoHandler.ListarPorCategoria).Methods("GET")
	r.HandleFunc("/productos/material/{materialId}", productoHandler.ListarPorMaterial).Methods("GET")

	// Tallas de un producto
	r.HandleFunc("/productos/{id}/tallas", productoHandler.AgregarTalla).Methods("POST")
	r.HandleFunc("/productos/{id}/tallas", productoHandler.ObtenerTallas).Methods("GET")
	r.HandleFunc("/productos/tallas/{tallaId}/stock", productoHandler.ActualizarStock).Methods("PATCH")
	r.HandleFunc("/productos/tallas/{tallaId}", productoHandler.EliminarTalla).Methods("DELETE")

	// ==========================================
	// Rutas de carrito
	// ==========================================
	r.HandleFunc("/carrito", carritoHandler.Abrir).Methods("POST")
	r.HandleFunc("/carrito/{id}/cerrar", carritoHandler.Cerrar).Methods("PATCH")
	r.HandleFunc("/carrito/{id}/items", carritoHandler.AgregarItem).Methods("POST")
	r.HandleFunc("/carrito/{id}/items", carritoHandler.ObtenerItems).Methods("GET")
	r.HandleFunc("/carrito/items/{itemId}", carritoHandler.ActualizarCantidadItem).Methods("PATCH")
	r.HandleFunc("/carrito/items/{itemId}", carritoHandler.QuitarItem).Methods("DELETE")

	// ==========================================
	// Rutas de órdenes
	// ==========================================
	r.HandleFunc("/ordenes", ordenHandler.Generar).Methods("POST")
	r.HandleFunc("/ordenes", ordenHandler.ListarTodas).Methods("GET")
	r.HandleFunc("/ordenes/{id}", ordenHandler.ObtenerPorID).Methods("GET")
	r.HandleFunc("/ordenes/{id}/estado", ordenHandler.ActualizarEstado).Methods("PATCH")
	r.HandleFunc("/ordenes/{id}/cancelar", ordenHandler.Cancelar).Methods("PATCH")
	r.HandleFunc("/ordenes/{id}/items", ordenHandler.AgregarItem).Methods("POST")
	r.HandleFunc("/ordenes/{id}/items", ordenHandler.ObtenerItems).Methods("GET")
	r.HandleFunc("/ordenes/usuario/{usuarioId}", ordenHandler.ListarPorUsuario).Methods("GET")

	// ==========================================
	// Rutas de pagos
	// ==========================================
	r.HandleFunc("/pagos", pagoHandler.Registrar).Methods("POST")
	r.HandleFunc("/pagos", pagoHandler.Listar).Methods("GET")
	r.HandleFunc("/pagos/orden/{ordenId}/estado", pagoHandler.ActualizarEstado).Methods("PATCH")
	r.HandleFunc("/pagos/orden/{ordenId}/anular", pagoHandler.Anular).Methods("PATCH")

	// Iniciar el servidor en el puerto 8080
	log.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error al iniciar el servidor:", err)
	}
}
