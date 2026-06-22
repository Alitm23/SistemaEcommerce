package repository

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/db"
	"github.com/Alitm23/SistemaEcommerce/models"
)

// UsuarioRepository ejecuta todas las operaciones SQL relacionadas con Usuario
type UsuarioRepository struct{}

// NuevoUsuarioRepository construye una instancia del repositorio
func NuevoUsuarioRepository() *UsuarioRepository {
	return &UsuarioRepository{}
}

// Insertar persiste un nuevo usuario en la base de datos
func (r *UsuarioRepository) Insertar(u *models.Usuario) error {
	query := `
		INSERT INTO usuario (nombre, apellido, email, password, rol, direccion, telefono)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, fecha_registro
	`
	return db.DB.QueryRow(
		query,
		u.Nombre, u.Apellido, u.Email, u.Password,
		u.Rol, u.Direccion, u.Telefono,
	).Scan(&u.ID, &u.FechaRegistro)
}

// Actualizar modifica los datos de un usuario existente
func (r *UsuarioRepository) Actualizar(u *models.Usuario) error {
	query := `
		UPDATE usuario
		SET nombre = $1, apellido = $2, email = $3,
		    password = $4, rol = $5, direccion = $6, telefono = $7
		WHERE id = $8
	`
	resultado, err := db.DB.Exec(
		query,
		u.Nombre, u.Apellido, u.Email, u.Password,
		u.Rol, u.Direccion, u.Telefono, u.ID,
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

// Eliminar borra un usuario de la base de datos por su ID
func (r *UsuarioRepository) Eliminar(id int) error {
	resultado, err := db.DB.Exec(`DELETE FROM usuario WHERE id = $1`, id)
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

// BuscarPorID recupera un usuario por su ID
func (r *UsuarioRepository) BuscarPorID(id int) (*models.Usuario, bool) {
	query := `
		SELECT id, nombre, apellido, email, password,
		       rol, direccion, telefono, fecha_registro
		FROM usuario
		WHERE id = $1
	`
	u := &models.Usuario{}
	err := db.DB.QueryRow(query, id).Scan(
		&u.ID, &u.Nombre, &u.Apellido, &u.Email, &u.Password,
		&u.Rol, &u.Direccion, &u.Telefono, &u.FechaRegistro,
	)
	if err != nil {
		return nil, false
	}
	return u, true
}

// BuscarPorEmail recupera un usuario por su correo electrónico
func (r *UsuarioRepository) BuscarPorEmail(email string) (*models.Usuario, bool) {
	query := `
		SELECT id, nombre, apellido, email, password,
		       rol, direccion, telefono, fecha_registro
		FROM usuario
		WHERE email = $1
	`
	u := &models.Usuario{}
	err := db.DB.QueryRow(query, email).Scan(
		&u.ID, &u.Nombre, &u.Apellido, &u.Email, &u.Password,
		&u.Rol, &u.Direccion, &u.Telefono, &u.FechaRegistro,
	)
	if err != nil {
		return nil, false
	}
	return u, true
}

// ListarTodos recupera todos los usuarios ordenados por ID
func (r *UsuarioRepository) ListarTodos() ([]models.Usuario, error) {
	query := `
		SELECT id, nombre, apellido, email,
		       rol, direccion, telefono, fecha_registro
		FROM usuario
		ORDER BY id ASC
	`
	filas, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer filas.Close()

	var usuarios []models.Usuario

	for filas.Next() {
		var u models.Usuario
		err := filas.Scan(
			&u.ID, &u.Nombre, &u.Apellido, &u.Email,
			&u.Rol, &u.Direccion, &u.Telefono, &u.FechaRegistro,
		)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}

	return usuarios, filas.Err()
}
