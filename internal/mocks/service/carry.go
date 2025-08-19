package mocks

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type CarryServiceDBMock struct {
	mock.Mock
}

func (m *CarryServiceDBMock) Create(carryAttributes models.CarryAttributes) (models.Carry, error) {
	args := m.Called(carryAttributes)
	return args.Get(0).(models.Carry), args.Error(1)
}
