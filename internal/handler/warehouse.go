package handler

import (
	"encoding/json"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type WarehouseHandler struct {
	sv service.WarehouseService
}

func NewWarehouseHandler(sv service.WarehouseService) WarehouseHandler {
	return WarehouseHandler{
		sv: sv,
	}
}
func (h WarehouseHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var warehouse models.WarehouseAttributes
		err := json.NewDecoder(r.Body).Decode(&warehouse)
		if err != nil {
			response.Error(w, 422, "Invalid JSON body")
			return
		}
		warehouseData, err := h.sv.Create(warehouse)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": warehouseData,
		})
	}
}

func (h WarehouseHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := h.sv.GetAll()
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

func (h WarehouseHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, 400, "Invalid Id format")
			return
		}

		warehouseData, err := h.sv.GetByID(id)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": warehouseData,
		})
	}
}

func (h WarehouseHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, 400, "Invalid Id format")
			return
		}

		var updatedWarehouse models.WarehouseAttributes
		err = json.NewDecoder(r.Body).Decode(&updatedWarehouse)
		if err != nil {
			response.Error(w, 422, "Invalid JSON body")
			return
		}

		warehouseData, err := h.sv.Update(id, updatedWarehouse)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": warehouseData,
		})
	}
}

func (h WarehouseHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, 400, "Invalid ID format")
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
