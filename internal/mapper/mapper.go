package mapper

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type EmployeeMapper interface {
	MapToSlice(employees map[int]models.Employee) []models.Employee
}
