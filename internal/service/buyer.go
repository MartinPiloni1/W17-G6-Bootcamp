package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type BuyerServiceDefault struct {
	rp repository.BuyerRepository
}

func NewBuyerServiceDefault(rp repository.BuyerRepository) BuyerService {
	return &BuyerServiceDefault{rp: rp}
}

func (s *BuyerServiceDefault) GetAll() (map[int]models.Buyer, error) {
	buyerData, err := s.rp.GetAll()
	return buyerData, err
}

func (s *BuyerServiceDefault) GetByID(id int) (models.Buyer, error) {
	buyer, err := s.rp.GetByID(id)
	return buyer, err
}

func (s *BuyerServiceDefault) Delete(id int) error {
	err := s.rp.Delete(id)
	return err
}

func (s *BuyerServiceDefault) Create(newBuyer models.BuyerAttributes) (models.Buyer, error) {
	// check if the cardNumber already exist
	exist, err := s.rp.CardNumberIdAlreadyExist(newBuyer.CardNumberId)
	if err != nil {
		return models.Buyer{}, err
	}

	if exist {
		return models.Buyer{}, httperrors.ConflictError{Message: "Buyer already exist"}
	}
	buyer, err := s.rp.Create(newBuyer)
	return buyer, err
}

func (s *BuyerServiceDefault) Update(id int, BuyerData models.BuyerPatchRequest) (models.Buyer, error) {
	buyer, err := s.rp.GetByID(id)
	if err != nil {
		return models.Buyer{}, err
	}
	if BuyerData.CardNumberId != nil {
		exist, err := s.rp.CardNumberIdAlreadyExist(*BuyerData.CardNumberId)
		if err != nil {
			return models.Buyer{}, err
		}

		if exist {
			return models.Buyer{}, httperrors.ConflictError{Message: "CardNumberId already in use"}
		}
		buyer.CardNumberId = *BuyerData.CardNumberId
	}

	if BuyerData.FirstName != nil {
		buyer.FirstName = *BuyerData.FirstName
	}
	if BuyerData.LastName != nil {
		buyer.LastName = *BuyerData.LastName
	}

	updatedBuyer, err := s.rp.Update(id, buyer)
	return updatedBuyer, err
}
