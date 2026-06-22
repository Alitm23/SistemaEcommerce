package models

import "time"

// Interfaz que define las operaciones básicas de gestión de usuarios
type ControlUsuario interface {
	Registrar() error
	Actualizar() error
	Eliminar() error
}

// Usuario representa la información de un usuario dentro del sistema
type Usuario struct {
	ID            int
	Nombre        string
	Apellido      string
	Email         string
	Password      string
	Rol           string
	Direccion     string
	Telefono      string
	FechaRegistro time.Time
}
