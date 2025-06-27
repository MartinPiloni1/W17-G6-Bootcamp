package handler

import (
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
)

type SellerHandler struct {
	sv service.SellerService
}

func NewSellerHandler(sv service.SellerService) *SellerHandler {
	return &SellerHandler{sv: sv}
}

func (h *SellerHandler) Create() http.HandlerFunc {
	return nil
}

func (h *SellerHandler) GetAll() http.HandlerFunc {
	return nil
}

func (h *SellerHandler) GetByID() http.HandlerFunc {
	return nil
}

func (h *SellerHandler) Update() http.HandlerFunc {
	return nil
}

func (h *SellerHandler) Delete() http.HandlerFunc {
	return nil
}
