package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type CarryServiceDefault struct {
	repo repository.CarryRepository
}

// NewCarryService creates a new CarryServiceDefault instance.
func NewCarryService(repo repository.CarryRepository) *CarryServiceDefault {
	return &CarryServiceDefault{repo: repo}
}

// Create validates and creates a new carry.
func (c *CarryServiceDefault) Create(carryAttributes models.CarryAttributes) (models.Carry, error) {
	if carryAttributes.Cid == "" {
		return models.Carry{}, httperrors.BadRequestError{Message: "the field Cid must not be empty"}
	}
	if carryAttributes.CompanyName == "" {
		return models.Carry{}, httperrors.BadRequestError{Message: "the field CompanyName must not be empty"}
	}
	if carryAttributes.Address == "" {
		return models.Carry{}, httperrors.BadRequestError{Message: "the field Address must not be empty"}
	}
	if carryAttributes.Telephone == "" {
		return models.Carry{}, httperrors.BadRequestError{Message: "the field Telephone must not be empty"}
	}
	if carryAttributes.LocalityId == "" {
		return models.Carry{}, httperrors.BadRequestError{Message: "the field LocalityId must not be empty"}
	}
	return c.repo.Create(carryAttributes)
}
