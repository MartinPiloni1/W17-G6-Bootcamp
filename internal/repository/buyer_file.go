package repository

import (
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
