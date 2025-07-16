package handler

import (
	"encoding/json"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-playground/validator"
	"net/http"
)

type InboundOrderHandler struct {
	sv service.InboundOrderService
}

func NewInboundOrderHandler(sv service.InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{sv: sv}
}

// Create returns an HTTP handler that creates a new inbound order.
// It decodes and validates the request body, then calls the service to persist the inbound order data.
// The handler responds with the created inbound order as JSON or returns an appropriate error message.
func (h *InboundOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.InboundOrderAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&req)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		data, err := h.sv.Create(req)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": data,
		})
	}
}
