package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type SectionHandler struct {
	sectionService service.SectionServiceInterface
}

func NewSectionHandler(sectionService service.SectionServiceInterface) SectionHandler {
	return SectionHandler{sectionService: sectionService}
}

func (h SectionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.sectionService.GetAll()
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data":    data,
		})
	}
}

func (h *SectionHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "ID inválido")
			return
		}

		data, err := h.sectionService.GetByID(id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *SectionHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "ID inválido")
			return
		}

		err = h.sectionService.Delete(id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusNoContent, map[string]any{
			"data": "Sección eliminada correctamente",
		})
	}
}

func (h *SectionHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var section models.CreateSectionRequest
		err := json.NewDecoder(r.Body).Decode(&section)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
			return
		}

		validator := validator.New()
		if err := validator.Struct(section); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "faltan campos o hay campos inválidos")
			return
		}

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

		createdSection, err := h.sectionService.Create(sectionModel)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": createdSection,
		})
	}
}