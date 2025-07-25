package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockInboundOrderRepository struct{ mock.Mock }

func (m *MockInboundOrderRepository) Create(order models.InboundOrder) (models.InboundOrder, error) {
	args := m.Called(order)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}

func (m *MockInboundOrderRepository) GetByOrderNumber(orderNumber string) (models.InboundOrder, error) {
	args := m.Called(orderNumber)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}

func (m *MockInboundOrderRepository) CountInboundOrdersForEmployees() (map[int]int, error) {
	args := m.Called()
	return args.Get(0).(map[int]int), args.Error(1)
}

func (m *MockInboundOrderRepository) CountInboundOrdersForEmployee(employeeID int) (int, error) {
	args := m.Called(employeeID)
	return args.Int(0), args.Error(1)
}
