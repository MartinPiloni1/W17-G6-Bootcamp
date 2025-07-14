package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

/*
ProductServiceDefault is the implementation of ProductService interface,
delegating persistence to a ProductRepository.
*/
type ProductServiceDefault struct {
	repo repository.ProductRepository
}

/*
NewProductServiceDefault constructs a ProductServiceDefault
with the given repository.
*/
func NewProductServiceDefault(repo repository.ProductRepository) ProductService {
	return &ProductServiceDefault{repo: repo}
}

/*
Create validates that no existing product shares the same ProductCode,
then delegates to the repository to persist a new product.
*/
func (svc *ProductServiceDefault) Create(product models.ProductAttributes) (models.Product, error) {
	products, err := svc.repo.GetAll()
	if err != nil {
		return models.Product{}, err
	}

	for _, p := range products {
		if p.ProductCode == product.ProductCode {
			return models.Product{},
				httperrors.ConflictError{Message: "A product with this product code already exists"}
		}
	}

	return svc.repo.Create(product)
}

/*
GetAll retrieves all products from the repository, returning a slice
containing every product or an error.
*/
func (svc *ProductServiceDefault) GetAll() ([]models.Product, error) {
	return svc.repo.GetAll()
}

/*
GetByID fetches a single product by its ID.
Returns the product or ErrNotFound if no such product exists.
*/
func (svc *ProductServiceDefault) GetByID(id int) (models.Product, error) {
	return svc.repo.GetByID(id)
}

/*
Update retrieves an existing Product by its ID, applies any non-nil fields
from the provided ProductPatchRequest, enforces uniqueness of ProductCode,
and then persists the updated Product via the repository.
Returns the updated Product or an error
*/
func (svc *ProductServiceDefault) Update(id int, productAttributes models.ProductPatchRequest) (models.Product, error) {
	product, err := svc.repo.GetByID(id)
	if err != nil {
		return models.Product{}, err
	}

	if productAttributes.ProductCode != nil {
		products, err := svc.repo.GetAll()
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

	return svc.repo.Update(id, product)
}

/*
Delete removes the Product with the given ID from the datastore.
Returns nil on success or an error if the repository fails.
*/
func (svc *ProductServiceDefault) Delete(id int) error {
	return svc.repo.Delete(id)
}
