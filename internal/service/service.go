package service

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type BuyerService interface {
	Create(newBuyer models.BuyerAttributes) (models.Buyer, error)
	GetAll() ([]models.Buyer, error)
	GetByID(id int) (models.Buyer, error)
	Update(id int, BuyerData models.BuyerPatchRequest) (models.Buyer, error)
	Delete(id int) error
}
