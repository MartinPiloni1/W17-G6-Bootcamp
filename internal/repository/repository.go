package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type WarehouseRepository interface {
	GetAll() (map[int]models.Warehouse, error)
}
