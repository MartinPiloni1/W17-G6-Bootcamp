package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

type BuyerHandler struct {
	sv service.BuyerService
}

func NewBuyerHandler(sv service.BuyerService) *BuyerHandler {
	return &BuyerHandler{sv: sv}
}

func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buyerData, err := h.sv.GetAll()
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}
		data := utils.MapToSlice(buyerData)
		slices.SortFunc(data, func(a, b models.Buyer) int {
			return a.Id - b.Id
		})
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *BuyerHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid Id")
			return
		}

		data, err := h.sv.GetByID(id)
		if err != nil {
			if errors.As(err, &httperrors.NotFoundError{}) {
				statusCode, msg := httperrors.GetErrorData(err)
				response.Error(w, statusCode, msg)
				return
			}
			http.Error(w, "Unexpected error at GetByID", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *BuyerHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil || id <= 0 {
			response.Error(w, http.StatusBadRequest, "Invalid Id")
			return
		}

		err = h.sv.Delete(id)
		if err != nil {
			if errors.As(err, &httperrors.NotFoundError{}) {
				statusCode, msg := httperrors.GetErrorData(err)
				response.Error(w, statusCode, msg)
				return
			}
			http.Error(w, "Unexpected error at Delete", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (h *BuyerHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newBuyer models.BuyerAttributes

		err := json.NewDecoder(r.Body).Decode(&newBuyer)
		if err != nil || newBuyer.CardNumberId <= 0 || newBuyer.FirstName == "" || newBuyer.LastName == "" {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		buyer, err := h.sv.Create(newBuyer)
		if err != nil {
			if errors.As(err, &httperrors.ConflictError{}) {
				statusCode, msg := httperrors.GetErrorData(err)
				response.Error(w, statusCode, msg)
				return
			}
			http.Error(w, "Unexpected error at Create", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": buyer,
		})
	}
}
