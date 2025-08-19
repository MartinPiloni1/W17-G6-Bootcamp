package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mock "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestLocalityHandler_Create_Success(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	input := models.Locality{ID: "10", LocalityName: "Nuñez"}
	mockService.On("Create", input).Return(input, nil)

	body, _ := json.Marshal(map[string]interface{}{
		"data": input,
	})

	req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var respBody map[string]models.Locality
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, input, respBody["data"])

	mockService.AssertExpectations(t)
}

func TestLocalityHandler_Create_InvalidJSON(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewBufferString("{malformed_json"))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestLocalityHandler_Create_ServiceError(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	input := models.Locality{ID: "11", LocalityName: "Caballito"}
	expectedErr := errors.New("service error")

	mockService.On("Create", input).Return(models.Locality{}, expectedErr)

	body, _ := json.Marshal(map[string]interface{}{
		"data": input,
	})
	req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetByID_Success(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	want := models.Locality{ID: "42", LocalityName: "Recoleta"}
	mockService.On("GetByID", "42").Return(want, nil)

	req := httptest.NewRequest(http.MethodGet, "/localities/42", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "42")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.GetByID()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]models.Locality
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, want, respBody["data"])

	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetByID_MissingID(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/localities/", nil)
	w := httptest.NewRecorder()
	h.GetByID()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLocalityHandler_GetByID_NotFound(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	mockService.On("GetByID", "99").Return(models.Locality{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/localities/99", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "99")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.GetByID()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetSellerReport_OK(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	reports := []models.SellerReport{
		{LocalityID: "1", LocalityName: "Palermo", SellersCount: 4},
		{LocalityID: "2", LocalityName: "Caballito", SellersCount: 2},
	}
	mockService.On("GetSellerReport", (*string)(nil)).Return(reports, nil)

	req := httptest.NewRequest(http.MethodGet, "/localities/sellers/report", nil)
	w := httptest.NewRecorder()
	h.GetSellerReport()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string][]models.SellerReport
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, reports, respBody["data"])

	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetSellerReport_WithID(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	id := "1"
	reports := []models.SellerReport{
		{LocalityID: "1", LocalityName: "Palermo", SellersCount: 3},
	}
	mockService.On("GetSellerReport", &id).Return(reports, nil)

	req := httptest.NewRequest(http.MethodGet, "/localities/sellers/report?id=1", nil)
	w := httptest.NewRecorder()
	h.GetSellerReport()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string][]models.SellerReport
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, reports, respBody["data"])

	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetSellerReport_Error(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	mockService.On("GetSellerReport", (*string)(nil)).Return([]models.SellerReport{}, errors.New("some error"))

	req := httptest.NewRequest(http.MethodGet, "/localities/sellers/report", nil)
	w := httptest.NewRecorder()
	h.GetSellerReport()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
func TestLocalityHandler_GetReportByLocalityId_OK(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	carryReports := []models.CarryReport{
		{LocalityId: "7", LocalityName: "Saavedra", CarriesCount: 2},
	}
	mockService.On("GetReportByLocalityId", "7").Return(carryReports, nil)

	req := httptest.NewRequest(http.MethodGet, "/localities/report_carry?id=7", nil)
	w := httptest.NewRecorder()
	h.GetReportByLocalityId()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string][]models.CarryReport
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, carryReports, respBody["data"])

	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetReportByLocalityId_NoIDAndNilResult(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	mockService.On("GetReportByLocalityId", "").Return([]models.CarryReport{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/localities/report_carry", nil)
	w := httptest.NewRecorder()
	h.GetReportByLocalityId()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string][]models.CarryReport
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(respBody["data"]))

	mockService.AssertExpectations(t)
}

func TestLocalityHandler_GetReportByLocalityId_Error(t *testing.T) {
	mockService := new(mock.MockLocalityService)
	h := handler.NewLocalityHandler(mockService)

	mockService.On("GetReportByLocalityId", "X").Return([]models.CarryReport{}, errors.New("db down"))

	req := httptest.NewRequest(http.MethodGet, "/localities/report_carry?id=X", nil)
	w := httptest.NewRecorder()
	h.GetReportByLocalityId()(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
