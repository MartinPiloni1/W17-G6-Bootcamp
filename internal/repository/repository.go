package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type BuyerRepository interface {
	GetAll() (map[int]models.Buyer, error)
	GetByID(id int) (models.Buyer, error)
	Delete(id int) error
	Create(newBuyer models.BuyerAttributes) (models.Buyer, error)
}
