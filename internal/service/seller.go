package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type SellerServiceDefault struct {
	rp repository.SellerRepository
}

func NewSellerService(repo repository.SellerRepository) SellerService {
	return &SellerServiceDefault{rp: repo}
}

func (s SellerServiceDefault) Create(seller models.SellerAttributes) (models.Seller, error) {
	if seller.CID <= 0 || seller.CompanyName == "" || seller.Address == "" || seller.Telephone == "" || seller.LocalityID == "" {
		return models.Seller{}, httperrors.BadRequestError{
			Message: "Invalid seller data",
		}
	}

	return s.rp.Create(seller)
}

func (s SellerServiceDefault) Delete(id int) error {
	return s.rp.Delete(id)
}

func (s SellerServiceDefault) GetAll() ([]models.Seller, error) {
	return s.rp.GetAll()
}

func (s SellerServiceDefault) GetByID(id int) (models.Seller, error) {
	return s.rp.GetByID(id)
}

func (s SellerServiceDefault) Update(id int, data *models.SellerAttributes) (models.Seller, error) {
	if data.CID != 0 {
		all, err := s.rp.GetAll()
		if err != nil {
			return models.Seller{}, err
		}
		for _, existing := range all {
			if existing.CID == data.CID && existing.ID != id {
				return models.Seller{}, httperrors.ConflictError{
					Message: "CID already exists",
				}
			}
		}
	}
	return s.rp.Update(id, data)
}
