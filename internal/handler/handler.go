package handler

import (
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
)

type ProductHandler struct {
	sv service.ProductServiceInterface
}

func NewProductHandler(sv service.ProductServiceInterface) ProductHandler {
	return ProductHandler{sv: sv}
}

func (h ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.sv.GetAll()
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
