package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type ProductRepository struct {
	filePath string
}

func NewProductRepository() ProductRepositoryInterface {
	filePath := os.Getenv("FILE_PATH_DEFAULT")
	return &ProductRepository{
		filePath: filePath,
	}
}

func (p *ProductRepository) Create(productAtribbutes models.ProductAtributtes) (models.Product, error) {
	data, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	newId, err := utils.GetNextID[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	newProduct := models.Product{
		ID:                newId,
		ProductAtributtes: productAtribbutes,
	}

	data[newId] = newProduct

	err = utils.Write(p.filePath, data)
	if err != nil {
		return models.Product{}, err
	}

	return newProduct, nil
}

func (p *ProductRepository) GetAll() (map[int]models.Product, error) {
	data, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *ProductRepository) GetByID(id int) (models.Product, error) {
	data, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	product, exists := data[id]
	if !exists {
		return models.Product{},
			httperrors.NotFoundError{Message: "No se encontró un producto con el ID proporcionado"}
	}
	return product, nil
}

func (p *ProductRepository) Update(id int, productAtributtes models.ProductAtributtes) (models.Product, error) {
	data, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return models.Product{}, err
	}

	if _, exists := data[id]; !exists {
		return models.Product{},
			httperrors.NotFoundError{Message: "No se encontró un producto con el ID proporcionado"}
	}

	updatedProduct := models.Product{
		ID:                id,
		ProductAtributtes: productAtributtes,
	}
	data[id] = updatedProduct

	utils.Write(p.filePath, data)
	return updatedProduct, nil
}

func (p *ProductRepository) Delete(id int) error {
	data, err := utils.Read[models.Product](p.filePath)
	if err != nil {
		return err
	}

	if _, exists := data[id]; !exists {
		return httperrors.NotFoundError{Message: "No se encontró un producto con el ID proporcionado"}
	}

	delete(data, id)
	utils.Write(p.filePath, data)

	return nil
}
