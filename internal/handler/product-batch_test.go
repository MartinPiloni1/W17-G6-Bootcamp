package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestProductBatchHandler_Create(t *testing.T) {
	okAtt := models.ProductBatchAttibutes{
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
	okBatch := models.ProductBatch{
		ID: 1,
		ProductBatchAttibutes: okAtt,
	}

	validRequestBody := `{
		"data": {
			"batch_number": 10001,
			"current_quantity": 100,
			"current_temperature": 20.5,
			"due_date": "2024-07-10",
			"initial_quantity": 120,
			"manufacturing_date": "2024-07-01",
			"manufacturing_hour": 9,
			"minimum_temperature": 10.5,
			"product_id": 1,
			"section_id": 2
		}
	}`

	invalidDateRequestBody := `{
		"data": {
			"batch_number": 10001,
			"current_quantity": 100,
			"current_temperature": 20.5,
			"due_date": "2024-13-99",
			"initial_quantity": 120,
			"manufacturing_date": "no-date",
			"manufacturing_hour": 9,
			"minimum_temperature": 10.5,
			"product_id": 1,
			"section_id": 2
		}
	}`

	missingFieldRequestBody := `{
		"data": {
			"current_quantity": 100,
			"current_temperature": 20.5,
			"due_date": "2024-07-10",
			"initial_quantity": 120,
			"manufacturing_date": "2024-07-01",
			"manufacturing_hour": 9,
			"minimum_temperature": 10.5,
			"product_id": 1,
			"section_id": 2
		}
	}`

	outOfRangeRequestBody := `{
		"data": {
			"batch_number": 0,
			"current_quantity": -8,
			"current_temperature": 20.5,
			"due_date": "2024-07-10",
			"initial_quantity": -99,
			"manufacturing_date": "2024-07-01",
			"manufacturing_hour": 9,
			"minimum_temperature": 10.5,
			"product_id": 0,
			"section_id": 0
		}
	}`

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
		mockInput    models.ProductBatchAttibutes
		mockOutput   models.ProductBatch
		mockError    error
		mockCalled   bool
	}{
		{
			name:         "Success: Create ProductBatch",
			body:         validRequestBody,
			expectedCode: http.StatusCreated,
			expectedBody: `{"data":{
				"id":1,
				"batch_number":10001,
				"current_quantity":100,
				"current_temperature":20.5,
				"due_date":"2024-07-10",
				"initial_quantity":120,
				"manufacturing_date":"2024-07-01",
				"manufacturing_hour":9,
				"minimum_temperature":10.5,
				"product_id":1,
				"section_id":2
			}}`,
			mockInput:  okAtt,
			mockOutput: okBatch,
			mockError:  nil,
			mockCalled: true,
		},
		{
			name:         "Fail: invalid JSON body (bad json)",
			body:         `{data: }`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"status":"Unprocessable Entity","message":"Invalid body"}`,
			mockCalled:   false,
		},
		{
			name:         "Fail: missing required field",
			body:         missingFieldRequestBody,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"status":"Unprocessable Entity","message":"Invalid body"}`,
			mockCalled:   false,
		},
		{
			name:         "Fail: field out of range",
			body:         outOfRangeRequestBody,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"status":"Unprocessable Entity","message":"Invalid body"}`,
			mockCalled:   false,
		},
		{
			name:         "Fail: invalid date format",
			body:         invalidDateRequestBody,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"status":"Unprocessable Entity","message":"Invalid body"}`,
			mockCalled:   false,
		},
		{
			name:         "Fail: conflict error from service",
			body:         validRequestBody,
			expectedCode: http.StatusConflict,
			expectedBody: `{"status":"Conflict","message":"Product batch already exists"}`,
			mockInput:    okAtt,
			mockOutput:   models.ProductBatch{},
			mockError:    httperrors.ConflictError{Message: "Product batch already exists"},
			mockCalled:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.ProductBatchServiceDefaultMock)
			if tt.mockCalled {
				mockService.On("Create", testifyMock.Anything, tt.mockInput).
					Return(tt.mockOutput, tt.mockError)
			}

			h := handler.NewProductBatchHandler(mockService)
			router := chi.NewRouter()
			router.Post("/api/v1/product-batches", h.Create())

			req := httptest.NewRequest("POST", "/api/v1/product-batches", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}