package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"

type ProductServiceInterface interface {
	Create(Product models.ProductAtributtes) (models.Product, error)
	GetAll() (map[int]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(id int, productAtributtes models.ProductAtributtes) (models.Product, error)
	Delete(id int) error
}
