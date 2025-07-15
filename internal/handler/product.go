package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

// ProductHandler handles HTTP requests for product resources.
type ProductHandler struct {
	svc service.ProductService
}

// NewProductHandler constructs a new ProductHandler with the given service.
func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// Create returns an http.HandlerFunc that decodes a JSON payload,
// validates it, delegates creation of a Product to the service layer, and writes
// the appropriate JSON response.
func (h ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var newProduct models.ProductAttributes

		// Decode JSON Body to a ProductAttributes struct
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&newProduct)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// Delete trailing whitespaces before validating
		newProduct.Description = strings.TrimSpace(newProduct.Description)
		newProduct.ProductCode = strings.TrimSpace(newProduct.ProductCode)

		// Validate the ProductAttributes struct
		validator := validator.New()
		err = validator.Struct(newProduct)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		// Delegate creation to the service layer
		productData, err := h.svc.Create(ctx, newProduct)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Write the appropiate JSON Response
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": productData,
		})
	}
}

// GetAll returns an http.HandlerFunc that fetches all products
// from the service layer and writes them as JSON.
func (h ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get all the products from the service layer
		productData, err := h.svc.GetAll(ctx)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Write the JSON Response
		response.JSON(w, http.StatusOK, map[string]any{
			"data": productData,
		})
	}
}

// GetById returns an http.HandlerFunc that parses the product ID
// from the URL, retrieves the product, and writes it as JSON.
func (h ProductHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse the product ID
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// Get the product from the service layer
		product, err := h.svc.GetByID(ctx, id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Write the JSON response
		response.JSON(w, http.StatusOK, map[string]any{
			"data": product,
		})
	}
}

// Update returns an http.HandlerFunc that parses the id URL parameter,
// decodes a partial-product JSON payload, validates it, delegates the update
// to the service layer, and responds with the updated product.
func (h ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse the product ID
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// Decode JSON Body to a ProductPatchRequest struct
		var updatedProduct models.ProductPatchRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&updatedProduct)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// Delete trailing whitespaces before validating
		if updatedProduct.Description != nil {
			*updatedProduct.Description = strings.TrimSpace(*updatedProduct.Description)
		}
		if updatedProduct.ProductCode != nil {
			*updatedProduct.ProductCode = strings.TrimSpace(*updatedProduct.ProductCode)
		}

		// Validate the ProductPatchRequest struct
		validator := validator.New()
		err = validator.Struct(updatedProduct)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		// Delegate update to the service layer
		product, err := h.svc.Update(ctx, id, updatedProduct)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Write the JSON response
		response.JSON(w, http.StatusOK, map[string]any{
			"data": product,
		})
	}
}

// Delete returns an http.HandlerFunc that parses the id URL parameter,
// delegates deletion to the service layer, and responds with no content.
func (h ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse the product ID
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// Delegate deletion to the service layer
		err = h.svc.Delete(ctx, id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		// Write the NoContent header after deleting the resource
		w.WriteHeader(http.StatusNoContent)
	}
}
