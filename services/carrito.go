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
	repo     *repository.CarritoRepository
	itemRepo *repository.ItemCarritoRepository
}

// NuevoCarritoService construye el service inyectando los repositorios necesarios.
// Se inyectan dos repositorios porque el carrito gestiona también sus ítems.
func NuevoCarritoService() *CarritoService {
	return &CarritoService{
		repo:     repository.NuevoCarritoRepository(),
		itemRepo: repository.NuevoItemCarritoRepository(),
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

	// Validar la transición de estado antes de persistir
	if c.Estado != "activo" {
		return errors.New("solo se puede cerrar un carrito en estado activo")
	}

	c.Estado = "cerrado"
	return s.repo.ActualizarEstado(c)
}

// AgregarItem agrega un producto al carrito con su cantidad y precio unitario.
// Valida que la cantidad sea positiva antes de registrar el ítem.
func (s *CarritoService) AgregarItem(carritoID, productoID, cantidad int, precioUnitario float64) (*models.ItemCarrito, error) {
	if cantidad <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a cero")
	}
	if precioUnitario <= 0 {
		return nil, errors.New("el precio unitario debe ser mayor a cero")
	}

	item := &models.ItemCarrito{
		CarritoID:      carritoID,
		ProductoID:     productoID,
		Cantidad:       cantidad,
		PrecioUnitario: precioUnitario,
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
// El total se calcula en el service porque es lógica de presentación de negocio,
// no una responsabilidad del repositorio ni del handler.
func (s *CarritoService) ObtenerItems(carritoID int) ([]models.ItemCarrito, float64, error) {
	items, err := s.itemRepo.ListarPorCarrito(carritoID)
	if err != nil {
		return nil, 0, err
	}

	// Calcular el total sumando el subtotal de cada ítem
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
