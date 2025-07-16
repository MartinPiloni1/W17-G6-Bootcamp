package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
)

// ProductRecordServiceDefault is the implementation of ProductRecordService interface,
// delegating persistence to a ProductRecordRepository.
type ProductRecordServiceDefault struct {
	repo repository.ProductRecordRepository
}

// NewProductRecordServiceDefault constructs a ProductRecordServiceDefault
// with the given repository.
func NewProductRecordServiceDefault(repo repository.ProductRecordRepository) ProductRecordService {
	return &ProductRecordServiceDefault{
		repo: repo,
	}
}

// Create creates a new product record entry in the repository using the supplied
// attributes. It returns the fully populated ProductRecord or an error
func (p *ProductRecordServiceDefault) Create(ctx context.Context, attributes models.ProductRecordAttributes) (models.ProductRecord, error) {
	return p.repo.Create(ctx, attributes)
}
