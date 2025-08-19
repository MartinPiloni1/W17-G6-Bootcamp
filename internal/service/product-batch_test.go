package service_test

import (
	"context"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"errors"
)

func TestProductBatchServiceDefault_Create(t *testing.T) {
	ctx := context.Background()

	input := models.ProductBatchAttibutes{
		BatchNumber:        10001,
		CurrentQuantity:    100,
		CurrentTemperature: 20.5,
		DueDate:            "2024-07-10",
		InitialQuantity:    120,
		ManufacturingDate:  "2024-07-01",
		ManufacturingHour:  9,
		MinimumTemperature: 10.5,
		ProductID:          1,
		SectionID:          2,
	}
	output := models.ProductBatch{
		ID:                   1,
		ProductBatchAttibutes: input,
	}

	tests := []struct {
		name        string
		mockOutput  models.ProductBatch
		mockError   error
		expectedOut models.ProductBatch
		expectedErr error
	}{
		{
			name:        "Success: creates product batch",
			mockOutput:  output,
			mockError:   nil,
			expectedOut: output,
			expectedErr: nil,
		},
		{
			name:        "Fail: repository error",
			mockOutput:  models.ProductBatch{},
			mockError:   errors.New("repository error"),
			expectedOut: models.ProductBatch{},
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mocks.ProductBatchRepositoryMock)

			repoMock.
				On("Create", testifyMock.Anything, input).
				Return(tt.mockOutput, tt.mockError)

			s := service.NewProductBatchServiceDefault(repoMock)

			result, err := s.Create(ctx, input)

			assert.Equal(t, tt.expectedOut, result)
			assert.Equal(t, tt.expectedErr, err)

			repoMock.AssertExpectations(t)
		})
	}
}