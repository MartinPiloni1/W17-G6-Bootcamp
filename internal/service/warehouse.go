package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type WarehouseServiceDefault struct {
	repo repository.WarehouseRepository
}

// NewWarehouseService creates a new warehouse service.
func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseServiceDefault {
	return &WarehouseServiceDefault{repo: repo}
}

// Create validates and adds a new warehouse.
func (w *WarehouseServiceDefault) Create(warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	if warehouseAttributes.WarehouseCode == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field WarehouseCode must not be empty"}
	}
	if warehouseAttributes.Address == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field Address must not be empty"}
	}
	if warehouseAttributes.Telephone == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field Telephone must not be empty"}
	}
	if warehouseAttributes.MinimunCapacity <= 0 {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field MinimunCapacity must not be zero or negative"}
	}
	return w.repo.Create(warehouseAttributes)
}

// GetAll returns all warehouses.
func (w *WarehouseServiceDefault) GetAll() ([]models.Warehouse, error) {
	result, err := w.repo.GetAll()
	if err != nil {
		return []models.Warehouse{}, err
	}
	return result, nil

}

// GetByID returns a warehouse by ID.
func (w *WarehouseServiceDefault) GetByID(id int) (models.Warehouse, error) {
	return w.repo.GetByID(id)
}

// Update modifies a warehouse by ID.
func (w *WarehouseServiceDefault) Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	warehouse, err := w.repo.GetByID(id)
	if err != nil {
		return models.Warehouse{}, err
	}
	errZeroVelue := utils.ApplyNonZero(&warehouse.WarehouseAttributes, warehouseAttributes)
	if errZeroVelue != nil {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the input fields are not valid"}
	}
	warehouses, err := w.repo.GetAll()
	if err != nil {
		return models.Warehouse{}, err
	}
	for _, w := range warehouses {
		if w.WarehouseCode == warehouse.WarehouseAttributes.WarehouseCode && w.Id != id {
			return models.Warehouse{}, httperrors.ConflictError{Message: "the WarehouseCode already exists"}
		}
	}
	return w.repo.Update(id, warehouse.WarehouseAttributes)
}

// Delete removes a warehouse by ID.
func (w *WarehouseServiceDefault) Delete(id int) error {
	return w.repo.Delete(id)
}
