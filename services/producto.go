package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// ProductoService gestiona la lógica de negocio relacionada con los productos.
// Centraliza las validaciones de precio, stock y nombre antes de delegar
// al repositorio, protegiendo la integridad de los datos del catálogo.
type ProductoService struct {
	repo *repository.ProductoRepository
}

// NuevoProductoService construye el service inyectando su repositorio correspondiente.
func NuevoProductoService() *ProductoService {
	return &ProductoService{
		repo: repository.NuevoProductoRepository(),
	}
}

// CrearProducto valida los datos del producto y lo persiste a través del repositorio.
// Las validaciones aseguran que ningún producto ingrese al sistema con datos incoherentes.
func (s *ProductoService) CrearProducto(categoriaID int, nombre, descripcion string, precio float64, stock int) (*models.Producto, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del producto no puede estar vacío")
	}
	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}
	if stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	p := &models.Producto{
		CategoriaID: categoriaID,
		Nombre:      nombre,
		Descripcion: descripcion,
		Precio:      precio,
		Stock:       stock,
	}

	if err := s.repo.Insertar(p); err != nil {
		return nil, err
	}

	return p, nil
}

// ActualizarProducto verifica la existencia del producto y aplica los nuevos valores.
func (s *ProductoService) ActualizarProducto(id, categoriaID int, nombre, descripcion string, precio float64, stock int) (*models.Producto, error) {
	if nombre == "" {
		return nil, errors.New("el nombre del producto no puede estar vacío")
	}
	if precio <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}
	if stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	// Verificar que el producto existe antes de modificarlo
	p, ok := s.repo.BuscarPorID(id)
	if !ok {
		return nil, errors.New("producto no encontrado")
	}

	p.CategoriaID = categoriaID
	p.Nombre = nombre
	p.Descripcion = descripcion
	p.Precio = precio
	p.Stock = stock

	if err := s.repo.Actualizar(p); err != nil {
		return nil, err
	}

	return p, nil
}

// ActualizarStock aplica la regla de negocio que impide que el stock sea negativo.
// Recibe un delta: positivo para aumentar existencias, negativo para reducirlas.
func (s *ProductoService) ActualizarStock(id, cantidad int) error {
	p, ok := s.repo.BuscarPorID(id)
	if !ok {
		return errors.New("producto no encontrado")
	}

	// Regla de negocio: el stock resultante nunca puede ser negativo
	if p.Stock+cantidad < 0 {
		return errors.New("stock insuficiente para realizar la operación")
	}

	p.Stock += cantidad
	return s.repo.Actualizar(p)
}

// EliminarProducto elimina un producto del sistema por su identificador.
func (s *ProductoService) EliminarProducto(id int) error {
	return s.repo.Eliminar(id)
}

// BuscarPorID recupera un producto por su identificador único.
func (s *ProductoService) BuscarPorID(id int) (*models.Producto, bool) {
	return s.repo.BuscarPorID(id)
}

// ListarProductos recupera todos los productos disponibles en el catálogo.
func (s *ProductoService) ListarProductos() ([]models.Producto, error) {
	return s.repo.ListarTodos()
}
