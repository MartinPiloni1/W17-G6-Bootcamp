package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type SellerRepositoryFile struct {
	filePath string
}

func NewSellerRepository() SellerRepository {
	return &SellerRepositoryFile{filePath: os.Getenv("FILE_PATH_SELLER")}
}

func (r SellerRepositoryFile) GetByID(id int) (models.Seller, error) {
	data, err := utils.Read[models.Seller](r.filePath)
	if err != nil {
		return models.Seller{}, err
	}
	seller, ok := data[id]
	if !ok {
		return models.Seller{}, httperrors.NotFoundError{Message: "Seller not found"}
	}
	return seller, nil
}

func (s SellerRepositoryFile) GetAll() (map[int]models.Seller, error) {
	data, err := utils.Read[models.Seller](s.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r SellerRepositoryFile) Create(seller models.SellerAttributes) (models.Seller, error) {
	data, err := utils.Read[models.Seller](r.filePath)
	if err != nil {
		return models.Seller{}, err
	}

	sellerId, err := utils.GetNextID[models.Seller](r.filePath)
	if err != nil {
		return models.Seller{}, err
	}
	for _, v := range data {
		if v.CID == seller.CID {
			return models.Seller{}, httperrors.ConflictError{Message: "CID already exists"}
		}
	}

	newSeller := models.Seller{
		ID:               sellerId,
		SellerAttributes: seller,
	}
	data[sellerId] = newSeller
	if err := utils.Write(r.filePath, data); err != nil {
		return models.Seller{}, err
	}
	return newSeller, nil
}

func (r SellerRepositoryFile) Delete(id int) error {
	data, err := utils.Read[models.Seller](r.filePath) // map[int]models.Seller
	if err != nil {
		return err
	}
	if _, exists := data[id]; !exists {
		return httperrors.NotFoundError{Message: "Seller not found"}
	}
	delete(data, id)
	if err := utils.Write(r.filePath, data); err != nil {
		return err
	}
	return nil
}

func (r SellerRepositoryFile) Update(id int, data *models.SellerAttributes) (models.Seller, error) {
	sellers, err := utils.Read[models.Seller](r.filePath)
	if err != nil {
		return models.Seller{}, err
	}
	seller, exists := sellers[id]
	if !exists {
		return models.Seller{}, httperrors.NotFoundError{Message: "Seller not found"}
	}

	//actualizar solo los campos no vacios
	if data.CID != 0 {
		seller.CID = data.CID
	}
	if data.CompanyName != "" {
		seller.CompanyName = data.CompanyName
	}
	if data.Address != "" {
		seller.Address = data.Address
	}
	if data.Telephone != "" {
		seller.Telephone = data.Telephone
	}

	sellers[id] = seller
	if err := utils.Write(r.filePath, sellers); err != nil {
		return models.Seller{}, err
	}
	return seller, nil
}
