package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-playground/validator"
)

type PurchaseOrderHandler struct {
	service service.PurchaseOrderService
}

func NewPurchaseOrderHandler(serviceInstance service.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: serviceInstance}
}

func (h *PurchaseOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var newPurchaseOrder models.PurchaseOrderAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // cannot send unexpected fields in request

		err := dec.Decode(&newPurchaseOrder)
		if err != nil {
			fmt.Println(err)
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		// delete leading and last white spacing before validator
		newPurchaseOrder.OrderNumber = strings.TrimSpace(newPurchaseOrder.OrderNumber)
		newPurchaseOrder.TrackingCode = strings.TrimSpace(newPurchaseOrder.TrackingCode)

		// validate the json body with validator
		validator := validator.New()
		err = validator.Struct(newPurchaseOrder)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		purchaseOrder, err := h.service.Create(ctx, newPurchaseOrder)
		if err != nil {
			fmt.Println(err)
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": purchaseOrder,
		})
	}
}
