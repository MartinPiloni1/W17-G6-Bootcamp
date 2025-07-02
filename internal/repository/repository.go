package repository

import "github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"

type BuyerRepository interface {
	Create(newBuyer models.BuyerAttributes) (models.Buyer, error)
	GetAll() (map[int]models.Buyer, error)
	GetByID(id int) (models.Buyer, error)
	Update(id int, updatedBuyer models.Buyer) (models.Buyer, error)
	Delete(id int) error
	CardNumberIdAlreadyExist(newCardNumberId int) (bool, error)
}
