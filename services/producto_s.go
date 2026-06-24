package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// ProductoService gestiona la lógica de negocio relacionada con los productos.
// Centraliza las validaciones de precio y nombre antes de delegar al repositorioy administra las tallas que determinan el stock real de cada joya.
type ProductoService struct {
	repo      *repository.ProductoRepository
	tallaRepo *repository.ProductoTallaRepository
}

// NuevoProductoService construye el service inyectando los repositorios necesarios como las tallas
func NuevoProductoService() *ProductoService {
	return &ProductoService{
		repo:      repository.NuevoProductoRepository(),
		tallaRepo: repository.NuevoProductoTallaRepository(),
	}
}

// CrearProducto valida los datos del producto y lo persiste a través del repositorio.
func (s *ProductoService) CrearProducto(categoriaID int, nombre, descripcion, material string, precio float64) (*models.Producto, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del producto no puede estar vacío")
	}
	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}

	p := &models.Producto{
		CategoriaID: categoriaID,
		Nombre:      nombre,
		Material:    material,
		Descripcion: descripcion,
		Precio:      precio,
	}

	if err := s.repo.Insertar(p); err != nil {
		return nil, err
	}

	return p, nil
}

// ActualizarProducto verifica la existencia del producto y aplica los nuevos valores
func (s *ProductoService) ActualizarProducto(id, categoriaID int, nombre, descripcion, material string, precio float64) (*models.Producto, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del producto no puede estar vacío")
	}
	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}

	// Verificar que el producto existe antes de modificarlo
	p, ok := s.repo.BuscarPorID(id)
	if !ok {
		return nil, errors.New("producto no encontrado")
	}

	p.CategoriaID = categoriaID
	p.Nombre = nombre
	p.Descripcion = descripcion
	p.Material = material
	p.Precio = precio

	if err := s.repo.Actualizar(p); err != nil {
		return nil, err
	}

	return p, nil
}

// AgregarTalla agrega una nueva talla con su stock inicial a un producto existente
func (s *ProductoService) AgregarTalla(productoID int, talla string, stock int) (*models.ProductoTalla, error) {
	if talla == "" {
		return nil, errors.New("la talla no puede estar vacía")
	}
	if stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	// Verificar que el producto existe antes de agregar la talla
	if _, ok := s.repo.BuscarPorID(productoID); !ok {
		return nil, errors.New("producto no encontrado")
	}

	pt := &models.ProductoTalla{
		ProductoID: productoID,
		Talla:      talla,
		Stock:      stock,
	}

	if err := s.tallaRepo.Insertar(pt); err != nil {
		return nil, err
	}

	return pt, nil
}

// ActualizarStockTalla aplica la regla de negocio que impide que el stock sea negativo.
func (s *ProductoService) ActualizarStockTalla(tallaID, delta int) error {
	pt, ok := s.tallaRepo.BuscarPorID(tallaID)
	if !ok {
		return errors.New("talla no encontrada")
	}

	// Regla de negocio: el stock resultante nunca puede ser negativo
	if pt.Stock+delta < 0 {
		return errors.New("stock insuficiente para realizar la operación")
	}

	pt.Stock += delta
	return s.tallaRepo.ActualizarStock(pt)
}

// EliminarTalla elimina una talla específica de un producto por su identificador
func (s *ProductoService) EliminarTalla(tallaID int) error {
	return s.tallaRepo.Eliminar(tallaID)
}

// ObtenerTallas recupera todas las tallas disponibles para un producto específico
func (s *ProductoService) ObtenerTallas(productoID int) ([]models.ProductoTalla, error) {
	return s.tallaRepo.ListarPorProducto(productoID)
}

// BuscarTallaPorID recupera una talla específica por su identificador único
func (s *ProductoService) BuscarTallaPorID(tallaID int) (*models.ProductoTalla, bool) {
	return s.tallaRepo.BuscarPorID(tallaID)
}

// EliminarProducto elimina un producto del sistema por su identificador.
func (s *ProductoService) EliminarProducto(id int) error {
	return s.repo.Eliminar(id)
}

// BuscarPorID recupera un producto por su identificador único
func (s *ProductoService) BuscarPorID(id int) (*models.Producto, bool) {
	return s.repo.BuscarPorID(id)
}

// ListarProductos recupera todos los productos disponibles en el catálogo
func (s *ProductoService) ListarProductos() ([]models.Producto, error) {
	return s.repo.ListarTodos()
}

// ListarPorCategoria recupera todos los productos de una categoría específica
func (s *ProductoService) ListarPorCategoria(categoriaID int) ([]models.Producto, error) {
	return s.repo.ListarPorCategoria(categoriaID)
}
