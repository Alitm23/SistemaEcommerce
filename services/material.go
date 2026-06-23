package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// MaterialService gestiona la lógica de negocio relacionada con los materiales.
// Los materiales son datos maestros que clasifican las joyas por composición
// (oro, plata, especiales) y se asignan al crear o editar un producto.
type MaterialService struct {
	repo *repository.MaterialRepository
}

// NuevoMaterialService construye el service inyectando su repositorio correspondiente
func NuevoMaterialService() *MaterialService {
	return &MaterialService{
		repo: repository.NuevoMaterialRepository(),
	}
}

// CrearMaterial valida que el nombre no esté vacío y delega la persistencia al repositorio
func (s *MaterialService) CrearMaterial(nombre string) (*models.Material, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del material no puede estar vacío")
	}

	m := &models.Material{Nombre: nombre}

	if err := s.repo.Insertar(m); err != nil {
		return nil, err
	}

	return m, nil
}

// ActualizarMaterial verifica la existencia del material y aplica el nuevo nombre
func (s *MaterialService) ActualizarMaterial(id int, nombre string) (*models.Material, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del material no puede estar vacío")
	}

	// Verificar que el material existe antes de modificarlo
	m, ok := s.repo.BuscarPorID(id)
	if !ok {
		return nil, errors.New("material no encontrado")
	}

	m.Nombre = nombre

	if err := s.repo.Actualizar(m); err != nil {
		return nil, err
	}

	return m, nil
}

// EliminarMaterial elimina un material por su identificador.
// Si tiene productos asociados, PostgreSQL retornará un error de integridad referencial.
func (s *MaterialService) EliminarMaterial(id int) error {
	return s.repo.Eliminar(id)
}

// BuscarPorID recupera un material por su identificador único
func (s *MaterialService) BuscarPorID(id int) (*models.Material, bool) {
	return s.repo.BuscarPorID(id)
}

// ListarMateriales recupera todos los materiales disponibles en el sistema
func (s *MaterialService) ListarMateriales() ([]models.Material, error) {
	return s.repo.ListarTodos()
}
