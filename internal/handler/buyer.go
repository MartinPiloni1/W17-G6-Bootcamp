package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

type BuyerHandler struct {
	sv service.BuyerServiceInterface
}

func NewBuyerHandler(sv service.BuyerServiceInterface) *BuyerHandler {
	return &BuyerHandler{sv: sv}
}

func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buyerData, err := h.sv.GetAll()
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}
		data := make([]models.Buyer, 0)
		for _, buyer := range buyerData {
			data = append(data, buyer)
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *BuyerHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid Id")
			return
		}

		data, err := h.sv.GetByID(id)
		if err != nil {
			if errors.As(err, &httperrors.NotFoundError{}) {
				statusCode, msg := httperrors.GetErrorData(err)
				response.Error(w, statusCode, msg)
				return
			}
			http.Error(w, "error at GetByID", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}
