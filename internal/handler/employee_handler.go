package handler

import (
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type EmployeeHandler struct {
	sv service.EmployeeService
}

func NewEmployeeHandler(sv service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{sv: sv}
}

func (h *EmployeeHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.sv.GetAll()
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *EmployeeHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		data, err := h.sv.GetByID(id)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *EmployeeHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee models.EmployeeAttributes
		err := request.JSON(r, &employee)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		data, err := h.sv.Create(employee)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *EmployeeHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee models.EmployeeAttributes
		err := request.JSON(r, &employee)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		idReq := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idReq)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		dbEmployee, err := h.sv.GetByID(id)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}
		if employee.CardNumberID != "" {
			dbEmployee.CardNumberID = employee.CardNumberID
		}
		if employee.FirstName != "" {
			dbEmployee.FirstName = employee.FirstName
		}
		if employee.LastName != "" {
			dbEmployee.LastName = employee.LastName
		}
		if employee.WarehouseID != 0 {
			dbEmployee.WarehouseID = employee.WarehouseID
		}

		data, err := h.sv.Update(id, dbEmployee.EmployeeAttributes)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *EmployeeHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		err = h.sv.Delete(id)
		if err != nil {
			httperrors.RespondError(w, err)
			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}
