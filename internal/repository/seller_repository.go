package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type SellerRepositoryImpl struct {
	filePath string
}

func NewSellerRepository() SellerRepository {
	return &SellerRepositoryImpl{filePath: os.Getenv("FILE_PATH_DEFAULT")}
}

// Create implements SellerRepository.
func (s *SellerRepositoryImpl) Create(seller models.Seller) (*models.Seller, error) {
	panic("unimplemented")
}

// Delete implements SellerRepository.
func (s *SellerRepositoryImpl) Delete(id int) error {
	panic("unimplemented")
}

// GetAll implements SellerRepository.
func (s *SellerRepositoryImpl) GetAll() (map[int]models.Seller, error) {
	data, err := utils.Read[models.Seller](s.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetByCID implements SellerRepository.
func (s *SellerRepositoryImpl) GetByCID(cid int) (models.Seller, error) {
	panic("unimplemented")
}

// GetByID implements SellerRepository.
func (s *SellerRepositoryImpl) GetByID(id int) (models.Seller, error) {
	panic("unimplemented")
}

// Update implements SellerRepository.
func (s *SellerRepositoryImpl) Update(id int, data models.Seller) (models.Seller, error) {
	panic("unimplemented")
}
