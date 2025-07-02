package repository

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"os"
)

type EmployeeRepositoryFile struct {
	filePath string
}

func NewEmployeeRepository() EmployeeRepository {
	return &EmployeeRepositoryFile{filePath: os.Getenv("FILE_PATH_EMPLOYEE")}
}

func (e *EmployeeRepositoryFile) Create(employee models.Employee) (models.Employee, error) {
	dataList, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, err
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

func (e *EmployeeRepositoryFile) GetAll() (map[int]models.Employee, error) {
	data, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *EmployeeRepositoryFile) GetByID(id int) (models.Employee, error) {
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

func (e *EmployeeRepositoryFile) Update(id int, employee models.Employee) (models.Employee, error) {
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

func (e *EmployeeRepositoryFile) Delete(id int) error {
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
