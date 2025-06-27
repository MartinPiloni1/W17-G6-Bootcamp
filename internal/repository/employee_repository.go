package repository

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"os"
)

type EmployeeRepositoryImpl struct {
	filePath string
}

func NewEmployeeRepository() EmployeeRepository {
	return &EmployeeRepositoryImpl{filePath: os.Getenv("FILE_PATH_DEFAULT")}
}

func (e EmployeeRepositoryImpl) Create(Employee models.Employee) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmployeeRepositoryImpl) GetAll() (map[int]models.Employee, error) {
	data, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (e EmployeeRepositoryImpl) GetByID(id int) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmployeeRepositoryImpl) Update(id int, data models.Employee) (models.Employee, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmployeeRepositoryImpl) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
