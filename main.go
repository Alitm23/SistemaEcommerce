/*
@Configuracion inicial del proyecto
@autor: Nataly Tituaña

*/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Archivos estáticos
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// Instanciar handlers — cada uno inyecta su service internamente
	usuarioHandler := handlers.NuevoUsuarioHandler()

	// Rutas de Usuario
	router.HandleFunc("/registro", usuarioHandler.MostrarRegistro).Methods("GET")
	router.HandleFunc("/registro", usuarioHandler.ProcesarRegistro).Methods("POST")
	router.HandleFunc("/login", usuarioHandler.MostrarLogin).Methods("GET")
	router.HandleFunc("/login", usuarioHandler.ProcesarLogin).Methods("POST")
	router.HandleFunc("/usuarios", usuarioHandler.ListarUsuarios).Methods("GET")
	router.HandleFunc("/usuarios/{id}/editar", usuarioHandler.MostrarEdicion).Methods("GET")
	router.HandleFunc("/usuarios/{id}/editar", usuarioHandler.ProcesarEdicion).Methods("POST")
	router.HandleFunc("/usuarios/{id}/eliminar", usuarioHandler.EliminarUsuario).Methods("POST")

	fmt.Println("Servidor corriendo levantado")
	log.Fatal(http.ListenAndServe(":8080", router))
}
