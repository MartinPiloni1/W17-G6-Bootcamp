package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type WarehouseServiceDefault struct {
	rp repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseServiceDefault {
	return &WarehouseServiceDefault{
		rp: repo,
	}
}

func (p *WarehouseServiceDefault) Create(product models.WarehouseAttributes) (models.Warehouse, error) {
	if product.WarehouseCode == "" {
		return models.Warehouse{}, httperrors.BadRequestError{Message: "the field WarehouseCode must not be empty"}
	}
	return p.rp.Create(product)
}

func (p *WarehouseServiceDefault) GetAll() (map[int]models.Warehouse, error) {
	result, err := p.rp.GetAll()
	if err != nil {
		return map[int]models.Warehouse{}, err
	}
	return result, nil
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
	return p.rp.Update(id, warehouse.WarehouseAttributes)
}

func (p *WarehouseServiceDefault) Delete(id int) error {
	return p.rp.Delete(id)
}
