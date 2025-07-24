package service

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
)

// SectionServiceDefault implements SectionService
type SectionServiceDefault struct {
	repository repository.SectionRepository
}

/*
NewSectionServiceDefault constructs a SectionServiceDefault
with the given repository.
*/
func NewSectionServiceDefault(repo repository.SectionRepository) SectionService {
	return &SectionServiceDefault{repository: repo}
}

// Create, creates a new section in the repository
func (service SectionServiceDefault) Create(ctx context.Context, section models.Section) (models.Section, error) {
	return service.repository.Create(ctx, section)
}

// Delete, deletes a section from the repository
func (service SectionServiceDefault) Delete(ctx context.Context, id int) error {
	return service.repository.Delete(ctx, id)
}

// GetAll returns all sections in the repository
func (service SectionServiceDefault) GetAll(ctx context.Context) ([]models.Section, error) {
	data, err := service.repository.GetAll(ctx)
	if err != nil {
		return []models.Section{}, err
	}
	return data, nil
}

// GetByID returns a section by its ID
func (service SectionServiceDefault) GetByID(ctx context.Context, id int) (models.Section, error) {
	return service.repository.GetByID(ctx, id)
}

// Update, updates a section in the repository
func (service *SectionServiceDefault) Update(ctx context.Context, id int, patchData models.UpdateSectionRequest) (models.Section, error) {
	section, err := service.repository.GetByID(ctx, id)
	if err != nil {
		return models.Section{}, err
	}

	service.applyChanges(&section, patchData)

	return service.repository.Update(ctx, id, section)
}

// applyChanges, applies the changes to the section
func (service *SectionServiceDefault) applyChanges(sectionToUpdate *models.Section, patchData models.UpdateSectionRequest) {
	// Update section number
	if patchData.SectionNumber != nil {
		sectionToUpdate.SectionNumber = *patchData.SectionNumber
	}
	// Update temperature
	if patchData.CurrentTemperature != nil {
		sectionToUpdate.CurrentTemperature = *patchData.CurrentTemperature
	}
	// Update minimum temperature
	if patchData.MinimumTemperature != nil {
		sectionToUpdate.MinimumTemperature = *patchData.MinimumTemperature
	}
	// Update current capacity
	if patchData.CurrentCapacity != nil {
		sectionToUpdate.CurrentCapacity = *patchData.CurrentCapacity
	}
	// Update minimum capacity
	if patchData.MinimumCapacity != nil {
		sectionToUpdate.MinimumCapacity = *patchData.MinimumCapacity
	}
	// Update maximum capacity
	if patchData.MaximumCapacity != nil {
		sectionToUpdate.MaximumCapacity = *patchData.MaximumCapacity
	}
	// Update warehouse ID
	if patchData.WarehouseID != nil {
		sectionToUpdate.WarehouseID = *patchData.WarehouseID
	}
	// Update product type ID
	if patchData.ProductTypeID != nil {
		sectionToUpdate.ProductTypeID = *patchData.ProductTypeID
	}

	return
}

// GetProductsReport calls the repository to get a single section report.
func (service SectionServiceDefault) GetProductsReport(ctx context.Context, id int) (models.SectionProductsReport, error) {
	return service.repository.GetProductsReport(ctx, id)
}

// GetAllProductsReport calls the repository to get all section reports.
func (service SectionServiceDefault) GetAllProductsReport(ctx context.Context) ([]models.SectionProductsReport, error) {
	return service.repository.GetAllProductsReport(ctx)
}
