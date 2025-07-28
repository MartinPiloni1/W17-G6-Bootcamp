package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Verifies the behavior of the HTTP handler responsible for creating a new Product. It covers:
// - Successful creation
// - Error when incomplete JSON bodies are given
// - Error when invalid JSON fields values are given
// - Error when unknown JSON fields are given
// - Error propagation from the service layer
func TestProductHandler_Create(t *testing.T) {
	// Define the payloads and products used in common by the test cases
	newProductAttributes := models.ProductAttributes{
		Description:                    "Yogurt helado",
		ExpirationRate:                 7,
		FreezingRate:                   2,
		Height:                         10.5,
		Length:                         20.0,
		Width:                          15.0,
		NetWeight:                      1.2,
		ProductCode:                    "YOG01",
		RecommendedFreezingTemperature: -5.0,
		ProductTypeID:                  3,
		SellerID:                       utils.Ptr(1),
	}

	newProduct := models.Product{
		ID:                1,
		ProductAttributes: newProductAttributes,
	}

	payload := `
		{
			"description": "Yogurt helado",
			"expiration_rate": 7,
			"freezing_rate": 2,
			"height": 10.5,
			"length": 20.0,
			"width": 15.0,
			"netweight": 1.2,
			"product_code": "YOG01",
			"recommended_freezing_temperature": -5.0,
			"product_type_id": 3,
			"seller_id": 1
		}
	`

	payloadWithMissingFields := `
		{
			"description": "Yogurt helado",
			"recommended_freezing_temperature": -5.0,
			"product_type_id": 3,
			"seller_id": 1
		}
	`

	payloadWithWrongValues := `
		{
			"description": "Yogurt helado",
			"expiration_rate": 7,
			"freezing_rate": 2,
			"height": -10,
			"length": -20.0,
			"width": -15.0,
			"netweight": -50.0,
			"product_code": "YOG01",
			"recommended_freezing_temperature": -5.0,
			"product_type_id": 3,
			"seller_id": -1
		}
	`

	payloadWithUnkownFields := `
		{
			"description": "Yogurt helado",
			"anUnkownField": 1,
			"expiration_rate": 7,
			"freezing_rate": 2,
			"height": -10,
			"length": -20.0,
			"width": -15.0,
			"netweight": -50.0,
			"product_code": "YOG01",
			"recommended_freezing_temperature": -5.0,
			"product_type_id": 3,
			"seller_id": -1
		}
	`

	// Each test case is constructed by:
	// testName            — human‐readable description
	// payload             — raw JSON payload sent in the HTTP request
	// isPayloadError      — whether we expect JSON validation to fail inside the handler
	// productAttributes   — Product attributes of the new product
	// serviceData         — the Product object returned by the mocked service
	// serviceError        — the error returned by the mocked service
	// expectedCode        — HTTP status code we expect the handler to produce
	// expectedHeaders     — HTTP headers we expect in the HTTP response
	// expectedBody        — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName          string
		payload           string
		isPayloadError    bool
		productAttributes models.ProductAttributes
		serviceData       models.Product
		serviceError      error
		expectedCode      int
		expectedHeaders   http.Header
		expectedBody      string
	}{
		{
			testName:          "Success: Create a new product",
			payload:           payload,
			isPayloadError:    false,
			productAttributes: newProductAttributes,
			serviceData:       newProduct,
			serviceError:      nil,
			expectedCode:      http.StatusCreated,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": {
					"id": 1,
					"description": "Yogurt helado",
					"expiration_rate": 7,
					"freezing_rate": 2,
					"height": 10.5,
					"length": 20.0,
					"width": 15.0,
					"netweight": 1.2,
					"product_code": "YOG01",
					"recommended_freezing_temperature": -5.0,
					"product_type_id": 3,
					"seller_id": 1
				}
			}`,
		},
		{
			testName:          "Error case: JSON with missing fields",
			payload:           payloadWithMissingFields,
			isPayloadError:    true,
			productAttributes: models.ProductAttributes{},
			serviceData:       models.Product{},
			serviceError:      nil,
			expectedCode:      http.StatusUnprocessableEntity,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:          "Error case: Wrong value in a JSON field",
			payload:           payloadWithWrongValues,
			isPayloadError:    true,
			productAttributes: models.ProductAttributes{},
			serviceData:       models.Product{},
			serviceError:      nil,
			expectedCode:      http.StatusUnprocessableEntity,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:          "Error case: Unknown JSON fields",
			payload:           payloadWithUnkownFields,
			isPayloadError:    true,
			productAttributes: models.ProductAttributes{},
			serviceData:       models.Product{},
			serviceError:      nil,
			expectedCode:      http.StatusBadRequest,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:          "Error case: Process an error from the service layer",
			payload:           payload,
			isPayloadError:    false,
			productAttributes: newProductAttributes,
			serviceData:       models.Product{},
			serviceError:      errors.New("db error"),
			expectedCode:      http.StatusInternalServerError,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
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
			serviceMock := &mocks.ProductServiceMock{}

			// If a validation error occurs the service method is not called
			if !tc.isPayloadError {
				serviceMock.
					On("Create", mock.Anything, tc.productAttributes).
					Return(tc.serviceData, tc.serviceError)
			}

			handler := handler.NewProductHandler(serviceMock)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// Act
			handler.Create()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedHeaders, response.Header())
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			serviceMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the HTTP handler responsible for retrieving all products.
// It covers:
// - Successful retrieval of multiple products
// - Successful retrieval when no products exist (empty list)
// - Error propagation from the service layer
func TestProductHandler_GetAll(t *testing.T) {
	// Define the products used in common by the test cases
	product1 := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Yogurt helado",
			ExpirationRate:                 7,
			FreezingRate:                   2,
			Height:                         10.5,
			Length:                         20.0,
			Width:                          15.0,
			NetWeight:                      1.2,
			ProductCode:                    "YOG01",
			RecommendedFreezingTemperature: -5.0,
			ProductTypeID:                  3,
			SellerID:                       utils.Ptr(1),
		},
	}

	product2 := models.Product{
		ID: 2,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Pechuga de pollo",
			ExpirationRate:                 3,
			FreezingRate:                   1,
			Height:                         5.0,
			Length:                         25.0,
			Width:                          12.5,
			NetWeight:                      0.8,
			ProductCode:                    "POLLO01",
			RecommendedFreezingTemperature: 0.0,
			ProductTypeID:                  7,
			SellerID:                       utils.Ptr(2),
		},
	}

	// Each test case is constructed by:
	// testName            — human‐readable description
	// serviceData         — the Products slice returned by the mocked service
	// serviceError        — the error returned by the mocked service
	// expectedCode        — HTTP status code we expect the handler to produce
	// expectedHeaders     — HTTP headers we expect in the HTTP response
	// expectedBody        — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName        string
		serviceData     []models.Product
		serviceError    error
		expectedCode    int
		expectedHeaders http.Header
		expectedBody    string
	}{
		{
			testName:        "Success: Get all products",
			serviceData:     []models.Product{product1, product2},
			serviceError:    nil,
			expectedCode:    http.StatusOK,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": [
					{
						"id": 1,
						"description": "Yogurt helado",
						"expiration_rate": 7,
						"freezing_rate": 2,
						"height": 10.5,
						"length": 20.0,
						"width": 15.0,
						"netweight": 1.2,
						"product_code": "YOG01",
						"recommended_freezing_temperature": -5.0,
						"product_type_id": 3,
						"seller_id": 1
					},
					{
						"id": 2,
						"description": "Pechuga de pollo",
						"expiration_rate": 3,
						"freezing_rate": 1,
						"height": 5.0,
						"length": 25.0,
						"width": 12.5,
						"netweight": 0.8,
						"product_code": "POLLO01",
						"recommended_freezing_temperature": 0.0,
						"product_type_id": 7,
						"seller_id": 2
					}
				]
			}`,
		},
		{
			testName:        "Success: Get an empty list if the DB is empty",
			serviceData:     []models.Product{},
			serviceError:    nil,
			expectedCode:    http.StatusOK,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"data": []
				}
			`,
		},
		{
			testName:        "Error case: Process an error from the service layer",
			serviceData:     nil,
			serviceError:    errors.New("db error"),
			expectedCode:    http.StatusInternalServerError,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
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
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.On("GetAll", mock.Anything).Return(tc.serviceData, tc.serviceError)
			handler := handler.NewProductHandler(serviceMock)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// Act
			handler.GetAll()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedHeaders, response.Header())
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			serviceMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the HTTP handler responsible for retrieving a single product.
// It covers:
// - Successful retrieval of a single product
// - Error when an invalid ID is given
// - Error propagation from the service layer
func TestProductHandler_GetById(t *testing.T) {
	// Define the product used in common by the test cases
	product := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Yogurt helado",
			ExpirationRate:                 7,
			FreezingRate:                   2,
			Height:                         10.5,
			Length:                         20.0,
			Width:                          15.0,
			NetWeight:                      1.2,
			ProductCode:                    "YOG01",
			RecommendedFreezingTemperature: -5.0,
			ProductTypeID:                  3,
			SellerID:                       utils.Ptr(1),
		},
	}

	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to retrieve
	// isIdError           — whether we expect ID validation to fail inside the handler
	// serviceData         — the Product object returned by the mocked service
	// serviceError        — the error returned by the mocked service
	// expectedCode        — HTTP status code we expect the handler to produce
	// expectedHeaders     — HTTP headers we expect in the HTTP response
	// expectedBody        — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName        string
		id              int
		isIdError       bool
		serviceData     models.Product
		serviceError    error
		expectedCode    int
		expectedHeaders http.Header
		expectedBody    string
	}{
		{
			testName:        "Success: Get product with ID 1",
			id:              1,
			isIdError:       false,
			serviceData:     product,
			serviceError:    nil,
			expectedCode:    http.StatusOK,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": {
					"id": 1,
					"description": "Yogurt helado",
					"expiration_rate": 7,
					"freezing_rate": 2,
					"height": 10.5,
					"length": 20.0,
					"width": 15.0,
					"netweight": 1.2,
					"product_code": "YOG01",
					"recommended_freezing_temperature": -5.0,
					"product_type_id": 3,
					"seller_id": 1
				}
			}`,
		},
		{
			testName:        "Error case: Invalid ID is given",
			id:              -1,
			isIdError:       true,
			serviceData:     models.Product{},
			serviceError:    httperrors.NotFoundError{Message: "Invalid ID"},
			expectedCode:    http.StatusBadRequest,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid ID"
				}
			`,
		},
		{
			testName:        "Error case: Process an error from the service layer",
			id:              1,
			isIdError:       false,
			serviceData:     models.Product{},
			serviceError:    errors.New("db error"),
			expectedCode:    http.StatusInternalServerError,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
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
			serviceMock := &mocks.ProductServiceMock{}

			// If a validation error occurs the service method is not called
			if !tc.isIdError {
				serviceMock.
					On("GetByID", mock.Anything, tc.id).
					Return(tc.serviceData, tc.serviceError)
			}
			handler := handler.NewProductHandler(serviceMock)

			url := fmt.Sprintf("/api/v1/products/%d", tc.id)
			request := httptest.NewRequest(http.MethodGet, url, nil)

			// Create chi context to pass the ID to the handler test
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", strconv.Itoa(tc.id))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()

			// Act
			handler.GetByID()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedHeaders, response.Header())
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			serviceMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the HTTP handler responsible for retrieving records per product.
// It covers:
// - Successful retrieval of the record count of every product
// - Successful retrieval of the record count of a single product
// - Error when an invalid ID is given
// - Error propagation from the service layer
func TestProductHandler_GetRecordsPerProduct(t *testing.T) {
	// Define the records used in common by the test cases
	record1 := models.ProductRecordCount{
		ProductID:    1,
		Description:  "Yogurt helado",
		RecordsCount: 3,
	}

	record2 := models.ProductRecordCount{
		ProductID:    2,
		Description:  "Pechuga de pollo",
		RecordsCount: 1,
	}

	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to retrieve as a pointer
	// isIdError           — whether we expect ID validation to fail inside the handler
	// serviceData         — the records per product slice returned by the mocked service
	// serviceError        — the error returned by the mocked service
	// expectedCode        — HTTP status code we expect the handler to produce
	// expectedHeaders     — HTTP headers we expect in the HTTP response
	// expectedBody        — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName        string
		id              *int
		isIdError       bool
		serviceData     []models.ProductRecordCount
		serviceError    error
		expectedCode    int
		expectedHeaders http.Header
		expectedBody    string
	}{
		{
			testName:        "Success: Get all product if no ID query param is given",
			id:              nil,
			isIdError:       false,
			serviceData:     []models.ProductRecordCount{record1, record2},
			serviceError:    nil,
			expectedCode:    http.StatusOK,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": [
					{
						"product_id": 1,
						"description": "Yogurt helado",
						"records_count": 3
					},
					{
						"product_id": 2,
						"description": "Pechuga de pollo",
						"records_count": 1
					}
				]
			}`,
		},
		{
			testName:        "Success: Get a single product if ID query param is given",
			id:              utils.Ptr(2),
			isIdError:       false,
			serviceData:     []models.ProductRecordCount{record2},
			serviceError:    nil,
			expectedCode:    http.StatusOK,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": [
					{
						"product_id": 2,
						"description": "Pechuga de pollo",
						"records_count": 1
					}
				]
			}`,
		},
		{
			testName:        "Error case: Invalid ID is given",
			id:              utils.Ptr(-1),
			isIdError:       true,
			serviceData:     nil,
			serviceError:    errors.New("invalid id"),
			expectedCode:    http.StatusBadRequest,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"status": "Bad Request",
				"message": "invalid id"
			}`,
		},
		{
			testName:        "Error case: Process an error from the service layer",
			id:              utils.Ptr(1),
			isIdError:       false,
			serviceData:     nil,
			serviceError:    errors.New("db error"),
			expectedCode:    http.StatusInternalServerError,
			expectedHeaders: http.Header{"Content-Type": []string{"application/json"}},
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
			serviceMock := &mocks.ProductServiceMock{}

			// If a validation error occurs the service method is not called
			if !tc.isIdError {
				serviceMock.
					On("GetRecordsPerProduct", mock.Anything, tc.id).
					Return(tc.serviceData, tc.serviceError)
			}
			handler := handler.NewProductHandler(serviceMock)

			url := "/api/v1/products/reportRecords"
			fmt.Println(url)
			if tc.id != nil {
				url += fmt.Sprintf("?id=%d", *tc.id)
			}

			request := httptest.NewRequest(http.MethodGet, url, nil)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// Act
			handler.GetRecordsPerProduct()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedHeaders, response.Header())
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			serviceMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the HTTP handler responsible for updating a product.
// It covers:
// - Successful update all fields of a product
// - Successful update a single field of a product
// - Error when an invalid ID is given
// - Error when invalid JSON fields values are given
// - Error when unknown JSON fields are given
// - Error propagation from the service layer
func TestProductHandler_Update(t *testing.T) {
	// Define the products and payloads used in common by the test cases
	productAttributes := models.ProductPatchRequest{
		Description:                    utils.Ptr("Pechuga de pollo"),
		ExpirationRate:                 utils.Ptr(6),
		FreezingRate:                   utils.Ptr(3),
		Height:                         utils.Ptr(11.5),
		Length:                         utils.Ptr(22.0),
		Width:                          utils.Ptr(13.0),
		NetWeight:                      utils.Ptr(10.0),
		ProductCode:                    utils.Ptr("POL01"),
		RecommendedFreezingTemperature: utils.Ptr(-3.0),
		ProductTypeID:                  utils.Ptr(2),
		SellerID:                       utils.Ptr(2),
	}

	singleProductAttribute := models.ProductPatchRequest{
		NetWeight: utils.Ptr(10.0),
	}

	updatedProduct := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Pechuga de pollo",
			ExpirationRate:                 6,
			FreezingRate:                   3,
			Height:                         11.5,
			Length:                         22.0,
			Width:                          13.0,
			NetWeight:                      10.0,
			ProductCode:                    "POL01",
			RecommendedFreezingTemperature: -3.0,
			ProductTypeID:                  2,
			SellerID:                       utils.Ptr(2),
		},
	}

	singleFieldUpdatedProduct := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Yogurt helado",
			ExpirationRate:                 7,
			FreezingRate:                   2,
			Height:                         10.5,
			Length:                         20.0,
			Width:                          15.0,
			NetWeight:                      10.0,
			ProductCode:                    "YOG01",
			RecommendedFreezingTemperature: -5.0,
			ProductTypeID:                  3,
			SellerID:                       utils.Ptr(1),
		},
	}

	payloadUpdatedProduct := `{
		"description": "Pechuga de pollo",
		"expiration_rate": 6,
		"freezing_rate": 3,
		"height": 11.5,
		"length": 22.0,
		"width": 13.0,
		"netweight": 10.0,
		"product_code": "POL01",
		"recommended_freezing_temperature": -3.0,
		"product_type_id": 2,
		"seller_id": 2
	}`

	payloadSingleFieldUpdatedProduct := `{
		"netweight": 10.0
	}`

	payloadWithUnkownFields := `{
		"description": "Pechuga de pollo",
		"anUnkownField": 1,
		"expiration_rate": 6,
		"freezing_rate": 3,
		"height": 11.5,
		"length": 22.0,
		"width": 13.0,
		"netweight": 10.0,
		"product_code": "POL01",
		"recommended_freezing_temperature": -3.0,
		"product_type_id": 2,
		"seller_id": 2
	}`

	payloadWithInvalidField := `{
		"description": "Pechuga de pollo",
		"expiration_rate": -1,
		"freezing_rate": -3,
		"height": -11.5,
		"length": -22.0,
		"width": -13.0,
		"netweight": 10.0,
		"product_code": "POL01",
		"recommended_freezing_temperature": -3.0,
		"product_type_id": 2,
		"seller_id": -2
	}`

	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to update
	// isValidationError   — whether we expect validation to fail inside the handler
	// serviceData         — the product returned by the mocked service
	// serviceError        — the error returned by the mocked service
	// expectedCode        — HTTP status code we expect the handler to produce
	// expectedHeaders     — HTTP headers we expect in the HTTP response
	// expectedBody        — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName          string
		payload           string
		productAttributes models.ProductPatchRequest
		id                int
		isValidationError bool
		serviceData       models.Product
		serviceError      error
		expectedCode      int
		expectedHeaders   http.Header
		expectedBody      string
	}{
		{
			testName:          "Success: Update all fields of product with ID 1",
			payload:           payloadUpdatedProduct,
			productAttributes: productAttributes,
			id:                1,
			isValidationError: false,
			serviceData:       updatedProduct,
			serviceError:      nil,
			expectedCode:      http.StatusOK,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": {
					"id": 1,
					"description": "Pechuga de pollo",
					"expiration_rate": 6,
					"freezing_rate": 3,
					"height": 11.5,
					"length": 22.0,
					"width": 13.0,
					"netweight": 10.0,
					"product_code": "POL01",
					"recommended_freezing_temperature": -3.0,
					"product_type_id": 2,
					"seller_id": 2
				}
			}`,
		},
		{
			testName:          "Success: Update a single field of product with ID 1",
			payload:           payloadSingleFieldUpdatedProduct,
			productAttributes: singleProductAttribute,
			id:                1,
			isValidationError: false,
			serviceData:       singleFieldUpdatedProduct,
			serviceError:      nil,
			expectedCode:      http.StatusOK,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"data": {
					"id": 1,
					"description": "Yogurt helado",
					"expiration_rate": 7,
					"freezing_rate": 2,
					"height": 10.5,
					"length": 20.0,
					"width": 15.0,
					"netweight": 10.0,
					"product_code": "YOG01",
					"recommended_freezing_temperature": -5.0,
					"product_type_id": 3,
					"seller_id": 1
				}
			}`,
		},
		{
			testName:          "Error case: Invalid ID is given",
			payload:           payloadUpdatedProduct,
			productAttributes: productAttributes,
			id:                -1,
			isValidationError: true,
			serviceData:       models.Product{},
			serviceError:      httperrors.NotFoundError{Message: "Invalid ID"},
			expectedCode:      http.StatusBadRequest,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"status": "Bad Request",
				"message": "Invalid ID"
			}
			`,
		},
		{
			testName:          "Error case: Unknown JSON fields",
			payload:           payloadWithUnkownFields,
			productAttributes: productAttributes,
			id:                1,
			isValidationError: true,
			serviceData:       models.Product{},
			serviceError:      httperrors.BadRequestError{Message: "Invalid JSON body"},
			expectedCode:      http.StatusBadRequest,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"status": "Bad Request",
				"message": "Invalid JSON body"
			}
			`,
		},
		{
			testName:          "Error case: Invalid values in JSON fields",
			payload:           payloadWithInvalidField,
			productAttributes: productAttributes,
			id:                1,
			isValidationError: true,
			serviceData:       models.Product{},
			serviceError:      httperrors.UnprocessableEntityError{Message: "Invalid JSON body"},
			expectedCode:      http.StatusUnprocessableEntity,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
			expectedBody: `
			{
				"status": "Unprocessable Entity",
				"message": "Invalid JSON body"
			}
			`,
		},
		{
			testName:          "Error case: Process an error from the service layer",
			payload:           payloadUpdatedProduct,
			productAttributes: productAttributes,
			id:                1,
			isValidationError: false,
			serviceData:       models.Product{},
			serviceError:      errors.New("db error"),
			expectedCode:      http.StatusInternalServerError,
			expectedHeaders:   http.Header{"Content-Type": []string{"application/json"}},
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
			serviceMock := &mocks.ProductServiceMock{}

			// If a validation error occurs the service method is not called
			if !tc.isValidationError {
				serviceMock.
					On("Update", mock.Anything, tc.id, tc.productAttributes).
					Return(tc.serviceData, tc.serviceError)
			}
			handler := handler.NewProductHandler(serviceMock)

			url := fmt.Sprintf("/api/v1/products/%d", tc.id)
			request := httptest.NewRequest(http.MethodPatch, url, strings.NewReader(tc.payload))

			// Create chi context to pass the ID to the handler test
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", strconv.Itoa(tc.id))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()

			// Act
			handler.Update()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedHeaders, response.Header())
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			serviceMock.AssertExpectations(t)
		})
	}
}

// Verifies the behavior of the HTTP handler responsible for deleting a product.
// It covers:
// - Successful deletion of a product
// - Error when an invalid ID is given
// - Error propagation from the service layer
func TestProductHandler_Delete(t *testing.T) {
	// Each test case is constructed by:
	// testName            — human‐readable description
	// id                  - ID of the product to delete
	// isIdError           — whether we expect ID validation to fail inside the handler
	// serviceError        — the error returned by the mocked service
	// expectedCode        — HTTP status code we expect the handler to produce
	// expectedBody        — JSON body (string) we expect in the HTTP response
	tests := []struct {
		testName     string
		id           int
		isIdError    bool
		serviceError error
		expectedCode int
		expectedBody string
	}{
		{
			testName:     "Success: Delete product with ID 1",
			id:           1,
			isIdError:    false,
			serviceError: nil,
			expectedCode: http.StatusNoContent,
			expectedBody: "null",
		},
		{
			testName:     "Error case: invalid ID is given",
			id:           -1,
			isIdError:    true,
			serviceError: httperrors.NotFoundError{Message: "Invalid ID"},
			expectedCode: http.StatusBadRequest,
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid ID"
				}
			`,
		},
		{
			testName:     "Error case: Process an error from the service layer",
			id:           1,
			isIdError:    false,
			serviceError: errors.New("db error"),
			expectedCode: http.StatusInternalServerError,
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
			serviceMock := &mocks.ProductServiceMock{}

			// If a validation error occurs the service method is not called
			if !tc.isIdError {
				serviceMock.
					On("Delete", mock.Anything, tc.id).
					Return(tc.serviceError)
			}
			handler := handler.NewProductHandler(serviceMock)

			url := fmt.Sprintf("/api/v1/products/%d", tc.id)
			request := httptest.NewRequest(http.MethodDelete, url, nil)

			// Create chi context to pass the ID to the handler test
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", strconv.Itoa(tc.id))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

			response := httptest.NewRecorder()

			// Act
			handler.Delete()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			if tc.expectedCode == http.StatusNoContent {
				require.Empty(t, response.Body.String())
			} else {
				require.JSONEq(t, tc.expectedBody, response.Body.String())
			}
			serviceMock.AssertExpectations(t)
		})
	}
}
