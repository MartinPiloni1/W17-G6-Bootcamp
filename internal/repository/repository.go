package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"

type ProductRepositoryInterface interface {
	Create(Product models.ProductAttributes) (models.Product, error)
	GetAll() (map[int]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, productAttributes models.ProductAttributes) (models.Product, error)
	Delete(id int) error
}
