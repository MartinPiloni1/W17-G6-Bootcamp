package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"

type ProductService interface {
	Create(product models.ProductAttributes) (models.Product, error)
	GetAll() ([]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, productAttributes models.ProductPatchRequest) (models.Product, error)
	Delete(id int) error
}
