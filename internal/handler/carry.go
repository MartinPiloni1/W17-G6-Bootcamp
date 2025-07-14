package handler

import (
	"encoding/json"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"net/http"
)

type CarryHandler struct {
	sv service.CarryService
}

func NewCarryHandler(sv service.CarryService) *CarryHandler {
	return &CarryHandler{
		sv: sv,
	}
}

func (h CarryHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var carry models.CarryAttributes
		err := json.NewDecoder(r.Body).Decode(&carry)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}
		carryData, err := h.sv.Create(carry)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": carryData,
		})
	}
}
