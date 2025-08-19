package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// ProductBatchServiceDefaultMock is a mock of ProductBatchServiceDefault
type ProductBatchServiceDefaultMock struct {
	mock.Mock
}

// NewProductBatchServiceDefaultMock returns an instance of ProductBatchServiceDefaultMock
func NewProductBatchServiceDefaultMock() *ProductBatchServiceDefaultMock {
	return new(ProductBatchServiceDefaultMock)
}

func (m *ProductBatchServiceDefaultMock) Create(ctx context.Context, productBatch models.ProductBatchAttibutes) (models.ProductBatch, error) {
	args := m.Called(ctx, productBatch)
	return args.Get(0).(models.ProductBatch), args.Error(1)
}
