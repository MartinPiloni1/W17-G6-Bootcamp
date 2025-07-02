package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"

type ProductRepository interface {
	Create(product models.ProductAttributes) (models.Product, error)
	GetAll() (map[int]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, product models.Product) (models.Product, error)
	Delete(id int) error
}
