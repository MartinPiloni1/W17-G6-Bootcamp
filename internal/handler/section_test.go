package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mock/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

// TestSectionHandler_Create tests the Create method of SectionHandler.
func TestSectionHandler_Create(t *testing.T) {
	sectionResponse := models.Section{
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

	tests := []struct {
		testName      string
		requestBody   string
		serviceInput  models.Section
		serviceOutput models.Section
		serviceError  error
		expectedCode  int
		expectedBody  string
	}{
		{
			testName: "Success: Create Section with valid JSON",
			requestBody: `{
                "section_number": "SEC-101", "current_temperature": 20, "minimum_temperature": 15,
                "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100,
                "warehouse_id": 1, "product_type_id": 1
            }`,
			serviceInput: models.Section{
				SectionNumber:      "SEC-101",
				CurrentTemperature: 20,
				MinimumTemperature: 15,
				CurrentCapacity:    50,
				MinimumCapacity:    10,
				MaximumCapacity:    100,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
			serviceOutput: sectionResponse,
			serviceError:  nil,
			expectedCode:  http.StatusCreated,
			expectedBody: `{
                "data": {
                    "id": 1, "section_number": "SEC-101", "current_temperature": 20, "minimum_temperature": 15,
                    "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100,
                    "warehouse_id": 1, "product_type_id": 1
                }
            }`,
		},
		{
			testName:      "Fail: Create Section with unprocessable entity",
			requestBody:   `{"section_number": "SEC-101"}`,
			serviceInput:  models.Section{},
			serviceOutput: models.Section{},
			serviceError:  nil,
			expectedCode:  http.StatusUnprocessableEntity,
			expectedBody:  `{"status": "Unprocessable Entity", "message": "Invalid JSON body"}`,
		},
		{
			testName:      "Fail: Create Section with invalid JSON",
			requestBody:   `{invalid json}`,
			serviceInput:  models.Section{},
			serviceOutput: models.Section{},
			serviceError:  nil,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"status": "Bad Request", "message": "Invalid body"}`,
		},
		{
			testName: "Fail: Create Section with conflict error from service",
			requestBody: `{
                "section_number": "SEC-101", "current_temperature": 20, "minimum_temperature": 15,
                "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100,
                "warehouse_id": 1, "product_type_id": 1
            }`,
			serviceInput: models.Section{
				SectionNumber:      "SEC-101",
				CurrentTemperature: 20,
				MinimumTemperature: 15,
				CurrentCapacity:    50,
				MinimumCapacity:    10,
				MaximumCapacity:    100,
				WarehouseID:        1,
				ProductTypeID:      1,
			},
			serviceOutput: models.Section{},
			serviceError:  httperrors.ConflictError{Message: "Section number already exists."},
			expectedCode:  http.StatusConflict,
			expectedBody:  `{"status": "Conflict", "message": "Section number already exists."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mock.SectionServiceMock)

			if tt.expectedCode == http.StatusCreated || tt.expectedCode == http.StatusConflict {
				mockService.On("Create", testifyMock.Anything, tt.serviceInput).Return(tt.serviceOutput, tt.serviceError)
			}

			handler := NewSectionHandler(mockService)
			router := chi.NewRouter()
			router.Post("/api/v1/sections", handler.Create())

			req := httptest.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedCode, response.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, response.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}


func TestSectionHandler_GetAll(t *testing.T) {
	sectionResponse := models.Section{
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

	tests := []struct {
		testName      string
		serviceOutput []models.Section
		serviceError  error
		expectedCode  int
		expectedBody  string
	}{
		{
			testName:      "Success: Get All Sections",
			serviceOutput: []models.Section{sectionResponse},
			serviceError:  nil,
			expectedCode:  http.StatusOK,
			expectedBody: `{
                "data": [
                    {
                        "id": 1, "section_number": "SEC-101", "current_temperature": 20, "minimum_temperature": 15,
                        "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100,
                        "warehouse_id": 1, "product_type_id": 1
                    }
                ]
            }`,
		},
		{
			testName:      "Fail: Service returns an error",
			serviceOutput: nil,
			serviceError:  errors.New("internal server error"),
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"status": "Internal Server Error", "message": "Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mock.SectionServiceMock)

			mockService.On("GetAll", testifyMock.Anything).Return(tt.serviceOutput, tt.serviceError)

			handler := NewSectionHandler(mockService)
			router := chi.NewRouter()
			router.Get("/api/v1/sections", handler.GetAll())

			req := httptest.NewRequest(http.MethodGet, "/api/v1/sections", nil)
			req.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedCode, response.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, response.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}


func TestSectionHandler_GetByID(t *testing.T) {
	sectionResponse := models.Section{
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

	tests := []struct {
		testName      string
		inputID       int    // ID que espera el mock
		requestURL    string // URL para la petición
		serviceOutput models.Section
		serviceError  error
		expectedCode  int
		expectedBody  string
	}{
		{
			testName:      "Success: Get Section by ID",
			inputID:       1,
			requestURL:    "/api/v1/sections/1",
			serviceOutput: sectionResponse,
			serviceError:  nil,
			expectedCode:  http.StatusOK,
			expectedBody: `{
                "data": {
                    "id": 1, "section_number": "SEC-101", "current_temperature": 20, "minimum_temperature": 15,
                    "current_capacity": 50, "minimum_capacity": 10, "maximum_capacity": 100,
                    "warehouse_id": 1, "product_type_id": 1
                }
            }`,
		},
		{
			testName:      "Fail: Invalid ID",
			inputID:       0,
			requestURL:    "/api/v1/sections/abc",
			serviceOutput: models.Section{},
			serviceError:  nil,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"status": "Bad Request", "message": "Invalid ID"}`,
		},
		{
			testName:      "Fail: Section not found",
			inputID:       99,
			requestURL:    "/api/v1/sections/99",
			serviceOutput: models.Section{},
			serviceError:  httperrors.NotFoundError{Message: "Section not found"},
			expectedCode:  http.StatusNotFound,
			expectedBody:  `{"status": "Not Found", "message": "Section not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mock.SectionServiceMock)

			if tt.expectedCode == http.StatusOK || tt.expectedCode == http.StatusNotFound {
				mockService.On("GetByID", testifyMock.Anything, tt.inputID).Return(tt.serviceOutput, tt.serviceError)
			}

			handler := NewSectionHandler(mockService)
			router := chi.NewRouter()
			router.Get("/api/v1/sections/{id}", handler.GetByID())

			req := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)
			req.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedCode, response.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, response.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}


func TestSectionHandler_Delete(t *testing.T) {
	tests := []struct {
		testName      string
		inputID       int
		requestURL    string
		serviceError  error
		expectedCode  int
		expectedBody  string
	}{
		{
			testName:      "Success: Delete Section",
			inputID:       1,
			requestURL:    "/api/v1/sections/1",
			serviceError:  nil,
			expectedCode:  http.StatusNoContent,
			expectedBody:  "",
		},
		{
			testName:      "Fail: Section not found",
			inputID:       99,
			requestURL:    "/api/v1/sections/99",
			serviceError:  httperrors.NotFoundError{Message: "Section not found"},
			expectedCode:  http.StatusNotFound,
			expectedBody:  `{"status": "Not Found", "message": "Section not found"}`,
		},
		{
			testName:      "Fail: Invalid ID",
			inputID:       0,
			requestURL:    "/api/v1/sections/abc",
			serviceError:  nil,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"status": "Bad Request", "message": "Invalid ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mock.SectionServiceMock)

			if tt.expectedCode == http.StatusNoContent || tt.expectedCode == http.StatusNotFound {
				mockService.On("Delete", testifyMock.Anything, tt.inputID).Return(tt.serviceError)
			}

			handler := NewSectionHandler(mockService)
			router := chi.NewRouter()
			router.Delete("/api/v1/sections/{id}", handler.Delete())

			req := httptest.NewRequest(http.MethodDelete, tt.requestURL, nil)
			
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedCode, response.Code)
			
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, response.Body.String())
			} else {
				assert.Empty(t, response.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Helper to create a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// Helper to create a pointer to an int.
func intPtr(i int) *int {
	return &i
}

func TestSectionHandler_Update(t *testing.T) {
	// Data for the request body
	updateRequest := models.UpdateSectionRequest{
		// Use the helper function to get a pointer
		SectionNumber:   stringPtr("SEC-102"),
		CurrentCapacity: intPtr(75),
	}

	// Expected service response on success
	updatedSectionResponse := models.Section{
		ID:                 1,
		SectionNumber:      "SEC-102",
		CurrentTemperature: 20,
		MinimumTemperature: 15,
		CurrentCapacity:    75,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	tests := []struct {
		testName      string
		inputID       int
		requestURL    string
		requestBody   string
		serviceInput  models.UpdateSectionRequest
		serviceOutput models.Section
		serviceError  error
		expectedCode  int
		expectedBody  string
	}{
		{
			testName:     "Success: Update Section",
			inputID:      1,
			requestURL:   "/api/v1/sections/1",
			requestBody:  `{"section_number": "SEC-102", "current_capacity": 75}`,
			serviceInput: updateRequest,
			serviceOutput:updatedSectionResponse,
			serviceError: nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"data": {"id": 1, "section_number": "SEC-102", "current_temperature": 20, "minimum_temperature": 15, "current_capacity": 75, "minimum_capacity": 10, "maximum_capacity": 100, "warehouse_id": 1, "product_type_id": 1}}`,
		},
		{
			testName:     "Fail: Section not found",
			inputID:      99,
			requestURL:   "/api/v1/sections/99",
			requestBody:  `{"section_number": "SEC-102"}`,
			serviceInput: models.UpdateSectionRequest{SectionNumber: stringPtr("SEC-102")},
			serviceError: httperrors.NotFoundError{Message: "Section not found"},
			expectedCode: http.StatusNotFound,
			expectedBody: `{"status": "Not Found", "message": "Section not found"}`,
		},
		{
			testName:     "Fail: Invalid ID",
			inputID:      0,
			requestURL:   "/api/v1/sections/abc",
			serviceError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status": "Bad Request", "message": "Invalid ID"}`,
		},
		{
			testName:     "Fail: Invalid JSON",
			inputID:      1,
			requestURL:   "/api/v1/sections/1",
			requestBody:  `{"section_number": "SEC-102", "current_capacity": "75"}`,
			serviceInput: models.UpdateSectionRequest{SectionNumber: stringPtr("SEC-102"), CurrentCapacity: intPtr(75)},
			serviceError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status": "Bad Request", "message": "Invalid body"}`,
		},
        {
            testName:     "Fail: Unprocessable entity",
            inputID:      1,
            requestURL:   "/api/v1/sections/1",
            requestBody:  `{"current_capacity": -5}`,
            serviceInput: models.UpdateSectionRequest{}, 
            serviceError: nil,
            expectedCode: http.StatusUnprocessableEntity,
            expectedBody: `{"status": "Unprocessable Entity", "message": "Invalid JSON body"}`,
        },
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mock.SectionServiceMock)

			if tt.expectedCode == http.StatusOK || tt.expectedCode == http.StatusNotFound {
				mockService.On("Update", testifyMock.Anything, tt.inputID, tt.serviceInput).Return(tt.serviceOutput, tt.serviceError)
			}

			handler := NewSectionHandler(mockService)
			router := chi.NewRouter()
			router.Patch("/api/v1/sections/{id}", handler.Update())

			req := httptest.NewRequest(http.MethodPatch, tt.requestURL, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedCode, response.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, response.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSectionHandler_GetProductsReport(t *testing.T) {
	report1 := models.SectionProductsReport{
		SectionID:     1,
		SectionNumber: "SEC-101",
		ProductsCount: 50,
	}
	report2 := models.SectionProductsReport{
		SectionID:     2,
		SectionNumber: "SEC-102",
		ProductsCount: 75,
	}

	tests := []struct {
		testName            string
		requestURL          string
		mockMethod          string
		inputID             int
		serviceOutputSingle models.SectionProductsReport
		serviceOutputAll    []models.SectionProductsReport
		serviceError        error
		expectedCode        int
		expectedBody        string
	}{
		{
			testName:            "Success: Get report for single section",
			requestURL:          "/api/v1/sections/reportProducts?id=1",
			mockMethod:          "GetProductsReport",
			inputID:             1,
			serviceOutputSingle: report1,
			serviceError:        nil,
			expectedCode:        http.StatusOK,
			expectedBody:        `{"data": [{"section_id":1, "section_number":"SEC-101", "products_count":50}]}`,
		},
		{
			testName:            "Fail: Section not found for single report",
			requestURL:          "/api/v1/sections/reportProducts?id=99",
			mockMethod:          "GetProductsReport",
			inputID:             99,
			serviceError:        httperrors.NotFoundError{Message: "Section not found"},
			expectedCode:        http.StatusNotFound,
			expectedBody:        `{"status": "Not Found", "message": "Section not found"}`,
		},
		{
			testName:     "Fail: Invalid ID for single report",
			requestURL:   "/api/v1/sections/reportProducts?id=abc",
			mockMethod:   "", 
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status": "Bad Request", "message": "Invalid section ID"}`,
		},
		{
			testName:         "Success: Get report for all sections",
			requestURL:       "/api/v1/sections/reportProducts",
			mockMethod:       "GetAllProductsReport",
			serviceOutputAll: []models.SectionProductsReport{report1, report2},
			serviceError:     nil,
			expectedCode:     http.StatusOK,
			expectedBody:     `{"data": [{"section_id":1, "section_number":"SEC-101", "products_count":50}, {"section_id":2, "section_number":"SEC-102", "products_count":75}]}`,
		},
		{
			testName:     "Fail: Service error on get all reports",
			requestURL:   "/api/v1/sections/reportProducts",
			mockMethod:   "GetAllProductsReport",
			serviceError: errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status": "Internal Server Error", "message": "Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mock.SectionServiceMock)

			switch tt.mockMethod {
			case "GetProductsReport":
				mockService.On(tt.mockMethod, testifyMock.Anything, tt.inputID).Return(tt.serviceOutputSingle, tt.serviceError)
			case "GetAllProductsReport":
				mockService.On(tt.mockMethod, testifyMock.Anything).Return(tt.serviceOutputAll, tt.serviceError)
			}

			handler := NewSectionHandler(mockService)
			router := chi.NewRouter()
			router.Get("/api/v1/sections/reportProducts", handler.GetProductsReport())

			req := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)

			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedCode, response.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, response.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}