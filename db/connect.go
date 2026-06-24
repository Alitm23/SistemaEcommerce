package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// estructurar todos los parametros hacia la base de datos

var DB *sql.DB

func Connect() (*sql.DB, error) {
	// variables de entorno
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// formato de dsn para postgresql
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//crear la conexion con postgresql
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	//establecer la conexión y hacer un ping permanente a la base de datos
	if err := db.Ping(); err != nil {
		return nil, err
	}
	DB = db

	log.Println("Se realizó la conexión a base de datos correctamente")
	return db, nil
}
