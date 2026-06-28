package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/utils"
	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Rol      string `json:"rol"`
}

// RegistrarUsuario encripta la contraseña e inserta al nuevo cliente (Servicio Web 2)
func RegistrarUsuario(u Usuario) error {
	if strings.TrimSpace(u.Nombre) == "" {
		return errors.New("el nombre es obligatorio")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("el email es obligatorio")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("la contrasena es obligatoria")
	}

	hash, err := utils.HashPassword(u.Password) // Encriptamos la contraseña
	if err != nil {
		return err
	}

	// Por defecto, todos los registros públicos son "cliente"
	_, err = db.DB.Exec(`
		INSERT INTO usuario (nombre, apellido, email, password, rol) 
		VALUES ($1, $2, $3, $4, 'cliente')`,
		u.Nombre, u.Apellido, u.Email, string(hash),
	)
	return err
}

// AutenticarUsuario verifica las credenciales y devuelve el usuario si son correctas (Servicio Web 1)
func AutenticarUsuario(email, password string) (*Usuario, error) {
	var u Usuario
	var hashPassword string

	query := "SELECT id, nombre, apellido, email, password, rol FROM usuario WHERE email = $1"
	err := db.DB.QueryRow(query, email).Scan(&u.ID, &u.Nombre, &u.Apellido, &u.Email, &hashPassword, &u.Rol)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("credenciales incorrectas")
		}
		return nil, err
	}

	// Comparamos el hash de la BD con la contraseña ingresada
	err = utils.CheckPassword(hashPassword, password)
	if err != nil {
		return nil, errors.New("credenciales incorrectas")
	}
	return &u, nil
}

// ObtenerUsuarioPorID busca un cliente por su identificador único
func ObtenerUsuarioPorID(id int) (*Usuario, error) {
	var u Usuario
	err := db.DB.QueryRow(
		"SELECT id, nombre, apellido, email, rol FROM usuario WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Nombre, &u.Apellido, &u.Email, &u.Rol)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, errors.New("error al obtener el usuario")
	}
	return &u, nil
}

// ListarUsuarios devuelve todos los clientes y administradores para el panel.
func ListarUsuarios() ([]Usuario, error) {
	filas, err := db.DB.Query("SELECT id, nombre, apellido, email, rol FROM usuario ORDER BY id")
	if err != nil {
		return nil, errors.New("error al listar usuarios")
	}
	defer filas.Close()

	var usuarios []Usuario
	for filas.Next() {
		var u Usuario
		if err := filas.Scan(&u.ID, &u.Nombre, &u.Apellido, &u.Email, &u.Rol); err != nil {
			return nil, errors.New("error al leer usuarios")
		}
		usuarios = append(usuarios, u)
	}
	if err := filas.Err(); err != nil {
		return nil, errors.New("error al recorrer usuarios")
	}
	return usuarios, nil
}

// ActualizarPerfil modifica los datos básicos del cliente
func ActualizarPerfil(u Usuario) error {
	_, err := db.DB.Exec(
		"UPDATE usuario SET nombre = $1, apellido = $2, email = $3 WHERE id = $4",
		u.Nombre, u.Apellido, u.Email, u.ID,
	)
	return err
}

// CambiarPassword encripta la nueva contraseña y la actualiza
func CambiarPassword(id int, nuevaPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(nuevaPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec("UPDATE usuario SET password = $1 WHERE id = $2", string(hash), id)
	return err
}

// ActualizarRol permite a un administrador dar permisos a otro usuario
func ActualizarRol(id int, nuevoRol string) error {
	_, err := db.DB.Exec("UPDATE usuario SET rol = $1 WHERE id = $2", nuevoRol, id)
	return err
}

// EliminarUsuario borra un registro de cliente de la base de datos
func EliminarUsuario(id int) error {
	_, err := db.DB.Exec("DELETE FROM usuario WHERE id = $1", id)
	return err
}
