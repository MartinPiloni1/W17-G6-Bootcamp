package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-playground/validator"
)

// ProductRecordHandler handles HTTP requests for product record resources.
type ProductRecordHandler struct {
	svc service.ProductRecordService
}

// NewProductRecordHandler constructs a new ProductRecordHandler with the given service.
func NewProductRecordHandler(svc service.ProductRecordService) *ProductRecordHandler {
	return &ProductRecordHandler{svc: svc}
}

// Create returns an http.HandlerFunc that decodes a JSON payload,
// validates it, delegates creation of a new product record to the service layer,
// and writes the appropriate JSON response.
func (h ProductRecordHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var newProductRecord models.ProductRecordAttributes

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&newProductRecord)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// Validates that LastUpdateDate is in a valid date format
		_, err = time.Parse("2006-01-02", newProductRecord.LastUpdateDate)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		validator := validator.New()
		err = validator.Struct(newProductRecord)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		productRecordData, err := h.svc.Create(ctx, newProductRecord)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": productRecordData,
		})
	}
}
