package service_test

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func Ptr[T any](v T) *T { return &v }

func TestProductService_GetAll(t *testing.T) {
	// Arrange
	products := []models.Product{
		{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				Description:                    "",
				ExpirationRate:                 1,
				FreezingRate:                   2,
				Height:                         1,
				Length:                         1,
				Width:                          1,
				NetWeight:                      1,
				ProductCode:                    "",
				RecommendedFreezingTemperature: 1,
				ProductTypeID:                  1,
				SellerID:                       Ptr(1),
			},
		},
		{
			ID: 2,
			ProductAttributes: models.ProductAttributes{
				Description:                    "",
				ExpirationRate:                 1,
				FreezingRate:                   2,
				Height:                         1,
				Length:                         1,
				Width:                          1,
				NetWeight:                      1,
				ProductCode:                    "",
				RecommendedFreezingTemperature: 1,
				ProductTypeID:                  1,
				SellerID:                       Ptr(1),
			},
		},
	}

	tests := []struct {
		testName        string
		repositoryData  []models.Product
		repositoryError error
		expectedResp    []models.Product
		expectedError   error
	}{
		{
			testName:        "Success: Should return all products",
			repositoryData:  products,
			repositoryError: nil,
			expectedResp:    products,
			expectedError:   nil,
		},
		{
			testName:        "Fail: should return a DB Error",
			repositoryData:  []models.Product{},
			repositoryError: errors.New("db error"),
			expectedResp:    []models.Product{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetAll", mock.Anything).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetAll(context.Background())

			// Assert
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestProductService_GetByID(t *testing.T) {
	// Arrange
	product := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "",
			ExpirationRate:                 1,
			FreezingRate:                   2,
			Height:                         1,
			Length:                         1,
			Width:                          1,
			NetWeight:                      1,
			ProductCode:                    "",
			RecommendedFreezingTemperature: 1,
			ProductTypeID:                  1,
			SellerID:                       Ptr(1),
		},
	}

	tests := []struct {
		testName        string
		repositoryData  models.Product
		repositoryError error
		idParam         int
		expectedResp    models.Product
		expectedError   error
	}{
		{
			testName:        "Success: Should return all products",
			repositoryData:  product,
			repositoryError: nil,
			idParam:         1,
			expectedResp:    product,
			expectedError:   nil,
		},
		{
			testName:        "Fail: should return a not found error",
			repositoryData:  models.Product{},
			repositoryError: httperrors.NotFoundError{Message: "Product not found"},
			expectedResp:    models.Product{},
			expectedError:   httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName:        "Fail: should return a DB Error",
			repositoryData:  models.Product{},
			repositoryError: errors.New("db error"),
			expectedResp:    models.Product{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetByID", mock.Anything, tc.idParam).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetByID(context.Background(), tc.idParam)

			// Assert
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}
