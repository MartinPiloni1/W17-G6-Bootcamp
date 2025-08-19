package service_test

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Verifies the behavior of the service layer responsible for creating a new ProductRecord. It covers:
// - Successful creation
// - Error propagation from the repository layer
func TestProductRecordService_Create(t *testing.T) {
	// Define the product record used in common by the test cases
	newProductRecordAttributes := models.ProductRecordAttributes{
		LastUpdateDate: "2024-06-01",
		PurchasePrice:  150,
		SalePrice:      200,
		ProductID:      1,
	}

	newProductRecord := models.ProductRecord{
		ID:                      1,
		ProductRecordAttributes: newProductRecordAttributes,
	}

	// Each test case is constructed by:
	// testName            — human‐readable description
	// productAttributes   — ProductRecord attributes of the new product record
	// repositoryData      — the Product record object returned by the mocked repository
	// repositoryError     — the error returned by the mocked repository
	// expectedData        — the data we expect the service to produce
	// expectedError       — the error we expect the service to produce
	tests := []struct {
		testName                string
		productRecordAttributes models.ProductRecordAttributes
		repositoryData          models.ProductRecord
		repositoryError         error
		expectedData            models.ProductRecord
		expectedError           error
	}{
		{
			testName:                "Success case: Should create a product",
			productRecordAttributes: newProductRecordAttributes,
			repositoryData:          newProductRecord,
			repositoryError:         nil,
			expectedData:            newProductRecord,
			expectedError:           nil,
		},
		{
			testName:                "Error case: Process an error from the repository layer",
			productRecordAttributes: newProductRecordAttributes,
			repositoryData:          models.ProductRecord{},
			repositoryError:         errors.New("A repository error"),
			expectedData:            models.ProductRecord{},
			expectedError:           errors.New("A repository error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			repositoryMock := mocks.ProductRecordRepositoryDBMock{}
			service := service.NewProductRecordServiceDefault(&repositoryMock)

			repositoryMock.
				On("Create", mock.Anything, tc.productRecordAttributes).
				Return(tc.repositoryData, tc.repositoryError)

			// Act
			result, err := service.Create(context.Background(), newProductRecordAttributes)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedData, result)
			repositoryMock.AssertExpectations(t)
		})
	}
}
