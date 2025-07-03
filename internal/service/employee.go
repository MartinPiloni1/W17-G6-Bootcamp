package service

import (
	"slices"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type EmployeeServiceDefault struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &EmployeeServiceDefault{repo: repo}
}

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

func (e EmployeeServiceDefault) GetAll() ([]models.Employee, error) {
	employees, err := e.repo.GetAll()
	if err != nil {
		return nil, err
	}

	slicedData := utils.MapToSlice(employees)
	slices.SortFunc(slicedData, func(a, b models.Employee) int {
		return a.Id - b.Id
	})
	return slicedData, nil
}

func (e EmployeeServiceDefault) GetByID(id int) (models.Employee, error) {
	return e.repo.GetByID(id)
}

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

func (e EmployeeServiceDefault) Delete(id int) error {
	return e.repo.Delete(id)
}
