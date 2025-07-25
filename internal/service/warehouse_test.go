package service_test

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mock/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/require"
	"testing"
)

var emptyWarehouse = models.Warehouse{}

// Test_Create tests the Create method of the WarehouseService.
func Test_Create(t *testing.T) {

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

	t.Run("create warehouse success (create_ok)", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("Create", warehouseAtt).Return(warehouse, nil)

		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		result, err := serviceTest.Create(warehouseAtt)

		// Assert
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("on error repository create (create_conflict)", func(t *testing.T) {
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

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Arrange
				mockRepository := new(repository.WarehouseRepositoryMock)
				serviceTest := service.NewWarehouseService(mockRepository)
				mockRepository.On("Create", warehouseAtt).Return(models.Warehouse{}, test.repositoryError)

				// Act
				result, err := serviceTest.Create(warehouseAtt)

				// Assert
				require.Equal(t, emptyWarehouse, result)
				require.Error(t, err)
				require.ErrorIs(t, err, test.repositoryError)
				mockRepository.AssertExpectations(t)
				mockRepository.AssertNumberOfCalls(t, "Create", 1)
			})
		}

	})

	t.Run("on invalid warehouse attributes (create_invalid)", func(t *testing.T) {
		mockRepository := new(repository.WarehouseRepositoryMock)
		serviceTest := service.NewWarehouseService(mockRepository)

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

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Arrange
				// base structure with all fields valid
				warehouseAtt := models.WarehouseAttributes{
					WarehouseCode:      "DHK",
					Address:            "123 Main St",
					Telephone:          "123-456-7890",
					MinimunCapacity:    100,
					MinimunTemperature: 15.0,
				}

				test.modifyStructFunc(&warehouseAtt)

				// Act
				result, err := serviceTest.Create(warehouseAtt)

				// Assert
				require.Equal(t, emptyWarehouse, result)
				require.Error(t, err)
				mockRepository.AssertExpectations(t)
				mockRepository.AssertNotCalled(t, "Create")
			})
		}
	})
}

// Test_GetAll tests the GetAll method of the WarehouseService.
func Test_GetAll(t *testing.T) {
	t.Run("on success (find_all)", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.WarehouseRepositoryMock)
		serviceTest := service.NewWarehouseService(mockRepository)

		warehouses := []models.Warehouse{
			{Id: 1, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Fake Street 123", Telephone: "123456789", MinimunCapacity: 10, MinimunTemperature: 5.0}},
			{Id: 2, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-002", Address: "Uwu 342", Telephone: "987654321", MinimunCapacity: 2, MinimunTemperature: 5.0}},
		}

		mockRepository.On("GetAll").Return(warehouses, nil)

		// Act
		result, err := serviceTest.GetAll()

		// Assert
		require.NoError(t, err)
		require.Equal(t, warehouses, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})

	t.Run("on error repository get all (find_all_error)", func(t *testing.T) {
		// Arrange
		errorReturned := httperrors.InternalServerError{Message: "error getting warehouses"}
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("GetAll").Return([]models.Warehouse{}, errorReturned)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		result, err := serviceTest.GetAll()

		// Assert
		require.Error(t, err)
		require.ErrorIs(t, errorReturned, err)
		require.Empty(t, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})

}

// Test_GetByID tests the GetByID method of the WarehouseService.
func Test_GetByID(t *testing.T) {
	t.Run("on success (find_by_id_existent)", func(t *testing.T) {
		mockRepository := new(repository.WarehouseRepositoryMock)
		serviceTest := service.NewWarehouseService(mockRepository)

		warehouse := models.Warehouse{
			Id: 1,
			WarehouseAttributes: models.WarehouseAttributes{
				WarehouseCode:      "WH-001",
				Address:            "Fake Street 123",
				Telephone:          "123456789",
				MinimunCapacity:    10,
				MinimunTemperature: 5.0},
		}

		mockRepository.On("GetByID", warehouse.Id).
			Return(warehouse, nil)

		// Act
		result, err := serviceTest.GetByID(warehouse.Id)

		// Assert
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)

	})

	t.Run("on error repository get by id (find_by_id_non_existent)", func(t *testing.T) {
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

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// arrange
				mockRepository := new(repository.WarehouseRepositoryMock)
				mockRepository.On("GetByID", idSearched).Return(models.Warehouse{}, test.repositoryError)
				serviceTest := service.NewWarehouseService(mockRepository)

				// act
				result, err := serviceTest.GetByID(idSearched)

				// assert
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
func Test_Update(t *testing.T) {
	warehouseAttributesTest := models.WarehouseAttributes{
		Address: "Fake Street 123 Updated",
	}

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

	// Update fields that will be updated
	wareHouseUpdated := warehouseTest
	warehouseAttUpdated := warehouseAttTest
	warehouseAttUpdated.Address = warehouseAttributesTest.Address

	warehouses := []models.Warehouse{
		{Id: 1, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Fake Street 123", Telephone: "123456789", MinimunCapacity: 10, MinimunTemperature: 5.0}},
		{Id: 2, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-002", Address: "Uwu 342", Telephone: "987654321", MinimunCapacity: 2, MinimunTemperature: 5.0}},
	}
	t.Run("on success", func(t *testing.T) {
		// Arrange

		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).
			Return(warehouseTest, nil)
		mockRepository.On("GetAll").
			Return(warehouses, nil)
		mockRepository.On("Update", warehouseTest.Id, warehouseAttUpdated).
			Return(wareHouseUpdated, nil)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttUpdated)

		// Assert
		require.NoError(t, err)
		require.Equal(t, wareHouseUpdated, result)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Update", 1)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
	})

	t.Run("on id not found (update_non_existent)", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).Return(models.Warehouse{}, httperrors.NotFoundError{Message: "warehouse not found"})
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttributesTest)

		// Assert
		require.Equal(t, emptyWarehouse, result)
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.NotFoundError{Message: "warehouse not found"})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
		mockRepository.AssertNotCalled(t, "Update")
	})

	t.Run("WarehouseCode duplicated", func(t *testing.T) {
		// Arrange

		warehousesTest := []models.Warehouse{
			{Id: 1, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Fake Street 123", Telephone: "123456789", MinimunCapacity: 10, MinimunTemperature: 5.0}},
			{Id: 2, WarehouseAttributes: models.WarehouseAttributes{WarehouseCode: "WH-001", Address: "Uwu 342", Telephone: "987654321", MinimunCapacity: 2, MinimunTemperature: 5.0}},
		}
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).
			Return(warehouseTest, nil)
		mockRepository.On("GetAll").
			Return(warehousesTest, nil)

		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttUpdated)

		// Assert
		require.Equal(t, emptyWarehouse, result)
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.ConflictError{
			Message: "the WarehouseCode already exists",
		})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})

	t.Run("WarehouseCode duplicated", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("GetByID", warehouseTest.Id).
			Return(warehouseTest, nil)
		mockRepository.On("GetAll").
			Return([]models.Warehouse{}, httperrors.InternalServerError{Message: "error reading warehouse data"})

		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		result, err := serviceTest.Update(warehouseTest.Id, warehouseAttUpdated)

		// Assert
		require.Equal(t, emptyWarehouse, result)
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.InternalServerError{Message: "error reading warehouse data"})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "GetByID", 1)
		mockRepository.AssertNumberOfCalls(t, "GetAll", 1)
	})
}

// Test_Delete tests the Delete method of the WarehouseService.
func Test_Delete(t *testing.T) {
	idSearched := 1
	t.Run("on success (delete_ok)", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("Delete", idSearched).
			Return(nil)
		serviceTest := service.NewWarehouseService(mockRepository)

		// Act
		err := serviceTest.Delete(idSearched)

		// Assert
		require.NoError(t, err)
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Delete", 1)
	})

	t.Run("on repository error (delete_non_existent)", func(t *testing.T) {
		// arrange
		mockRepository := new(repository.WarehouseRepositoryMock)
		mockRepository.On("Delete", idSearched).Return(httperrors.NotFoundError{Message: "warehouse not found"})
		serviceTest := service.NewWarehouseService(mockRepository)

		// act
		err := serviceTest.Delete(idSearched)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, httperrors.NotFoundError{Message: "warehouse not found"})
		mockRepository.AssertExpectations(t)
		mockRepository.AssertNumberOfCalls(t, "Delete", 1)
	})
}
