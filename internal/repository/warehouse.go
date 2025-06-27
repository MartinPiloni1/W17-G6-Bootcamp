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
