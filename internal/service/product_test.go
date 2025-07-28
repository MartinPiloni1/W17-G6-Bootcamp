package service_test

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Verifies the behavior of the HTTP handler responsible for creating a new Product. It covers:
// - Successful creation
// - Error propagation from the repository layer
func TestProductService_Create(t *testing.T) {
	// Define the products used in common by the test cases
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
		SellerID:                       utils.Ptr(1),
	}

	createdProduct := models.Product{
		ID:                1,
		ProductAttributes: newProduct,
	}

	// Each test case is constructed by:
	// testName            — human‐readable description
	// repositoryData      — the Product object returned by the mocked repository
	// repositoryError     — the error returned by the mocked repository
	// expectedData        — the data we expect the service to produce
	// expectedError       — the error we expect the service to produce
	tests := []struct {
		testName        string
		repositoryData  models.Product
		repositoryError error
		expectedData    models.Product
		expectedError   error
	}{
		{
			testName:        "Success case: Should create a product",
			repositoryData:  createdProduct,
			repositoryError: nil,
			expectedData:    createdProduct,
			expectedError:   nil,
		},
		{
			testName:        "Error case: Process an error from the repository layer",
			repositoryData:  models.Product{},
			repositoryError: errors.New("db error"),
			expectedData:    models.Product{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("Create", mock.Anything, newProduct).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.Create(context.Background(), newProduct)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedData, result)
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
				SellerID:                       utils.Ptr(1),
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
				SellerID:                       utils.Ptr(1),
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
			testName:        "Success case: should return all products",
			repositoryData:  products,
			repositoryError: nil,
			expectedResp:    products,
			expectedError:   nil,
		},
		{
			testName:        "Error case: should return an error",
			repositoryData:  []models.Product{},
			repositoryError: errors.New("db error"),
			expectedResp:    []models.Product{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetAll", mock.Anything).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetAll(context.Background())

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
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
			SellerID:                       utils.Ptr(1),
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
			testName:        "Error case: should return an Error",
			repositoryData:  models.Product{},
			repositoryError: errors.New("db error"),
			expectedResp:    models.Product{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetByID", mock.Anything, tc.idParam).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetByID(context.Background(), tc.idParam)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestProductService_GetRecordsPerProduct(t *testing.T) {
	// Arrange
	recordsPerProduct := []models.ProductRecordCount{
		{
			ProductID:    1,
			Description:  "Yogurt helado",
			RecordsCount: 3,
		},
		{
			ProductID:    2,
			Description:  "Pechuga de pollo",
			RecordsCount: 1,
		},
	}

	singleRecord := []models.ProductRecordCount{
		{
			ProductID:    2,
			Description:  "Pechuga de pollo",
			RecordsCount: 1,
		},
	}

	tests := []struct {
		testName        string
		repositoryData  []models.ProductRecordCount
		repositoryError error
		idParam         *int
		expectedResp    []models.ProductRecordCount
		expectedError   error
	}{
		{
			testName:        "Success: Should return product record count of every product",
			repositoryData:  recordsPerProduct,
			repositoryError: nil,
			idParam:         nil,
			expectedResp:    recordsPerProduct,
			expectedError:   nil,
		},
		{
			testName:        "Success: Should return product record count of a single product",
			repositoryData:  singleRecord,
			repositoryError: nil,
			idParam:         utils.Ptr(2),
			expectedResp:    singleRecord,
			expectedError:   nil,
		},
		{
			testName:        "Error case: should return an Error",
			repositoryData:  []models.ProductRecordCount{},
			repositoryError: errors.New("db error"),
			idParam:         nil,
			expectedResp:    []models.ProductRecordCount{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetRecordsPerProduct", mock.Anything, tc.idParam).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetRecordsPerProduct(context.Background(), tc.idParam)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestProductService_Update(t *testing.T) {
	sellerID := utils.Ptr(1) // Should be declared in a variable to avoid different memory addreses issues

	originalProduct := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Yogurt helado",
			ExpirationRate:                 7,
			FreezingRate:                   2,
			Height:                         10.5,
			Length:                         20.0,
			Width:                          15.0,
			NetWeight:                      1.2,
			ProductCode:                    "YOG01",
			RecommendedFreezingTemperature: -5.0,
			ProductTypeID:                  3,
			SellerID:                       sellerID,
		},
	}

	updatePayload := models.ProductPatchRequest{
		Description:                    utils.Ptr(""),
		ExpirationRate:                 utils.Ptr(1),
		FreezingRate:                   utils.Ptr(2),
		Height:                         utils.Ptr(1.0),
		Length:                         utils.Ptr(1.0),
		Width:                          utils.Ptr(1.0),
		NetWeight:                      utils.Ptr(1.0),
		ProductCode:                    utils.Ptr(""),
		RecommendedFreezingTemperature: utils.Ptr(1.0),
		ProductTypeID:                  utils.Ptr(1),
		SellerID:                       sellerID,
	}

	updatedProduct := models.Product{
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
			SellerID:                       sellerID,
		},
	}

	singleFieldUpdatePayload := models.ProductPatchRequest{
		Description: utils.Ptr("Yogur helado"),
	}

	singleFieldUpdatedProduct := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Yogur helado",
			ExpirationRate:                 7,
			FreezingRate:                   2,
			Height:                         10.5,
			Length:                         20.0,
			Width:                          15.0,
			NetWeight:                      1.2,
			ProductCode:                    "YOG01",
			RecommendedFreezingTemperature: -5.0,
			ProductTypeID:                  3,
			SellerID:                       sellerID,
		},
	}

	tests := []struct {
		testName             string
		payload              models.ProductPatchRequest
		updatedProduct       models.Product
		repositoryGetData    models.Product
		repositoryGetError   error
		repositoryPatchData  models.Product
		repositoryPatchError error
		idParam              int
		expectedResp         models.Product
		expectedError        error
	}{
		{
			testName:             "Success: should update all fields of a product",
			payload:              updatePayload,
			updatedProduct:       updatedProduct,
			repositoryGetData:    originalProduct,
			repositoryGetError:   nil,
			repositoryPatchData:  updatedProduct,
			repositoryPatchError: nil,
			idParam:              1,
			expectedResp:         updatedProduct,
			expectedError:        nil,
		},
		{
			testName:             "Success: should update a single field of a product",
			payload:              singleFieldUpdatePayload,
			updatedProduct:       singleFieldUpdatedProduct,
			repositoryGetData:    originalProduct,
			repositoryGetError:   nil,
			repositoryPatchData:  singleFieldUpdatedProduct,
			repositoryPatchError: nil,
			idParam:              1,
			expectedResp:         singleFieldUpdatedProduct,
			expectedError:        nil,
		},
		{
			testName:           "Error case: GetByID fails",
			payload:            updatePayload,
			repositoryGetData:  models.Product{},
			repositoryGetError: httperrors.NotFoundError{Message: "Product not found"},
			idParam:            1,
			expectedResp:       models.Product{},
			expectedError:      httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName:             "Error case: should return an error",
			payload:              updatePayload,
			updatedProduct:       updatedProduct,
			repositoryGetData:    originalProduct,
			repositoryGetError:   nil,
			repositoryPatchData:  models.Product{},
			repositoryPatchError: errors.New("db error"),
			idParam:              1,
			expectedResp:         models.Product{},
			expectedError:        errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetByID", mock.Anything, tc.idParam).
				Return(tc.repositoryGetData, tc.repositoryGetError)

			if tc.repositoryGetError == nil {
				repositoryMock.
					On("Update", mock.Anything, tc.idParam, tc.updatedProduct).
					Return(tc.repositoryPatchData, tc.repositoryPatchError)
			}

			// Act
			result, err := service.Update(context.Background(), tc.idParam, tc.payload)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
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
			testName:        "Fail: should return an Error",
			idParam:         1,
			repositoryError: errors.New("db error"),
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("Delete", mock.Anything, tc.idParam).
				Return(tc.repositoryError)

			// Act
			err := service.Delete(context.Background(), tc.idParam)

			// Assert
			require.Equal(t, tc.expectedError, err)
			repositoryMock.AssertExpectations(t)
		})
	}
}
