package repository

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
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

func (e *EmployeeRepositoryImpl) Create(employee models.Employee) (models.Employee, error) {
	dataList, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, err
	}

	if _, exist := dataList[employee.Id]; exist {
		return models.Employee{}, httperrors.ConflictError{Message: "already exist"}
	}

	employee.Id, err = utils.GetNextID[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, httperrors.ConflictError{Message: "error generating sequential id"}
	}
	dataList[employee.Id] = employee

	if err := utils.Write[models.Employee](e.filePath, dataList); err != nil {
		return models.Employee{}, err
	}

	return employee, nil
}

func (e *EmployeeRepositoryImpl) GetAll() (map[int]models.Employee, error) {
	data, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *EmployeeRepositoryImpl) GetByID(id int) (models.Employee, error) {
	data, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, err
	}
	emp, exist := data[id]
	if !exist {
		return models.Employee{}, httperrors.NotFoundError{Message: "employee not found"}
	}
	return emp, nil
}

func (e *EmployeeRepositoryImpl) Update(id int, employee models.Employee) (models.Employee, error) {
	dataList, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, err
	}

	if _, exist := dataList[id]; !exist {
		return models.Employee{}, httperrors.NotFoundError{Message: "employee not found"}
	}

	employee.Id = id
	dataList[id] = employee
	if err := utils.Write[models.Employee](e.filePath, dataList); err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}

func (e *EmployeeRepositoryImpl) Delete(id int) error {
	dataList, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return err
	}

	if _, exist := dataList[id]; !exist {
		return httperrors.NotFoundError{Message: "employee not found"}
	}

	delete(dataList, id)
	if err := utils.Write[models.Employee](e.filePath, dataList); err != nil {
		return err
	}
	return nil
}
