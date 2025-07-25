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
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: Move this function to an utils package
func Ptr[T any](v T) *T { return &v }

func TestProductHandler_Create(t *testing.T) {
	// TODO: Generate random products with values in a valid range
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
		SellerID:                       Ptr(1),
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

	tests := []struct {
		testName     string
		serviceData  models.Product
		serviceError error
		payload      string
		expectedCode int
		expectedBody string
	}{
		{
			testName:     "Success: Create a new product",
			serviceData:  newProduct,
			serviceError: nil,
			payload:      payload,
			expectedCode: http.StatusCreated,
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
			testName:     "Fail: Unprocessable entity when payload with missing fields is given",
			serviceData:  models.Product{},
			serviceError: httperrors.ConflictError{Message: "A product with the given product code already exists"},
			payload:      payloadWithMissingFields,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:     "Fail: Unprocessable entity when payload with wrong values is given",
			serviceData:  models.Product{},
			serviceError: httperrors.ConflictError{Message: "A product with the given product code already exists"},
			payload:      payloadWithWrongValues,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `
				{
					"status": "Unprocessable Entity",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:     "Fail: Bad request when payload with unkown fields is given",
			serviceData:  models.Product{},
			serviceError: httperrors.ConflictError{Message: "A product with the given product code already exists"},
			payload:      payloadWithUnkownFields,
			expectedCode: http.StatusBadRequest,
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid JSON body"
				}
			`,
		},
		{
			testName:     "Fail: Conflict error when the given product code already exists",
			serviceData:  models.Product{},
			serviceError: httperrors.ConflictError{Message: "A product with the given product code already exists"},
			payload:      payload,
			expectedCode: http.StatusConflict,
			expectedBody: `
				{
					"status": "Conflict",
					"message": "A product with the given product code already exists"
				}
			`,
		},
		{
			testName:     "Fail: Internal server error after a DB Error",
			serviceData:  models.Product{},
			serviceError: errors.New("db error"),
			payload:      payload,
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
			// Run tests parallel
			t.Parallel()

			// Arrange
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.On("Create", mock.Anything, newProductAttributes).Return(tc.serviceData, tc.serviceError)
			handler := handler.NewProductHandler(serviceMock)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
			response := httptest.NewRecorder()

			// Act
			handler.Create()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}

func TestProductHandler_GetAll(t *testing.T) {
	// TODO: Generate random products with values in a valid range
	p1 := models.Product{
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
			SellerID:                       Ptr(1),
		},
	}

	p2 := models.Product{
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
			SellerID:                       Ptr(2),
		},
	}

	tests := []struct {
		testName     string
		serviceData  []models.Product
		serviceError error
		expectedCode int
		expectedBody string
	}{
		{
			testName:     "Success: Get all products",
			serviceData:  []models.Product{p1, p2},
			serviceError: nil,
			expectedCode: http.StatusOK,
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
			testName:     "Success: Get an empty list if the DB is empty",
			serviceData:  []models.Product{},
			serviceError: nil,
			expectedCode: http.StatusOK,
			expectedBody: `
				{
					"data": []
				}
			`,
		},
		{
			testName:     "Fail: Internal server error after a DB Error",
			serviceData:  nil,
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
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Run tests parallel
			t.Parallel()

			// Arrange
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.On("GetAll", mock.Anything).Return(tt.serviceData, tt.serviceError)
			handler := handler.NewProductHandler(serviceMock)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			// Act
			handler.GetAll()(response, request)

			// Assert
			require.Equal(t, tt.expectedCode, response.Code)
			require.JSONEq(t, tt.expectedBody, response.Body.String())
		})
	}
}

func TestProductHandler_GetById(t *testing.T) {
	// TODO: Generate random products with values in a valid range
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
			SellerID:                       Ptr(1),
		},
	}

	tests := []struct {
		testName     string
		serviceData  models.Product
		serviceError error
		idParam      int
		expectedCode int
		expectedBody string
	}{
		{
			testName:     "Success: Get product with ID 1",
			serviceData:  product,
			serviceError: nil,
			idParam:      1,
			expectedCode: http.StatusOK,
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
			testName:     "Fail: Not found when giving a non existant ID",
			serviceData:  models.Product{},
			serviceError: httperrors.NotFoundError{Message: "Product not found"},
			idParam:      10000,
			expectedCode: http.StatusNotFound,
			expectedBody: `
				{
					"status": "Not Found",
					"message": "Product not found"
				}
			`,
		},
		{
			testName:     "Fail: Bad request when giving an invalid ID",
			serviceData:  models.Product{},
			serviceError: httperrors.NotFoundError{Message: "Invalid ID"},
			idParam:      -1,
			expectedCode: http.StatusBadRequest,
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid ID"
				}
			`,
		},
		{
			testName:     "Fail: Internal server error after a DB Error",
			serviceData:  models.Product{},
			serviceError: errors.New("db error"),
			idParam:      1,
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
			// Run tests parallel
			t.Parallel()

			// Arrange
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.On("GetByID", mock.Anything, tc.idParam).Return(tc.serviceData, tc.serviceError)
			handler := handler.NewProductHandler(serviceMock)

			url := fmt.Sprintf("/api/v1/products/%d", tc.idParam)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", strconv.Itoa(tc.idParam))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

			response := httptest.NewRecorder()

			// Act
			handler.GetByID()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}

func TestProductHandler_GetRecordsPerProduct(t *testing.T) {
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

	tests := []struct {
		testName     string
		serviceData  []models.ProductRecordCount
		serviceError error
		idParam      *int
		expectedCode int
		expectedBody string
	}{
		{
			testName:     "Success: Get all product if no ID query param is given",
			serviceData:  []models.ProductRecordCount{record1, record2},
			serviceError: nil,
			idParam:      nil,
			expectedCode: http.StatusOK,
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
			testName:     "Success: Get a single product if ID query param is given",
			serviceData:  []models.ProductRecordCount{record2},
			serviceError: nil,
			idParam:      Ptr(2),
			expectedCode: http.StatusOK,
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
			testName:     "Fail: Invalid ID is given",
			serviceData:  nil,
			serviceError: errors.New("invalid id"),
			idParam:      Ptr(-1),
			expectedCode: http.StatusBadRequest,
			expectedBody: `
			{
				"status": "Bad Request",
				"message": "invalid id"
			}`,
		},
		{
			testName:     "Fail: Internal server error after a DB Error",
			serviceData:  nil,
			serviceError: errors.New("db error"),
			idParam:      Ptr(1),
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
			// Run tests parallel
			t.Parallel()

			// Arrange
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.On("GetRecordsPerProduct", mock.Anything, tc.idParam).Return(tc.serviceData, tc.serviceError)
			handler := handler.NewProductHandler(serviceMock)

			url := "/api/v1/products/reportRecords"
			fmt.Println(url)
			if tc.idParam != nil {
				url += fmt.Sprintf("?id=%d", *tc.idParam)
			}

			request := httptest.NewRequest(http.MethodGet, url, nil)
			response := httptest.NewRecorder()

			// Act
			handler.GetRecordsPerProduct()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}

func TestProductHandler_Update(t *testing.T) {
	productAttributes := models.ProductPatchRequest{
		Description:                    Ptr("Pechuga de pollo"),
		ExpirationRate:                 Ptr(6),
		FreezingRate:                   Ptr(3),
		Height:                         Ptr(11.5),
		Length:                         Ptr(22.0),
		Width:                          Ptr(13.0),
		NetWeight:                      Ptr(10.0),
		ProductCode:                    Ptr("POL01"),
		RecommendedFreezingTemperature: Ptr(-3.0),
		ProductTypeID:                  Ptr(2),
		SellerID:                       Ptr(2),
	}

	singleProductAttribute := models.ProductPatchRequest{
		NetWeight: Ptr(10.0),
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
			SellerID:                       Ptr(2),
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
			SellerID:                       Ptr(1),
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

	tests := []struct {
		testName          string
		serviceData       models.Product
		serviceError      error
		payload           string
		idParam           int
		productAttributes models.ProductPatchRequest
		expectedCode      int
		expectedBody      string
	}{
		{
			testName:          "Success: Update all fields of product with ID 1",
			serviceData:       updatedProduct,
			serviceError:      nil,
			payload:           payloadUpdatedProduct,
			idParam:           1,
			productAttributes: productAttributes,
			expectedCode:      http.StatusOK,
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
			serviceData:       singleFieldUpdatedProduct,
			serviceError:      nil,
			payload:           payloadSingleFieldUpdatedProduct,
			idParam:           1,
			productAttributes: singleProductAttribute,
			expectedCode:      http.StatusOK,
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
			testName:          "Fail: Not found when giving a non existant ID",
			serviceData:       models.Product{},
			serviceError:      httperrors.NotFoundError{Message: "Product not found"},
			payload:           payloadUpdatedProduct,
			idParam:           10000,
			productAttributes: productAttributes,
			expectedCode:      http.StatusNotFound,
			expectedBody: `
			{
				"status": "Not Found",
				"message": "Product not found"
			}
			`,
		},
		{
			testName:          "Fail: Bad request when giving an invalid ID",
			serviceData:       models.Product{},
			serviceError:      httperrors.NotFoundError{Message: "Invalid ID"},
			payload:           payloadUpdatedProduct,
			idParam:           -1,
			productAttributes: productAttributes,
			expectedCode:      http.StatusBadRequest,
			expectedBody: `
			{
				"status": "Bad Request",
				"message": "Invalid ID"
			}
			`,
		},
		{
			testName:          "Fail: Bad request when body contains unknown fields",
			serviceData:       models.Product{},
			serviceError:      httperrors.BadRequestError{Message: "Invalid JSON body"},
			payload:           payloadWithUnkownFields,
			idParam:           1,
			productAttributes: productAttributes,
			expectedCode:      http.StatusBadRequest,
			expectedBody: `
			{
				"status": "Bad Request",
				"message": "Invalid JSON body"
			}
			`,
		},
		{
			testName:          "Fail: Bad request when body contains an invalid value in some field",
			serviceData:       models.Product{},
			serviceError:      httperrors.UnprocessableEntityError{Message: "Invalid JSON body"},
			payload:           payloadWithInvalidField,
			idParam:           1,
			productAttributes: productAttributes,
			expectedCode:      http.StatusUnprocessableEntity,
			expectedBody: `
			{
				"status": "Unprocessable Entity",
				"message": "Invalid JSON body"
			}
			`,
		},
		{
			testName:          "Fail: Internal server error after a DB Error",
			serviceData:       models.Product{},
			serviceError:      errors.New("db error"),
			payload:           payloadUpdatedProduct,
			idParam:           1,
			productAttributes: productAttributes,
			expectedCode:      http.StatusInternalServerError,
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
			// Run tests parallel
			t.Parallel()

			// Arrange
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.
				On("Update", mock.Anything, tc.idParam, tc.productAttributes).
				Return(tc.serviceData, tc.serviceError)
			handler := handler.NewProductHandler(serviceMock)

			url := fmt.Sprintf("/api/v1/products/%d", tc.idParam)
			request := httptest.NewRequest(http.MethodPatch, url, strings.NewReader(tc.payload))
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", strconv.Itoa(tc.idParam))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

			response := httptest.NewRecorder()

			// Act
			handler.Update()(response, request)

			// Assert
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	tests := []struct {
		testName     string
		serviceError error
		idParam      int
		expectedCode int
		expectedBody string
	}{
		{
			testName:     "Success: Delete product with ID 1",
			serviceError: nil,
			idParam:      1,
			expectedCode: http.StatusNoContent,
			expectedBody: "null",
		},
		{
			testName:     "Fail: Not found when giving a non existant ID",
			serviceError: httperrors.NotFoundError{Message: "Product not found"},
			idParam:      10000,
			expectedCode: http.StatusNotFound,
			expectedBody: `
				{
					"status": "Not Found",
					"message": "Product not found"
				}
			`,
		},
		{
			testName:     "Fail: Bad request when giving an invalid ID",
			serviceError: httperrors.NotFoundError{Message: "Invalid ID"},
			idParam:      -1,
			expectedCode: http.StatusBadRequest,
			expectedBody: `
				{
					"status": "Bad Request",
					"message": "Invalid ID"
				}
			`,
		},
		{
			testName:     "Fail: Internal server error after a DB Error",
			serviceError: errors.New("db error"),
			idParam:      1,
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
			// Run tests parallel
			t.Parallel()

			// Arrange
			serviceMock := &mocks.ProductServiceMock{}
			serviceMock.On("Delete", mock.Anything, tc.idParam).Return(tc.serviceError)
			handler := handler.NewProductHandler(serviceMock)

			url := fmt.Sprintf("/api/v1/products/%d", tc.idParam)
			request := httptest.NewRequest(http.MethodDelete, url, nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", strconv.Itoa(tc.idParam))
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
		})
	}
}
