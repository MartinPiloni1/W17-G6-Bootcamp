package mock

import (
	"context"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
)

type SectionServiceMock struct {
	mock.Mock
}

func (m *SectionServiceMock) Create(ctx context.Context, section models.Section) (models.Section, error){
	args := m.Called(ctx, section)
	return args.Get(0).(models.Section), args.Error(1)
}
func (m *SectionServiceMock) GetAll(ctx context.Context) ([]models.Section, error){
	args := m.Called(ctx)
	return args.Get(0).([]models.Section), args.Error(1)
}
func (m *SectionServiceMock) GetByID(ctx context.Context, id int) (models.Section, error){
	args := m.Called(ctx, id)
	return args.Get(0).(models.Section), args.Error(1)
}
func (m *SectionServiceMock) Update(ctx context.Context, id int, data models.UpdateSectionRequest) (models.Section, error){
	args := m.Called(ctx, id, data)
	return args.Get(0).(models.Section), args.Error(1)
}
func (m *SectionServiceMock) Delete(ctx context.Context, id int) error{
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *SectionServiceMock) GetProductsReport(ctx context.Context, id int) (models.SectionProductsReport, error){
	args := m.Called(ctx, id)
	return args.Get(0).(models.SectionProductsReport), args.Error(1)
}
func (m *SectionServiceMock) GetAllProductsReport(ctx context.Context) ([]models.SectionProductsReport, error){
	args := m.Called(ctx)
	return args.Get(0).([]models.SectionProductsReport), args.Error(1)
}