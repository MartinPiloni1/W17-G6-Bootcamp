package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type SellerRepositoryInterface interface {
	Create(seller models.Seller) (*models.Seller, error)
	GetAll() (map[int]models.Seller, error)
	GetByID(id int) (models.Seller, error)
	Update(id int, data *models.Seller) (*models.Seller, error)
	Delete(id int) error
}
