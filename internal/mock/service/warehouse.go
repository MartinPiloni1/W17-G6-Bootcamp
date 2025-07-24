package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

// WarehouseServiceMock is a mock implementation of the warehouse service for testing.
type WarehouseServiceMock struct {
	mock.Mock
}

// GetAll retrieves all warehouses from the mock service.
// Receives: nothing.
// Returns: a slice of Warehouse and an error.
func (w *WarehouseServiceMock) GetAll() ([]models.Warehouse, error) {
	args := w.Called()
	return args.Get(0).([]models.Warehouse), args.Error(1)
}

// Create adds a new warehouse to the mock service.
// Receives: WarehouseAttributes with the data for the new warehouse.
// Returns: the created Warehouse and an error.
func (w *WarehouseServiceMock) Create(warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	args := w.Called(warehouseAttributes)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

// Update modifies an existing warehouse in the mock service.
// Receives: an int (id) identifying the warehouse, and WarehouseAttributes with updated data.
// Returns: the updated Warehouse and an error.
func (w *WarehouseServiceMock) Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	args := w.Called(id, warehouseAttributes)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

// GetById retrieves a warehouse by its id from the mock service.
// Receives: an int (id) identifying the warehouse.
// Returns: the Warehouse and an error.
func (w *WarehouseServiceMock) GetById(id int) (models.Warehouse, error) {
	args := w.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

// Delete removes a warehouse by its id from the mock service.
// Receives: an int (id) identifying the warehouse.
// Returns: an error if the operation fails.
func (w *WarehouseServiceMock) Delete(id int) error {
	args := w.Called(id)
	return args.Error(0)
}
