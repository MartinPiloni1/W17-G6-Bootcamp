package service

import (
	"errors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type EmployeeServiceDefault struct {
	repo          repository.EmployeeRepository
	warehouseRepo repository.WarehouseRepository
}

func NewEmployeeService(repo repository.EmployeeRepository, warehouseRepo repository.WarehouseRepository) EmployeeService {
	return &EmployeeServiceDefault{repo: repo, warehouseRepo: warehouseRepo}
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

	_, err = e.warehouseRepo.GetByID(employee.WarehouseID)
	if err != nil {
		var notFoundError httperrors.NotFoundError
		if errors.As(err, &notFoundError) {
			return models.Employee{}, httperrors.UnprocessableEntityError{Message: "warehouse_id does not exist"}
		}
		return models.Employee{}, err
	}

	newEmployee := models.Employee{EmployeeAttributes: employee}
	return e.repo.Create(newEmployee)
}

func (e EmployeeServiceDefault) GetAll() ([]models.Employee, error) {
	employees, err := e.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return utils.MapToSlice[models.Employee](employees), nil
}

func (e EmployeeServiceDefault) GetByID(id int) (models.Employee, error) {
	return e.repo.GetByID(id)
}

func (e EmployeeServiceDefault) Update(id int, employee models.EmployeeAttributes) (models.Employee, error) {
	existing, err := e.repo.GetAll()
	if err != nil {
		return models.Employee{}, err
	}
	for _, emp := range existing {
		if emp.CardNumberID == employee.CardNumberID && emp.Id != id {
			return models.Employee{}, httperrors.ConflictError{Message: "duplicated card number"}
		}
	}
	_, err = e.warehouseRepo.GetByID(employee.WarehouseID)
	if err != nil {
		var notFoundError httperrors.NotFoundError
		if errors.As(err, &notFoundError) {
			return models.Employee{}, httperrors.UnprocessableEntityError{Message: "warehouse_id does not exist"}
		}
		return models.Employee{}, err
	}

	modifiedEmployee := models.Employee{EmployeeAttributes: employee}
	return e.repo.Update(id, modifiedEmployee)
}

func (e EmployeeServiceDefault) Delete(id int) error {
	return e.repo.Delete(id)
}
