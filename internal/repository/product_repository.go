package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type ProductRepository struct {
	filePath string
}

func NewProductRepository() ProductRepositoryInterface {
	return &ProductRepository{filePath: os.Getenv("FILE_PATH_DEFAULT")}
}

func (p *ProductRepository) Create(Product models.Product) (*models.Product, error) {
	panic("unimplemented")
}

func (p *ProductRepository) GetAll() (map[int]models.Product, error) {
	data, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetByID implements ProductRepositoryInterface.
func (p *ProductRepository) GetByID(id int) (models.Product, error) {
	panic("unimplemented")
}

func (p *ProductRepository) Delete(id int) error {
	panic("unimplemented")
}

// Update implements ProductRepositoryInterface.
func (p *ProductRepository) Update(id int, data models.Product) (models.Product, error) {
	panic("unimplemented")
}
