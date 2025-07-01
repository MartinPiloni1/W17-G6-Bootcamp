package service

import (
	"slices"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type ProductServiceDefault struct {
	rp repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceDefault{rp: repo}
}

func (p *ProductServiceDefault) Create(product models.ProductAttributes) (models.Product, error) {
	return p.rp.Create(product)
}

func (p *ProductServiceDefault) GetAll() ([]models.Product, error) {
	data, err := p.rp.GetAll()
	if err != nil {
		return []models.Product{}, err
	}

	slicedData := utils.MapToSlice(data)
	slices.SortFunc(slicedData, func(a, b models.Product) int {
		return a.ID - b.ID
	})
	return slicedData, nil
}

func (p *ProductServiceDefault) GetByID(id int) (models.Product, error) {
	return p.rp.GetByID(id)
}

func (p *ProductServiceDefault) Update(id int, productAttributes models.ProductAttributes) (models.Product, error) {
	return p.rp.Update(id, productAttributes)
}

func (p *ProductServiceDefault) Delete(id int) error {
	return p.rp.Delete(id)
}
