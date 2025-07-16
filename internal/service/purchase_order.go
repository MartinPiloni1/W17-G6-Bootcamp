package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
)

type PurchaseOrderDefault struct {
	repository repository.PurchaseOrderRepository
}

func NewPurchaseOrderDefault(repositoryInstance repository.PurchaseOrderRepository) PurchaseOrderService {
	return &PurchaseOrderDefault{repository: repositoryInstance}
}

func (s *PurchaseOrderDefault) Create(
	ctx context.Context,
	newPurchaseOrder models.PurchaseOrderAttributes) (models.PurchaseOrder, error) {
	purchaseOrder, err := s.repository.Create(ctx, newPurchaseOrder)
	return purchaseOrder, err
}
