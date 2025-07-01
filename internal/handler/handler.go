package handler

import (
	"errors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/models"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func NewEmployeeHandler(sv service.EmployeeService) EmployeeHandler {
	return EmployeeHandler{sv: sv}
}

type EmployeeHandler struct {
	// sv is the service that will be used by the handler
	sv service.EmployeeService
}

func (h *EmployeeHandler) GetAll() http.HandlerFunc {
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

func (h *EmployeeHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		data, err := h.sv.GetByID(id)
		if err != nil {
			if errors.As(err, &httperrors.NotFoundError{}) {
				statusCode, msg := httperrors.GetErrorData(err)
				response.Error(w, statusCode, msg)
				return
			}
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EmployeeHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee models.Employee
		err := request.JSON(r, &employee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		data, err := h.sv.Create(employee)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EmployeeHandler) Patch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee models.Employee
		err := request.JSON(r, &employee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		idReq := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		dbEmployee, err := h.sv.GetByID(id)
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

		data, err := h.sv.Update(id, dbEmployee)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EmployeeHandler) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.sv.Delete(id)
		if err != nil {
			if errors.As(err, &httperrors.NotFoundError{}) {
				statusCode, msg := httperrors.GetErrorData(err)
				response.Error(w, statusCode, msg)
				return
			}
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    nil,
		})
	}
}
