package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
)

type LocalityHandler struct {
	service service.LocalityService
}

func NewLocalityHandler(service service.LocalityService) *LocalityHandler {
	return &LocalityHandler{service: service}
}

// This function create a new locality
// It expects a JSON body with the locality data
// Returns the created locality with a 201 status code
func (h *LocalityHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Data models.Locality `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		created, err := h.service.Create(body.Data)
		if err != nil {
			status, msg := httperrors.GetErrorData(err)
			response.Error(w, status, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": created,
		})
	}
}

// This function retrieves a locality by its ID
// It expects the ID as a URL parameter
func (h *LocalityHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			response.Error(w, http.StatusBadRequest, "Missing id")
			return
		}

		locality, err := h.service.GetByID(id)
		if err != nil {
			status, msg := httperrors.GetErrorData(err)
			response.Error(w, status, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": locality,
		})
	}
}

// This function retrieves a report of sellers by locality
// If an ID is provided, it returns the report for that specific locality
func (h *LocalityHandler) GetSellerReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id *string
		idStr := r.URL.Query().Get("id")
		if idStr != "" {
			id = &idStr
		}

		reports, err := h.service.GetSellerReport(id)
		if err != nil {
			status, msg := httperrors.GetErrorData(err)
			response.Error(w, status, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": reports,
		})
	}
}

// GetReportByLocalityId handles the retrieval of a report by locality ID.
func (h *LocalityHandler) GetReportByLocalityId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localityId := r.URL.Query().Get("id")
		result, err := h.service.GetReportByLocalityId(localityId)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}
		if result == nil && localityId == "" {
			result = make([]models.CarryReport, 0)
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": result,
		})
	}
}
