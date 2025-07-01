package service

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
)

type SectionServiceImpl struct {
	repo repository.SectionRepositoryInterface
}

func NewSectionService(repo repository.SectionRepositoryInterface) SectionServiceInterface {
	return &SectionServiceImpl{repo: repo}
}

// Create implements SectionServiceInterface.
func (s *SectionServiceImpl) Create(section models.Section) (*models.Section, error) {
	allSections, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	for _, sections := range allSections {
		if sections.SectionNumber == section.SectionNumber {
			return nil, httperrors.ConflictError{Message: "el section_number ya existe"}
		}
	}

	return s.repo.Create(section)
}

// Delete implements SectionServiceInterface.
func (s *SectionServiceImpl) Delete(id int) error {
	return s.repo.Delete(id)
}

// GetAll implements SectionServiceInterface.
func (s *SectionServiceImpl) GetAll() (map[int]models.Section, error) {
	return s.repo.GetAll()
}

// GetByID implements SectionServiceInterface.
func (s *SectionServiceImpl) GetByID(id int) (models.Section, error) {
	return s.repo.GetByID(id)
}

// Update implements SectionServiceInterface.
func (s *SectionServiceImpl) Update(id int, data models.Section) (models.Section, error) {
	panic("unimplemented")
}