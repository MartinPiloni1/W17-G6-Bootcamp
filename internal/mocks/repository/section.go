package mocks

import (
    "context"

    "github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
    "github.com/stretchr/testify/mock"
)

type SectionRepositoryDBMock struct {
    mock.Mock
}

func (m *SectionRepositoryDBMock) Create(ctx context.Context, section models.Section) (models.Section, error) {
    args := m.Called(ctx, section)
    if args.Get(0) == nil {
        return models.Section{}, args.Error(1)
    }
    return args.Get(0).(models.Section), args.Error(1)
}

func (m *SectionRepositoryDBMock) GetAll(ctx context.Context) ([]models.Section, error) {
    args := m.Called(ctx)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]models.Section), args.Error(1)
}

func (m *SectionRepositoryDBMock) GetByID(ctx context.Context, id int) (models.Section, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return models.Section{}, args.Error(1)
    }
    return args.Get(0).(models.Section), args.Error(1)
}

func (m *SectionRepositoryDBMock) Update(ctx context.Context, id int, data models.Section) (models.Section, error) {
    args := m.Called(ctx, id, data)
    if args.Get(0) == nil {
        return models.Section{}, args.Error(1)
    }
    return args.Get(0).(models.Section), args.Error(1)
}

func (m *SectionRepositoryDBMock) Delete(ctx context.Context, id int) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *SectionRepositoryDBMock) GetProductsReport(ctx context.Context, id int) (models.SectionProductsReport, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return models.SectionProductsReport{}, args.Error(1)
    }
    return args.Get(0).(models.SectionProductsReport), args.Error(1)
}

func (m *SectionRepositoryDBMock) GetAllProductsReport(ctx context.Context) ([]models.SectionProductsReport, error) {
    args := m.Called(ctx)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]models.SectionProductsReport), args.Error(1)
}