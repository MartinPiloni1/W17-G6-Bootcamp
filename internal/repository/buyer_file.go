package repository

import (
	"fmt"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type BuyerRepositoryFile struct {
	filePath string
}

func NewBuyerRepositoryFile(filePath string) BuyerRepositoryInterface {
	return &BuyerRepositoryFile{filePath: filePath}
}

func (r *BuyerRepositoryFile) GetAll() (map[int]models.Buyer, error) {
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return nil, err
	}
	return buyersData, nil
}

func (r *BuyerRepositoryFile) GetByID(id int) (models.Buyer, error) {
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return models.Buyer{}, err
	}
	buyer, ok := buyersData[id]
	if !ok {
		return models.Buyer{},
			httperrors.NotFoundError{Message: fmt.Sprintf("Buyer %d, not found", id)}
	}
	return buyer, nil
}
