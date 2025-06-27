package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type WarehouseService interface {
	GetAll() (map[int]models.Warehouse, error)
}
