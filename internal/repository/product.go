package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type ProductRepositoryFile struct {
	filePath string
}

func NewProductRepositoryFile() ProductRepository {
	filePath := os.Getenv("FILE_PATH_PRODUCTS")
	return &ProductRepositoryFile{
		filePath: filePath,
	}
}

func (p *ProductRepositoryFile) Create(productAttribbutes models.ProductAttributes) (models.Product, error) {
	productData, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	newId, err := utils.GetNextID[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	newProduct := models.Product{
		ID:                newId,
		ProductAttributes: productAttribbutes,
	}

	productData[newId] = newProduct

	err = utils.Write(p.filePath, productData)
	if err != nil {
		return models.Product{}, err
	}

	return newProduct, nil
}

func (p *ProductRepositoryFile) GetAll() (map[int]models.Product, error) {
	productData, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return nil, err
	}
	return productData, nil
}

func (p *ProductRepositoryFile) GetByID(id int) (models.Product, error) {
	productData, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	product, exists := productData[id]
	if !exists {
		return models.Product{},
			httperrors.NotFoundError{Message: "Product not found"}
	}
	return product, nil
}

func (p *ProductRepositoryFile) Update(id int, product models.Product) (models.Product, error) {
	productData, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}
	productData[id] = product
	err = utils.Write(p.filePath, productData)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (p *ProductRepositoryFile) Delete(id int) error {
	productData, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return err
	}

	if _, exists := productData[id]; !exists {
		return httperrors.NotFoundError{Message: "Product not found"}
	}

	delete(productData, id)

	err = utils.Write(p.filePath, productData)
	if err != nil {
		return err
	}

	return nil
}
