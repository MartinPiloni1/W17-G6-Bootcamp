package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"slices"
)

type WarehouseServiceDefault struct {
	rp repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseServiceDefault {
	return &WarehouseServiceDefault{
		rp: repo,
	}
}

func (p *WarehouseServiceDefault) Create(warehouse models.WarehouseAttributes) (models.Warehouse, error) {
	if warehouse.WarehouseCode == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field WarehouseCode must not be empty"}
	}
	if warehouse.Address == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field Address must not be empty"}
	}
	if warehouse.Telephone == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field Telephone must not be empty"}
	}
	if warehouse.MinimunCapacity <= 0 {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field MinimunCapacity must not be zero or negative"}
	}
	warehouses, err := p.rp.GetAll()
	if err != nil {
		return models.Warehouse{}, err
	}
	for _, w := range warehouses {
		if w.WarehouseCode == warehouse.WarehouseCode {
			return models.Warehouse{}, httperrors.ConflictError{Message: "the WarehouseCode already exists"}
		}
	}
	return p.rp.Create(warehouse)
}

func (p *WarehouseServiceDefault) GetAll() ([]models.Warehouse, error) {
	result, err := p.rp.GetAll()
	if err != nil {
		return []models.Warehouse{}, err
	}
	slicedData := utils.MapToSlice(result)
	slices.SortFunc(slicedData, func(a, b models.Warehouse) int {
		return a.Id - b.Id
	})
	return slicedData, nil

}

func (p *WarehouseServiceDefault) GetByID(id int) (models.Warehouse, error) {
	return p.rp.GetByID(id)
}

func (p *WarehouseServiceDefault) Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	warehouse, err := p.rp.GetByID(id)
	if err != nil {
		return models.Warehouse{}, err
	}
	errZeroVelue := utils.ApplyNonZero(&warehouse.WarehouseAttributes, warehouseAttributes)
	if errZeroVelue != nil {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the input fields are not valid"}
	}
	warehouses, err := p.rp.GetAll()
	if err != nil {
		return models.Warehouse{}, err
	}
	for _, w := range warehouses {
		if w.WarehouseCode == warehouse.WarehouseAttributes.WarehouseCode {
			return models.Warehouse{}, httperrors.ConflictError{Message: "the WarehouseCode already exists"}
		}
	}
	return p.rp.Update(id, warehouse.WarehouseAttributes)
}

func (p *WarehouseServiceDefault) Delete(id int) error {
	return p.rp.Delete(id)
}
