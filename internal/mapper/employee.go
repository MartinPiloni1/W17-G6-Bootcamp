package mapper

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type EmployeeMapperImpl struct{}

func NewEmployeeMapper() EmployeeMapper {
	return &EmployeeMapperImpl{}
}

func (e EmployeeMapperImpl) MapToSlice(employees map[int]models.Employee) []models.Employee {
	var res []models.Employee
	for _, employee := range employees {
		res = append(res, employee)
	}
	return res
}
