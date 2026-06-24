package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
	"github.com/Alitm23/SistemaEcommerce/utils"
)

// UsuarioService gestiona la lógica de negocio relacionada con usuarios
type UsuarioService struct {
	repo *repository.UsuarioRepository
}

// NuevoUsuarioService construye el service inyectando el repositorio
func NuevoUsuarioService() *UsuarioService {
	return &UsuarioService{
		repo: repository.NuevoUsuarioRepository(),
	}
}

// RegistrarUsuario valida los datos, hashea la contraseña y persiste el usuario
func (s *UsuarioService) RegistrarUsuario(nombre, apellido, email, password, rol, direccion, telefono string) (*models.Usuario, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if email == "" {
		return nil, errors.New("el email no puede estar vacío")
	}
	if password == "" {
		return nil, errors.New("la contraseña no puede estar vacía")
	}
	if rol == "" {
		rol = "cliente"
	}

	// Hashear la contraseña antes de construir el struct
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	u := &models.Usuario{
		Nombre:    nombre,
		Apellido:  apellido,
		Email:     email,
		Password:  hash,
		Rol:       rol,
		Direccion: direccion,
		Telefono:  telefono,
	}

	// Delegar la persistencia al repositorio
	if err := s.repo.Insertar(u); err != nil {
		return nil, err
	}

	return u, nil
}

// ActualizarUsuario aplica los cambios sobre un usuario existente
func (s *UsuarioService) ActualizarUsuario(id int, nombre, apellido, email, direccion, telefono string) (*models.Usuario, error) {
	// Verificar que el usuario existe antes de modificarlo
	u, ok := s.repo.BuscarPorID(id)
	if !ok {
		return nil, errors.New("usuario no encontrado")
	}

	u.Nombre = nombre
	u.Apellido = apellido
	u.Email = email
	u.Direccion = direccion
	u.Telefono = telefono

	if err := s.repo.Actualizar(u); err != nil {
		return nil, err
	}

	return u, nil
}

// CambiarPassword hashea y actualiza la contraseña del usuario
func (s *UsuarioService) CambiarPassword(id int, nuevaPassword string) error {
	if nuevaPassword == "" {
		return errors.New("la contraseña no puede estar vacía")
	}

	u, ok := s.repo.BuscarPorID(id)
	if !ok {
		return errors.New("usuario no encontrado")
	}

	hash, err := utils.HashPassword(nuevaPassword)
	if err != nil {
		return err
	}

	u.Password = hash
	return s.repo.Actualizar(u)
}

// CambiarRol valida y actualiza el rol del usuario
func (s *UsuarioService) CambiarRol(id int, rol string) error {
	if rol != "cliente" && rol != "admin" {
		return errors.New("rol inválido: debe ser 'cliente' o 'admin'")
	}

	u, ok := s.repo.BuscarPorID(id)
	if !ok {
		return errors.New("usuario no encontrado")
	}

	u.Rol = rol
	return s.repo.Actualizar(u)
}

// EliminarUsuario elimina un usuario por su ID
func (s *UsuarioService) EliminarUsuario(id int) error {
	return s.repo.Eliminar(id)
}

// Autenticar verifica las credenciales del usuario
func (s *UsuarioService) Autenticar(email, password string) (*models.Usuario, error) {
	u, ok := s.repo.BuscarPorEmail(email)
	if !ok {
		return nil, errors.New("credenciales inválidas")
	}

	if err := utils.CheckPassword(u.Password, password); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return u, nil
}

// BuscarPorID delega la búsqueda al repositorio
func (s *UsuarioService) BuscarPorID(id int) (*models.Usuario, bool) {
	return s.repo.BuscarPorID(id)
}

// ListarUsuarios recupera todos los usuarios
func (s *UsuarioService) ListarUsuarios() ([]models.Usuario, error) {
	return s.repo.ListarTodos()
}
