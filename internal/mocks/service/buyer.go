package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// BuyerServiceDefaultMock is a mock of BuyerServiceDefault
type BuyerServiceDefaultMock struct {
	mock.Mock
}

// NewBuyerServiceDefaultMock retuns an instance of BuyerServiceDefaultMock
func NewBuyerServiceDefaultMock() *BuyerServiceDefaultMock {
	return new(BuyerServiceDefaultMock)
}

func (m *BuyerServiceDefaultMock) Create(ctx context.Context, newBuyer models.BuyerAttributes) (models.Buyer, error) {
	args := m.Called(ctx, newBuyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *BuyerServiceDefaultMock) GetAll(ctx context.Context) ([]models.Buyer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Buyer), args.Error(1)
}

func (m *BuyerServiceDefaultMock) GetByID(ctx context.Context, id int) (models.Buyer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *BuyerServiceDefaultMock) Update(ctx context.Context, id int, BuyerData models.BuyerPatchRequest) (models.Buyer, error) {
	args := m.Called(ctx, id, BuyerData)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *BuyerServiceDefaultMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *BuyerServiceDefaultMock) GetWithPurchaseOrdersCount(ctx context.Context, id *int) ([]models.BuyerWithPurchaseOrdersCount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.BuyerWithPurchaseOrdersCount), args.Error(1)
}
