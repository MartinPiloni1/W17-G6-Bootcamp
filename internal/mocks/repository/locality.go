package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type LocalityRepositoryMock struct {
	mock.Mock
}

func (m *LocalityRepositoryMock) Create(l models.Locality) (models.Locality, error) {
	args := m.Called(l)
	return args.Get(0).(models.Locality), args.Error(1)
}
func (m *LocalityRepositoryMock) GetByID(id string) (models.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(models.Locality), args.Error(1)
}
func (m *LocalityRepositoryMock) GetSellerReport(id *string) ([]models.SellerReport, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SellerReport), args.Error(1)
}
func (m *LocalityRepositoryMock) GetReportByLocalityId(localityId string) ([]models.CarryReport, error) {
	args := m.Called(localityId)

	var res []models.CarryReport
	if args.Get(0) != nil {
		res = args.Get(0).([]models.CarryReport)
	} else {
		res = nil
	}
	return res, args.Error(1)
}
