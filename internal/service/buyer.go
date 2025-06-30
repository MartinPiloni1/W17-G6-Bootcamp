package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type BuyerService struct {
	rp repository.BuyerRepositoryInterface
}

func NewBuyerService(rp repository.BuyerRepositoryInterface) BuyerServiceInterface {
	return &BuyerService{rp: rp}
}

func (s *BuyerService) GetAll() (map[int]models.Buyer, error) {
	buyerData, err := s.rp.GetAll()
	return buyerData, err
}

func (s *BuyerService) GetByID(id int) (models.Buyer, error) {
	buyer, err := s.rp.GetByID(id)
	return buyer, err
}
