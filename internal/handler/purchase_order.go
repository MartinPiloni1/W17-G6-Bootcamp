package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
)

type PurchaseOrderHandler struct {
	service service.PurchaseOrderService
}

func NewPurchaseOrderHandler(serviceInstance service.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: serviceInstance}
}

func (h *PurchaseOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()

		var purchaseOrder models.PurchaseOrderAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // cannot send unexpected fields in request

		err := dec.Decode(&purchaseOrder)
		if err != nil {

		}

	}
}
