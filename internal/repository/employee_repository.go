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

	for _, empIteration := range dataList {
		if empIteration.CardNumberID == employee.CardNumberID {
			return models.Employee{}, httperrors.UnprocessableEntityError{Message: "duplicate card number"}
		}
	}

	employee.Id, err = utils.GetNextID[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, httperrors.ConflictError{Message: "error generating sequential id"}
	}
	dataList[employee.Id] = employee

	utils.Write[models.Employee](e.filePath, dataList)
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
	for _, emp := range data {
		if emp.Id == id {
			return emp, nil
		}
	}
	return models.Employee{}, httperrors.NotFoundError{Message: "employee not found"}
}

func (e *EmployeeRepositoryImpl) Update(id int, employee models.Employee) (models.Employee, error) {
	dataList, err := utils.Read[models.Employee](e.filePath)
	if err != nil {
		return models.Employee{}, err
	}

	if _, exist := dataList[id]; !exist {
		return models.Employee{}, httperrors.NotFoundError{Message: "employee not found"}
	}

	for _, empIteration := range dataList {
		if empIteration.CardNumberID == employee.CardNumberID && empIteration.Id != id {
			return models.Employee{}, httperrors.UnprocessableEntityError{Message: "duplicated card number"}
		}
	}

	dataList[id] = employee
	utils.Write[models.Employee](e.filePath, dataList)
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
	utils.Write[models.Employee](e.filePath, dataList)
	return nil
}
