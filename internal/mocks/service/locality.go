package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) Create(l models.Locality) (models.Locality, error) {
	args := m.Called(l)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetByID(id string) (models.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetSellerReport(id *string) ([]models.SellerReport, error) {
	args := m.Called(id)
	return args.Get(0).([]models.SellerReport), args.Error(1)
}

func (m *MockLocalityService) GetReportByLocalityId(id string) ([]models.CarryReport, error) {
	args := m.Called(id)
	return args.Get(0).([]models.CarryReport), args.Error(1)
}
