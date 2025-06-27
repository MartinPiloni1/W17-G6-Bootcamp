package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type ProductRepositoryInterface interface {
	Create(Product models.Product) (*models.Product, error)
	GetAll() (map[int]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, data models.Product) (models.Product, error)
	Delete(id int) error
}
