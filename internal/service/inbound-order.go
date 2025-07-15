package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

// InboundOrderServiceDefault implements the InboundOrderService interface.
type InboundOrderServiceDefault struct {
	repo         repository.InboundOrderRepository
	employeeRepo repository.EmployeeRepository
}

// NewInboundOrderService returns a new InboundOrderService.
func NewInboundOrderService(repo repository.InboundOrderRepository, employeeRepo repository.EmployeeRepository) InboundOrderService {
	return &InboundOrderServiceDefault{
		repo:         repo,
		employeeRepo: employeeRepo,
	}
}

// Create creates a new inbound order with the provided attributes.
// It checks for duplicate OrderNumber before persisting the new inbound order,
// and verifies that the given EmployeeID exists.
// Returns the created inbound order or an error if the operation fails.
func (s InboundOrderServiceDefault) Create(attrs models.InboundOrderAttributes) (models.InboundOrder, error) {
	// Check for duplicate order number (unique)
	existing, err := s.repo.GetByOrderNumber(attrs.OrderNumber)
	if err != nil {
		return models.InboundOrder{}, err
	}
	if existing.ID != 0 {
		return models.InboundOrder{}, httperrors.ConflictError{Message: "duplicate order number"}
	}

	// Check if employee exists
	employee, err := s.employeeRepo.GetByID(attrs.EmployeeID)
	if err != nil || employee.Id == 0 {
		return models.InboundOrder{}, httperrors.ConflictError{Message: "employee does not exist"}
	}

	newInboundOrder := models.InboundOrder{InboundOrderAttributes: attrs}
	return s.repo.Create(newInboundOrder)
}
