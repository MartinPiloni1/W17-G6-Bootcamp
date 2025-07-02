package service

import (
	"slices"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type SectionServiceDefault struct {
	repo             repository.SectionRepository
	warehouseService WarehouseService 
}

func NewSectionService(repo repository.SectionRepository, warehouseService WarehouseService) SectionService {
	return &SectionServiceDefault{
		repo:             repo,
		warehouseService: warehouseService,
	}
}

func (s SectionServiceDefault) Create(section models.Section) (*models.Section, error) {
	_, err := s.warehouseService.GetByID(section.WarehouseID)
	if err != nil {
		return nil, httperrors.BadRequestError{Message: "warehouse not found"}
	}

	allSections, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	for _, sections := range allSections {
		if sections.SectionNumber == section.SectionNumber {
			return nil, httperrors.ConflictError{Message: "section number already exists"}
		}
	}

	return s.repo.Create(section)
}

func (s SectionServiceDefault) Delete(id int) error {
	return s.repo.Delete(id)
}
func (s SectionServiceDefault) GetAll() ([]models.Section, error) {
	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	slicedData := utils.MapToSlice(data)
	slices.SortFunc(slicedData, func(a, b models.Section) int {
		return a.ID - b.ID
	})
	return slicedData, nil
}

func (s SectionServiceDefault) GetByID(id int) (models.Section, error) {
	return s.repo.GetByID(id)
}

func (s *SectionServiceDefault) Update(id int, patchData models.UpdatePatchSectionRequest) (models.Section, error) {
	_, err := s.warehouseService.GetByID(id)
	if err != nil {
		return models.Section{}, httperrors.BadRequestError{Message: "warehouse not found"}
	}
	sectionToUpdate, err := s.repo.GetByID(id)
	if err != nil {
		return models.Section{}, err
	}
	if err := s.applyChanges(&sectionToUpdate, patchData); err != nil {
		return models.Section{}, err
	}

	return s.repo.Update(id, sectionToUpdate)
}

func (s *SectionServiceDefault) applyChanges(sectionToUpdate *models.Section, patchData models.UpdatePatchSectionRequest) error {
	if patchData.SectionNumber != nil {
		allSections, err := s.repo.GetAll()
		if err != nil {
			return err
		}
		for _, section := range allSections {
			if section.ID != sectionToUpdate.ID && section.SectionNumber == *patchData.SectionNumber {
				return httperrors.ConflictError{Message: "section number already exists in another section"}
			}
		}
		sectionToUpdate.SectionNumber = *patchData.SectionNumber
	}

	if patchData.CurrentTemperature != nil {
		sectionToUpdate.CurrentTemperature = *patchData.CurrentTemperature
	}
	if patchData.MinimumTemperature != nil {
		sectionToUpdate.MinimumTemperature = *patchData.MinimumTemperature
	}
	if patchData.CurrentCapacity != nil {
		sectionToUpdate.CurrentCapacity = *patchData.CurrentCapacity
	}
	if patchData.MinimumCapacity != nil {
		sectionToUpdate.MinimumCapacity = *patchData.MinimumCapacity
	}
	if patchData.MaximumCapacity != nil {
		sectionToUpdate.MaximumCapacity = *patchData.MaximumCapacity
	}
	if patchData.WarehouseID != nil {
		sectionToUpdate.WarehouseID = *patchData.WarehouseID
	}
	if patchData.ProductTypeID != nil {
		sectionToUpdate.ProductTypeID = *patchData.ProductTypeID
	}

	return nil
}