package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type CarryServiceDefault struct {
	rp repository.CarriesRepository
}

// NewCarryService creates a new CarryServiceDefault instance.
func NewCarryService(repo repository.CarriesRepository) *CarryServiceDefault {
	return &CarryServiceDefault{rp: repo}
}

// Create validates and creates a new carry.
func (p *CarryServiceDefault) Create(carryAttributes models.CarryAttributes) (models.Carry, error) {
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
	return p.rp.Create(carryAttributes)
}

// GetReportByLocalityId retrieves a report of carries by locality ID.
func (p *CarryServiceDefault) GetReportByLocalityId(localityId string) ([]models.CarryReport, error) {
	result, err := p.rp.GetReportByLocalityId(localityId)
	if result == nil {
		return nil, httperrors.NotFoundError{Message: "no carries found for the given locality ID"}
	}
	return result, err
}
