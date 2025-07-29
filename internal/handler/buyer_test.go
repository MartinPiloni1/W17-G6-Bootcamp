package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// similar to the error response in reponse.Error()
type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// allows to mock ONE path params to the request
func addChiURLParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
}

func TestBuyerHandler_Create(t *testing.T) {
	t.Parallel()

	// simulates the success response of the method
	type DataResponseCreate struct {
		Data models.Buyer `json:"data"`
	}
	t.Run("passing unkonwn fields in the body will return a StatusBadRequest", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		bodyWithUnknownFields := strings.NewReader(`{
			"first_name": "Juan",
			"last_name": "Perez",
			"card_number_id": 23456789
			"unknown_field": "thisFieldIsUnknown"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/", bodyWithUnknownFields) // r
		rec := httptest.NewRecorder()                                           // w

		// act
		buyerHandler.Create().ServeHTTP(rec, req)

		// assert
		expectedMessage := "Invalid JSON body"
		expectedCode := http.StatusBadRequest

		assert.Equal(t, expectedCode, rec.Code)
		serviceMock.AssertNotCalled(t, "Create")

		// validate string appears in response
		assert.Contains(t, rec.Body.String(), expectedMessage)

		// unmarshal response to check the exact response
		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, expectedMessage, actualResponseError.Message)
	})

	t.Run("passing a body that has invalid values returning StatusUnprocessableEntity", func(t *testing.T) {
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)
		invalidRequestBody := `{"first_name":"Juan","last_name":"","card_number_id":23456789}` // empty last name

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidRequestBody))
		rec := httptest.NewRecorder()

		buyerHandler.Create().ServeHTTP(rec, req)

		expectedMessage := "Invalid JSON body"
		expectedCode := http.StatusUnprocessableEntity
		assert.Equal(t, expectedCode, rec.Code)
		serviceMock.AssertNotCalled(t, "Create")

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, expectedMessage, actualResponseError.Message)
	})

	t.Run("Create with a card_number_id that exist will return a ConflictError 409", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		newBuyerWithDuplicateCardNumberId := models.BuyerAttributes{
			CardNumberId: 12345678,
			FirstName:    "Juan",
			LastName:     "Perez",
		}

		bodyRequest := `{"first_name":"Juan","last_name":"Perez","card_number_id":12345678}`

		conflictErr := httperrors.ConflictError{Message: "CardNumberId already in use"}

		ctx := mock.Anything // this avoid comparison problems
		serviceMock.On("Create", ctx, newBuyerWithDuplicateCardNumberId).
			Return(models.Buyer{}, conflictErr)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyRequest))
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Create().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusConflict
		assert.Equal(t, expectedCode, rec.Code)

		require.True(t, serviceMock.AssertExpectations(t))

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, conflictErr.Message, actualResponseError.Message)
	})

	t.Run("successfully created a buyer returning StatusCodeOK and new buyer data with id", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		newBuyerValid := models.BuyerAttributes{
			CardNumberId: 12345678,
			FirstName:    "Juan",
			LastName:     "Perez",
		}
		expectedBuyer := models.Buyer{
			Id:              1,
			BuyerAttributes: newBuyerValid,
		}

		bodyRequest := fmt.Sprintf(
			`{"first_name":"%s","last_name":"%s","card_number_id":%d}`,
			newBuyerValid.FirstName,
			newBuyerValid.LastName,
			newBuyerValid.CardNumberId,
		)

		ctx := mock.Anything
		serviceMock.On("Create", ctx, newBuyerValid).Return(expectedBuyer, nil)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyRequest))
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Create().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusCreated
		assert.Equal(t, expectedCode, rec.Code)

		var resp DataResponseCreate
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, expectedBuyer, resp.Data)

		serviceMock.AssertExpectations(t)
	})

	t.Run("service responds with a random error from Repository, should return Internal Server Error", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		newBuyerValid := models.BuyerAttributes{
			CardNumberId: 12345678,
			FirstName:    "Juan",
			LastName:     "Perez",
		}

		bodyRequest := fmt.Sprintf(
			`{"first_name":"%s","last_name":"%s","card_number_id":%d}`,
			newBuyerValid.FirstName,
			newBuyerValid.LastName,
			newBuyerValid.CardNumberId,
		)

		ctx := mock.Anything
		randomErrorFromRepo := errors.New("the database is busy")
		serviceMock.On("Create", ctx, newBuyerValid).Return(models.Buyer{}, randomErrorFromRepo)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyRequest))
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Create().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusInternalServerError
		assert.Equal(t, expectedCode, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, "Internal Server Error", actualResponseError.Message)

		require.True(t, serviceMock.AssertExpectations(t))
	})
}

func TestBuyerHandler_GetAll(t *testing.T) {
	t.Parallel()

	// simulates the success response of the method
	type DataResponseGetAll struct {
		Data []models.Buyer `json:"data"`
	}
	t.Run("the service responds with a random error from repository, returns InternalServerError", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		randomErrorFromRepo := errors.New("database is busy")
		serviceMock.On("GetAll", mock.Anything).Return([]models.Buyer{}, randomErrorFromRepo)

		// act
		buyerHandler.GetAll().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusInternalServerError
		assert.Equal(t, expectedCode, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, "Internal Server Error", actualResponseError.Message)

		require.True(t, serviceMock.AssertExpectations(t))
	})

	t.Run("successfully fetch 2 buyers and returns a statusOK", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		expectedBuyers := []models.Buyer{
			{
				Id: 1,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberId: 12345678,
					LastName:     "Perez",
					FirstName:    "Juan",
				},
			},
			{
				Id: 2,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberId: 22345678,
					LastName:     "Esteban",
					FirstName:    "Pablo",
				},
			},
		}
		ctx := mock.Anything
		serviceMock.On("GetAll", ctx).Return(expectedBuyers, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// act
		buyerHandler.GetAll().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusOK
		assert.Equal(t, expectedCode, rec.Code)

		var actualResponse DataResponseGetAll
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
		require.NoError(t, err)

		assert.Equal(t, expectedBuyers, actualResponse.Data)
		assert.Equal(t, 2, len(actualResponse.Data))
		serviceMock.AssertExpectations(t)
	})
}

func TestBuyerHandler_GetByID(t *testing.T) {
	t.Parallel()

	// simulates the success response of the method
	type DataResponseGetByID struct {
		Data models.Buyer `json:"data"`
	}
	t.Run("passed and invalid id param returns StatusBadRequest Invalid Id", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		invalidID := "invalidID"
		req := httptest.NewRequest(http.MethodGet, "/"+invalidID, nil) // add id to the param (not fully nescesary)
		rec := httptest.NewRecorder()

		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", invalidID)

		// act
		buyerHandler.GetByID().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedMessage := "Invalid Id"

		assert.Equal(t, expectedCode, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, expectedMessage, actualResponseError.Message)
		serviceMock.AssertNotCalled(t, "GetByID")
	})

	t.Run("passed and invalid id (0) param returns StatusBadRequest Invalid Id", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		invalidID := "0"
		req := httptest.NewRequest(http.MethodGet, "/"+invalidID, nil) // add id to the param (not fully nescesary)
		rec := httptest.NewRecorder()

		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", invalidID)

		// act
		buyerHandler.GetByID().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedMessage := "Invalid Id"

		assert.Equal(t, expectedCode, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, expectedMessage, actualResponseError.Message)
		serviceMock.AssertNotCalled(t, "GetByID")
	})

	t.Run("service return NotFoundError the handler return StatusNotFound", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		validID := "1"
		req := httptest.NewRequest(http.MethodGet, "/"+validID, nil) // add id to the param (not fully nescesary)
		rec := httptest.NewRecorder()

		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", validID)

		notFoundErr := httperrors.NotFoundError{Message: "Buyer Not Found"}
		serviceMock.On("GetByID", mock.Anything, 1).Return(models.Buyer{}, notFoundErr)

		// act
		buyerHandler.GetByID().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedMessage := "Buyer Not Found"
		assert.Equal(t, expectedCode, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, expectedMessage, actualResponseError.Message)
		serviceMock.AssertExpectations(t)
	})

	t.Run("successfully returns the buyer with statusOK", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		validID := "1"
		req := httptest.NewRequest(http.MethodGet, "/"+validID, nil) // add id to the param (not fully nescesary)
		rec := httptest.NewRecorder()

		expectedBuyer := models.Buyer{
			Id: 1, // equal to validID
			BuyerAttributes: models.BuyerAttributes{
				CardNumberId: 12345678,
				LastName:     "Juan",
				FirstName:    "Perez",
			},
		}
		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", validID)

		serviceMock.On("GetByID", mock.Anything, expectedBuyer.Id).
			Return(expectedBuyer, nil)

		// act
		buyerHandler.GetByID().ServeHTTP(rec, req)

		// assert
		expectedCode := http.StatusOK
		assert.Equal(t, expectedCode, rec.Code)

		var actualResponse DataResponseGetByID
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
		require.NoError(t, err)

		assert.Equal(t, expectedBuyer, actualResponse.Data)
		serviceMock.AssertExpectations(t)
	})
}

func TestBuyerHandler_Delete(t *testing.T) {
	t.Parallel()

	t.Run("executing the delete without passing the id return StatusBadRequest", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Delete().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)
		assert.Equal(t, "Invalid Id", actualResponseError.Message)

		serviceMock.AssertNotCalled(t, "Delete")
	})

	t.Run("passing as an id an invalid Integer return StatusBadRequest", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		invalidID := "invalidID"
		req := httptest.NewRequest(http.MethodDelete, "/"+invalidID, nil)
		rec := httptest.NewRecorder()

		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", invalidID)

		// act
		buyerHandler.Delete().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)
		assert.Equal(t, "Invalid Id", actualResponseError.Message)

		serviceMock.AssertNotCalled(t, "Delete")
	})

	t.Run("passing an id of a buyer that does not exist in the db returns NotFound", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		validID := "1"
		req := httptest.NewRequest(http.MethodDelete, "/"+validID, nil)
		rec := httptest.NewRecorder()

		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", validID)

		notFoundErr := httperrors.NotFoundError{Message: "Buyer Not Found"}
		serviceMock.On("Delete", mock.Anything, 1).Return(notFoundErr)

		// act
		buyerHandler.Delete().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)
		assert.Equal(t, notFoundErr.Message, actualResponseError.Message)

	})

	t.Run("successfully deletes the buyer returning StatusNoContent", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		validID := "1"
		req := httptest.NewRequest(http.MethodDelete, "/"+validID, nil)
		rec := httptest.NewRecorder()

		// boiler plate to add pathParams to the route
		req = addChiURLParam(req, "id", validID)

		serviceMock.On("Delete", mock.Anything, 1).Return(nil)

		// act
		buyerHandler.Delete().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

func TestBuyerHandler_GetWithPurchaseOrdersCount(t *testing.T) {
	t.Run("passing an invalidID as param returns StatusBadRequest", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		req := httptest.NewRequest(http.MethodGet, "/reportPurchaseOrders?id=invalidID", nil)
		rec := httptest.NewRecorder()

		// act
		buyerHandler.GetWithPurchaseOrdersCount().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)
		assert.Equal(t, "Invalid Id", actualResponseError.Message)
	})

	t.Run("the service returns ErrNotFound, handler returns NotFound", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		validID := "1"
		req := httptest.NewRequest(http.MethodGet, "/reportPurchaseOrders?id="+validID, nil)
		rec := httptest.NewRecorder()

		id, _ := strconv.Atoi(validID)
		notFoundErr := httperrors.NotFoundError{Message: "Buyer Not Found"}
		serviceMock.On("GetWithPurchaseOrdersCount", mock.Anything, &id).
			Return([]models.BuyerWithPurchaseOrdersCount{}, notFoundErr)

		// act
		buyerHandler.GetWithPurchaseOrdersCount().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var actualResponseError errorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &actualResponseError)
		require.NoError(t, err)

		assert.Equal(t, "Buyer Not Found", actualResponseError.Message)
	})

	t.Run("successfully fetchs a reportPurchaseOrders with statusOK", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		expectedBuyersWIthPurchaseOrdersCount := []models.BuyerWithPurchaseOrdersCount{
			{
				PurchaseOrdersCount: 2,
				Buyer: models.Buyer{
					Id: 1,
					BuyerAttributes: models.BuyerAttributes{
						CardNumberId: 12345678,
						LastName:     "Pedro",
						FirstName:    "Perez",
					},
				},
			},
		}
		validID := "1"
		req := httptest.NewRequest(http.MethodGet, "/reportPurchaseOrders?id="+validID, nil)
		rec := httptest.NewRecorder()

		id, _ := strconv.Atoi(validID)
		serviceMock.On("GetWithPurchaseOrdersCount", mock.Anything, &id).
			Return(expectedBuyersWIthPurchaseOrdersCount, nil)

		// act
		buyerHandler.GetWithPurchaseOrdersCount().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestBuyerHandler_Update(t *testing.T) {
	t.Run("invalid id returns StatusBadRequest", func(t *testing.T) {
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		h := handler.NewBuyerHandler(serviceMock)

		invalidID := "invalidID"
		req := httptest.NewRequest(http.MethodPatch, "/buyers/"+invalidID, nil)
		req = addChiURLParam(req, "id", invalidID)
		rec := httptest.NewRecorder()

		h.Update().ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Invalid id", resp.Message)
		serviceMock.AssertNotCalled(t, "Update")
	})

	t.Run("malformed JSON returns StatusBadRequest", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		body := bytes.NewBufferString(`{first_name:"Juan"`) // invalid json
		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", body)
		req = addChiURLParam(req, "id", "1")
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Update().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// validate the response body
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, "Invalid JSON body", resp.Message)
		serviceMock.AssertNotCalled(t, "Update")
	})

	t.Run("unknown fields in JSON returns StatusBadRequest", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		body := bytes.NewBufferString(`{"unknown":"how knows"}`)
		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", body)
		req = addChiURLParam(req, "id", "1")
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Update().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// decode response for validation
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, "Invalid JSON body", resp.Message)
		serviceMock.AssertNotCalled(t, "Update")
	})

	t.Run("validation error returns StatusUnprocessableEntity", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		body := bytes.NewBufferString(`{"first_name":""}`)
		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", body)
		req = addChiURLParam(req, "id", "1")
		rec := httptest.NewRecorder()

		// act
		buyerHandler.Update().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		// decode response for validation
		var resp errorResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, "Invalid JSON body", resp.Message)
		serviceMock.AssertNotCalled(t, "Update")
	})

	t.Run("service returns error NotFoundError that can be returned as it is", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		newFirstName := "Juan"
		buyerNewData := models.BuyerPatchRequest{FirstName: &newFirstName}
		b, err := json.Marshal(&buyerNewData)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewReader(b))
		req = addChiURLParam(req, "id", "1")
		rec := httptest.NewRecorder()

		notFoundErr := httperrors.NotFoundError{Message: "Buyer not Found"}
		serviceMock.On("Update", mock.Anything, 1, buyerNewData).
			Return(models.Buyer{}, notFoundErr)

		// act
		buyerHandler.Update().ServeHTTP(rec, req)

		// assert
		assert.NotEqual(t, http.StatusOK, rec.Code)
		serviceMock.AssertExpectations(t)
	})

	t.Run("successfully updates the buyer", func(t *testing.T) {
		// arrange
		serviceMock := mocks.NewBuyerServiceDefaultMock()
		buyerHandler := handler.NewBuyerHandler(serviceMock)

		newFirstName := "Carlos"
		patch := models.BuyerPatchRequest{FirstName: &newFirstName}
		b, _ := json.Marshal(&patch)
		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewReader(b))
		req = addChiURLParam(req, "id", "1")
		rec := httptest.NewRecorder()

		expected := models.Buyer{
			Id: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberId: 12345678,
				FirstName:    "Carlos",
				LastName:     "Perez",
			},
		}
		serviceMock.On("Update", mock.Anything, 1, patch).Return(expected, nil)

		// act
		buyerHandler.Update().ServeHTTP(rec, req)

		// assert
		assert.Equal(t, http.StatusOK, rec.Code)

		// internal structure thas similar to the actual response
		var resp struct {
			Data models.Buyer `json:"data"`
		}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, expected, resp.Data)
		serviceMock.AssertExpectations(t)
	})
}
