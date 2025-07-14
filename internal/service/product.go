package service

import (
	"slices"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type ProductServiceDefault struct {
	repository repository.ProductRepository
}

func NewProductServiceDefault(repository repository.ProductRepository) ProductService {
	return &ProductServiceDefault{repository: repository}
}

func (service *ProductServiceDefault) Create(product models.ProductAttributes) (models.Product, error) {
	products, err := service.repository.GetAll()
	if err != nil {
		return models.Product{}, err
	}

	for _, p := range products {
		if p.ProductCode == product.ProductCode {
			return models.Product{},
				httperrors.ConflictError{Message: "A product with this product code already exists"}
		}
	}

	return service.repository.Create(product)
}

func (service *ProductServiceDefault) GetAll() ([]models.Product, error) {
	data, err := service.repository.GetAll()
	if err != nil {
		return []models.Product{}, err
	}

	slices.SortFunc(data, func(a, b models.Product) int {
		return a.ID - b.ID
	})
	return data, nil
}

func (service *ProductServiceDefault) GetByID(id int) (models.Product, error) {
	return service.repository.GetByID(id)
}

func (service *ProductServiceDefault) Update(id int, productAttributes models.ProductPatchRequest) (models.Product, error) {
	product, err := service.repository.GetByID(id)
	if err != nil {
		return models.Product{}, err
	}

	if productAttributes.ProductCode != nil {
		products, err := service.repository.GetAll()
		if err != nil {
			return models.Product{}, err
		}

		product.ProductCode = *productAttributes.ProductCode
		for _, p := range products {
			if p.ProductCode == product.ProductCode && p.ID != product.ID {
				return models.Product{},
					httperrors.ConflictError{Message: "A product with this product code already exists"}
			}
		}
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

	return service.repository.Update(id, product)
}

func (service *ProductServiceDefault) Delete(id int) error {
	return service.repository.Delete(id)
}
