/*
@Configuracion inicial del proyecto
@autor: Nataly Tituaña

*/

package main

import (
	"fmt"
	"log"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

func TestAgregarNuevoProducto() {
	fmt.Println("Iniciando prueba de registro...")
	// Crear categoría
	categoria, err := models.NuevaCategoria("Aretes")
	if err != nil {
		log.Printf("Error al crear la categoría: %v\n", err)
		return
	}

	err = categoria.Registrar()
	if err != nil {
		log.Printf("Error al registrar la categoría: %v\n", err)
		return
	}

	fmt.Println("\n=== CATEGORÍA REGISTRADA ===")
	fmt.Printf("ID: %d\n", categoria.ID)
	fmt.Printf("Nombre: %s\n", categoria.Nombre)

	// Crear producto utilizando el ID de la categoría
	producto, err := models.NuevoProducto(
		categoria.ID,
		"Cristal",
		"Elegantes aretes artesanales que destacan por su delicado brillo y diseño sofisticado",
		15.99,
		2,
	)

	if err != nil {
		log.Printf("Error al crear el producto: %v\n", err)
		return
	}

	err = producto.Registrar()
	if err != nil {
		log.Printf("Error al registrar el producto: %v\n", err)
		return
	}

	fmt.Println("\n=== PRODUCTO REGISTRADO ===")
	fmt.Printf("ID: %d\n", producto.ID)
	fmt.Printf("Categoría ID: %d\n", producto.CategoriaID)
	fmt.Printf("Nombre: %s\n", producto.Nombre)
	fmt.Printf("Descripción: %s\n", producto.Descripcion)
	fmt.Printf("Precio: %.2f\n", producto.Precio)
	fmt.Printf("Stock: %d\n", producto.Stock)
	fmt.Printf("Fecha Ingreso: %v\n", producto.FechaIngreso)
}

func main() {
	// Verificar si se ejecuta el ping de la conexion
	DB, err := db.Connect()
	if err != nil {
		// Fatalf nos permite inyectar el error dentro del texto
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	fmt.Println("Conectado exitosamente a la base de datos")
	defer DB.Close()

	TestAgregarNuevoProducto()

	log.Println("Servidor corriendo correctamente")

	// Levantamos el servidor
	//if err := http.ListenAndServe(":8081", nil); err != nil {
	//log.Fatalf("Error al levantar el servidor: %v", err)
	//}
}
