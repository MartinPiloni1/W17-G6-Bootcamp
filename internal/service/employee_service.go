package service

import (
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

func (e EmployeeServiceImpl) GetAll() (map[int]models.Employee, error) {
	return e.repo.GetAll()
}

func (e EmployeeServiceImpl) GetByID(id int) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmployeeServiceImpl) Update(id int, employee models.Employee) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmployeeServiceImpl) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
