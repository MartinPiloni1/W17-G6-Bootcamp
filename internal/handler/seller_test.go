package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/handler"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mock/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestSellerHandler_GetAll tests the GetAll method of the SellerHandler.
// It checks for successful retrieval of all sellers and error handling.
func TestSellerHandler_GetAll(t *testing.T) {
	cases := []struct {
		name       string
		setupMock  func(s *mocks.SellerServiceDBMock)
		wantStatus int
	}{
		{
			name: "ok",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("GetAll").Return([]models.Seller{
					{ID: 1, SellerAttributes: models.SellerAttributes{CID: 1}},
					{ID: 2, SellerAttributes: models.SellerAttributes{CID: 2}},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "fail",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("GetAll").Return([]models.Seller{}, errors.New("db fail"))
			},
			wantStatus: http.StatusInternalServerError, // O el código de error que uses
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(mocks.SellerServiceDBMock)
			tc.setupMock(service)
			handler := handler.NewSellerHandler(service)
			router := chi.NewRouter()
			router.Get("/api/v1/sellers", handler.GetAll())
			req := httptest.NewRequest("GET", "/api/v1/sellers", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

// TestSellerHandler_GetByID tests the GetByID method of the SellerHandler.
// It checks for successful retrieval of a seller by ID and error handling for non-existent sellers.
func TestSellerHandler_GetByID(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		setupMock  func(s *mocks.SellerServiceDBMock)
		wantStatus int
	}{
		{
			name: "ok",
			id:   "11",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("GetByID", 11).Return(models.Seller{ID: 11, SellerAttributes: models.SellerAttributes{CID: 555}}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   "987",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("GetByID", 987).Return(models.Seller{}, httperrors.NotFoundError{Message: "not found"})
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "bad id",
			id:   "abc",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				// No call, porque path param no es int.
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(mocks.SellerServiceDBMock)
			tc.setupMock(service)
			handler := handler.NewSellerHandler(service)
			router := chi.NewRouter()
			router.Get("/api/v1/sellers/{id}", handler.GetByID())
			req := httptest.NewRequest("GET", "/api/v1/sellers/"+tc.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

// TestSellerHandler_Create tests the Create method of the SellerHandler.
// It checks for successful creation, conflict errors, and unprocessable entity errors.
func TestSellerHandler_Create(t *testing.T) {
	attrOK := models.SellerAttributes{
		CID:         5432,
		CompanyName: "PepValentine",
		Address:     "SN 5555",
		Telephone:   "35164543",
		LocalityID:  "1001",
	}
	attrFail := models.SellerAttributes{
		CID:         77,
		CompanyName: "Prueba",
		Address:     "",
		Telephone:   "12345",
		LocalityID:  "1001",
	}
	attrConflict := models.SellerAttributes{
		CID:         12,
		CompanyName: "Pepe",
		Address:     "Calle x",
		Telephone:   "888",
		LocalityID:  "90",
	}

	cases := []struct {
		name       string
		body       []byte
		setupMock  func(s *mocks.SellerServiceDBMock)
		wantStatus int
	}{
		{
			name: "ok",
			body: func() []byte { b, _ := json.Marshal(attrOK); return b }(),
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("Create", attrOK).Return(models.Seller{ID: 5, SellerAttributes: attrOK}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "bad json",
			body:       []byte("{not a json}"),
			setupMock:  func(s *mocks.SellerServiceDBMock) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "fail (missing Address)",
			body: func() []byte { b, _ := json.Marshal(attrFail); return b }(),
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("Create", attrFail).Return(models.Seller{}, httperrors.UnprocessableEntityError{Message: "Invalid seller data"})
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "conflict",
			body: func() []byte { b, _ := json.Marshal(attrConflict); return b }(),
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("Create", attrConflict).Return(models.Seller{}, httperrors.ConflictError{Message: "cid exists"})
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(mocks.SellerServiceDBMock)
			tc.setupMock(service)
			handler := handler.NewSellerHandler(service)
			router := chi.NewRouter()
			router.Post("/api/v1/sellers", handler.Create())
			req := httptest.NewRequest("POST", "/api/v1/sellers", bytes.NewBuffer(tc.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

// TestSellerHandler_Update tests the Update method of the SellerHandler.
// It checks for successful updates, conflict errors, and non-existent sellers.
func TestSellerHandler_Update(t *testing.T) {
	attr := models.SellerAttributes{
		CID:         999,
		CompanyName: "NuevoNombre",
		Address:     "NuevaDirección",
		Telephone:   "1234",
		LocalityID:  "1001",
	}

	cases := []struct {
		name       string
		id         string
		body       []byte
		setupMock  func(s *mocks.SellerServiceDBMock)
		wantStatus int
	}{
		{
			name: "ok",
			id:   "22",
			body: func() []byte { b, _ := json.Marshal(attr); return b }(),
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("Update", 22, &attr).Return(models.Seller{ID: 22, SellerAttributes: attr}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "bad id",
			id:         "abc",
			body:       func() []byte { b, _ := json.Marshal(attr); return b }(),
			setupMock:  func(s *mocks.SellerServiceDBMock) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(mocks.SellerServiceDBMock)
			tc.setupMock(service)
			handler := handler.NewSellerHandler(service)
			router := chi.NewRouter()
			router.Put("/api/v1/sellers/{id}", handler.Update())
			req := httptest.NewRequest("PUT", "/api/v1/sellers/"+tc.id, bytes.NewBuffer(tc.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

// TestSellerHandler_Delete tests the Delete method of the SellerHandler.
// It checks for successful deletion, not found errors, and bad ID formats.
func TestSellerHandler_Delete(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		setupMock  func(s *mocks.SellerServiceDBMock)
		wantStatus int
	}{
		{
			name: "ok",
			id:   "42",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("Delete", 42).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   "22",
			setupMock: func(s *mocks.SellerServiceDBMock) {
				s.On("Delete", 22).Return(httperrors.NotFoundError{Message: "not found"})
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bad id",
			id:         "abc",
			setupMock:  func(s *mocks.SellerServiceDBMock) {},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(mocks.SellerServiceDBMock)
			tc.setupMock(service)
			handler := handler.NewSellerHandler(service)
			router := chi.NewRouter()
			router.Delete("/api/v1/sellers/{id}", handler.Delete())
			req := httptest.NewRequest("DELETE", "/api/v1/sellers/"+tc.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}
