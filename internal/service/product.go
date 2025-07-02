package service

import (
	"slices"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type ProductServiceDefault struct {
	rp repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceDefault{rp: repo}
}

func (p *ProductServiceDefault) Create(product models.ProductAttributes) (models.Product, error) {
	products, err := p.rp.GetAll()
	if err != nil {
		return models.Product{}, err
	}

	for _, p := range products {
		if p.ProductCode == product.ProductCode {
			return models.Product{}, httperrors.ConflictError{Message: "A product with this product code already exists"}
		}
	}

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

func (p *ProductServiceDefault) Update(id int, productAttributes models.ProductPatchRequest) (models.Product, error) {
	product, err := p.rp.GetByID(id)
	if err != nil {
		return models.Product{}, err
	}

	if productAttributes.ProductCode != nil {
		products, err := p.rp.GetAll()
		if err != nil {
			return models.Product{}, err
		}

		for _, p := range products {
			if p.ProductCode == product.ProductCode {
				return models.Product{}, httperrors.ConflictError{Message: "A product with this product code already exists"}
			}
		}
		product.ProductCode = *productAttributes.ProductCode
	}

	if productAttributes.Description != nil {
		product.Description = *productAttributes.Description
	}
	if productAttributes.ExpirationRate != nil {
		product.ExpirationRate = *productAttributes.ExpirationRate
	}
	if productAttributes.FreezingRate != nil {
		product.FreezingRate = *productAttributes.FreezingRate
	}
	if productAttributes.Height != nil {
		product.Height = *productAttributes.Height
	}
	if productAttributes.Length != nil {
		product.Length = *productAttributes.Length
	}
	if productAttributes.Width != nil {
		product.Width = *productAttributes.Width
	}
	if productAttributes.NetWeight != nil {
		product.NetWeight = *productAttributes.NetWeight
	}
	if productAttributes.RecommendedFreezingTemperature != nil {
		product.RecommendedFreezingTemperature = *productAttributes.RecommendedFreezingTemperature
	}
	if productAttributes.ProductTypeID != nil {
		product.ProductTypeID = *productAttributes.ProductTypeID
	}
	if productAttributes.SellerID != nil {
		product.SellerID = *productAttributes.SellerID
	}

	return p.rp.Update(id, product)
}

func (p *ProductServiceDefault) Delete(id int) error {
	return p.rp.Delete(id)
}
