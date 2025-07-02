package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type BuyerRepositoryFile struct {
	filePath string
}

func NewBuyerRepositoryFile() BuyerRepository {
	return &BuyerRepositoryFile{filePath: os.Getenv("FILE_PATH_BUYER")}
}

func (r *BuyerRepositoryFile) Create(newBuyer models.BuyerAttributes) (models.Buyer, error) {
	nextId, err := utils.GetNextID[models.Buyer](r.filePath)
	if err != nil {
		return models.Buyer{}, err
	}

	buyer := models.Buyer{
		Id:              nextId,
		BuyerAttributes: newBuyer,
	}

	// fetch complete file and add the new buyer
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return models.Buyer{}, err
	}
	buyersData[nextId] = buyer

	err = utils.Write(r.filePath, buyersData)
	if err != nil {
		return models.Buyer{}, err
	}

	return buyer, nil
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
			httperrors.NotFoundError{Message: "Buyer not found"}
	}
	return buyer, nil
}

func (r *BuyerRepositoryFile) Update(id int, updatedBuyer models.Buyer) (models.Buyer, error) {
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return models.Buyer{}, err
	}

	buyersData[id] = updatedBuyer

	err = utils.Write(r.filePath, buyersData)
	if err != nil {
		return models.Buyer{}, err
	}

	return updatedBuyer, nil
}

func (r *BuyerRepositoryFile) Delete(id int) error {
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return err
	}

	_, ok := buyersData[id]
	if !ok {
		return httperrors.NotFoundError{Message: "Buyer not found"}
	}

	delete(buyersData, id)

	err = utils.Write(r.filePath, buyersData)
	if err != nil {
		return err
	}

	return nil
}

func (r *BuyerRepositoryFile) CardNumberIdAlreadyExist(newCardNumberId int) (bool, error) {
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return false, err
	}

	for _, buyer := range buyersData {
		if buyer.CardNumberId == newCardNumberId {
			return true, nil
		}
	}
	return false, nil

}
