package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// BuyerRepositoryDBMock is a mock of BuyerRepositoryDB
type BuyerRepositoryDBMock struct {
	mock.Mock
}

// NewBuyerRepositoryDBMock creates an instance of BuyerRepositoryDBMock
func NewBuyerRepositoryDBMock() *BuyerRepositoryDBMock {
	return new(BuyerRepositoryDBMock)
}

func (m *BuyerRepositoryDBMock) Create(ctx context.Context, newBuyer models.BuyerAttributes) (models.Buyer, error) {
	args := m.Called(ctx, newBuyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *BuyerRepositoryDBMock) GetAll(ctx context.Context) ([]models.Buyer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Buyer), args.Error(1)
}

func (m *BuyerRepositoryDBMock) GetByID(ctx context.Context, id int) (models.Buyer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *BuyerRepositoryDBMock) Update(ctx context.Context, id int, updatedBuyer models.Buyer) (models.Buyer, error) {
	args := m.Called(ctx, id, updatedBuyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *BuyerRepositoryDBMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *BuyerRepositoryDBMock) GetWithPurchaseOrdersCount(ctx context.Context, id *int) ([]models.BuyerWithPurchaseOrdersCount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.BuyerWithPurchaseOrdersCount), args.Error(1)
}
