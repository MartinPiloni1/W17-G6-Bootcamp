package mocks

import (
	"context"

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

// Si tu service tiene más métodos agrégalos aquí. Por ejemplo (puedes eliminar si no los usás):

func (m *InboundOrderServiceDefaultMock) GetAll(ctx context.Context) ([]models.InboundOrder, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.InboundOrder), args.Error(1)
}

func (m *InboundOrderServiceDefaultMock) GetByID(ctx context.Context, id int) (models.InboundOrder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}

func (m *InboundOrderServiceDefaultMock) Update(ctx context.Context, id int, attrs models.InboundOrderAttributes) (models.InboundOrder, error) {
	args := m.Called(ctx, id, attrs)
	return args.Get(0).(models.InboundOrder), args.Error(1)
}

func (m *InboundOrderServiceDefaultMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
