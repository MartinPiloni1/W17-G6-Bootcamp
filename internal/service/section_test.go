package service

import (
	"context"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mock/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

// Helper para crear un puntero a un string.
func stringPtr(s string) *string {
	return &s
}

// Helper para crear un puntero a un int.
func intPtr(i int) *int {
	return &i
}


func TestSectionService_Update(t *testing.T) {
	// Arrange: Preparamos los datos de prueba.

	originalSection := models.Section{
		ID:                 1,
		SectionNumber:      "SEC-101",
		CurrentTemperature: 20,
		MinimumTemperature: 15,
		CurrentCapacity:    50,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	patchDataPartial := models.UpdateSectionRequest{
		SectionNumber:   stringPtr("SEC-102"),
		CurrentCapacity: intPtr(75),
	}
	
	updatedSectionPartial := models.Section{
		ID:                 1,
		SectionNumber:      "SEC-102", // Actualizado
		CurrentTemperature: 20,
		MinimumTemperature: 15,
		CurrentCapacity:    75, // Actualizado
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	patchDataAllFields := models.UpdateSectionRequest{
		SectionNumber:      stringPtr("SEC-ALL"),
		CurrentTemperature: float64Ptr(25.5),
		MinimumTemperature: float64Ptr(18.5),
		CurrentCapacity:    intPtr(99),
		MinimumCapacity:    intPtr(22),
		MaximumCapacity:    intPtr(111),
		WarehouseID:        intPtr(3),
		ProductTypeID:      intPtr(4),
	}
	
	updatedSectionAllFields := models.Section{
		ID:                 1,
		SectionNumber:      "SEC-ALL",
		CurrentTemperature: 25.5,
		MinimumTemperature: 18.5,
		CurrentCapacity:    99,
		MinimumCapacity:    22,
		MaximumCapacity:    111,
		WarehouseID:        3,
		ProductTypeID:      4,
	}


	tests := []struct {
		testName         string
		inputID          int
		patchData        models.UpdateSectionRequest
		mockGetByIDResp  models.Section
		mockGetByIDError error
		mockUpdateResp   models.Section
		mockUpdateError  error
		expectedResp     models.Section
		expectedError    error
	}{
		{
			testName:         "Success: Should update partial fields of a section",
			inputID:          1,
			patchData:        patchDataPartial,
			mockGetByIDResp:  originalSection,
			mockGetByIDError: nil,
			mockUpdateResp:   updatedSectionPartial,
			mockUpdateError:  nil,
			expectedResp:     updatedSectionPartial,
			expectedError:    nil,
		},
		{
			testName:         "Success: Should update all fields of a section",
			inputID:          1,
			patchData:        patchDataAllFields,
			mockGetByIDResp:  originalSection,
			mockGetByIDError: nil,
			mockUpdateResp:   updatedSectionAllFields,
			mockUpdateError:  nil,
			expectedResp:     updatedSectionAllFields,
			expectedError:    nil,
		},
		{
			testName:         "Fail: Should return error when section not found",
			inputID:          99,
			patchData:        patchDataPartial,
			mockGetByIDResp:  models.Section{},
			mockGetByIDError: httperrors.NotFoundError{Message: "Section not found"},
			expectedResp:     models.Section{},
			expectedError:    httperrors.NotFoundError{Message: "Section not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockRepo := new(mock.SectionRepositoryDBMock)
			service := NewSectionServiceDefault(mockRepo)
			
			mockRepo.On("GetByID", testifyMock.Anything, tt.inputID).Return(tt.mockGetByIDResp, tt.mockGetByIDError)

			if tt.mockGetByIDError == nil {
				mockRepo.On("Update", testifyMock.Anything, tt.inputID, testifyMock.AnythingOfType("models.Section")).Return(tt.mockUpdateResp, tt.mockUpdateError)
			}

			result, err := service.Update(context.Background(), tt.inputID, tt.patchData)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}