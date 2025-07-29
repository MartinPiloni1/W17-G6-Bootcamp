package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
)

// BuyerServiceDefault contains a repository for buyer
type BuyerServiceDefault struct {
	repository repository.BuyerRepository
}

// NewBuyerServiceDefault returns an instance of BuyerServiceDefault
func NewBuyerServiceDefault(repositoryInstance repository.BuyerRepository) BuyerService {
	return &BuyerServiceDefault{repository: repositoryInstance}
}

// creates a buyers if it has an unique CardNumberId that dosent already exist in the db
func (s *BuyerServiceDefault) Create(ctx context.Context, newBuyer models.BuyerAttributes) (models.Buyer, error) {
	buyer, err := s.repository.Create(ctx, newBuyer)
	return buyer, err
}

// Get all buyers
func (s *BuyerServiceDefault) GetAll(ctx context.Context) ([]models.Buyer, error) {
	buyerData, err := s.repository.GetAll(ctx)
	return buyerData, err
}

// Get one user based on his id
func (s *BuyerServiceDefault) GetByID(ctx context.Context, id int) (models.Buyer, error) {
	buyer, err := s.repository.GetByID(ctx, id)
	return buyer, err
}

// Update updates the Buyer identified by the given ID with the provided patch data.
// It retrieves the existing Buyer from the repository, applies any non-nil fields
// from BuyerData, and checks for CardNumberId conflicts before saving.
//
// Returns the updated Buyer struct if successful. If the Buyer is not found,
// a not found error is returned. If the provided CardNumberId is already in use
// by another Buyer, a conflict error is returned. Any other database or internal
// errors are also returned.
func (s *BuyerServiceDefault) Update(ctx context.Context, id int, BuyerData models.BuyerPatchRequest) (models.Buyer, error) {
	buyer, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return models.Buyer{}, err
	}
	if BuyerData.CardNumberId != nil {
		buyer.CardNumberId = *BuyerData.CardNumberId
	}

	if BuyerData.FirstName != nil {
		buyer.FirstName = *BuyerData.FirstName
	}
	if BuyerData.LastName != nil {
		buyer.LastName = *BuyerData.LastName
	}

	updatedBuyer, err := s.repository.Update(ctx, id, buyer)
	return updatedBuyer, err
}

// Delete Buyer based on his id
func (s *BuyerServiceDefault) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	return err
}

// if id is given validate if the buyer exist by id, if exist, obtains the BuyerWithPurchaseOrdersCount
// no id, GetAll with BuyerWithPurchaseOrdersCount
func (s *BuyerServiceDefault) GetWithPurchaseOrdersCount(
	ctx context.Context, id *int) ([]models.BuyerWithPurchaseOrdersCount, error) {
	if id != nil {
		_, err := s.repository.GetByID(ctx, *id)
		if err != nil {
			return nil, err
		}
	}
	return s.repository.GetWithPurchaseOrdersCount(ctx, id)
}
