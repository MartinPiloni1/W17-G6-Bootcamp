package service

import (
	"slices"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type ProductService struct {
	repo repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{repo: repo}
}

func (p *ProductService) Create(product models.ProductAttributes) (models.Product, error) {
	return p.repo.Create(product)
}

func (p *ProductService) GetAll() ([]models.Product, error) {
	data, err := p.repo.GetAll()
	if err != nil {
		return []models.Product{}, err
	}

	slicedData := utils.MapToSlice(data)
	slices.SortFunc(slicedData, func(a, b models.Product) int {
		return a.ID - b.ID
	})
	return slicedData, nil
}

func (p *ProductService) GetByID(id int) (models.Product, error) {
	return p.repo.GetByID(id)
}

func (p *ProductService) Update(id int, productAttributes models.ProductAttributes) (models.Product, error) {
	return p.repo.Update(id, productAttributes)
}

func (p *ProductService) Delete(id int) error {
	return p.repo.Delete(id)
}
