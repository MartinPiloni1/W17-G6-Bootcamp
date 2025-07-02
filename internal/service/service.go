package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type ProductService interface {
	Create(product models.ProductAttributes) (models.Product, error)
	GetAll() ([]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, productAttributes models.ProductPatchRequest) (models.Product, error)
	Delete(id int) error
}

type BuyerService interface {
	Create(newBuyer models.BuyerAttributes) (models.Buyer, error)
	GetAll() ([]models.Buyer, error)
	GetByID(id int) (models.Buyer, error)
	Update(id int, BuyerData models.BuyerPatchRequest) (models.Buyer, error)
	Delete(id int) error
}
