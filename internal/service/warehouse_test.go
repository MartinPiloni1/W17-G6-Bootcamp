// Package service_test provides comprehensive testing for the warehouse business logic service.
//
// This test suite covers all CRUD operations for warehouse business logic,
// ensuring proper data validation, error handling, and repository interactions.
// Tests use mocked repositories to isolate service logic from data access layer.
//
// Test Categories:
//   - Create: Tests warehouse creation with validation and conflict handling
//   - GetAll: Tests retrieval of all warehouses from repository
//   - GetByID: Tests retrieval of specific warehouses by ID
//   - Update: Tests warehouse modification with validation and conflict detection
//   - Delete: Tests warehouse deletion operations
//
// Each test follows the AAA pattern (Arrange, Act, Assert) and uses testify
// for assertions and mocking. Business logic validation is thoroughly tested
// to ensure data integrity and proper error responses.
package service_test

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/require"
)

// emptyWarehouse represents an empty warehouse model used for error case assertions
var emptyWarehouse = models.Warehouse{}

// Test_Create tests the Create method of the WarehouseService.
//
// This test suite verifies that the warehouse creation business logic:
//   - Successfully creates warehouses with valid input data
//   - Validates all required fields before repository interaction
//   - Properly handles repository-level conflicts (duplicate warehouse codes)
//   - Returns appropriate business logic errors for invalid data
//   - Maintains data integrity through validation rules
//
// Test scenarios:
//   - Success case: Valid warehouse attributes return created warehouse
//   - Repository errors: Conflict and internal server errors are properly propagated
//   - Validation errors: Missing or invalid fields return validation errors
func Test_Create(t *testing.T) {
	// Test data: valid warehouse attributes for successful creation
	warehouseAtt := models.WarehouseAttributes{
		WarehouseCode:      "DHK",
		Address:            "123 Main St",
		Telephone:          "123-456-7890",
		MinimunCapacity:    100,
		MinimunTemperature: 15.0,
	}
	warehouse := models.Warehouse{
		Id:                  1,
		WarehouseAttributes: warehouseAtt,
	}

	// Success scenario: create warehouse with valid attributes
	t.Run("create warehouse success (create_ok)", func(t *testing.T) {
		// Arrange: mock repository to return successful creation
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("Create", warehouseAtt).Return(warehouse, nil)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: execute warehouse creation through service
		result, err := serviceTest.Create(warehouseAtt)

		// Assert: verify successful creation with correct warehouse data
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Create", 1)
	})

	// Error scenarios: repository-level errors during creation
	t.Run("on error repository create (create_conflict)", func(t *testing.T) {
		// Table-driven test for different repository error scenarios
		tests := []struct {
			name            string
			repositoryError error
		}{
			{
				name:            "on conflict duplicate warehouse code error",
				repositoryError: httperrors.ConflictError{Message: "the WarehouseCode already exists"},
			},
			{
				name:            "on internal server error",
				repositoryError: httperrors.InternalServerError{Message: "error creating warehouse"},
			},
		}

		// Execute each repository error scenario
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Arrange: mock repository to return specific error
				mockRepository := new(mocks.WarehouseRepositoryMock)
				serviceTest := service.NewWarehouseService(mockRepository)
				mockRepository.On("Create", warehouseAtt).Return(models.Warehouse{}, test.repositoryError)

				// Act: attempt warehouse creation that will result in repository error
				result, err := serviceTest.Create(warehouseAtt)

				// Assert: verify error propagation and empty result
				require.Equal(t, emptyWarehouse, result)
				require.Error(t, err)
				require.ErrorIs(t, err, test.repositoryError)
				mockRepository.AssertExpectations(t)
				mockRepository.AssertNumberOfCalls(t, "Create", 1)
			})
		}
	})

	// Validation scenarios: invalid warehouse attributes
	t.Run("on invalid warehouse attributes (create_invalid)", func(t *testing.T) {
		// Setup: common test infrastructure
		mockRepository := new(mocks.WarehouseRepositoryMock)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Table-driven test for different validation scenarios
		tests := []struct {
			name              string
			modifyStructFunc  func(w *models.WarehouseAttributes)
			expectedErrorText string
		}{
			{
				name: "Missing warehouse code",
				modifyStructFunc: func(w *models.WarehouseAttributes) {
					w.WarehouseCode = ""
				},
				expectedErrorText: "the field WarehouseCode must not be empty",
			},
			{
				name: "Missing address",
				modifyStructFunc: func(w *models.WarehouseAttributes) {
					w.Address = ""
				},
				expectedErrorText: "the field Address must not be empty",
			},
			{
				name: "Missing telephone",
				modifyStructFunc: func(w *models.WarehouseAttributes) {
					w.Telephone = ""
				},
				expectedErrorText: "the field Telephone must not be empty",
			},
			{
				name: "Missing minimum capacity",
				modifyStructFunc: func(w *models.WarehouseAttributes) {
					w.MinimunCapacity = 0
				},
				expectedErrorText: "the field MinimunCapacity must not be zero or negative",
			},
		}

		// Execute each validation scenario
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Arrange: create base structure with all fields valid, then modify specific field
				warehouseAtt := models.WarehouseAttributes{
					WarehouseCode:      "DHK",
					Address:            "123 Main St",
					Telephone:          "123-456-7890",
					MinimunCapacity:    100,
					MinimunTemperature: 15.0,
				}
				test.modifyStructFunc(&warehouseAtt)

				// Act: attempt creation with invalid attributes
				result, err := serviceTest.Create(warehouseAtt)

				// Assert: verify validation error and no repository interaction
				require.Equal(t, emptyWarehouse, result)
				require.Error(t, err)
				mockRepository.AssertExpectations(t)
				mockRepository.AssertNotCalled(t, "Create")
			})
		}
	})
}

// Test_GetAll tests the GetAll method of the WarehouseService.
//
// This test suite verifies that the warehouse retrieval business logic:
//   - Successfully retrieves all warehouses from repository
//   - Properly handles empty warehouse collections
//   - Propagates repository errors appropriately
//   - Returns warehouse collections in expected format
//
// Test scenarios:
//   - Success case: Multiple warehouses returned from repository
//   - Repository error: Internal server errors are properly propagated
func Test_GetAll(t *testing.T) {
	// Success scenario: retrieve multiple warehouses
	t.Run("on success (find_all)", func(t *testing.T) {
		// Arrange: mock repository to return warehouse collection
		mockRepository := new(mocks.WarehouseRepositoryMock)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Test data: collection of warehouses with different attributes
		warehouses := []models.Warehouse{
			{Id: 1, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Fake Street 123", Telephone: "123456789", MinimunCapacity: 10, MinimunTemperature: 5.0}},
			{Id: 2, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-002", Address: "Uwu 342", Telephone: "987654321", MinimunCapacity: 2, MinimunTemperature: 5.0}},
		}

		mockRepository.On("GetAll").Return(warehouses, nil)

		// Act: execute warehouse retrieval through service
		result, err := serviceTest.GetAll()

		// Assert: verify successful retrieval with correct warehouse collection
		require.NoError(t, err)
		require.Equal(t, warehouses, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})

	// Error scenario: repository fails during warehouse retrieval
	t.Run("on error repository get all (find_all_error)", func(t *testing.T) {
		// Arrange: mock repository to return internal server error
		errorReturned := httperrors.InternalServerError{Message: "error getting warehouses"}
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("GetAll").Return([]models.Warehouse{}, errorReturned)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: attempt warehouse retrieval that will result in repository error
		result, err := serviceTest.GetAll()

		// Assert: verify error propagation and empty result
		require.Error(t, err)
		require.ErrorIs(t, errorReturned, err)
		require.Empty(t, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})
}

// Test_GetByID tests the GetByID method of the WarehouseService.
//
// This test suite verifies that the warehouse retrieval by ID business logic:
//   - Successfully retrieves warehouses when valid ID is provided
//   - Returns appropriate errors when warehouse doesn't exist
//   - Properly handles repository-level errors during retrieval
//   - Maintains consistent error handling across different scenarios
//
// Test scenarios:
//   - Success case: Valid ID returns corresponding warehouse
//   - Not found: Non-existent ID returns not found error
//   - Repository error: Internal server errors are properly propagated
func Test_GetByID(t *testing.T) {
	// Success scenario: retrieve existing warehouse by valid ID
	t.Run("on success (find_by_id_existent)", func(t *testing.T) {
		// Arrange: mock repository to return existing warehouse
		mockRepository := new(mocks.WarehouseRepositoryMock)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Test data: warehouse with complete attributes
		warehouse := models.Warehouse{
			Id: 1,
			WarehouseAttributes: models.WarehouseAttributes{
				WarehouseCode:      "WH-001",
				Address:            "Fake Street 123",
				Telephone:          "123456789",
				MinimunCapacity:    10,
				MinimunTemperature: 5.0},
		}

		mockRepository.On("GetByID", warehouse.Id).Return(warehouse, nil)

		// Act: execute warehouse retrieval by ID through service
		result, err := serviceTest.GetByID(warehouse.Id)

		// Assert: verify successful retrieval with correct warehouse data
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
	})

	// Error scenarios: repository-level errors during retrieval
	t.Run("on error repository get by id (find_by_id_non_existent)", func(t *testing.T) {
		// Table-driven test for different repository error scenarios
		tests := []struct {
			name            string
			repositoryError error
		}{
			{
				name:            "on not found error",
				repositoryError: httperrors.NotFoundError{Message: "warehouse not found"},
			},
			{
				name: "on internal server error",
				repositoryError: httperrors.InternalServerError{
					Message: "error obtaining warehouse by ID",
				},
			},
		}

		idSearched := 1

		// Execute each repository error scenario
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Arrange: mock repository to return specific error
				mockRepository := new(mocks.WarehouseRepositoryMock)
				mockRepository.On("GetByID", idSearched).Return(models.Warehouse{}, test.repositoryError)
				serviceTest := service.NewWarehouseService(mockRepository)

				// Act: attempt retrieval that will result in repository error
				result, err := serviceTest.GetByID(idSearched)

				// Assert: verify error propagation and empty result
				require.Equal(t, emptyWarehouse, result)
				require.Error(t, err)
				require.ErrorIs(t, test.repositoryError, err)
				mockRepository.AssertExpectations(t)
				mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
			})
		}
	})
}

// Test_Update tests the Update method of the WarehouseService.
//
// This test suite verifies that the warehouse update business logic:
//   - Successfully updates warehouses with valid data and existing ID
//   - Validates warehouse existence before applying updates
//   - Detects and prevents duplicate warehouse codes during updates
//   - Properly handles repository-level errors during update operations
//   - Maintains data consistency through validation and conflict detection
//
// Test scenarios:
//   - Success case: Valid data and existing ID return updated warehouse
//   - Not found: Non-existent ID returns not found error
//   - Conflict: Duplicate warehouse code returns conflict error
//   - Repository error: Internal server errors are properly propagated
func Test_Update(t *testing.T) {
	// Test data: attributes for partial update
	warehouseAttributesTest := models.WarehouseAttributes{
		Address: "Fake Street 123 Updated",
	}

	// Test data: complete warehouse attributes for existing warehouse
	warehouseAttTest := models.WarehouseAttributes{
		WarehouseCode:      "WH-001",
		Address:            "Calle Falsa 123",
		Telephone:          "123456789",
		MinimunCapacity:    100,
		MinimunTemperature: 4.5,
	}

	warehouseTest := models.Warehouse{
		Id:                  1,
		WarehouseAttributes: warehouseAttTest,
	}

	// Updated warehouse with merged attributes
	wareHouseUpdated := warehouseTest
	warehouseAttUpdated := warehouseAttTest
	warehouseAttUpdated.Address = warehouseAttributesTest.Address

	// Test data: collection of warehouses for conflict detection
	warehouses := []models.Warehouse{
		{Id: 1, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Fake Street 123", Telephone: "123456789", MinimunCapacity: 10, MinimunTemperature: 5.0}},
		{Id: 2, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-002", Address: "Uwu 342", Telephone: "987654321", MinimunCapacity: 2, MinimunTemperature: 5.0}},
	}

	// Success scenario: update existing warehouse with valid data
	t.Run("on success", func(t *testing.T) {
		// Arrange: mock repository for successful update workflow
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).Return(warehouseTest, nil)
		mockRepository.On("GetAll").Return(warehouses, nil)
		mockRepository.On("Update", warehouseTest.Id, warehouseAttUpdated).Return(wareHouseUpdated, nil)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: execute warehouse update through service
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttUpdated)

		// Assert: verify successful update with correct data
		require.NoError(t, err)
		require.Equal(t, wareHouseUpdated, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Update", 1)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
	})

	// Error scenario: attempt to update non-existent warehouse
	t.Run("on id not found (update_non_existent)", func(t *testing.T) {
		// Arrange: mock repository to return not found error during existence check
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).Return(models.Warehouse{}, httperrors.NotFoundError{Message: "warehouse not found"})
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: attempt update of non-existent warehouse
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttributesTest)

		// Assert: verify not found error and no update operation
		require.Equal(t, emptyWarehouse, result)
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.NotFoundError{Message: "warehouse not found"})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
		mockRepository.AssertNotCalled(t, "Update")
	})

	// Error scenario: duplicate warehouse code conflict
	t.Run("WarehouseCode duplicated", func(t *testing.T) {
		// Arrange: setup warehouses with duplicate codes for conflict detection
		warehousesTest := []models.Warehouse{
			{Id: 1, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Fake Street 123", Telephone: "123456789", MinimunCapacity: 10, MinimunTemperature: 5.0}},
			{Id: 2, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Uwu 342", Telephone: "987654321", MinimunCapacity: 2, MinimunTemperature: 5.0}},
		}
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).Return(warehouseTest, nil)
		mockRepository.On("GetAll").Return(warehousesTest, nil)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: attempt update that would create duplicate warehouse code
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttUpdated)

		// Assert: verify conflict error and no update operation
		require.Equal(t, emptyWarehouse, result)
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.ConflictError{
			Message: "the WarehouseCode already exists",
		})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})

	// Error scenario: repository error during warehouse collection retrieval
	t.Run("Repository error during validation", func(t *testing.T) {
		// Arrange: mock repository to return error during GetAll operation
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).Return(warehouseTest, nil)
		mockRepository.On("GetAll").Return([]models.Warehouse{}, httperrors.InternalServerError{Message: "error reading warehouse data"})
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: attempt update that fails during validation
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttUpdated)

		// Assert: verify internal server error propagation
		require.Equal(t, emptyWarehouse, result)
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.InternalServerError{Message: "error reading warehouse data"})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})
}

// Test_Delete tests the Delete method of the WarehouseService.
//
// This test suite verifies that the warehouse deletion business logic:
//   - Successfully deletes existing warehouses through repository
//   - Properly handles attempts to delete non-existent warehouses
//   - Propagates repository errors appropriately
//   - Maintains referential integrity through proper error handling
//
// Test scenarios:
//   - Success case: Existing warehouse ID results in successful deletion
//   - Not found: Non-existent warehouse ID returns not found error
func Test_Delete(t *testing.T) {
	idSearched := 1

	// Success scenario: delete existing warehouse
	t.Run("on success (delete_ok)", func(t *testing.T) {
		// Arrange: mock repository to return successful deletion
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("Delete", idSearched).Return(nil)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: execute warehouse deletion through service
		err := serviceTest.Delete(idSearched)

		// Assert: verify successful deletion with no error
		require.NoError(t, err)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Delete", 1)
	})

	// Error scenario: attempt to delete non-existent warehouse
	t.Run("on repository error (delete_non_existent)", func(t *testing.T) {
		// Arrange: mock repository to return not found error
		mockRepository := new(mocks.WarehouseRepositoryMock)
		mockRepository.On("Delete", idSearched).Return(httperrors.NotFoundError{Message: "warehouse not found"})
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act: attempt deletion of non-existent warehouse
		err := serviceTest.Delete(idSearched)

		// Assert: verify not found error propagation
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.NotFoundError{Message: "warehouse not found"})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Delete", 1)
	})
}
