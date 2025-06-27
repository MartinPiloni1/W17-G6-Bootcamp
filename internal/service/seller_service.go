package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type SellerServiceImpl struct {
	repo repository.SellerRepository
}

func NewSellerService(repo repository.SellerRepository) SellerService {
	return &SellerServiceImpl{repo: repo}
}

// Create implements SellerService.
func (s *SellerServiceImpl) Create(seller models.Seller) (models.Seller, error) {
	panic("unimplemented")
}

// Delete implements SellerService.
func (s *SellerServiceImpl) Delete(id int) error {
	panic("unimplemented")
}

// GetAll implements SellerService.
func (s *SellerServiceImpl) GetAll() (map[int]models.Seller, error) {
	return s.repo.GetAll()
}

// GetByID implements SellerService.
func (s *SellerServiceImpl) GetByID(id int) (models.Seller, error) {
	panic("unimplemented")
}

// Update implements SellerService.
func (s *SellerServiceImpl) Update(id int, data models.Seller) (models.Seller, error) {
	panic("unimplemented")
}
