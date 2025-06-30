package handler

import (
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
)

type BuyerHandler struct {
	sv service.BuyerServiceInterface
}

func NewBuyerHandler(sv service.BuyerServiceInterface) *BuyerHandler {
	return &BuyerHandler{sv: sv}
}

func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
