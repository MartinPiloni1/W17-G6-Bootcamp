package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type SellerRepositoryInterface interface {
	Create(seller models.SellerAttributes) (models.Seller, error)
	GetAll() (map[int]models.Seller, error)
	GetByID(id int) (models.Seller, error)
	Update(id int, data *models.SellerAttributes) (models.Seller, error)
	Delete(id int) error
}
