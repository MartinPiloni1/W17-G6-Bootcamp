package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
)

// ProductBatchServiceDefault implements ProductBatchService
type ProductBatchServiceDefault struct {
	repository repository.ProductBatchRepository
}

/*
NewProductBatchServiceDefault constructs a ProductBatchServiceDefault
with the given repository.
*/
func NewProductBatchServiceDefault(repo repository.ProductBatchRepository) ProductBatchService {
	return &ProductBatchServiceDefault{repository: repo}
}

// Create creates a new product batch in the repository
func (service ProductBatchServiceDefault) Create(ctx context.Context, productBatch models.ProductBatchAttibutes) (models.ProductBatch, error) {
	return service.repository.Create(ctx, productBatch)
}