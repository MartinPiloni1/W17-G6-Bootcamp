package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type ProductHandler struct {
	sv service.ProductService
}

func NewProductHandler(sv service.ProductService) *ProductHandler {
	return &ProductHandler{sv: sv}
}

func (h ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newProduct models.ProductAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&newProduct)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		validator := validator.New()
		err = validator.Struct(newProduct)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		productData, err := h.sv.Create(newProduct)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": productData,
		})
	}
}

func (h ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productData, err := h.sv.GetAll()
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": productData,
		})
	}
}

func (h ProductHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		product, err := h.sv.GetByID(id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": product,
		})
	}
}

func (h ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		var updatedProduct models.ProductPatchRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err = dec.Decode(&updatedProduct)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		validator := validator.New()
		err = validator.Struct(updatedProduct)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		product, err := h.sv.Update(id, updatedProduct)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": product,
		})
	}
}

func (h ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		err = h.sv.Delete(id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
