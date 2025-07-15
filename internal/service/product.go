package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
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

// Create creates a new product entry in the repository using the supplied
// attributes. It returns the fully populated Product model or an error
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

// GetRecordsPerProduct returns a list of ProductRecordCount entries.
// If id is nil, it retrieves counts for all products; otherwise it returns
// the record count for the single product identified by id.
func (s *ProductServiceDefault) GetRecordsPerProduct(ctx context.Context, id *int) ([]models.ProductRecordCount, error) {
	return s.repo.GetRecordsPerProduct(ctx, id)
}

// Update retrieves an existing Product by its ID, applies any non-nil fields
// from the provided ProductPatchRequest, and then persists the updated Product
// via the repository. Returns the updated Product or an error
func (s *ProductServiceDefault) Update(ctx context.Context, id int, productAttributes models.ProductPatchRequest) (models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Product{}, err
	}

	// Apply non-nil fields to the product
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
