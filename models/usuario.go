package models

import (
	"errors"
	"time"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/utils"
)

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
	Email         string
	password      string
	Rol           string
	Direccion     string
	Telefono      string
	FechaRegistro time.Time
}

// Función para validar los datos recibidos y crea una instancia de Usuario con la contraseña almacenada como hash
func RegistrarUsuario(nombre, email, passwordhash, rol, direccion, telefono string) (*Usuario, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	if email == "" {
		return nil, errors.New("el email no puede estar vacío")
	}

	if passwordhash == "" {
		return nil, errors.New("la contraseña no puede estar vacía")
	}

	// Genera el hash de la contraseña antes de almacenarla
	hash, err := utils.HashPassword(passwordhash)
	if err != nil {
		return nil, err
	}

	return &Usuario{
		Nombre:    nombre,
		Email:     email,
		password:  hash,
		Rol:       rol,
		Direccion: direccion,
		Telefono:  telefono,
	}, nil
}

// Métodos para gestionar el acceso y modificación de la contraseña

func (u *Usuario) ObtenerPw() string {
	return u.password
}

func (u *Usuario) AsignarPw(hash string) {
	u.password = hash
}

func (u *Usuario) CambiarPw(nueva string) error {
	if nueva == "" {
		return errors.New("la contraseña no puede estar vacía")
	}

	hash, err := utils.HashPassword(nueva)
	if err != nil {
		return err
	}

	u.password = hash
	return nil
}

func (u *Usuario) CambiarRol(rol string) error {
	if rol != "cliente" && rol != "admin" {
		return errors.New("rol inválido: debe ser 'cliente' o 'admin'")
	}

	u.Rol = rol
	return nil
}

func (u *Usuario) Autenticar(passwordhash string) error {
	return utils.CheckPassword(u.password, passwordhash)
}

// Registrar almacena un nuevo usuario en la base de datos
func (u *Usuario) Registrar() error {
	query := `
		INSERT INTO usuario (nombre, email, contrasenia, rol, direccion, telefono)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, fecha_registro
	`

	return db.DB.QueryRow(
		query,
		u.Nombre,
		u.Email,
		u.password,
		u.Rol,
		u.Direccion,
		u.Telefono,
	).Scan(&u.ID, &u.FechaRegistro)
}

// Actualizar modifica la información de un usuario existente
func (u *Usuario) Actualizar() error {
	query := `
		UPDATE usuario
		SET nombre      = $1,
		    email       = $2,
		    contrasenia = $3,
		    rol         = $4,
		    direccion   = $5,
		    telefono    = $6
		WHERE id = $7
	`

	resultado, err := db.DB.Exec(
		query,
		u.Nombre,
		u.Email,
		u.password,
		u.Rol,
		u.Direccion,
		u.Telefono,
		u.ID,
	)

	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if filas == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}

// Eliminar un usuario de la base de datos
func (u *Usuario) Eliminar() error {
	resultado, err := db.DB.Exec(
		`DELETE FROM usuario WHERE id = $1`,
		u.ID,
	)

	if err != nil {
		return err
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}

	if filas == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}

// Funciones de consulta y búsqueda de usuarios

func BuscarPorID(id int) (*Usuario, bool) {
	query := `
		SELECT id, nombre, email, contrasenia, rol, direccion, telefono, fecha_registro
		FROM usuario
		WHERE id = $1
	`

	u := &Usuario{}
	var hash string

	err := db.DB.QueryRow(query, id).Scan(
		&u.ID,
		&u.Nombre,
		&u.Email,
		&hash,
		&u.Rol,
		&u.Direccion,
		&u.Telefono,
		&u.FechaRegistro,
	)

	if err != nil {
		return nil, false
	}

	u.AsignarPw(hash)
	return u, true
}

// consulta para obtener un usuario utilizando su dirección de correo electrónico.
func BuscarPorEmail(email string) (*Usuario, bool) {
	query := `
		SELECT id, nombre, email, contrasenia, rol, direccion, telefono, fecha_registro
		FROM usuario
		WHERE email = $1
	`
	//instancia vacía donde se almacenarán los datos obtenidos de la base de datos
	u := &Usuario{}
	var hash string //// Variable temporal para almacenar la contraseña encriptada

	//asignar cada columna recuperada de la consulta a los atributos correspondientes de la estructura Usuario
	err := db.DB.QueryRow(query, email).Scan(
		&u.ID,
		&u.Nombre,
		&u.Email,
		&hash,
		&u.Rol,
		&u.Direccion,
		&u.Telefono,
		&u.FechaRegistro,
	)

	if err != nil {
		return nil, false
	}

	u.AsignarPw(hash)
	return u, true
}

// Listar todos los usuarios registrados y organizarlos en una lista para su posterior utilización.
func ListarUsuarios() ([]Usuario, error) {
	//consulta todos los usuarios de la base de datos para ordenarlos por su id
	query := `
		SELECT id, nombre, email, rol, direccion, telefono, fecha_registro
		FROM usuario
		ORDER BY id ASC
	`
	// Ejecuta la consulta y obtiene el conjunto de registros resultantes
	filas, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer filas.Close()
	// Slice que almacenará los usuarios recuperados
	var usuarios []Usuario

	for filas.Next() {
		// Variable temporal donde se almacenará la información de cada usuario recuperado
		var u Usuario
		// Asigna los valores de cada columna a los atributos correspondientes de la estructura Usuario
		err := filas.Scan(
			&u.ID,
			&u.Nombre,
			&u.Email,
			&u.Rol,
			&u.Direccion,
			&u.Telefono,
			&u.FechaRegistro,
		)

		if err != nil {
			return nil, err
		}

		usuarios = append(usuarios, u)
	}

	return usuarios, filas.Err()
}
