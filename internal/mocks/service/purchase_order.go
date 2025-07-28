package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type PurchaseOrderDefaultMock struct {
	mock.Mock
}

func NewPurchaseOrderDefaultMock() *PurchaseOrderDefaultMock {
	return new(PurchaseOrderDefaultMock)
}

func (m *PurchaseOrderDefaultMock) Create(ctx context.Context, newPurchaseOrder models.PurchaseOrderAttributes) (models.PurchaseOrder, error) {
	args := m.Called(ctx, newPurchaseOrder)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}
