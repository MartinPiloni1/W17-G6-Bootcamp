package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/bootcamp-go/web/response"
	"github.com/go-playground/validator"
)

// ProductBatchHandler handles product batch operations
type ProductBatchHandler struct {
	productBatchService service.ProductBatchService
}

// NewProductBatchHandler returns a new ProductBatchHandler
func NewProductBatchHandler(productBatchService service.ProductBatchService) *ProductBatchHandler {
	return &ProductBatchHandler{productBatchService: productBatchService}
}

// Create creates a new product batch in the repository
// @Summary Create a new product batch
// @Description Create a new product batch
// @Tags product-batches
// @Accept json
// @Produce json
// @Param product-batch body models.CreateProductBatchRequest true "Product batch to create"
// @Success 201 {object} models.ProductBatch
// @Router /product-batches [post]
func (handler *ProductBatchHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse request body
		var productBatch models.CreateProductBatchRequest
		// Initialize decoder
		dec := json.NewDecoder(r.Body)
		// Disallow unknown fields
		dec.DisallowUnknownFields()

		// Decode request body
		err := dec.Decode(&productBatch)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		// Validate date format
		validator := validator.New()
		validator.RegisterValidation("date_format", utils.ValidateDateFormat)

		// Validate product batch
		if err := validator.Struct(productBatch.Data); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		// Convert request body to product batch model
		productBatchModel := models.ProductBatchAttibutes{
			BatchNumber:        productBatch.Data.BatchNumber,
			CurrentQuantity:    productBatch.Data.CurrentQuantity,
			CurrentTemperature: productBatch.Data.CurrentTemperature,
			DueDate:            productBatch.Data.DueDate,
			InitialQuantity:    productBatch.Data.InitialQuantity,
			ManufacturingDate:  productBatch.Data.ManufacturingDate,
			ManufacturingHour:  productBatch.Data.ManufacturingHour,
			MinimumTemperature: productBatch.Data.MinimumTemperature,
			ProductID:          productBatch.Data.ProductID,
			SectionID:          productBatch.Data.SectionID,
		}

		// Create product batch in service
		createdProductBatch, err := handler.productBatchService.Create(ctx, productBatchModel)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Return created product batch
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": createdProductBatch,
		})
	}
}