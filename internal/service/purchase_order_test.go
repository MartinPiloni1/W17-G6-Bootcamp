package service_test

import (
	"context"
	"testing"
	"time"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseOrderDefault_Create(t *testing.T) {
	// arrange
	repoMock := mocks.NewPurchaseOrderRepositoryDBMock()
	serviceDefault := service.NewPurchaseOrderDefault(repoMock)

	purchaseOrderAtt := models.PurchaseOrderAttributes{
		OrderNumber:     "ORD-1",
		OrderDate:       time.Date(2025, 4, 4, 15, 4, 5, 0, time.UTC),
		TrackingCode:    "abc123asd",
		BuyerId:         1,
		ProductRecordId: 1,
	}

	expectedPurchaseOrder := models.PurchaseOrder{
		Id:                      1,
		PurchaseOrderAttributes: purchaseOrderAtt,
	}

	ctx := context.Background()
	repoMock.On("Create", ctx, purchaseOrderAtt).
		Return(expectedPurchaseOrder, nil).
		Once()

	// act
	got, err := serviceDefault.Create(ctx, purchaseOrderAtt)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, got, expectedPurchaseOrder)
	repoMock.AssertNumberOfCalls(t, "Create", 1)
}
