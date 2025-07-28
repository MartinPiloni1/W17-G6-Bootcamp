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

// Verifies the behavior of the service layer responsible for creating a new Product. It covers:
// - Successful creation
// - Error propagation from the repository layer
func TestProductService_Create(t *testing.T) {
	// Define the product used in common by the test cases
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

// Verifies the behavior of the service layer responsible for retrieving all products.
// It covers:
// - Successful retrieval of multiple products
// - Error propagation from the repository layer
func TestProductService_GetAll(t *testing.T) {
	// Define the products used in common by the test cases
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

	// Each test case is constructed by:
	// testName            — human‐readable description
	// repositoryData      — the Products slice returned by the mocked repository
	// repositoryError     — the error returned by the mocked repository
	// expectedData        — the data we expect the service to produce
	// expectedError       — the error we expect the service to produce
	tests := []struct {
		testName        string
		repositoryData  []models.Product
		repositoryError error
		expectedData    []models.Product
		expectedError   error
	}{
		{
			testName:        "Success case: Should return all products",
			repositoryData:  products,
			repositoryError: nil,
			expectedData:    products,
			expectedError:   nil,
		},
		{
			testName:        "Error case: Process an error from the repository layer",
			repositoryData:  []models.Product{},
			repositoryError: errors.New("db error"),
			expectedData:    []models.Product{},
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
			require.Equal(t, tc.expectedData, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the service layer responsible for retrieving a single product.
// It covers:
// - Successful retrieval of a single product
// - Error propagation from the repository layer
func TestProductService_GetByID(t *testing.T) {
	// Define the product used in common by the test cases
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

	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to retrieve
	// repositoryData      — the Product returned by the mocked repository
	// repositoryError     — the error returned by the mocked repository
	// expectedData        — the data we expect the service to produce
	// expectedError       — the error we expect the service to produce
	tests := []struct {
		testName        string
		id              int
		repositoryData  models.Product
		repositoryError error
		expectedData    models.Product
		expectedError   error
	}{
		{
			testName:        "Success: Should return a single product",
			id:              1,
			repositoryData:  product,
			repositoryError: nil,
			expectedData:    product,
			expectedError:   nil,
		},
		{
			testName:        "Error case: Process an error from the repository layer",
			id:              1,
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
				On("GetByID", mock.Anything, tc.id).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetByID(context.Background(), tc.id)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedData, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

// Verifies the service layer responsible for retrieving records per product.
// It covers:
// - Successful retrieval of the record count of every product
// - Successful retrieval of the record count of a single product
// - Error propagation from the repository layer
func TestProductService_GetRecordsPerProduct(t *testing.T) {
	// Define the records used in common by the test cases
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

	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to retrieve as a pointer
	// repositoryData      — the Product returned by the mocked repository
	// repositoryError     — the error returned by the mocked repository
	// expectedData        — the data we expect the service to produce
	// expectedError       — the error we expect the service to produce
	tests := []struct {
		testName        string
		id              *int
		repositoryData  []models.ProductRecordCount
		repositoryError error
		expectedData    []models.ProductRecordCount
		expectedError   error
	}{
		{
			testName:        "Success: Should return product record count of every product",
			id:              nil,
			repositoryData:  recordsPerProduct,
			repositoryError: nil,
			expectedData:    recordsPerProduct,
			expectedError:   nil,
		},
		{
			testName:        "Success: Should return product record count of a single product",
			id:              utils.Ptr(2),
			repositoryData:  singleRecord,
			repositoryError: nil,
			expectedData:    singleRecord,
			expectedError:   nil,
		},
		{
			testName:        "Error case: Process an error from the repository layer",
			id:              nil,
			repositoryData:  []models.ProductRecordCount{},
			repositoryError: errors.New("db error"),
			expectedData:    []models.ProductRecordCount{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetRecordsPerProduct", mock.Anything, tc.id).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.GetRecordsPerProduct(context.Background(), tc.id)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedData, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of service layer responsible for updating a product.
// It covers:
// - Successful update all fields of a product
// - Successful update a single field of a product
// - Error when GetByID call fails
// - Error propagation from the repository layer
func TestProductService_Update(t *testing.T) {
	// Define the products used in common by the test cases

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

	// Each test case is constructed by:
	// testName              — human‐readable description
	// payload               - The ProductPatchRequest to update a product
	// id                    - ID of the product to update
	// updatedProduct        - The Product with the updated fields
	// repositoryGetData     — the Product returned by the get call on the mocked repository
	// repositoryGetError    — the error returned by get call on the mocked repository
	// repositoryPatchData   — the Product returned by the patch call on the mocked repository
	// repositoryPatchError  — the error returned by the patch call mocked repository
	// expectedData          — the data we expect the service to produce
	// expectedError         — the error we expect the service to produce
	tests := []struct {
		testName             string
		payload              models.ProductPatchRequest
		id                   int
		updatedProduct       models.Product
		repositoryGetData    models.Product
		repositoryGetError   error
		repositoryPatchData  models.Product
		repositoryPatchError error
		expectedData         models.Product
		expectedError        error
	}{
		{
			testName:             "Success: should update all fields of a product",
			payload:              updatePayload,
			id:                   1,
			updatedProduct:       updatedProduct,
			repositoryGetData:    originalProduct,
			repositoryGetError:   nil,
			repositoryPatchData:  updatedProduct,
			repositoryPatchError: nil,
			expectedData:         updatedProduct,
			expectedError:        nil,
		},
		{
			testName:             "Success: should update a single field of a product",
			payload:              singleFieldUpdatePayload,
			updatedProduct:       singleFieldUpdatedProduct,
			id:                   1,
			repositoryGetData:    originalProduct,
			repositoryGetError:   nil,
			repositoryPatchData:  singleFieldUpdatedProduct,
			repositoryPatchError: nil,
			expectedData:         singleFieldUpdatedProduct,
			expectedError:        nil,
		},
		{
			testName:           "Error case: GetByID fails",
			payload:            updatePayload,
			id:                 1,
			updatedProduct:     models.Product{},
			repositoryGetData:  models.Product{},
			repositoryGetError: httperrors.NotFoundError{Message: "Product not found"},
			expectedData:       models.Product{},
			expectedError:      httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName:             "Error case: Process an error from the repository layer",
			payload:              updatePayload,
			updatedProduct:       updatedProduct,
			id:                   1,
			repositoryGetData:    originalProduct,
			repositoryGetError:   nil,
			repositoryPatchData:  models.Product{},
			repositoryPatchError: errors.New("db error"),
			expectedData:         models.Product{},
			expectedError:        errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRepositoryDBMock{}
			service := service.NewProductServiceDefault(&repositoryMock)

			repositoryMock.
				On("GetByID", mock.Anything, tc.id).
				Return(tc.repositoryGetData, tc.repositoryGetError)

			// Check if the GetByID call failed
			if tc.repositoryGetError == nil {
				repositoryMock.
					On("Update", mock.Anything, tc.id, tc.updatedProduct).
					Return(tc.repositoryPatchData, tc.repositoryPatchError)
			}

			// Act
			result, err := service.Update(context.Background(), tc.id, tc.payload)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedData, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the service layer responsible for deleting a product.
// It covers:
// - Successful deletion of a product
// - Error propagation from the repository layer
func TestProductService_Delete(t *testing.T) {
	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to delete
	// repositoryError     — the error returned by the mocked repository
	// expectedError       — the error we expect the service to produce
	tests := []struct {
		testName        string
		id              int
		repositoryError error
		expectedError   error
	}{
		{
			testName:        "Success: should delete the product with the given ID",
			id:              1,
			repositoryError: nil,
			expectedError:   nil,
		},
		{
			testName:        "Error case: Process an error from the repository layer",
			id:              1,
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
				On("Delete", mock.Anything, tc.id).
				Return(tc.repositoryError)

			// Act
			err := service.Delete(context.Background(), tc.id)

			// Assert
			require.Equal(t, tc.expectedError, err)
			repositoryMock.AssertExpectations(t)
		})
	}
}
