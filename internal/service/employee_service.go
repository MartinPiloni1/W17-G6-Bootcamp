package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mapper"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type EmployeeServiceImpl struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &EmployeeServiceImpl{repo: repo}
}

func (e EmployeeServiceImpl) Create(Employee models.Employee) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
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

func (e EmployeeServiceImpl) Update(id int, employee models.Employee) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmployeeServiceImpl) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
