package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) Create(ctx context.Context, newProduct models.ProductAttributes) (models.Product, error) {
	args := m.Called(ctx, newProduct)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *ProductServiceMock) GetAll(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductServiceMock) GetByID(ctx context.Context, id int) (models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *ProductServiceMock) GetRecordsPerProduct(ctx context.Context, id *int) ([]models.ProductRecordCount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.ProductRecordCount), args.Error(1)
}

func (m *ProductServiceMock) Update(ctx context.Context, id int, updatedProduct models.ProductPatchRequest) (models.Product, error) {
	args := m.Called(ctx, id, updatedProduct)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *ProductServiceMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
