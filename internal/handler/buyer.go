package handler

import (
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
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

		response.JSON(w, http.StatusOK, map[string]any{
			"data": buyerData,
		})
	}
}
