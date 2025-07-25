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

func TestProductService_Create(t *testing.T) {
	// Arrange
	newProduct := models.ProductAttributes{
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
	}

	createdProduct := models.Product{
		ID:                1,
		ProductAttributes: newProduct,
	}

	tests := []struct {
		testName        string
		repositoryData  models.Product
		repositoryError error
		expectedResp    models.Product
		expectedError   error
	}{
		{
			testName:        "Success: Should create a product",
			repositoryData:  createdProduct,
			repositoryError: nil,
			expectedResp:    createdProduct,
			expectedError:   nil,
		},
		{
			testName:        "Fail: should return a Conflict Error if a product with the given product_code already exists ",
			repositoryData:  models.Product{},
			repositoryError: httperrors.ConflictError{Message: "A product with the given product code already exists"},
			expectedResp:    models.Product{},
			expectedError:   httperrors.ConflictError{Message: "A product with the given product code already exists"},
		},
		{
			testName:        "Fail: should return a Conflict Error if the given seller id doesn't exists",
			repositoryData:  models.Product{},
			repositoryError: httperrors.ConflictError{Message: "The given seller id does not exists"},
			expectedResp:    models.Product{},
			expectedError:   httperrors.ConflictError{Message: "The given seller id does not exists"},
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
				On("Create", mock.Anything, newProduct).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.Create(context.Background(), newProduct)

			// Assert
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResp, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

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

func TestProductService_Delete(t *testing.T) {
	// Arrange
	tests := []struct {
		testName        string
		idParam         int
		repositoryError error
		expectedError   error
	}{
		{
			testName:        "Success: should delete the product with the given ID",
			idParam:         1,
			repositoryError: nil,
			expectedError:   nil,
		},
		{
			testName:        "Fail: should return a not found error",
			idParam:         1,
			repositoryError: httperrors.NotFoundError{Message: "Product not found"},
			expectedError:   httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName:        "Fail: should return a not found error",
			idParam:         1,
			repositoryError: httperrors.NotFoundError{Message: "Product not found"},
			expectedError:   httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName:        "Fail: should return a conflict error",
			idParam:         1,
			repositoryError: httperrors.ConflictError{Message: "The product to delete is still referenced by some product records"},
			expectedError:   httperrors.ConflictError{Message: "The product to delete is still referenced by some product records"},
		},
		{
			testName:        "Fail: should return a DB Error",
			idParam:         1,
			repositoryError: errors.New("db error"),
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
				On("Delete", mock.Anything, tc.idParam).
				Return(tc.repositoryError)

			// Act
			err := service.Delete(context.Background(), tc.idParam)

			// Assert
			assert.Equal(t, tc.expectedError, err)
			repositoryMock.AssertExpectations(t)
		})
	}
}
