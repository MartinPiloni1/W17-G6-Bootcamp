package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

// SectionHandler handles section operations
type SectionHandler struct {
	sectionService service.SectionService
}

// NewSectionHandler returns a new SectionHandler
func NewSectionHandler(sectionService service.SectionService) *SectionHandler {
	return &SectionHandler{sectionService: sectionService}
}

// GetAll returns all sections
// @Summary Get all sections
// @Description Get all sections
// @Tags sections
// @Accept json
// @Produce json
// @Success 200 {array} models.Section
// @Router /sections [get]
func (handler *SectionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Get all sections from the repository
		data, err := handler.sectionService.GetAll(ctx)
		// Check for errors
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Return all sections
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

// GetByID returns a section by its ID
// @Summary Get a section by ID
// @Description Get a section by ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 200 {object} models.Section
// @Router /sections/{id} [get]
func (handler *SectionHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Get section ID from URL parameter
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		// Check for invalid ID
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// Get section by ID from service
		data, err := handler.sectionService.GetByID(ctx, id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Return section
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

// Delete deletes a section from the repository
// @Summary Delete a section
// @Description Delete a section
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 204
// @Router /sections/{id} [delete]
func (handler *SectionHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Get section ID from URL parameter
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		// Check for invalid ID
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// Delete section from service
		err = handler.sectionService.Delete(ctx, id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Create creates a new section in the repository
// @Summary Create a new section
// @Description Create a new section
// @Tags sections
// @Accept json
// @Produce json
// @Param section body models.CreateSectionRequest true "Section to create"
// @Success 201 {object} models.Section
// @Router /sections [post]
func (handler *SectionHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse request body
		var section models.CreateSectionRequest
		// Initialize decoder
		dec := json.NewDecoder(r.Body)
		// Disallow unknown fields
		dec.DisallowUnknownFields()

		// Decode request body
		err := dec.Decode(&section)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid body")
			return
		}

		// Validate section
		validator := validator.New()
		if err := validator.Struct(section); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		// Convert request body to section model
		sectionModel := models.Section{
			SectionNumber:      section.SectionNumber,
			CurrentTemperature: section.CurrentTemperature,
			MinimumTemperature: section.MinimumTemperature,
			CurrentCapacity:    section.CurrentCapacity,
			MinimumCapacity:    section.MinimumCapacity,
			MaximumCapacity:    section.MaximumCapacity,
			WarehouseID:        section.WarehouseID,
			ProductTypeID:      section.ProductTypeID,
		}

		// Create section in service
		createdSection, err := handler.sectionService.Create(ctx, sectionModel)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Return created section
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": createdSection,
		})
	}
}

// Update updates a section in the repository
// @Summary Update a section
// @Description Update a section
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Param section body models.UpdateSectionRequest true "Section to update"
// @Success 200 {object} models.Section
// @Router /sections/{id} [patch]
func (handler *SectionHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// Parse request body
		var req models.UpdateSectionRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		// Decode request body
		err = dec.Decode(&req)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid body")
			return
		}

		// Validate request body
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		// Update section in service
		updatedSection, err := handler.sectionService.Update(ctx, id, req)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Return updated section
		response.JSON(w, http.StatusOK, map[string]any{
			"data": updatedSection,
		})
	}
}
