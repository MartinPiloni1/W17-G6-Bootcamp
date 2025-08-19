package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Verifies the behavior of the HTTP handler responsible for creating a new ProductRecord.
func TestProductRecordHandler_Create(t *testing.T) {
	// Define the payloads and product records used in common by the test cases
	newProductRecordAttributes := models.ProductRecordAttributes{
		LastUpdateDate: "2024-06-01",
		PurchasePrice:  150,
		SalePrice:      200,
		ProductID:      1,
	}

	newProductRecord := models.ProductRecord{
		ID:                      1,
		ProductRecordAttributes: newProductRecordAttributes,
	}

	payload := `
		{
			"last_update_date": "2024-06-01",
			"purchase_price": 150.00,
			"sale_price": 200.00,
			"product_id": 1
		}
	`

	payloadWithWrongValues := `
		{
			"last_update_date": "2024-06-01",
			"purchase_price": -10.0,
			"sale_price": 150.00,
			"product_id": 1
		}
	`

	payloadWithUnknownFields := `
		{
			"last_update_date": "2024-06-01",
			"anUnknownField": 10.0,
			"purchase_price": -10.0,
			"sale_price": 150.00,
			"product_id": 1
		}
	`

	payloadWithFutureDate := `
		{
			"last_update_date": "9999-06-01",
			"purchase_price": -10.0,
			"sale_price": 150.00,
			"product_id": 1
		}
	`

	payloadWithInvalidDateFormat := `
		{
			"last_update_date": "2024-06-01T13:00:00Z",
			"purchase_price": -10.0,
			"sale_price": 150.00,
			"product_id": 1
		}
	`

	// Each test case is constructed by:
	// testName                  — human‐readable description
	// payload                   — raw JSON payload sent in the HTTP request
	// isPayloadError            — whether we expect JSON validation to fail inside the handler
	// productRecordAttributes   — ProductRecord attributes of the new product
	// serviceData               — the ProductRecord object returned by the mocked service
	// serviceError              — the error returned by the mocked service
	// expectedCode              — HTTP status code we expect the handler to produce
	// expectedHeaders           — HTTP headers we expect in the HTTP response
	// expectedBody              — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName                string
		payload                 string
		isPayloadError          bool
		productRecordAttributes models.ProductRecordAttributes
		serviceData             models.ProductRecord
		serviceError            error
		expectedCode            int
		expectedHeaders         http.Header
		expectedBody            string
	}{
		{
			testName:                "Success: Create a new product",
			payload:                 payload,
			isPayloadError:          false,
			productRecordAttributes: newProductRecordAttributes,
			serviceData:             newProductRecord,
			serviceError:            nil,
			expectedCode:            http.StatusCreated,
			expectedHeaders:         http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": {
					"id": 1,
					"last_update_date": "2024-06-01",
					"purchase_price": 150.00,
					"sale_price": 200.00,
					"product_id": 1
				}
			}`,
		},
		{
			testName:                "Error case: Wrong value in a JSON field",
			payload:                 payloadWithWrongValues,
			isPayloadError:          true,
			productRecordAttributes: models.ProductRecordAttributes{},
			serviceData:             models.ProductRecord{},
			serviceError:            nil,
			expectedCode:            http.StatusUnprocessableEntity,
			expectedHeaders:         http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:                "Error case: Unknown JSON fields",
			payload:                 payloadWithUnknownFields,
			isPayloadError:          true,
			productRecordAttributes: models.ProductRecordAttributes{},
			serviceData:             models.ProductRecord{},
			serviceError:            nil,
			expectedCode:            http.StatusBadRequest,
			expectedHeaders:         http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:                "Error case: Future date time",
			payload:                 payloadWithFutureDate,
			isPayloadError:          true,
			productRecordAttributes: models.ProductRecordAttributes{},
			serviceData:             models.ProductRecord{},
			serviceError:            nil,
			expectedCode:            http.StatusUnprocessableEntity,
			expectedHeaders:         http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "invalid date: cannot be in the future"
				}
			`,
		},
		{
			testName:                "Error case: Invalid date format",
			payload:                 payloadWithInvalidDateFormat,
			isPayloadError:          true,
			productRecordAttributes: models.ProductRecordAttributes{},
			serviceData:             models.ProductRecord{},
			serviceError:            nil,
			expectedCode:            http.StatusUnprocessableEntity,
			expectedHeaders:         http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "Invalid date format"
				}
			`,
		},
		{
			testName:                "Error case: Process an error from the service layer",
			payload:                 payload,
			isPayloadError:          false,
			productRecordAttributes: newProductRecordAttributes,
			serviceData:             models.ProductRecord{},
			serviceError:            errors.New("A service layer error"),
			expectedCode:            http.StatusInternalServerError,
			expectedHeaders:         http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Internal Server Error",
					"message": "Internal Server Error"
				}
			`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			serviceMock := &mocks.ProductRecordServiceMock{}

			// If a validation error occurs the service method is not called
			if !tc.isPayloadError {
				serviceMock.
					On("Create", mock.Anything, tc.productRecordAttributes).
					Return(tc.serviceData, tc.serviceError)
			}

			handler := handler.NewProductRecordHandler(serviceMock)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// Act
			handler.Create().ServeHTTP(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedHeaders, response.Header())
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			serviceMock.AssertExpectations(t)
		})
	}
}
