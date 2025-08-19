package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// InboundOrderServiceDefaultMock is a mock of InboundOrderService
type InboundOrderServiceDefaultMock struct {
	mock.Mock
}

// NewInboundOrderServiceDefaultMock returns an instance of InboundOrderServiceDefaultMock
func NewInboundOrderServiceDefaultMock() *InboundOrderServiceDefaultMock {
	return new(InboundOrderServiceDefaultMock)
}

func (m *InboundOrderServiceDefaultMock) Create(attrs models.InboundOrderAttributes) (models.InboundOrder, error) {
	args := m.Called(attrs)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}
