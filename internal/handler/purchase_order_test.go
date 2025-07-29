package handler_test

import (
	"context"
	"encoding/json"
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

func TestPurchaseOrderHandler_Create(t *testing.T) {

	// simulates the success response of the method
	type DataResponseCreate struct {
		Data models.PurchaseOrder `json:"data"`
	}
	t.Run("passing an invalid JSON returns StatusUnprocessableEntity", func(t *testing.T) {
		serviceMock := mocks.NewPurchaseOrderDefaultMock()
		purchaseOrderHandler := handler.NewPurchaseOrderHandler(serviceMock)

		invalidBody := strings.NewReader(`{
			"order_number": "ORD-1",
			"order_date": "2023-04-04T15:04:05",
			"tracking_code": "abc123asd",
			"buyer_id": 1,
			"product_record_id": 1
		}`) // order_date must have the the format "2025-04-04T15:04:05Z"

		req := httptest.NewRequest(http.MethodPost, "/", invalidBody)
		rec := httptest.NewRecorder()

		purchaseOrderHandler.Create().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		serviceMock.AssertNumberOfCalls(t, "Create", 0)
	})

	t.Run("passing a JSON boby with invalid field values returns StatusUnprocessableEntity", func(t *testing.T) {
		serviceMock := mocks.NewPurchaseOrderDefaultMock()
		purchaseOrderHandler := handler.NewPurchaseOrderHandler(serviceMock)

		invalidBody := strings.NewReader(`{
			"order_number": "ORD-1",
			"order_date": "3000-04-04T15:04:05",
			"tracking_code": "abc123asd",
			"buyer_id": 1,
			"product_record_id": 1
		}`) // order_date is setted into the future

		req := httptest.NewRequest(http.MethodPost, "/", invalidBody)
		rec := httptest.NewRecorder()

		// act
		purchaseOrderHandler.Create().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		serviceMock.AssertNumberOfCalls(t, "Create", 0)
	})

	t.Run("the service return error ProductRecordID/BuyerID not found, handler returning a NotFound", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewPurchaseOrderDefaultMock()
		purchaseOrderHandler := handler.NewPurchaseOrderHandler(serviceMock)

		purchaseOrderDate := time.Date(2023, 4, 4, 15, 4, 5, 0, time.UTC)
		purchaseOrderAtt := models.PurchaseOrderAttributes{
			OrderNumber:     "ORD-1",
			OrderDate:       purchaseOrderDate,
			TrackingCode:    "abc123asd",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		validBody := strings.NewReader(`{
			"order_number": "ORD-1",
			"order_date": "2023-04-04T15:04:05Z",
			"tracking_code": "abc123asd",
			"buyer_id": 1,
			"product_record_id": 1
		}`)

		req := httptest.NewRequest(http.MethodPost, "/", validBody)
		rec := httptest.NewRecorder()

		ctx := context.Background()
		notFoundErr := httperrors.NotFoundError{Message: "ProductRecordID/BuyerID not Found"}
		serviceMock.On("Create", ctx, purchaseOrderAtt).Return(models.PurchaseOrder{}, notFoundErr)

		// act
		purchaseOrderHandler.Create().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusNotFound, rec.Code)
		serviceMock.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("successfully creates a purchaseOrder return statusCreated", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewPurchaseOrderDefaultMock()
		purchaseOrderHandler := handler.NewPurchaseOrderHandler(serviceMock)

		purchaseOrderDate := time.Date(2023, 4, 4, 15, 4, 5, 0, time.UTC)
		purchaseOrderAtt := models.PurchaseOrderAttributes{
			OrderNumber:     "ORD-1",
			OrderDate:       purchaseOrderDate,
			TrackingCode:    "abc123asd",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		expectedPurchaseOrder := models.PurchaseOrder{
			Id:                      1,
			PurchaseOrderAttributes: purchaseOrderAtt,
		}

		validBody := strings.NewReader(`{
			"order_number": "ORD-1",
			"order_date": "2023-04-04T15:04:05Z",
			"tracking_code": "abc123asd",
			"buyer_id": 1,
			"product_record_id": 1
		}`)
		req := httptest.NewRequest(http.MethodPost, "/", validBody)
		rec := httptest.NewRecorder()

		ctx := context.Background()
		serviceMock.On("Create", ctx, purchaseOrderAtt).Return(expectedPurchaseOrder, nil)

		purchaseOrderHandler.Create().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusCreated, rec.Code)

		var actualResponse DataResponseCreate
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
		require.NoError(t, err)

		assert.Equal(t, expectedPurchaseOrder, actualResponse.Data)
		serviceMock.AssertNumberOfCalls(t, "Create", 1)
		serviceMock.AssertExpectations(t)
	})
}
