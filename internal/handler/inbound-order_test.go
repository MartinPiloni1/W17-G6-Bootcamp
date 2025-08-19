package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dataResponse struct {
	Data models.InboundOrder `json:"data"`
}

func TestInboundOrderHandler_Create(t *testing.T) {
	t.Parallel()

	validOrderDate := time.Date(2024, 6, 13, 15, 0, 0, 0, time.UTC)
	validOrderDateStr := validOrderDate.Format(time.RFC3339)
	// JSON body for valid requests (OrderDate RFC3339 encoded)
	validBody := `{
		"order_number": "123",
		"order_date": "` + validOrderDateStr + `",
		"employee_id": 1,
		"warehouse_id": 2,
		"product_batch_id": 3
	}`

	t.Run("unknown fields in the body should return StatusBadRequest", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)
		body := `{
			"order_number": "123",
			"order_date": "` + validOrderDateStr + `",
			"employee_id": 1,
			"warehouse_id": 2,
			"product_batch_id": 3,
			"unknown_field": 999
		}`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "Create")

		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Invalid JSON body", resp.Message)
	})

	t.Run("invalid json returns StatusBadRequest", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)
		body := `{"order_number":`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "Create")

		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Invalid JSON body", resp.Message)
	})

	t.Run("validation error returns StatusUnprocessableEntity", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)
		// employee_id = 0, no pasa la validación gt=0
		body := `{
			"order_number": "123",
			"order_date": "` + validOrderDateStr + `",
			"employee_id": 0,
			"warehouse_id": 2,
			"product_batch_id": 3
		}`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		serviceMock.AssertNotCalled(t, "Create")

		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Invalid JSON body", resp.Message)
	})

	t.Run("service returns conflict error, returns StatusConflict", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)

		var attrs models.InboundOrderAttributes
		err := json.Unmarshal([]byte(validBody), &attrs)
		require.NoError(t, err)

		expectedErr := httperrors.ConflictError{Message: "OrderNumber already exists"}
		serviceMock.On("Create", attrs).Return(models.InboundOrder{}, expectedErr)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validBody))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, expectedErr.Message, resp.Message)
		serviceMock.AssertExpectations(t)
	})

	t.Run("service returns not found error, returns StatusNotFound", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)

		var attrs models.InboundOrderAttributes
		err := json.Unmarshal([]byte(validBody), &attrs)
		require.NoError(t, err)

		expectedErr := httperrors.NotFoundError{Message: "Related Warehouse not found"}
		serviceMock.On("Create", attrs).Return(models.InboundOrder{}, expectedErr)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validBody))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, expectedErr.Message, resp.Message)
		serviceMock.AssertExpectations(t)
	})

	t.Run("service returns other error, returns StatusInternalServerError", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)

		var attrs models.InboundOrderAttributes
		err := json.Unmarshal([]byte(validBody), &attrs)
		require.NoError(t, err)

		expectedErr := errors.New("random db error")
		serviceMock.On("Create", attrs).Return(models.InboundOrder{}, expectedErr)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validBody))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Internal Server Error", resp.Message)
		serviceMock.AssertExpectations(t)
	})

	t.Run("successfully creates inbound order returns StatusCreated", func(t *testing.T) {
		serviceMock := mocks.NewInboundOrderServiceDefaultMock()
		handler := handler.NewInboundOrderHandler(serviceMock)

		var attrs models.InboundOrderAttributes
		err := json.Unmarshal([]byte(validBody), &attrs)
		require.NoError(t, err)

		expected := models.InboundOrder{
			ID:                     1,
			InboundOrderAttributes: attrs,
		}
		serviceMock.On("Create", attrs).Return(expected, nil)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validBody))
		rec := httptest.NewRecorder()

		handler.Create().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		var resp dataResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, expected, resp.Data)
		serviceMock.AssertExpectations(t)
	})
}
