package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type SellerRepositoryDBMock struct {
	mock.Mock
}

func (m *SellerRepositoryDBMock) GetByID(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *SellerRepositoryDBMock) GetAll() ([]models.Seller, error) {
	args := m.Called()
	return args.Get(0).([]models.Seller), args.Error(1)
}

func (m *SellerRepositoryDBMock) Create(attr models.SellerAttributes) (models.Seller, error) {
	args := m.Called(attr)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *SellerRepositoryDBMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *SellerRepositoryDBMock) Update(id int, attr *models.SellerAttributes) (models.Seller, error) {
	args := m.Called(id, attr)
	return args.Get(0).(models.Seller), args.Error(1)
}
