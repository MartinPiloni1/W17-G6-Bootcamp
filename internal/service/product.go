package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

// ProductServiceDefault is the implementation of ProductService interface,
// delegating persistence to a ProductRepository.
type ProductServiceDefault struct {
	repo repository.ProductRepository
}

// NewProductServiceDefault constructs a ProductServiceDefault
// with the given repository.
func NewProductServiceDefault(repo repository.ProductRepository) ProductService {
	return &ProductServiceDefault{repo: repo}
}

// Create validates that no existing product shares the same ProductCode,
// then delegates to the repository to persist a new product.
func (s *ProductServiceDefault) Create(ctx context.Context, product models.ProductAttributes) (models.Product, error) {
	return s.repo.Create(ctx, product)
}

// GetAll retrieves all products from the repository, returning a slice
// containing every product or an error.
func (s *ProductServiceDefault) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

// GetByID fetches a single product by its ID.
// Returns the product or ErrNotFound if no such product exists.
func (s *ProductServiceDefault) GetByID(ctx context.Context, id int) (models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// Update retrieves an existing Product by its ID, applies any non-nil fields
// from the provided ProductPatchRequest, and then persists the updated Product
// via the repository. Returns the updated Product or an error
func (s *ProductServiceDefault) Update(ctx context.Context, id int, productAttributes models.ProductPatchRequest) (models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Product{}, err
	}

	if productAttributes.ProductCode != nil {
		products, err := s.repo.GetAll(ctx)
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
		product.SellerID = productAttributes.SellerID
	}

	return s.repo.Update(ctx, id, product)
}

// Delete removes the Product with the given ID from the database.
// Returns nil on success or an error if the repository fails.
func (s *ProductServiceDefault) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
