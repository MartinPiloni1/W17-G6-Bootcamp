package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mapper"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type EmployeeServiceImpl struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &EmployeeServiceImpl{repo: repo}
}

func (e EmployeeServiceImpl) Create(employee models.EmployeeAttributes) (models.Employee, error) {
	existing, err := e.repo.GetAll()
	if err != nil {
		return models.Employee{}, err
	}
	for _, emp := range existing {
		if emp.CardNumberID == employee.CardNumberID {
			return models.Employee{}, httperrors.UnprocessableEntityError{Message: "duplicate card number"}
		}
	}

	newEmployee := models.Employee{EmployeeAttributes: employee}
	return e.repo.Create(newEmployee)
}

func (e EmployeeServiceImpl) GetAll() ([]models.Employee, error) {
	employees, err := e.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return mapper.NewEmployeeMapper().MapToSlice(employees), nil
}

func (e EmployeeServiceImpl) GetByID(id int) (models.Employee, error) {
	return e.repo.GetByID(id)
}

func (e EmployeeServiceImpl) Update(id int, employee models.EmployeeAttributes) (models.Employee, error) {
	existing, err := e.repo.GetAll()
	if err != nil {
		return models.Employee{}, err
	}
	for _, emp := range existing {
		if emp.CardNumberID == employee.CardNumberID && emp.Id != id {
			return models.Employee{}, httperrors.UnprocessableEntityError{Message: "duplicated card number"}
		}
	}
	modifiedEmployee := models.Employee{EmployeeAttributes: employee}
	return e.repo.Update(id, modifiedEmployee)
}

func (e EmployeeServiceImpl) Delete(id int) error {
	return e.repo.Delete(id)
}
