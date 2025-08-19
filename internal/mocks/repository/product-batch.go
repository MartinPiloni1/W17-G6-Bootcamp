package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// ProductBatchRepositoryMock mocks the ProductBatchRepository interface
type ProductBatchRepositoryMock struct {
	mock.Mock
}

func (m *ProductBatchRepositoryMock) Create(ctx context.Context, batch models.ProductBatchAttibutes) (models.ProductBatch, error) {
	args := m.Called(ctx, batch)
	return args.Get(0).(models.ProductBatch), args.Error(1)
}