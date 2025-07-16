package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type LocalityServiceDefault struct {
	repository repository.LocalityRepository
}

func NewLocalityService(repo repository.LocalityRepository) LocalityService {
	return &LocalityServiceDefault{repository: repo}
}

// This function creates a new locality
// It checks if all required fields are provided and returns an error if any are missing
func (s *LocalityServiceDefault) Create(locality models.Locality) (models.Locality, error) {
	if locality.ID == "" || locality.LocalityName == "" || locality.ProvinceName == "" || locality.CountryName == "" {
		return models.Locality{}, httperrors.BadRequestError{
			Message: "Invalid locality data",
		}
	}
	return s.repository.Create(locality)
}

// This function retrieves a locality by its ID
func (s *LocalityServiceDefault) GetByID(id string) (models.Locality, error) {
	return s.repository.GetByID(id)
}

// This function retrieves the seller report by locality ID
func (s *LocalityServiceDefault) GetSellerReport(id *string) ([]models.SellerReport, error) {
	return s.repository.GetSellerReport(id)
}

// GetReportByLocalityId retrieves a report of carries by locality ID.
func (s *LocalityServiceDefault) GetReportByLocalityId(localityId string) ([]models.CarryReport, error) {
	result, err := s.repository.GetReportByLocalityId(localityId)
	if result == nil && localityId != "" {
		return nil, httperrors.NotFoundError{Message: "locality not found"}
	}
	return result, err
}
