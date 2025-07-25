package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockEmployeeRepository struct{ mock.Mock }

func (m *MockEmployeeRepository) Create(e models.Employee) (models.Employee, error) {
	a := m.Called(e)
	return a.Get(0).(models.Employee), a.Error(1)
}
func (m *MockEmployeeRepository) GetAll() ([]models.Employee, error) {
	args := m.Called()
	v := args.Get(0)
	if v == nil {
		return nil, args.Error(1)
	}
	return v.([]models.Employee), args.Error(1)
}
func (m *MockEmployeeRepository) GetByID(id int) (models.Employee, error) {
	a := m.Called(id)
	return a.Get(0).(models.Employee), a.Error(1)
}
func (m *MockEmployeeRepository) Update(id int, e models.Employee) (models.Employee, error) {
	a := m.Called(id, e)
	return a.Get(0).(models.Employee), a.Error(1)
}
func (m *MockEmployeeRepository) Delete(id int) error {
	a := m.Called(id)
	return a.Error(0)
}
