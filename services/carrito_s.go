package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// CarritoService gestiona la lógica de negocio del carrito de compras.
// Controla el ciclo de vida del carrito y la gestión de sus ítems,
// aplicando reglas como la unicidad del carrito activo por usuario.
type CarritoService struct {
	repo      *repository.CarritoRepository
	itemRepo  *repository.ItemCarritoRepository
	tallaRepo *repository.ProductoTallaRepository
}

// NuevoCarritoService construye el service inyectando los repositorios necesarios.
func NuevoCarritoService() *CarritoService {
	return &CarritoService{
		repo:      repository.NuevoCarritoRepository(),
		itemRepo:  repository.NuevoItemCarritoRepository(),
		tallaRepo: repository.NuevoProductoTallaRepository(),
	}
}

// AbrirCarrito crea un nuevo carrito en estado 'activo' para el usuario indicado.
// Verifica previamente que el usuario no tenga ya un carrito activo.
func (s *CarritoService) AbrirCarrito(usuarioID int) (*models.Carrito, error) {
	if usuarioID <= 0 {
		return nil, errors.New("identificador de usuario inválido")
	}

	// Regla de negocio: un usuario solo puede tener un carrito activo a la vez
	if _, ok := s.repo.BuscarActivoPorUsuario(usuarioID); ok {
		return nil, errors.New("el usuario ya tiene un carrito activo")
	}

	c := &models.Carrito{
		UsuarioID: usuarioID,
		Estado:    "activo",
	}

	if err := s.repo.Insertar(c); err != nil {
		return nil, err
	}

	return c, nil
}

// CerrarCarrito cambia el estado del carrito a 'cerrado'.
// Solo los carritos en estado 'activo' pueden cerrarse.
func (s *CarritoService) CerrarCarrito(id int) error {
	c, ok := s.repo.BuscarPorID(id)
	if !ok {
		return errors.New("carrito no encontrado")
	}

	if c.Estado != "activo" {
		return errors.New("solo se puede cerrar un carrito en estado activo")
	}

	c.Estado = "cerrado"
	return s.repo.ActualizarEstado(c)
}

// AgregarItem agrega una talla de producto al carrito con su cantidad y precio unitario.
// Valida que exista stock suficiente en la talla seleccionada antes de agregar el ítem.
func (s *CarritoService) AgregarItem(carritoID, productoTallaID, cantidad int, precioUnitario float64) (*models.ItemCarrito, error) {
	if cantidad <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a cero")
	}
	if precioUnitario <= 0 {
		return nil, errors.New("el precio unitario debe ser mayor a cero")
	}

	// Regla de negocio: verificar que la talla existe y tiene stock suficiente
	pt, ok := s.tallaRepo.BuscarPorID(productoTallaID)
	if !ok {
		return nil, errors.New("la talla seleccionada no existe")
	}
	if pt.Stock < cantidad {
		return nil, errors.New("stock insuficiente para la talla seleccionada")
	}

	item := &models.ItemCarrito{
		CarritoID:       carritoID,
		ProductoTallaID: productoTallaID,
		Cantidad:        cantidad,
		PrecioUnitario:  precioUnitario,
	}

	if err := s.itemRepo.Insertar(item); err != nil {
		return nil, err
	}

	return item, nil
}

// ActualizarCantidadItem modifica la cantidad de un ítem existente en el carrito.
func (s *CarritoService) ActualizarCantidadItem(itemID, cantidad int) error {
	if cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}

	item := &models.ItemCarrito{ID: itemID, Cantidad: cantidad}
	return s.itemRepo.ActualizarCantidad(item)
}

// QuitarItem elimina un ítem del carrito por su identificador.
func (s *CarritoService) QuitarItem(itemID int) error {
	return s.itemRepo.Eliminar(itemID)
}

// ObtenerItems recupera todos los ítems de un carrito y calcula el total acumulado.
// El total se calcula en el service porque es lógica de negocio, no responsabilidad del repositorio.
func (s *CarritoService) ObtenerItems(carritoID int) ([]models.ItemCarrito, float64, error) {
	items, err := s.itemRepo.ListarPorCarrito(carritoID)
	if err != nil {
		return nil, 0, err
	}

	var total float64
	for _, item := range items {
		total += float64(item.Cantidad) * item.PrecioUnitario
	}

	return items, total, nil
}

// BuscarPorID recupera un carrito por su identificador único.
func (s *CarritoService) BuscarPorID(id int) (*models.Carrito, bool) {
	return s.repo.BuscarPorID(id)
}
