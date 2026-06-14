/*
@Configuracion inicial del proyecto
@autor: Nataly Tituaña

*/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/db"
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

	log.Println("Servidor corriendo correctamente")

	// Levantamos el servidor
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Error al levantar el servidor: %v", err)
	}
}
