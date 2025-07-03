package repository

import (
	"os"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
)

type SectionRepositoryFile struct {
	filePath string
}

func NewSectionRepository() SectionRepository {
	return &SectionRepositoryFile{filePath: os.Getenv("FILE_PATH_SECTIONS")}
}

func (s *SectionRepositoryFile) Create(section models.Section) (models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return models.Section{}, err
	}

	newId, err := utils.GetNextID[models.Section](s.filePath)
	if err != nil {
		return models.Section{}, err
	}

	section.ID = newId
	data[newId] = section

	return section, utils.Write(s.filePath, data)
}

func (s SectionRepositoryFile) Delete(id int) error {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return err
	}

	_, ok := data[id]
	if !ok {
		return httperrors.NotFoundError{Message: "Section not found"}
	}

	delete(data, id)

	return utils.Write(s.filePath, data)
}

func (s SectionRepositoryFile) GetAll() (map[int]models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s SectionRepositoryFile) GetByID(id int) (models.Section, error) {
	data, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return models.Section{}, err
	}

	section, ok := data[id]
	if !ok {
		return models.Section{}, httperrors.NotFoundError{Message: "Section not found"}
	}

	return section, nil
}

func (s SectionRepositoryFile) Update(id int, data models.Section) (models.Section, error) {
	dataMap, err := utils.Read[models.Section](s.filePath)
	if err != nil {
		return models.Section{}, err
	}

	_, ok := dataMap[id]
	if !ok {
		return models.Section{}, httperrors.NotFoundError{Message: "Section not found"}
	}

	dataMap[id] = data

	return dataMap[id], utils.Write(s.filePath, dataMap)
}
