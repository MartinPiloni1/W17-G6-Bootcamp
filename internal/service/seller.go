package service

import (
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type SellerServiceImpl struct {
	rp repository.SellerRepositoryInterface
}

func NewSellerService(repo repository.SellerRepositoryInterface) SellerServiceInterface {
	return &SellerServiceImpl{rp: repo}
}

func (s *SellerServiceImpl) Create(seller models.Seller) (*models.Seller, error) {
	if seller.CID == 0 || seller.CompanyName == "" || seller.Address == "" || seller.Telephone == "" {
		return &models.Seller{}, errors.New("missing required fields")
	}
	all, err := s.rp.GetAll()
	if err != nil {
		return &models.Seller{}, err
	}
	for _, v := range all {
		if v.CID == seller.CID {
			return &models.Seller{}, errors.New("cid already exists")
		}
	}
	return s.rp.Create(seller)
}

func (s *SellerServiceImpl) Delete(id int) error {
	return s.rp.Delete(id)
}

func (s *SellerServiceImpl) GetAll() (map[int]models.Seller, error) {
	return s.rp.GetAll()
}

func (s *SellerServiceImpl) GetByID(id int) (models.Seller, error) {
	return s.rp.GetByID(id)
}

func (s *SellerServiceImpl) Update(id int, data *models.Seller) (*models.Seller, error) {
	//no cambiar un seller inactivo, el CID debe seguir siendo único, etc.
	if data.CID != 0 {
		all, err := s.rp.GetAll()
		if err != nil {
			return nil, err
		}
		for _, existing := range all {
			if existing.CID == data.CID && existing.ID != id {
				return nil, errors.New("cid already exists")
			}
		}
	}
	return s.rp.Update(id, data)
}
