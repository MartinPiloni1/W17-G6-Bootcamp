package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type EmployeeService interface {
	Create(Employee models.EmployeeAttributes) (models.Employee, error)
	GetAll() ([]models.Employee, error)
	GetByID(id int) (models.Employee, error)
	Update(id int, employee models.EmployeeAttributes) (models.Employee, error)
	Delete(id int) error
}
