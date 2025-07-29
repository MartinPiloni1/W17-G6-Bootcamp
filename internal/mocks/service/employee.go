package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeService is a mock of EmployeeService
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) Create(employee models.EmployeeAttributes) (models.Employee, error) {
	args := m.Called(employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) GetAll() ([]models.Employee, error) {
	args := m.Called()
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *MockEmployeeService) GetByID(id int) (models.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) Update(id int, employee models.EmployeeAttributes) (models.Employee, error) {
	args := m.Called(id, employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeService) ReportInboundOrders(employeeID int) ([]models.EmployeeWithInboundCount, error) {
	args := m.Called(employeeID)
	return args.Get(0).([]models.EmployeeWithInboundCount), args.Error(1)
}
