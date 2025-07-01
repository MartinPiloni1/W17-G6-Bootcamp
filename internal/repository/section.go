package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type SectionRepositoryFile struct {
	filePath string
}

func NewSectionRepository() SectionRepositoryInterface {
	return &SectionRepositoryFile{filePath: os.Getenv("FILE_PATH_SECTIONS")}
}

// Create implements SectionRepositoryInterface.
func (s *SectionRepositoryFile) Create(section models.Section) (*models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return nil, err
	}

	id := len(data) + 1
	section.ID = id
	data[id] = section

	return &section, utils.Write(s.filePath, data)
}

// Delete implements SectionRepositoryInterface.
func (s SectionRepositoryFile) Delete(id int) error {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return err
	}

	_, ok := data[id]
	if !ok {
		return httperrors.NotFoundError{Message: "sección no encontrada"}
	}

	delete(data, id)

	return utils.Write(s.filePath, data)
}

// GetAll implements SectionRepositoryInterface.
func (s SectionRepositoryFile) GetAll() (map[int]models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetByID implements SectionRepository.
func (s SectionRepositoryFile) GetByID(id int) (models.Section, error) {
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
func (s SectionRepositoryFile) Update(id int, data models.Section) (models.Section, error) {
	dataMap, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return models.Section{}, err
	}

	_, ok := dataMap[id]
	if !ok {
		return models.Section{}, httperrors.NotFoundError{Message: "sección no encontrada"}
	}

	dataMap[id] = data

	return dataMap[id], utils.Write(s.filePath, dataMap)
}
