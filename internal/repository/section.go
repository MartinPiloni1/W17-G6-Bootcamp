package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)


type SectionRepositoryImpl struct {
	filePath string
}

func NewSectionRepository() SectionRepositoryInterface {
	return &SectionRepositoryImpl{filePath: os.Getenv("FILE_PATH_DEFAULT")}
}
// Create implements SectionRepositoryInterface.
func (s *SectionRepositoryImpl) Create(section models.Section) (*models.Section, error) {
	panic("unimplemented")
}

// Delete implements SectionRepositoryInterface.
func (s *SectionRepositoryImpl) Delete(id int) error {
	panic("unimplemented")
}

// GetAll implements SectionRepositoryInterface.
func (s *SectionRepositoryImpl) GetAll() (map[int]models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetByID implements SectionRepository.
func (s *SectionRepositoryImpl) GetByID(id int) (models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return models.Section{}, err
	}

	section, ok := data[id]
	if !ok {
		return models.Section{}, httperrors.NotFoundError{Message: "sección no encontrada"}
	} 

	return section, nil
}

// Update implements SectionRepositoryInterface.
func (s *SectionRepositoryImpl) Update(id int, data models.Section) (models.Section, error) {
	panic("unimplemented")
}