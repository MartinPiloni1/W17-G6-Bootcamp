package handler

import (
	"net/http"
	"strconv"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
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