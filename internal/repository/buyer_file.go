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

func NewBuyerRepositoryFile(filePath string) BuyerRepository {
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

func (r *BuyerRepositoryFile) Delete(id int) error {
	buyersData, err := utils.Read[models.Buyer](r.filePath)
	if err != nil {
		return err
	}

	_, ok := buyersData[id]
	if !ok {
		return httperrors.NotFoundError{Message: fmt.Sprintf("Buyer %d, not found", id)}
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

func (r *BuyerRepositoryFile) Create(newBuyer models.BuyerAttributes) (models.Buyer, error) {
	// check if the cardNumber already exist
	exist, err := r.CardNumberIdAlreadyExist(newBuyer.CardNumberId)
	if err != nil {
		return models.Buyer{}, err
	}

	if exist {
		return models.Buyer{}, httperrors.ConflictError{Message: "Buyer already exist"}
	}

	nextId, err := utils.GetNextID[models.Buyer](r.filePath)
	if err != nil {
		return models.Buyer{}, err
	}

	buyer := models.Buyer{
		Id: nextId,
		BuyerAttributes: models.BuyerAttributes{
			CardNumberId: newBuyer.CardNumberId,
			FirstName:    newBuyer.FirstName,
			LastName:     newBuyer.LastName,
		},
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
