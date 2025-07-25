package repository

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// WarehouseRepositoryMock is a mock implementation of the warehouse repository for testing.
type WarehouseRepositoryMock struct {
	mock.Mock
}

// GetAll retrieves all warehouses from the mock repository.
// Receives: nothing.
// Returns: a slice of Warehouse and an error.
func (w *WarehouseRepositoryMock) GetAll() ([]models.Warehouse, error) {
	args := w.Called()
	return args.Get(0).([]models.Warehouse), args.Error(1)
}

// Create adds a new warehouse to the mock repository.
// Receives: WarehouseAttributes with the data for the new warehouse.
// Returns: the created Warehouse and an error.
func (w *WarehouseRepositoryMock) Create(warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	args := w.Called(warehouseAttributes)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

// Update modifies an existing warehouse in the mock repository.
// Receives: an int (id) identifying the warehouse, and WarehouseAttributes with updated data.
// Returns: the updated Warehouse and an error.
func (w *WarehouseRepositoryMock) Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	args := w.Called(id, warehouseAttributes)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

// GetByID retrieves a warehouse by its id from the mock repository.
// Receives: an int (id) identifying the warehouse.
// Returns: the Warehouse and an error.
func (w *WarehouseRepositoryMock) GetByID(id int) (models.Warehouse, error) {
	args := w.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

// Delete removes a warehouse by its id from the mock repository.
// Receives: an int (id) identifying the warehouse.
// Returns: an error if the operation fails.
func (w *WarehouseRepositoryMock) Delete(id int) error {
	args := w.Called(id)
	return args.Error(0)
}
