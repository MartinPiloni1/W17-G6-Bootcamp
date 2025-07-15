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

type BuyerHandler struct {
	service service.BuyerService
}

func NewBuyerHandler(serviceInstance service.BuyerService) *BuyerHandler {
	return &BuyerHandler{service: serviceInstance}
}

func (h *BuyerHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var newBuyer models.BuyerAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // cannot send unexpected fields in request

		err := dec.Decode(&newBuyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// validate the json body with validator
		validator := validator.New()
		err = validator.Struct(newBuyer)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		buyer, err := h.service.Create(ctx, newBuyer)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": buyer,
		})
	}
}

func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buyerData, err := h.service.GetAll(ctx)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": buyerData,
		})
	}
}

func (h *BuyerHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		// id must be uint
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid Id")
			return
		}

		data, err := h.service.GetByID(ctx, id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *BuyerHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)

		// id must be uint
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		var patchReq models.BuyerPatchRequest

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // cannot send unexpected fields in request
		err = dec.Decode(&patchReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// validate the json body with validator
		validate := validator.New()
		err = validate.Struct(patchReq)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		buyer, err := h.service.Update(ctx, id, patchReq)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": buyer,
		})
	}
}

func (h *BuyerHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		// id must be uint
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid Id")
			return
		}

		err = h.service.Delete(ctx, id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}
