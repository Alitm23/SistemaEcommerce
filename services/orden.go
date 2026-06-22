package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// OrdenService gestiona la lógica de negocio relacionada con las órdenes de compra.
// Implementa la máquina de estados que controla las transiciones válidas
// entre los distintos estados del ciclo de vida de una orden.
type OrdenService struct {
	repo     *repository.OrdenRepository
	itemRepo *repository.ItemOrdenRepository
}

// NuevoOrdenService construye el service inyectando los repositorios necesarios.
func NuevoOrdenService() *OrdenService {
	return &OrdenService{
		repo:     repository.NuevoOrdenRepository(),
		itemRepo: repository.NuevoItemOrdenRepository(),
	}
}

// GenerarOrden crea una nueva orden con estado inicial 'pendiente'.
// Valida que el total sea positivo antes de persistir.
func (s *OrdenService) GenerarOrden(usuarioID int, total float64) (*models.Orden, error) {
	if total <= 0 {
		return nil, errors.New("el total de la orden debe ser mayor a cero")
	}

	o := &models.Orden{
		UsuarioID: usuarioID,
		Total:     total,
		Estado:    "pendiente",
	}

	if err := s.repo.Insertar(o); err != nil {
		return nil, err
	}

	return o, nil
}

// ActualizarEstado aplica la máquina de estados de la orden.
// Una orden entregada no puede cambiar de estado — regla de negocio irreversible.
func (s *OrdenService) ActualizarEstado(id int, nuevoEstado string) error {
	estados := map[string]bool{
		"pendiente":  true,
		"procesando": true,
		"enviada":    true,
		"entregada":  true,
		"cancelada":  true,
	}

	if !estados[nuevoEstado] {
		return errors.New("estado de orden inválido")
	}

	o, ok := s.repo.BuscarPorID(id)
	if !ok {
		return errors.New("orden no encontrada")
	}

	// Regla de negocio: una orden entregada es un estado terminal irreversible
	if o.Estado == "entregada" {
		return errors.New("una orden entregada no puede cambiar de estado")
	}

	o.Estado = nuevoEstado
	return s.repo.ActualizarEstado(o)
}

// CancelarOrden establece el estado de la orden como 'cancelada'.
// Reutiliza ActualizarEstado para aplicar las mismas validaciones.
func (s *OrdenService) CancelarOrden(id int) error {
	return s.ActualizarEstado(id, "cancelada")
}

// AgregarItem agrega un producto a la orden con su precio histórico de compra.
// El precio se guarda en el ítem para preservar el valor al momento de la transacción.
func (s *OrdenService) AgregarItem(ordenID, productoID, cantidad int, precioCompra float64) (*models.ItemOrden, error) {
	if cantidad <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a cero")
	}
	if precioCompra <= 0 {
		return nil, errors.New("el precio de compra debe ser mayor a cero")
	}

	item := &models.ItemOrden{
		OrdenID:      ordenID,
		ProductoID:   productoID,
		Cantidad:     cantidad,
		PrecioCompra: precioCompra,
	}

	if err := s.itemRepo.Insertar(item); err != nil {
		return nil, err
	}

	return item, nil
}

// ObtenerItems recupera todos los ítems asociados a una orden específica.
func (s *OrdenService) ObtenerItems(ordenID int) ([]models.ItemOrden, error) {
	return s.itemRepo.ListarPorOrden(ordenID)
}

// BuscarPorID recupera una orden por su identificador único.
func (s *OrdenService) BuscarPorID(id int) (*models.Orden, bool) {
	return s.repo.BuscarPorID(id)
}

// ListarPorUsuario recupera todas las órdenes de un usuario específico.
func (s *OrdenService) ListarPorUsuario(usuarioID int) ([]models.Orden, error) {
	return s.repo.ListarPorUsuario(usuarioID)
}

// ListarTodas recupera todas las órdenes del sistema.
func (s *OrdenService) ListarTodas() ([]models.Orden, error) {
	return s.repo.ListarTodas()
}
