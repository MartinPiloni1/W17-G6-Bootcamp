package service_test

import (
	"errors"
	"testing"
	"time"

	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
)

func TestInboundOrderService_Create(t *testing.T) {
	// Prepare valid attributes for testing
	validAttrs := models.InboundOrderAttributes{
		OrderNumber:    "X123",
		OrderDate:      time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
		EmployeeID:     1,
		WarehouseID:    2,
		ProductBatchID: 3,
	}
	expectedCreated := models.InboundOrder{
		ID:                     11,
		InboundOrderAttributes: validAttrs,
	}
	validEmployee := models.Employee{Id: 1}
	validWarehouse := models.Warehouse{Id: 2}

	tests := []struct {
		name          string
		setupMocks    func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock)
		input         models.InboundOrderAttributes
		expectedOrder models.InboundOrder
		expectedError error
	}{
		{
			name: "successfully creates a new inbound order",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, nil).Once()
				empRepo.On("GetByID", validAttrs.EmployeeID).
					Return(validEmployee, nil).Once()
				whRepo.On("GetByID", validAttrs.WarehouseID).
					Return(validWarehouse, nil).Once()
				repo.On("Create", models.InboundOrder{InboundOrderAttributes: validAttrs}).
					Return(expectedCreated, nil).Once()
			},
			input:         validAttrs,
			expectedOrder: expectedCreated,
			expectedError: nil,
		},
		{
			name: "returns error if GetByOrderNumber fails",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, errors.New("db error")).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: errors.New("db error"),
		},
		{
			name: "returns conflict error if order number already exists",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				duplicate := models.InboundOrder{ID: 99}
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(duplicate, nil).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: httperrors.ConflictError{},
		},
		{
			name: "returns conflict error if fetching employee returns error",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, nil).Once()
				empRepo.On("GetByID", validAttrs.EmployeeID).
					Return(models.Employee{}, errors.New("employee error")).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: httperrors.ConflictError{},
		},
		{
			name: "returns conflict error if employee not found",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, nil).Once()
				empRepo.On("GetByID", validAttrs.EmployeeID).
					Return(models.Employee{Id: 0}, nil).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: httperrors.ConflictError{},
		},
		{
			name: "returns conflict error if fetching warehouse returns error",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, nil).Once()
				empRepo.On("GetByID", validAttrs.EmployeeID).
					Return(validEmployee, nil).Once()
				whRepo.On("GetByID", validAttrs.WarehouseID).
					Return(models.Warehouse{}, errors.New("warehouse error")).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: httperrors.ConflictError{},
		},
		{
			name: "returns conflict error if warehouse not found",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, nil).Once()
				empRepo.On("GetByID", validAttrs.EmployeeID).
					Return(validEmployee, nil).Once()
				whRepo.On("GetByID", validAttrs.WarehouseID).
					Return(models.Warehouse{Id: 0}, nil).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: httperrors.ConflictError{},
		},
		{
			name: "returns error if repository Create fails",
			setupMocks: func(repo *mocks.MockInboundOrderRepository, empRepo *mocks.MockEmployeeRepository, whRepo *mocks.WarehouseRepositoryMock) {
				repo.On("GetByOrderNumber", validAttrs.OrderNumber).
					Return(models.InboundOrder{}, nil).Once()
				empRepo.On("GetByID", validAttrs.EmployeeID).
					Return(validEmployee, nil).Once()
				whRepo.On("GetByID", validAttrs.WarehouseID).
					Return(validWarehouse, nil).Once()
				repo.On("Create", models.InboundOrder{InboundOrderAttributes: validAttrs}).
					Return(models.InboundOrder{}, errors.New("create failure")).Once()
			},
			input:         validAttrs,
			expectedOrder: models.InboundOrder{},
			expectedError: errors.New("create failure"),
		},
	}

	for _, tc := range tests {
		tt := tc // capture iteration variable for parallel subtests
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := new(mocks.MockInboundOrderRepository)
			empRepo := new(mocks.MockEmployeeRepository)
			whRepo := new(mocks.WarehouseRepositoryMock)
			s := service.NewInboundOrderService(repo, empRepo, whRepo)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, empRepo, whRepo)
			}

			// Act
			result, err := s.Create(tt.input)

			// Assert
			assert.Equal(t, tt.expectedOrder, result)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				// Use assert.IsType for error type assertion
				assert.IsType(t, tt.expectedError, err)
			}

			repo.AssertExpectations(t)
			empRepo.AssertExpectations(t)
			whRepo.AssertExpectations(t)
		})
	}
}
