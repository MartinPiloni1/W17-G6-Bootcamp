package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type WarehouseServiceDefault struct {
	rp repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseServiceDefault {
	return &WarehouseServiceDefault{
		rp: repo,
	}
}

func (p *WarehouseServiceDefault) GetAll() (map[int]models.Warehouse, error) {
	result, err := p.rp.GetAll()
	if err != nil {
		return map[int]models.Warehouse{}, err
	}
	return result, nil
}
