package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	serviceMocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-chi/chi/v5"
)

// getValidEmployeeToCreate returns a valid set of employee attributes for creation tests.
func getValidEmployeeToCreate() models.EmployeeAttributes {
	return models.EmployeeAttributes{
		CardNumberID: "10000000",
		FirstName:    "Thomas",
		LastName:     "Shelby",
		WarehouseID:  1,
	}
}

// getValidEmployeeCreated returns a valid Employee with the given ID and attributes.
func getValidEmployeeCreated(id int, attrs models.EmployeeAttributes) models.Employee {
	return models.Employee{
		Id:                 id,
		EmployeeAttributes: attrs,
	}
}

// createRequestWithBody constructs an HTTP request with a marshaled JSON body.
func createRequestWithBody(method, url string, data interface{}) *http.Request {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// createRequestWithContext adds a URL parameter to the request context for routing.
func createRequestWithContext(method, url, paramKey, paramValue string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(paramKey, paramValue)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}

// createRequestWithBodyAndContext builds a request with both a JSON body and route context.
func createRequestWithBodyAndContext(method, url string, data interface{}, paramKey, paramValue string) *http.Request {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(paramKey, paramValue)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}

// TestEmployeeHandler_Create tests the Create handler for various scenarios.
func TestEmployeeHandler_Create(t *testing.T) {
	// Tests successful creation of an employee.
	t.Run("create_ok: success to create employee", func(t *testing.T) {
		// arrange
		employee := getValidEmployeeToCreate()
		expectedEmployee := getValidEmployeeCreated(1, employee)
		expectedEmployeeJson, err := json.Marshal(expectedEmployee)
		require.NoError(t, err)
		expectedCode := http.StatusCreated
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}

		mockService := &serviceMocks.MockEmployeeService{}
		mockService.On("Create", employee).Return(expectedEmployee, nil)
		h := handler.NewEmployeeHandler(mockService)
		rr := httptest.NewRecorder()
		req := createRequestWithBody(http.MethodPost, "/employees", employee)

		// act
		h.Create()(rr, req)

		// assert
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), string(expectedEmployeeJson))
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertExpectations(t)
	})

	// Tests that creating with missing required fields returns 422 and does not call the service.
	t.Run("create_fail: fail to create employee - missing required field", func(t *testing.T) {
		// arrange
		invalidEmployee := models.EmployeeAttributes{
			CardNumberID: "10000000",
			FirstName:    "Thomas",
			// missing LastName
			WarehouseID: 1,
		}

		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		rr := httptest.NewRecorder()
		req := createRequestWithBody(http.MethodPost, "/employees", invalidEmployee)

		// act
		h.Create()(rr, req)

		// assert
		expectedCode := http.StatusUnprocessableEntity
		expectedSubBody := "Invalid JSON body"
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), expectedSubBody)
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertNotCalled(t, "Create")
	})

	// Tests that attempting to create a conflicting card number returns 409 Conflict.
	t.Run("create_conflict: fail to create employee - conflict card number", func(t *testing.T) {
		// arrange
		employee := getValidEmployeeToCreate()
		conflictErr := httperrors.ConflictError{Message: "Card number already exists"}

		mockService := &serviceMocks.MockEmployeeService{}
		mockService.On("Create", employee).Return(models.Employee{}, conflictErr)
		h := handler.NewEmployeeHandler(mockService)
		rr := httptest.NewRecorder()
		req := createRequestWithBody(http.MethodPost, "/employees", employee)

		// act
		h.Create()(rr, req)

		// assert
		expectedCode := http.StatusConflict
		expectedSubBody := "Card number already exists"
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), expectedSubBody)
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertExpectations(t)
	})

	// Tests that an invalid JSON body results in a 400 Bad Request.
	t.Run("create_bad_request: invalid json body", func(t *testing.T) {
		// arrange
		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString("{bad json"))
		req.Header.Set("Content-Type", "application/json")

		// act
		h.Create()(rr, req)

		// assert
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Contains(t, rr.Body.String(), "Invalid JSON body")
	})
}

// TestEmployeeHandler_Get tests retrieval (find) endpoints for employees.
func TestEmployeeHandler_Get(t *testing.T) {
	// Tests retrieval of all employees returns status OK and correct data.
	t.Run("find_all: success get all employees", func(t *testing.T) {
		// arrange
		employeeAFields := getValidEmployeeToCreate()
		employeeA := getValidEmployeeCreated(1, employeeAFields)
		employeeBFields := getValidEmployeeToCreate()
		employeeBFields.FirstName = "Arthur"
		employeeBFields.CardNumberID = "10000001"
		employeeB := getValidEmployeeCreated(1, employeeBFields)
		expectedEmployees := []models.Employee{employeeA, employeeB}
		expectedJsonA, err := json.Marshal(employeeA)
		require.NoError(t, err)
		expectedJsonB, err := json.Marshal(employeeB)
		require.NoError(t, err)

		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("GetAll").Return(expectedEmployees, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees", nil)
		rr := httptest.NewRecorder()

		// act
		h.GetAll()(rr, req)

		// assert
		expectedCode := http.StatusOK
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Equal(t, expectedHeaders, rr.Header())
		require.Contains(t, rr.Body.String(), string(expectedJsonA))
		require.Contains(t, rr.Body.String(), string(expectedJsonB))
		mockService.AssertExpectations(t)
	})

	// Tests retrieval of a single employee by an existing ID returns correct data.
	t.Run("find_by_id_existent: success get employee by existent id", func(t *testing.T) {
		// arrange
		expectedEmployee := getValidEmployeeCreated(1, getValidEmployeeToCreate())
		expectedJsonBody, err := json.Marshal(expectedEmployee)
		require.NoError(t, err)
		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("GetByID", 1).Return(expectedEmployee, nil)
		req := createRequestWithContext(http.MethodGet, "/employees/1", "id", "1")
		rr := httptest.NewRecorder()

		// act
		h.GetById()(rr, req)

		// assert
		expectedCode := http.StatusOK
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), string(expectedJsonBody))
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertExpectations(t)
	})

	// Tests retrieval by a non-existent ID returns 404 Not Found.
	t.Run("find_by_id_non_existent: fail - get employee by non-existent id", func(t *testing.T) {
		// arrange
		notFoundError := httperrors.NotFoundError{Message: "Employee not found"}

		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("GetByID", 999).Return(models.Employee{}, notFoundError)
		req := createRequestWithContext(http.MethodGet, "/employees/999", "id", "999")
		rr := httptest.NewRecorder()

		// act
		h.GetById()(rr, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedSubBody := "Employee not found"
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), expectedSubBody)
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertExpectations(t)
	})
}

// TestEmployeeHandler_Update tests employee update logic for success and not-found cases.
func TestEmployeeHandler_Update(t *testing.T) {
	// Tests successful update of an employee.
	t.Run("update_ok: success update employee", func(t *testing.T) {
		// arrange
		employee := getValidEmployeeToCreate()
		expectedEmployee := getValidEmployeeCreated(1, employee)
		expectedJsonBody, err := json.Marshal(expectedEmployee)
		require.NoError(t, err)

		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("Update", 1, employee).Return(expectedEmployee, nil)
		req := createRequestWithBodyAndContext(http.MethodPut, "/employees/1", employee, "id", "1")
		rr := httptest.NewRecorder()

		// act
		h.Update()(rr, req)

		// assert
		expectedCode := http.StatusOK
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), string(expectedJsonBody))
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertExpectations(t)
	})

	// Tests that updating a non-existent employee returns 404 Not Found.
	t.Run("update_non_existent: fail - update non-existent employee", func(t *testing.T) {
		// arrange
		employee := getValidEmployeeToCreate()
		notFoundError := httperrors.NotFoundError{Message: "Employee not found"}

		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("Update", 999, employee).Return(models.Employee{}, notFoundError)
		req := createRequestWithBodyAndContext(http.MethodPut, "/employees/999", employee, "id", "999")
		rr := httptest.NewRecorder()

		// act
		h.Update()(rr, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedSubBody := "Employee not found"
		expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), expectedSubBody)
		require.Equal(t, expectedHeaders, rr.Header())
		mockService.AssertExpectations(t)
	})
}

// TestEmployeeHandler_Delete tests the Delete handler for both successful and not-found cases.
func TestEmployeeHandler_Delete(t *testing.T) {
	// Tests successful deletion of an employee results in 204 No Content.
	t.Run("delete_ok: success delete employee", func(t *testing.T) {
		// arrange
		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("Delete", 1).Return(nil)
		req := createRequestWithContext(http.MethodDelete, "/employees/1", "id", "1")
		rr := httptest.NewRecorder()

		// act
		h.Delete()(rr, req)

		// assert
		expectedCode := http.StatusNoContent
		require.Equal(t, expectedCode, rr.Code)
		require.Empty(t, rr.Body.String())
		mockService.AssertExpectations(t)
	})

	// Tests deletion of a non-existent employee returns 404 Not Found.
	t.Run("delete_non_existent: fail - delete non-existent employee", func(t *testing.T) {
		// arrange
		notFoundError := httperrors.NotFoundError{Message: "Employee not found"}

		mockService := &serviceMocks.MockEmployeeService{}
		h := handler.NewEmployeeHandler(mockService)
		mockService.On("Delete", 999).Return(notFoundError)
		req := createRequestWithContext(http.MethodDelete, "/employees/999", "id", "999")
		rr := httptest.NewRecorder()

		// act
		h.Delete()(rr, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedSubBody := "Employee not found"
		require.Equal(t, expectedCode, rr.Code)
		require.Contains(t, rr.Body.String(), expectedSubBody)
		require.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, rr.Header())
		mockService.AssertExpectations(t)
	})
}

// TestEmployeeHandler_GetInboundOrderReport tests the report handler for inbound orders report generation.
func TestEmployeeHandler_GetInboundOrderReport(t *testing.T) {
	// Tests report generation for all employees (no 'id' query param).
	t.Run("report_all: success - sin id", func(t *testing.T) {
		// arrange
		mockService := &serviceMocks.MockEmployeeService{}
		expectedReport := []models.EmployeeWithInboundCount{
			{
				Id:                 1,
				CardNumberID:       "10000000",
				FirstName:          "Thomas",
				LastName:           "Shelby",
				WarehouseID:        1,
				InboundOrdersCount: 3,
			},
			{
				Id:                 2,
				CardNumberID:       "10000001",
				FirstName:          "Arthur",
				LastName:           "Shelby",
				WarehouseID:        1,
				InboundOrdersCount: 1,
			},
		}
		mockService.On("ReportInboundOrders", 0).Return(expectedReport, nil)
		h := handler.NewEmployeeHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/employees/report-inbound-orders", nil)
		rr := httptest.NewRecorder()

		// act
		h.GetInboundOrderReport()(rr, req)

		// assert
		require.Equal(t, http.StatusOK, rr.Code)
		require.Contains(t, rr.Body.String(), `"id":1`)
		require.Contains(t, rr.Body.String(), `"id":2`)
		require.Contains(t, rr.Body.String(), `"inbound_orders_count":3`)
		require.Contains(t, rr.Body.String(), `"inbound_orders_count":1`)
		require.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, rr.Header())
		mockService.AssertExpectations(t)
	})

	// Tests report generation for a specific employee by valid 'id'.
	t.Run("report_by_id: success - con id válido", func(t *testing.T) {
		// arrange
		mockService := &serviceMocks.MockEmployeeService{}
		expectedReport := []models.EmployeeWithInboundCount{
			{
				Id:                 1,
				CardNumberID:       "10000000",
				FirstName:          "Thomas",
				LastName:           "Shelby",
				WarehouseID:        1,
				InboundOrdersCount: 3,
			},
		}
		mockService.On("ReportInboundOrders", 1).Return(expectedReport, nil)
		req := httptest.NewRequest(http.MethodGet, "/employees/report-inbound-orders?id=1", nil)
		rr := httptest.NewRecorder()
		h := handler.NewEmployeeHandler(mockService)

		// act
		h.GetInboundOrderReport()(rr, req)

		// assert
		require.Equal(t, http.StatusOK, rr.Code)
		require.Contains(t, rr.Body.String(), `"id":1`)
		require.Contains(t, rr.Body.String(), `"inbound_orders_count":3`)
		require.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, rr.Header())
		mockService.AssertExpectations(t)
	})

	// Tests that an invalid 'id' query returns 400 Bad Request.
	t.Run("report_invalid_id: fail - id inválido", func(t *testing.T) {
		mockService := &serviceMocks.MockEmployeeService{}
		req := httptest.NewRequest(http.MethodGet, "/employees/report-inbound-orders?id=foo", nil)
		rr := httptest.NewRecorder()
		h := handler.NewEmployeeHandler(mockService)

		// act
		h.GetInboundOrderReport()(rr, req)

		// assert
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Contains(t, rr.Body.String(), "Invalid employee_id")
	})

	// Tests report generation returns a service error (404 Not Found).
	t.Run("report_service_error: fail - error del service", func(t *testing.T) {
		notFoundErr := httperrors.NotFoundError{Message: "Employee/s not found"}
		mockService := &serviceMocks.MockEmployeeService{}
		mockService.On("ReportInboundOrders", 1).Return([]models.EmployeeWithInboundCount{}, notFoundErr)
		req := httptest.NewRequest(http.MethodGet, "/employees/report-inbound-orders?id=1", nil)
		rr := httptest.NewRecorder()
		h := handler.NewEmployeeHandler(mockService)

		// act
		h.GetInboundOrderReport()(rr, req)

		// assert
		require.Equal(t, http.StatusNotFound, rr.Code)
		require.Contains(t, rr.Body.String(), "Employee/s not found")
		require.Equal(t, http.Header{"Content-Type": []string{"application/json"}}, rr.Header())
		mockService.AssertExpectations(t)
	})
}
