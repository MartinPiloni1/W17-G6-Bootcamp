package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mock "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestCarryHandler_Create tests the Create method of the CarryHandler.
// It checks for successful creation, validation errors, and various error scenarios.
func TestCarryHandler_Create(t *testing.T) {
	// Test data
	validCarry := models.CarryAttributes{
		Cid:         "12345",
		CompanyName: "Test Company",
		Address:     "123 Test Street",
		Telephone:   "1234567890",
		LocalityId:  "1001",
	}

	invalidCarry := models.CarryAttributes{
		Cid:         "",
		CompanyName: "Test Company",
		Address:     "123 Test Street",
		Telephone:   "1234567890",
		LocalityId:  "1001",
	}

	conflictCarry := models.CarryAttributes{
		Cid:         "99999",
		CompanyName: "Conflict Company",
		Address:     "456 Conflict Ave",
		Telephone:   "0987654321",
		LocalityId:  "1002",
	}

	cases := []struct {
		name       string
		body       []byte
		setupMock  func(s *mock.CarryServiceDBMock)
		wantStatus int
	}{
		{
			name: "successful creation",
			body: func() []byte { b, _ := json.Marshal(validCarry); return b }(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				expectedCarry := models.Carry{
					Id:              1,
					CarryAttributes: validCarry,
				}
				s.On("Create", validCarry).Return(expectedCarry, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid JSON body",
			body:       []byte("{not a valid json}"),
			setupMock:  func(s *mock.CarryServiceDBMock) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "validation error - empty CID",
			body: func() []byte { b, _ := json.Marshal(invalidCarry); return b }(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				s.On("Create", invalidCarry).Return(models.Carry{}, httperrors.BadRequestError{Message: "the field Cid must not be empty"})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty CompanyName",
			body: func() []byte {
				invalidData := validCarry
				invalidData.CompanyName = ""
				b, _ := json.Marshal(invalidData)
				return b
			}(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				invalidData := validCarry
				invalidData.CompanyName = ""
				s.On("Create", invalidData).Return(models.Carry{}, httperrors.BadRequestError{Message: "the field CompanyName must not be empty"})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty Address",
			body: func() []byte {
				invalidData := validCarry
				invalidData.Address = ""
				b, _ := json.Marshal(invalidData)
				return b
			}(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				invalidData := validCarry
				invalidData.Address = ""
				s.On("Create", invalidData).Return(models.Carry{}, httperrors.BadRequestError{Message: "the field Address must not be empty"})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty Telephone",
			body: func() []byte {
				invalidData := validCarry
				invalidData.Telephone = ""
				b, _ := json.Marshal(invalidData)
				return b
			}(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				invalidData := validCarry
				invalidData.Telephone = ""
				s.On("Create", invalidData).Return(models.Carry{}, httperrors.BadRequestError{Message: "the field Telephone must not be empty"})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty LocalityId",
			body: func() []byte {
				invalidData := validCarry
				invalidData.LocalityId = ""
				b, _ := json.Marshal(invalidData)
				return b
			}(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				invalidData := validCarry
				invalidData.LocalityId = ""
				s.On("Create", invalidData).Return(models.Carry{}, httperrors.BadRequestError{Message: "the field LocalityId must not be empty"})
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "conflict error",
			body: func() []byte { b, _ := json.Marshal(conflictCarry); return b }(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				s.On("Create", conflictCarry).Return(models.Carry{}, httperrors.ConflictError{Message: "carry with this CID already exists"})
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal server error",
			body: func() []byte { b, _ := json.Marshal(validCarry); return b }(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				s.On("Create", validCarry).Return(models.Carry{}, errors.New("database connection failed"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "not found error",
			body: func() []byte { b, _ := json.Marshal(validCarry); return b }(),
			setupMock: func(s *mock.CarryServiceDBMock) {
				s.On("Create", validCarry).Return(models.Carry{}, httperrors.NotFoundError{Message: "locality not found"})
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			service := new(mock.CarryServiceDBMock)
			tc.setupMock(service)
			handler := handler.NewCarryHandler(service)

			// Create router and register handler
			router := chi.NewRouter()
			router.Post("/api/v1/carries", handler.Create())

			// Create request
			req := httptest.NewRequest("POST", "/api/v1/carries", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.wantStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

// TestCarryHandler_Create_ResponseBody tests the response body structure for successful creation
func TestCarryHandler_Create_ResponseBody(t *testing.T) {
	// Arrange
	validCarry := models.CarryAttributes{
		Cid:         "12345",
		CompanyName: "Test Company",
		Address:     "123 Test Street",
		Telephone:   "1234567890",
		LocalityId:  "1001",
	}

	expectedCarry := models.Carry{
		Id:              1,
		CarryAttributes: validCarry,
	}

	service := new(mock.CarryServiceDBMock)
	service.On("Create", validCarry).Return(expectedCarry, nil)

	handler := handler.NewCarryHandler(service)
	router := chi.NewRouter()
	router.Post("/api/v1/carries", handler.Create())

	body, _ := json.Marshal(validCarry)
	req := httptest.NewRequest("POST", "/api/v1/carries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response structure
	assert.Contains(t, response, "data")
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "12345", data["cid"])
	assert.Equal(t, "Test Company", data["company_name"])
	assert.Equal(t, "123 Test Street", data["address"])
	assert.Equal(t, "1234567890", data["telephone"])
	assert.Equal(t, "1001", data["locality_id"])

	service.AssertExpectations(t)
}

// TestCarryHandler_Create_EmptyBody tests the handler with an empty request body
func TestCarryHandler_Create_EmptyBody(t *testing.T) {
	// Arrange
	service := new(mock.CarryServiceDBMock)
	handler := handler.NewCarryHandler(service)
	router := chi.NewRouter()
	router.Post("/api/v1/carries", handler.Create())

	req := httptest.NewRequest("POST", "/api/v1/carries", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// TestCarryHandler_Create_MalformedJSON tests the handler with malformed JSON
func TestCarryHandler_Create_MalformedJSON(t *testing.T) {
	// Arrange
	service := new(mock.CarryServiceDBMock)
	handler := handler.NewCarryHandler(service)
	router := chi.NewRouter()
	router.Post("/api/v1/carries", handler.Create())

	malformedJSON := []byte(`{"cid": "12345", "company_name": "Test", "address": "123 Street", "telephone": "1234567890", "locality_id": "1001"`)
	req := httptest.NewRequest("POST", "/api/v1/carries", bytes.NewBuffer(malformedJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
