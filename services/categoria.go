package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// CategoriaService gestiona la lógica de negocio relacionada con las categorías.
// Actúa como intermediario entre los handlers HTTP y el repositorio de datos,
// aplicando validaciones y reglas antes de delegar la persistencia.
type CategoriaService struct {
	repo *repository.CategoriaRepository
}

// NuevoCategoriaService construye el service inyectando su repositorio correspondiente.
// La inyección de dependencias permite que el service opere sin conocer
// los detalles de implementación del acceso a datos.
func NuevoCategoriaService() *CategoriaService {
	return &CategoriaService{
		repo: repository.NuevoCategoriaRepository(),
	}
}

// CrearCategoria valida que el nombre no esté vacío y delega la persistencia al repositorio.
// La validación en el service garantiza que el repositorio reciba únicamente datos coherentes.
func (s *CategoriaService) CrearCategoria(nombre string) (*models.Categoria, error) {
	if nombre == "" {
		return nil, errors.New("el nombre de la categoría no puede estar vacío")
	}

	c := &models.Categoria{Nombre: nombre}

	// Delegar la operación de inserción al repositorio
	if err := s.repo.Insertar(c); err != nil {
		return nil, err
	}

	return c, nil
}

// ActualizarCategoria verifica la existencia del registro y aplica los nuevos valores.
// Se busca primero la entidad para garantizar que la actualización opera sobre datos reales.
func (s *CategoriaService) ActualizarCategoria(id int, nombre string) (*models.Categoria, error) {
	if nombre == "" {
		return nil, errors.New("el nombre de la categoría no puede estar vacío")
	}

	// Verificar que la categoría existe antes de intentar modificarla
	c, ok := s.repo.BuscarPorID(id)
	if !ok {
		return nil, errors.New("categoría no encontrada")
	}

	c.Nombre = nombre

	if err := s.repo.Actualizar(c); err != nil {
		return nil, err
	}

	return c, nil
}

// EliminarCategoria elimina una categoría por su identificador.
// Si la categoría tiene productos asociados, PostgreSQL retornará
// un error de integridad referencial que se propagará al handler.
func (s *CategoriaService) EliminarCategoria(id int) error {
	return s.repo.Eliminar(id)
}

// BuscarPorID delega la búsqueda al repositorio y retorna el resultado idiomático (entidad, bool).
func (s *CategoriaService) BuscarPorID(id int) (*models.Categoria, bool) {
	return s.repo.BuscarPorID(id)
}

// ListarCategorias recupera todas las categorías disponibles en el sistema.
func (s *CategoriaService) ListarCategorias() ([]models.Categoria, error) {
	return s.repo.ListarTodas()
}
