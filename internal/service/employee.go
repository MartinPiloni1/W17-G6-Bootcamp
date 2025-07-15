package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type EmployeeServiceDefault struct {
	repo             repository.EmployeeRepository
	inboundOrderRepo repository.InboundOrderRepository
}

func NewEmployeeService(repo repository.EmployeeRepository, inboundOrderRepo repository.InboundOrderRepository) EmployeeService {
	return &EmployeeServiceDefault{repo: repo, inboundOrderRepo: inboundOrderRepo}
}

// Create creates a new employee with the provided attributes.
// It checks for duplicate CardNumberID before persisting the new employee.
// Returns the created employee or an error if the operation fails.
func (e EmployeeServiceDefault) Create(employee models.EmployeeAttributes) (models.Employee, error) {
	existing, err := e.repo.GetAll()
	if err != nil {
		return models.Employee{}, err
	}
	for _, emp := range existing {
		if emp.CardNumberID == employee.CardNumberID {
			return models.Employee{}, httperrors.ConflictError{Message: "duplicate card number"}
		}
	}

	newEmployee := models.Employee{EmployeeAttributes: employee}
	return e.repo.Create(newEmployee)
}

// GetAll retrieves all employees from the repository.
// The returned list is sorted by employee ID in ascending order.
// Returns the list of employees or an error if the operation fails.
func (e EmployeeServiceDefault) GetAll() ([]models.Employee, error) {
	employees, err := e.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return employees, nil
}

// GetByID retrieves a single employee by its unique ID from the repository.
// Returns the employee or an error if not found.
func (e EmployeeServiceDefault) GetByID(id int) (models.Employee, error) {
	return e.repo.GetByID(id)
}

// Update updates an existing employee identified by id with the provided attributes.
// It updates only the non-zero fields and checks for CardNumberID duplication.
// Returns the updated employee or an error if the operation fails.
func (e EmployeeServiceDefault) Update(id int, attrs models.EmployeeAttributes) (models.Employee, error) {
	dbEmployee, err := e.repo.GetByID(id)
	if err != nil {
		return models.Employee{}, err
	}

	if attrs.CardNumberID != "" {
		dbEmployee.CardNumberID = attrs.CardNumberID
	}
	if attrs.FirstName != "" {
		dbEmployee.FirstName = attrs.FirstName
	}
	if attrs.LastName != "" {
		dbEmployee.LastName = attrs.LastName
	}
	if attrs.WarehouseID != 0 {
		dbEmployee.WarehouseID = attrs.WarehouseID
	}

	if attrs.CardNumberID != "" {
		existing, err := e.repo.GetAll()
		if err != nil {
			return models.Employee{}, err
		}
		for _, emp := range existing {
			if emp.CardNumberID == dbEmployee.CardNumberID && emp.Id != id {
				return models.Employee{}, httperrors.ConflictError{Message: "duplicated card number"}
			}
		}
	}
	return e.repo.Update(id, dbEmployee)
}

// Delete removes the employee identified by id from the repository.
// Returns an error if the employee does not exist or the operation fails.
func (e EmployeeServiceDefault) Delete(id int) error {
	return e.repo.Delete(id)
}

func (s *EmployeeServiceDefault) ReportInboundOrders(employeeID int) ([]models.EmployeeWithInboundCount, error) {
	// Check for one employee
	if employeeID != 0 {
		emp, err := s.repo.GetByID(employeeID)
		if err != nil || emp.Id == 0 {
			return nil, httperrors.NotFoundError{Message: "employee not found"}
		}
		count, err := s.inboundOrderRepo.CountInboundOrdersForEmployee(employeeID)
		if err != nil {
			return nil, err
		}
		result := models.EmployeeWithInboundCount{
			Id: emp.Id, CardNumberID: emp.CardNumberID, FirstName: emp.FirstName, LastName: emp.LastName, WarehouseID: emp.WarehouseID, InboundOrdersCount: count,
		}
		return []models.EmployeeWithInboundCount{result}, nil
	}

	// Check for all employees
	employees, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	counts, err := s.inboundOrderRepo.CountInboundOrdersForEmployees()
	if err != nil {
		return nil, err
	}
	response := make([]models.EmployeeWithInboundCount, len(employees))
	for i, emp := range employees {
		response[i] = models.EmployeeWithInboundCount{
			Id: emp.Id, CardNumberID: emp.CardNumberID, FirstName: emp.FirstName, LastName: emp.LastName, WarehouseID: emp.WarehouseID,
			InboundOrdersCount: counts[emp.Id], // 0 by default
		}
	}
	return response, nil
}
