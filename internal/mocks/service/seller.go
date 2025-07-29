package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type SellerServiceDBMock struct {
	mock.Mock
}

func (m *SellerServiceDBMock) Create(attr models.SellerAttributes) (models.Seller, error) {
	args := m.Called(attr)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *SellerServiceDBMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *SellerServiceDBMock) GetAll() ([]models.Seller, error) {
	args := m.Called()
	return args.Get(0).([]models.Seller), args.Error(1)
}

func (m *SellerServiceDBMock) GetByID(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *SellerServiceDBMock) Update(id int, data *models.SellerAttributes) (models.Seller, error) {
	args := m.Called(id, data)
	return args.Get(0).(models.Seller), args.Error(1)
}
