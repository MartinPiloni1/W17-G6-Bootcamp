package repository

import (
	"errors"
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type SellerRepositoryImpl struct {
	filePath string
}

func NewSellerRepository() SellerRepositoryInterface {
	return &SellerRepositoryImpl{filePath: os.Getenv("FILE_PATH_DEFAULT")}
}

func (r *SellerRepositoryImpl) GetByID(id int) (models.Seller, error) {
	data, err := utils.Read[models.Seller](r.filePath)
	if err != nil {
		return models.Seller{}, err
	}
	seller, ok := data[id]
	if !ok {
		return models.Seller{}, errors.New("seller not found")
	}
	return seller, nil
}

func (s *SellerRepositoryImpl) GetAll() (map[int]models.Seller, error) {
	data, err := utils.Read[models.Seller](s.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *SellerRepositoryImpl) Create(seller models.Seller) (*models.Seller, error) {
	data, err := utils.Read[models.Seller](r.filePath)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		if v.CID == seller.CID {
			return nil, errors.New("cid already exists")
		}
	}
	maxID := 0
	for id := range data {
		if id > maxID {
			maxID = id
		}
	}
	seller.ID = maxID + 1
	data[seller.ID] = seller
	if err := utils.Write(r.filePath, data); err != nil {
		return nil, err
	}
	return &seller, nil
}

func (r *SellerRepositoryImpl) Delete(id int) error {
	data, err := utils.Read[models.Seller](r.filePath) // map[int]models.Seller
	if err != nil {
		return err
	}
	if _, exists := data[id]; !exists {
		return errors.New("seller not found")
	}
	delete(data, id)
	if err := utils.Write(r.filePath, data); err != nil {
		return err
	}
	return nil
}

func (r *SellerRepositoryImpl) Update(id int, data *models.Seller) (*models.Seller, error) {
	sellers, err := utils.Read[models.Seller](r.filePath)
	if err != nil {
		return nil, err
	}
	seller, exists := sellers[id]
	if !exists {
		return nil, errors.New("seller not found")
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
		return nil, err
	}
	return &seller, nil
}
