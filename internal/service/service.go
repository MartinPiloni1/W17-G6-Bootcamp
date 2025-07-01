package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type WarehouseService interface {
	GetAll() (map[int]models.Warehouse, error)
	Create(Product models.WarehouseAttributes) (models.Warehouse, error)
	GetByID(id int) (models.Warehouse, error)
	Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error)
	Delete(id int) error
}
