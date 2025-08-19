package mocks

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRecordRepositoryDBMock struct {
	mock.Mock
}

func (m *ProductRecordRepositoryDBMock) Create(ctx context.Context, attributes models.ProductRecordAttributes) (models.ProductRecord, error) {
	args := m.Called(ctx, attributes)
	return args.Get(0).(models.ProductRecord), args.Error(1)
}
