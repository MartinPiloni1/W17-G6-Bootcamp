package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
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
	buyer, err := s.rp.Create(newBuyer)
	return buyer, err
}
