package services

import (
	"errors"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/repository"
)

// PagoService gestiona la lógica de negocio relacionada con los pagos.
// Controla las transiciones de estado del pago aplicando reglas que
// protegen la integridad financiera del sistema.
type PagoService struct {
	repo *repository.PagoRepository
}

// NuevoPagoService construye el service inyectando su repositorio correspondiente.
func NuevoPagoService() *PagoService {
	return &PagoService{
		repo: repository.NuevoPagoRepository(),
	}
}

// RegistrarPago crea un nuevo pago en estado inicial 'pendiente'.
// Valida que el monto sea positivo antes de persistir.
func (s *PagoService) RegistrarPago(ordenID int, monto float64) (*models.Pago, error) {
	if monto <= 0 {
		return nil, errors.New("el monto del pago debe ser mayor a cero")
	}

	p := &models.Pago{
		OrdenID: ordenID,
		Monto:   monto,
		Estado:  "pendiente",
	}

	if err := s.repo.Insertar(p); err != nil {
		return nil, err
	}

	return p, nil
}

// ActualizarEstado aplica las reglas de transición de estado del pago.
// Un pago completado no puede volver a estado pendiente — regla de integridad financiera.
func (s *PagoService) ActualizarEstado(ordenID int, nuevoEstado string) error {
	estados := map[string]bool{
		"pendiente":  true,
		"completado": true,
		"fallido":    true,
	}

	if !estados[nuevoEstado] {
		return errors.New("estado de pago inválido")
	}

	p, ok := s.repo.BuscarPorOrden(ordenID)
	if !ok {
		return errors.New("pago no encontrado")
	}

	// Regla de negocio: un pago completado no puede revertirse a pendiente
	if p.Estado == "completado" && nuevoEstado == "pendiente" {
		return errors.New("un pago completado no puede volver a estado pendiente")
	}

	p.Estado = nuevoEstado
	return s.repo.ActualizarEstado(p)
}

// AnularPago marca el pago como 'fallido' si aún no fue completado.
// Reutiliza ActualizarEstado para aplicar las mismas validaciones de transición.
func (s *PagoService) AnularPago(ordenID int) error {
	p, ok := s.repo.BuscarPorOrden(ordenID)
	if !ok {
		return errors.New("pago no encontrado")
	}

	// Regla de negocio: un pago completado no puede ser anulado
	if p.Estado == "completado" {
		return errors.New("no se puede anular un pago que ya fue completado")
	}

	return s.ActualizarEstado(ordenID, "fallido")
}

// ListarPagos recupera todos los pagos registrados en el sistema.
func (s *PagoService) ListarPagos() ([]models.Pago, error) {
	return s.repo.ListarTodos()
}
