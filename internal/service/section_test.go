package service

import (
	"context"
	"errors"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

// Helper to create a string pointer
func stringPtr(s string) *string {
	return &s
}

// Helper para crear un puntero a un int.
func intPtr(i int) *int {
	return &i
}

// TestSectionService_Update tests the Update method of SectionService
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
			mockRepo := new(mocks.SectionRepositoryDBMock)
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

// Helper to create a float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}

// TestSectionService_GetAll tests the GetAll method of SectionService
func TestSectionService_GetAll(t *testing.T) {
	// Arrange
	expectedSections := []models.Section{
		{ID: 1, SectionNumber: "SEC-101"},
		{ID: 2, SectionNumber: "SEC-102"},
	}
	expectedError := errors.New("database error")

	tests := []struct {
		testName      string
		mockResp      []models.Section
		mockErr       error
		expectedResp  []models.Section
		expectedError error
	}{
		{
			testName:      "Success: should return all sections",
			mockResp:      expectedSections,
			mockErr:       nil,
			expectedResp:  expectedSections,
			expectedError: nil,
		},
		{
			testName:      "Fail: should return an error",
			mockResp:      nil,
			mockErr:       expectedError,
			expectedResp:  []models.Section{},
			expectedError: expectedError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockRepo := new(mocks.SectionRepositoryDBMock)
			service := NewSectionServiceDefault(mockRepo)

			mockRepo.On("GetAll", testifyMock.Anything).Return(tt.mockResp, tt.mockErr)

			// Act
			result, err := service.GetAll(context.Background())

			// Assert
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestSectionService_GetByID tests the GetByID method of SectionService
func TestSectionService_GetByID(t *testing.T) {
	expectedSection := models.Section{ID: 1, SectionNumber: "SEC-101"}
	inputID := 1
	mockRepo := new(mocks.SectionRepositoryDBMock)
	service := NewSectionServiceDefault(mockRepo)

	mockRepo.On("GetByID", testifyMock.Anything, inputID).Return(expectedSection, nil)

	result, err := service.GetByID(context.Background(), inputID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSection, result)
	mockRepo.AssertExpectations(t)
}

// TestSectionService_Delete tests the Delete method of SectionService
func TestSectionService_Delete(t *testing.T) {
	inputID := 1
	mockRepo := new(mocks.SectionRepositoryDBMock)
	service := NewSectionServiceDefault(mockRepo)

	mockRepo.On("Delete", testifyMock.Anything, inputID).Return(nil)

	err := service.Delete(context.Background(), inputID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestSectionService_Create tests the Create method of SectionService
func TestSectionService_Create(t *testing.T) {
	sectionToCreate := models.Section{SectionNumber: "SEC-NEW"}
	expectedSection := models.Section{ID: 1, SectionNumber: "SEC-NEW"}
	mockRepo := new(mocks.SectionRepositoryDBMock)
	service := NewSectionServiceDefault(mockRepo)
	mockRepo.On("Create", testifyMock.Anything, sectionToCreate).Return(expectedSection, nil)

	result, err := service.Create(context.Background(), sectionToCreate)

	assert.NoError(t, err)
	assert.Equal(t, expectedSection, result)
	mockRepo.AssertExpectations(t)
}

// TestSectionService_GetProductsReport tests the GetProductsReport method of SectionService
func TestSectionService_GetProductsReport(t *testing.T) {
	expectedReport := models.SectionProductsReport{SectionID: 1, ProductsCount: 10}
	inputID := 1
	mockRepo := new(mocks.SectionRepositoryDBMock)
	service := NewSectionServiceDefault(mockRepo)
	mockRepo.On("GetProductsReport", testifyMock.Anything, inputID).Return(expectedReport, nil)

	result, err := service.GetProductsReport(context.Background(), inputID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReport, result)
	mockRepo.AssertExpectations(t)
}

// TestSectionService_GetAllProductsReport tests the GetAllProductsReport method of SectionService
func TestSectionService_GetAllProductsReport(t *testing.T) {
	expectedReports := []models.SectionProductsReport{
		{SectionID: 1, ProductsCount: 10},
		{SectionID: 2, ProductsCount: 20},
	}
	mockRepo := new(mocks.SectionRepositoryDBMock)
	service := NewSectionServiceDefault(mockRepo)

	mockRepo.On("GetAllProductsReport", testifyMock.Anything).Return(expectedReports, nil)

	result, err := service.GetAllProductsReport(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedReports, result)
	mockRepo.AssertExpectations(t)
}