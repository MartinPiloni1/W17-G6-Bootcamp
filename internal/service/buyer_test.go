package service_test

import (
	"context"
	"testing"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuyerServiceDefault_Create(t *testing.T) {
	t.Run("successfull creation", func(t *testing.T) {
		// arrange
		expectedAttrs := models.BuyerAttributes{
			CardNumberId: 34567890,
			FirstName:    "Juan",
			LastName:     "Lopez",
		}

		expectedBuyer := models.Buyer{
			Id:              1,
			BuyerAttributes: expectedAttrs,
		}

		ctx := context.TODO() // dummy context
		repoMock := mocks.NewBuyerRepositoryDBMock()
		repoMock.On("Create", ctx, expectedAttrs).
			Return(expectedBuyer, nil)

		srv := service.NewBuyerServiceDefault(repoMock)

		// act
		buyer, err := srv.Create(ctx, expectedAttrs)

		// assert
		require.NoError(t, err)
		assert.Equal(t, expectedBuyer, buyer)
		repoMock.AssertExpectations(t)
		repoMock.AssertNumberOfCalls(t, "Create", 1)
	})
}

func TestBuyerServiceDefault_Update(t *testing.T) {
	t.Run("successfully updates all fields", func(t *testing.T) {
		// Arrange
		buyerID := 1
		oldBuyer := models.Buyer{
			Id: buyerID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberId: 12345678,
				FirstName:    "Old First Name",
				LastName:     "Old Last Name",
			},
		}

		newCardNumberId := 12233445
		newFirstName := "New First Name"
		newLastName := "New Last Name"

		patch := models.BuyerPatchRequest{
			CardNumberId: &newCardNumberId,
			FirstName:    &newFirstName,
			LastName:     &newLastName,
		}

		updatedBuyer := models.Buyer{
			Id: buyerID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberId: newCardNumberId,
				FirstName:    newFirstName,
				LastName:     newLastName,
			},
		}

		repoMock := mocks.NewBuyerRepositoryDBMock()
		svc := service.NewBuyerServiceDefault(repoMock)

		ctx := context.Background()
		repoMock.On("GetByID", ctx, buyerID).
			Return(oldBuyer, nil).
			Once()
		repoMock.On("Update", ctx, buyerID, updatedBuyer).
			Return(updatedBuyer, nil).
			Once()

		// Act
		got, err := svc.Update(ctx, buyerID, patch)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedBuyer, got)
		repoMock.AssertExpectations(t)
		repoMock.AssertNumberOfCalls(t, "Update", 1)
		repoMock.AssertNumberOfCalls(t, "GetByID", 1)
	})

	t.Run("not found returns error", func(t *testing.T) {
		// arrange
		ctx := context.Background() // dummy context

		buyerID := 999
		expectedErr := httperrors.NotFoundError{Message: "Buyer Not Found"}

		repoMock := mocks.NewBuyerRepositoryDBMock()
		repoMock.On("GetByID", ctx, buyerID).
			Return(models.Buyer{}, expectedErr).
			Once()

		svc := service.NewBuyerServiceDefault(repoMock)

		// act
		got, err := svc.Update(ctx, buyerID, models.BuyerPatchRequest{})

		// assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, models.Buyer{}, got)
		repoMock.AssertNumberOfCalls(t, "GetByID", 1)
		repoMock.AssertNumberOfCalls(t, "Update", 0)
		repoMock.AssertExpectations(t)
	})
}

func TestBuyerServiceDefault_Delete(t *testing.T) {
	// assert
	repoMock := mocks.NewBuyerRepositoryDBMock()
	serviceDefault := service.NewBuyerServiceDefault(repoMock)

	ctx := context.Background()

	repoMock.On("Delete", ctx, 1).Return(nil).Once()

	// act
	actual := serviceDefault.Delete(ctx, 1)

	// assert
	assert.Nil(t, actual)
	repoMock.AssertNumberOfCalls(t, "Delete", 1)
}

func TestBuyerServiceDefault_GetWithPurchaseOrdersCount(t *testing.T) {

	t.Run("passing an id that does not exist returns err", func(t *testing.T) {

		repoMock := mocks.NewBuyerRepositoryDBMock()
		serviceDefault := service.NewBuyerServiceDefault(repoMock)

		ctx := context.Background()
		buyerID := 9999
		notFoundErr := httperrors.NotFoundError{Message: "Buyer Not Found"}

		repoMock.On("GetByID", ctx, 9999).Return(models.Buyer{}, notFoundErr).Once()

		got, err := serviceDefault.GetWithPurchaseOrdersCount(ctx, &buyerID)

		// assert
		assert.Nil(t, got)
		assert.Equal(t, notFoundErr, err)
		repoMock.AssertNumberOfCalls(t, "GetByID", 1)
	})

	t.Run("not passing buyerID retuns all buyers WithPurchaseOrdersCount, not calling GetByID", func(t *testing.T) {
		repoMock := mocks.NewBuyerRepositoryDBMock()
		serviceDefault := service.NewBuyerServiceDefault(repoMock)

		ctx := context.Background()
		var nilPtrID *int

		repoMock.On("GetWithPurchaseOrdersCount", ctx, nilPtrID).Return([]models.BuyerWithPurchaseOrdersCount{}, nil).Once()

		_, err := serviceDefault.GetWithPurchaseOrdersCount(ctx, nilPtrID)

		// assert
		assert.Nil(t, err)
		repoMock.AssertNumberOfCalls(t, "GetByID", 0)
		repoMock.AssertNumberOfCalls(t, "GetWithPurchaseOrdersCount", 1)
	})
}

func TestBuyerServiceDefault_GetByID(t *testing.T) {
	repoMock := mocks.NewBuyerRepositoryDBMock()
	serviceDefault := service.NewBuyerServiceDefault(repoMock)

	buyer := models.Buyer{
		Id: 1,
		BuyerAttributes: models.BuyerAttributes{
			CardNumberId: 12312312,
			FirstName:    "Juan",
			LastName:     "Perez",
		},
	}
	ctx := context.Background()
	repoMock.On("GetByID", ctx, buyer.Id).Return(buyer, nil)

	// act
	got, err := serviceDefault.GetByID(ctx, buyer.Id)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, buyer, got)
	repoMock.AssertNumberOfCalls(t, "GetByID", 1)
}
