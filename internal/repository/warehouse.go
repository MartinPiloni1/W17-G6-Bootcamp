package repository

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"os"
)

type WarehouseRepositoryFile struct {
	filePath string
}

func NewWarehouseRepository() *WarehouseRepositoryFile {
	return &WarehouseRepositoryFile{filePath: os.Getenv("FILE_PATH_DEFAULT")}
}

func (p *WarehouseRepositoryFile) Create(warehouseAtribbutes models.WarehouseAttributes) (models.Warehouse, error) {
	warehouseData, err := utils.Read[models.Warehouse](p.filePath)
	if err != nil {
		return models.Warehouse{}, err
	}

	newId, err := utils.GetNextID[models.Warehouse](p.filePath)
	if err != nil {
		return models.Warehouse{}, err
	}

	newWarehouse := models.Warehouse{
		Id:                  newId,
		WarehouseAttributes: warehouseAtribbutes,
	}

	warehouseData[newId] = newWarehouse

	err = utils.Write(p.filePath, warehouseData)
	if err != nil {
		return models.Warehouse{}, err
	}

	return newWarehouse, nil
}

func (p *WarehouseRepositoryFile) GetAll() (map[int]models.Warehouse, error) {
	data, err := utils.Read[models.Warehouse](p.filePath)
	if err != nil {
		return map[int]models.Warehouse{}, err
	}
	if len(data) == 0 {
		return map[int]models.Warehouse{}, httperrors.NotFoundError{Message: "data not find"}
	}
	return data, nil
}

func (p *WarehouseRepositoryFile) GetByID(id int) (models.Warehouse, error) {
	warehouseData, err := utils.Read[models.Warehouse](p.filePath)
	if err != nil {
		return models.Warehouse{}, err
	}
	warehouse, exists := warehouseData[id]
	if !exists {
		return models.Warehouse{},
			httperrors.NotFoundError{Message: "Warehouse not found"}
	}
	return warehouse, nil
}

func (p *WarehouseRepositoryFile) Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	WarehouseData, err := utils.Read[models.Warehouse](p.filePath)
	if err != nil {
		return models.Warehouse{}, err
	}

	updatedWarehouse := models.Warehouse{
		Id:                  id,
		WarehouseAttributes: warehouseAttributes,
	}
	WarehouseData[id] = updatedWarehouse

	utils.Write(p.filePath, WarehouseData)
	return updatedWarehouse, nil
}

func (p *WarehouseRepositoryFile) Delete(id int) error {
	warehouseData, err := utils.Read[models.Warehouse](p.filePath)
	if err != nil {
		return err
	}

	if _, exists := warehouseData[id]; !exists {
		return httperrors.NotFoundError{Message: "Warehouse not found"}
	}

	delete(warehouseData, id)

	err = utils.Write(p.filePath, warehouseData)
	if err != nil {
		return err
	}

	return nil
}
