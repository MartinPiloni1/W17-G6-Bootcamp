// Package handler_test provides comprehensive testing for the warehouse HTTP handlers.
//
// This test suite covers all CRUD operations for the warehouse resource,
// ensuring proper HTTP status codes, response formats, and error handling.
// Tests use mocked services to isolate handler logic from business logic.
//
// Test Categories:
//   - Create: Tests warehouse creation with various input scenarios
//   - GetAll: Tests retrieval of all warehouses
//   - GetById: Tests retrieval of specific warehouses by ID
//   - Update: Tests warehouse modification operations
//   - Delete: Tests warehouse deletion operations
//
// Each test follows the AAA pattern (Arrange, Act, Assert) and uses testify
// for assertions and mocking.
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/mock/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// expectedContentType defines the expected Content-Type header for all warehouse API responses
const expectedContentType = "application/json"

// Test_Create tests the Create method of the WarehouseHandler.
//
// This test suite verifies that the warehouse creation endpoint:
//   - Successfully creates warehouses with valid input data
//   - Properly handles malformed JSON requests
//   - Returns appropriate error responses for business logic conflicts
//   - Sets correct HTTP status codes and headers
//
// Test scenarios:
//   - Success case: Valid warehouse data returns 201 Created
//   - Invalid JSON: Malformed request body returns 422 Unprocessable Entity
//   - Conflict: Duplicate warehouse code returns 409 Conflict
func Test_Create(t *testing.T) {
	// Success scenario: warehouse creation with valid data
	t.Run("on success (create_ok)", func(t *testing.T) {
		// Arrange: prepare valid warehouse input and expected output
		warehouseInput := models.WarehouseAttributes{
			WarehouseCode:      "WH-001",
			Address:            "Calle Falsa 123",
			Telephone:          "123456789",
			MinimunCapacity:    100,
			MinimunTemperature: 4.5,
		}

		warehuseUotput := models.Warehouse{
			Id:                  1,
			WarehouseAttributes: warehouseInput,
		}

		// Mock service to return successful creation
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("Create", warehouseInput).
			Return(warehuseUotput, nil)
		handlerTest := handler.NewWarehouseHandler(mockService)
		body, err := json.Marshal(warehouseInput)
		require.NoError(t, err)

		// Act: execute HTTP request to create warehouse
		req := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handlerTest.Create()(res, req)

		// Assert: verify response status, body, and headers
		expectedStatus := http.StatusCreated
		expectedBody, err := json.Marshal(map[string]any{
			"data": warehuseUotput,
		})

		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "Create", 1)
	})

	// Error scenario: malformed JSON request body
	t.Run("Invalid JSON body (create_fail)", func(t *testing.T) {
		// Arrange: prepare invalid JSON with wrong data type
		mockService := new(service_mocks.WarehouseServiceMock)
		handlerTest := handler.NewWarehouseHandler(mockService)
		bodyInput := map[string]any{
			"minimun_temperature": "this is not a float", // Invalid: string instead of float
		}
		body, err := json.Marshal(bodyInput)
		require.NoError(t, err)

		// Act: execute request with malformed JSON
		req := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewReader(body))
		res := httptest.NewRecorder()
		handlerTest.Create()(res, req)

		// Assert: verify 422 status and error message, service should not be called
		expectedStatus := http.StatusUnprocessableEntity
		expectedBody := `{"status":"Unprocessable Entity","message":"Invalid JSON body"}`

		require.NoError(t, err)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "Create")
	})

	// Error scenario: business logic conflict (duplicate warehouse code)
	t.Run("On service error (create_conflict)", func(t *testing.T) {
		// Arrange: prepare valid input but simulate service conflict error
		warehouseInput := models.WarehouseAttributes{
			WarehouseCode:      "WH-001",
			Address:            "Calle Falsa 123",
			Telephone:          "123456789",
			MinimunCapacity:    100,
			MinimunTemperature: 4.5,
		}
		body, err := json.Marshal(warehouseInput)
		require.NoError(t, err)

		// Mock service to return conflict error
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("Create", warehouseInput).
			Return(models.Warehouse{}, httperrors.ConflictError{Message: "the WarehouseCode already exists"})
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute request that will result in conflict
		req := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handlerTest.Create()(res, req)

		// Assert: verify 409 status and conflict error message
		expectedBody := `{"status":"Conflict","message":"the WarehouseCode already exists"}`
		expectedStatus := http.StatusConflict

		require.NoError(t, err)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "Create", 1)
	})
}

// Test_GetAll tests the GetAll method of the WarehouseHandler.
//
// This test suite verifies that the warehouse listing endpoint:
//   - Successfully returns all warehouses when they exist
//   - Properly handles service errors during data retrieval
//   - Returns appropriate HTTP status codes and response format
//
// Test scenarios:
//   - Success case: Multiple warehouses returned with 200 OK
//   - Service error: Internal server error returns 500 status
func Test_GetAll(t *testing.T) {
	// Success scenario: retrieve multiple warehouses
	t.Run("on success (find_all)", func(t *testing.T) {
		// Arrange: prepare test warehouses data
		warehouseTest1 := models.Warehouse{
			Id: 1,
			WarehouseAttributes: models.WarehouseAttributes{
				WarehouseCode:      "WH-001",
				Address:            "Calle Falsa 123",
				Telephone:          "123456789",
				MinimunCapacity:    100,
				MinimunTemperature: 4.5,
			},
		}
		warehouseTest2 := models.Warehouse{
			Id: 2,
			WarehouseAttributes: models.WarehouseAttributes{
				WarehouseCode:      "WH-002",
				Address:            "Calle 123343",
				Telephone:          "16820823",
				MinimunCapacity:    90,
				MinimunTemperature: 4.1,
			},
		}

		// Mock service to return warehouse collection
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("GetAll").
			Return([]models.Warehouse{
				warehouseTest1,
				warehouseTest2,
			}, nil)
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute GET request for all warehouses
		req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
		res := httptest.NewRecorder()
		handlerTest.GetAll()(res, req)

		// Assert: verify successful response with warehouse collection
		expectedStatus := http.StatusOK
		expectedBody, err := json.Marshal(map[string]any{
			"data": []models.Warehouse{
				warehouseTest1,
				warehouseTest2,
			},
		})
		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "GetAll", 1)
	})

	// Error scenario: service fails during warehouse retrieval
	t.Run("on service error", func(t *testing.T) {
		// Arrange: mock service to return internal server error
		serviceError := httperrors.InternalServerError{Message: "error obtaining warehouses"}
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("GetAll").
			Return([]models.Warehouse{}, serviceError)
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute request that will result in service error
		req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
		res := httptest.NewRecorder()
		handlerTest.GetAll()(res, req)

		// Assert: verify 500 status and error message
		expectedStatus := http.StatusInternalServerError
		expectedBody := `{"status":"Internal Server Error","message":"error obtaining warehouses"}`

		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "GetAll", 1)
	})

	// Note: Commented test case for empty warehouse collection
	// This scenario might be implemented in the future to handle empty results
	//t.Run("There are not warehouses", func(t *testing.T) {
	// Arrange
	//	mockService := new(service.WarehouseServiceMock)
	//	mockService.On("GetAll").
	//		Return(nil, nil)
	//	handlerTest := handler.NewWarehouseHandler(mockService)

	// Act
	//	req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
	//	res := httptest.NewRecorder()
	//	handlerTest.GetAll()(res, req)

	// Assert
	//	expectedBody := `{"data":[]}`
	//	expectedStatus := http.StatusOK

	//	require.JSONEq(t, expectedBody, res.Body.String())
	//	require.Equal(t, expectedStatus, res.Code)
	//	require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
	//	mockService.AssertExpectations(t)
	//	mockService.AssertNumberOfCalls(t, "GetAll", 1)
	//})

}

// Test_GetById tests the GetById method of the WarehouseHandler.
//
// This test suite verifies that the warehouse retrieval by ID endpoint:
//   - Successfully returns warehouse when valid ID is provided
//   - Returns 404 when warehouse doesn't exist
//   - Returns 400 when invalid ID format is provided
//   - Properly extracts ID from URL parameters using chi router context
//
// Test scenarios:
//   - Success case: Valid ID returns warehouse with 200 OK
//   - Not found: Non-existent ID returns 404 Not Found
//   - Invalid ID: Non-numeric ID returns 400 Bad Request
func Test_GetById(t *testing.T) {
	// Test data used across multiple test cases
	warehouseTest := models.Warehouse{
		Id: 1,
		WarehouseAttributes: models.WarehouseAttributes{
			WarehouseCode:      "WH-001",
			Address:            "Calle Falsa 123",
			Telephone:          "123456789",
			MinimunCapacity:    100,
			MinimunTemperature: 4.5,
		},
	}

	// Success scenario: retrieve existing warehouse by valid ID
	t.Run("on success (find_by_id_existent)", func(t *testing.T) {
		// Arrange: mock service to return existing warehouse
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("GetByID", warehouseTest.Id).
			Return(warehouseTest, nil)
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute GET request with valid warehouse ID
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
		// Set up chi router context with ID parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		handlerTest.GetById()(res, req)

		// Assert: verify successful response with warehouse data
		expectedStatus := http.StatusOK
		expectedBody, err := json.Marshal(map[string]any{
			"data": warehouseTest,
		})

		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "GetByID", 1)
	})

	// Error scenario: attempt to retrieve non-existent warehouse
	t.Run("on id not exists (find_by_id_non_existent)", func(t *testing.T) {
		// Arrange: mock service to return not found error
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("GetByID", 2).
			Return(models.Warehouse{}, httperrors.NotFoundError{Message: "warehouse not found"})

		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute request with non-existent warehouse ID
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/warehouses/", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		handlerTest.GetById()(res, req)

		// Assert: verify 404 status and not found error message
		expectedStatus := http.StatusNotFound
		expectedBody, err := json.Marshal(map[string]any{
			"status":  http.StatusText(expectedStatus),
			"message": httperrors.NotFoundError{Message: "warehouse not found"}.Message,
		})
		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "GetByID", 1)
	})

	// Error scenario: invalid ID format provided
	t.Run("on id invalid", func(t *testing.T) {
		// Arrange: setup handler without mocking service calls
		mockService := new(service_mocks.WarehouseServiceMock)
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute request with non-numeric ID
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "letter") // Invalid: non-numeric ID
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		handlerTest.GetById()(res, req)

		// Assert: verify 400 status and invalid format error
		expectedStatus := http.StatusBadRequest
		expectedBody, err := json.Marshal(map[string]any{
			"status":  http.StatusText(expectedStatus),
			"message": httperrors.BadRequestError{Message: "Invalid Id format"}.Message,
		})

		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "GetByID", 1)
	})
}

// Test_Update tests the Update method of the WarehouseHandler.
//
// This test suite verifies that the warehouse update endpoint:
//   - Successfully updates warehouses with valid data and ID
//   - Returns 404 when attempting to update non-existent warehouse
//   - Handles malformed JSON requests appropriately
//   - Properly handles various service-layer errors (conflict, internal errors)
//   - Validates ID format before processing
//
// Test scenarios:
//   - Success case: Valid data and ID returns updated warehouse with 200 OK
//   - Not found: Non-existent ID returns 404 Not Found
//   - Invalid JSON: Malformed request body returns 422 Unprocessable Entity
//   - Service errors: Various business logic errors return appropriate status codes
//   - Invalid ID: Non-numeric ID returns 400 Bad Request
func Test_Update(t *testing.T) {
	// Test data used across multiple test cases
	warehouseAttInput := models.WarehouseAttributes{
		WarehouseCode:      "WH-001",
		Address:            "Calle Falsa 123",
		Telephone:          "123456789",
		MinimunCapacity:    100,
		MinimunTemperature: 4.5,
	}
	warehouseOutput := models.Warehouse{
		Id:                  1,
		WarehouseAttributes: warehouseAttInput,
	}

	// Success scenario: update existing warehouse with valid data
	t.Run("on success (update_ok)", func(t *testing.T) {
		// Arrange: mock service to return successful update
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("Update", 1, warehouseAttInput).
			Return(warehouseOutput, nil)
		handlerTest := handler.NewWarehouseHandler(mockService)
		body, err := json.Marshal(warehouseAttInput)
		require.NoError(t, err)

		// Act: execute PATCH request with valid data and ID
		req := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handlerTest.Update()(res, req)

		// Assert: verify successful update response
		expectedStatus := http.StatusOK
		expectedBody, err := json.Marshal(map[string]any{
			"data": warehouseOutput,
		})

		require.NoError(t, err)
		require.Equal(t, expectedStatus, res.Code)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "Update", 1)
	})

	// Error scenario: attempt to update non-existent warehouse
	t.Run("on id not exists (update_non_existent)", func(t *testing.T) {
		// Arrange: mock service to return not found error
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("Update", 2, warehouseAttInput).
			Return(models.Warehouse{}, httperrors.NotFoundError{Message: "warehouse not found"})
		handlerTest := handler.NewWarehouseHandler(mockService)
		body, err := json.Marshal(warehouseAttInput)
		require.NoError(t, err)

		// Act: execute update request for non-existent warehouse
		req := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handlerTest.Update()(res, req)

		// Assert: verify 404 status and not found error
		expectedStatus := http.StatusNotFound
		expectedBody, err := json.Marshal(map[string]any{
			"status":  http.StatusText(expectedStatus),
			"message": httperrors.NotFoundError{Message: "warehouse not found"}.Message,
		})

		require.NoError(t, err)
		require.Equal(t, expectedStatus, res.Code)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "Update", 1)
	})

	// Error scenario: malformed JSON in request body
	t.Run("on json malformed", func(t *testing.T) {
		// Arrange: prepare invalid JSON with wrong data type
		mockService := new(service_mocks.WarehouseServiceMock)
		handlerTest := handler.NewWarehouseHandler(mockService)
		bodyInput := map[string]any{
			"minimun_temperature": "this is not a float", // Invalid: string instead of float
		}
		body, err := json.Marshal(bodyInput)
		require.NoError(t, err)

		// Act: execute update request with malformed JSON
		req := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		res := httptest.NewRecorder()
		handlerTest.Update()(res, req)

		// Assert: verify 422 status and invalid JSON error
		expectedStatus := http.StatusUnprocessableEntity
		expectedBody, err := json.Marshal(map[string]any{
			"message": httperrors.BadRequestError{Message: "Invalid JSON body"}.Message,
			"status":  http.StatusText(expectedStatus),
		})

		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "Update")
	})

	// Error scenarios: various service-layer errors using table-driven tests
	t.Run("On service error", func(t *testing.T) {
		// Table-driven test for different service error scenarios
		tests := []struct {
			name           string
			expectedStatus int
			serviceError   error
		}{
			{
				name:           "On conflict",
				expectedStatus: http.StatusConflict,
				serviceError:   httperrors.ConflictError{Message: "the WarehouseCode already exists"},
			},
			{
				name:           "On internal server error",
				expectedStatus: http.StatusInternalServerError,
				serviceError:   httperrors.InternalServerError{Message: "error updating warehouse"},
			},
		}

		body, err := json.Marshal(warehouseAttInput)
		require.NoError(t, err)

		// Execute each error scenario
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Arrange: mock service to return specific error
				mockService := new(service_mocks.WarehouseServiceMock)
				mockService.On("Update", 1, warehouseAttInput).
					Return(models.Warehouse{}, test.serviceError)
				handlerTest := handler.NewWarehouseHandler(mockService)

				// Act: execute update request that will result in service error
				req := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewReader(body))
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("id", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				handlerTest.Update()(res, req)

				// Assert: verify expected status and error message
				expectedBody, err := json.Marshal(map[string]any{
					"status":  http.StatusText(test.expectedStatus),
					"message": test.serviceError.Error(),
				})

				require.NoError(t, err)
				require.Equal(t, test.expectedStatus, res.Code)
				require.JSONEq(t, string(expectedBody), res.Body.String())
				require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
				mockService.AssertExpectations(t)
				mockService.AssertNumberOfCalls(t, "Update", 1)
			})
		}
	})

	// Error scenario: invalid ID format provided
	t.Run("on id invalid", func(t *testing.T) {
		// Arrange: setup with invalid ID parameter
		mockService := new(service_mocks.WarehouseServiceMock)
		handlerTest := handler.NewWarehouseHandler(mockService)
		body, err := json.Marshal(warehouseAttInput)
		require.NoError(t, err)

		// Act: execute update request with non-numeric ID
		req := httptest.NewRequest(http.MethodPatch, "/warehouses/letter", bytes.NewReader(body))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "letter") // Invalid: non-numeric ID
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		res := httptest.NewRecorder()
		handlerTest.Update()(res, req)

		// Assert: verify 400 status and invalid ID format error
		expectedStatus := http.StatusBadRequest
		expectedBody, err := json.Marshal(map[string]any{
			"status":  http.StatusText(expectedStatus),
			"message": httperrors.BadRequestError{Message: "Invalid Id format"}.Message,
		})

		require.NoError(t, err)
		require.Equal(t, expectedStatus, res.Code)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "Update")
	})

}

// Test_Delete tests the Delete method of the WarehouseHandler.
//
// This test suite verifies that the warehouse deletion endpoint:
//   - Successfully deletes existing warehouses
//   - Returns 404 when attempting to delete non-existent warehouse
//   - Validates ID format before processing
//   - Returns appropriate HTTP status codes (204 for success, 404 for not found)
//
// Test scenarios:
//   - Success case: Valid ID deletes warehouse and returns 204 No Content
//   - Not found: Non-existent ID returns 404 Not Found
//   - Invalid ID: Non-numeric ID returns 400 Bad Request
func Test_Delete(t *testing.T) {
	// Success scenario: delete existing warehouse
	t.Run("on success (delete_ok)", func(t *testing.T) {
		// Arrange: mock service to return successful deletion
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("Delete", 1).
			Return(nil)
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute DELETE request with valid warehouse ID
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		handlerTest.Delete()(res, req)

		// Assert: verify successful deletion with 204 status and empty body
		expectedStatus := http.StatusNoContent

		require.Equal(t, expectedStatus, res.Code)
		require.Empty(t, res.Body.String())
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "Delete", 1)
	})

	// Error scenario: attempt to delete non-existent warehouse
	t.Run("Id not exists (delete_non_existent)", func(t *testing.T) {
		// Arrange: mock service to return not found error
		mockService := new(service_mocks.WarehouseServiceMock)
		mockService.On("Delete", 1).
			Return(httperrors.NotFoundError{Message: "warehouse not found"})
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute delete request for non-existent warehouse
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		handlerTest.Delete()(res, req)

		// Assert: verify 404 status and not found error message
		expectedStatus := http.StatusNotFound
		expectedBody, err := json.Marshal(map[string]any{
			"status":  http.StatusText(expectedStatus),
			"message": httperrors.NotFoundError{Message: "warehouse not found"}.Message,
		})

		require.NoError(t, err)
		require.Equal(t, expectedStatus, res.Code)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		mockService.AssertExpectations(t)
		mockService.AssertNumberOfCalls(t, "Delete", 1)
	})

	// Error scenario: invalid ID format provided
	t.Run("on id invalid", func(t *testing.T) {
		// Arrange: setup with invalid ID parameter
		mockService := new(service_mocks.WarehouseServiceMock)
		handlerTest := handler.NewWarehouseHandler(mockService)

		// Act: execute delete request with non-numeric ID
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "letter") // Invalid: non-numeric ID
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		handlerTest.Delete()(res, req)

		// Assert: verify 400 status and invalid ID format error
		expectedStatus := http.StatusBadRequest
		expectedBody, err := json.Marshal(map[string]any{
			"status":  http.StatusText(expectedStatus),
			"message": httperrors.BadRequestError{Message: "Invalid ID format"}.Message,
		})

		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), res.Body.String())
		require.Equal(t, expectedStatus, res.Code)
		require.Equal(t, expectedContentType, res.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "Delete")
	})

}
