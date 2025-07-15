package handler

import (
	"net/http"
	"strconv"

	"encoding/json"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type EmployeeHandler struct {
	sv service.EmployeeService
}

func NewEmployeeHandler(sv service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{sv: sv}
}

// GetAll returns an HTTP handler that retrieves all employees from the service.
// It responds with a JSON object containing a list of employees or an error message if the operation fails.
func (h *EmployeeHandler) GetAll() http.HandlerFunc {
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

// GetById returns an HTTP handler that retrieves a single employee by its ID from the service.
// It parses the ID from the URL, validates it, and responds with a JSON object containing the employee data or an error message if the operation fails.
func (h *EmployeeHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		data, err := h.sv.GetByID(id)
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

// Create returns an HTTP handler that creates a new employee.
// It decodes and validates the request body, then calls the service to persist the employee data.
// The handler responds with the created employee as JSON or returns an appropriate error message.
func (h *EmployeeHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee models.EmployeeAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&employee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		validator := validator.New()
		err = validator.Struct(employee)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		data, err := h.sv.Create(employee)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": data,
		})
	}
}

// Update returns an HTTP handler that updates an existing employee by its ID.
// It decodes and validates the request body, parses the ID from the URL, and updates the employee using the service.
// The handler responds with the updated employee as JSON or an appropriate error message if the operation fails.
func (h *EmployeeHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee models.EmployeeAttributes

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&employee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		validator := validator.New()
		err = validator.Struct(employee)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid JSON body")
			return
		}

		idReq := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		data, err := h.sv.Update(id, employee)
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

// Delete returns an HTTP handler that deletes an employee by its ID.
// It parses the ID from the URL, calls the service to delete the employee, and responds with a 204 No Content status on success or an error message if the operation fails.
func (h *EmployeeHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idReq)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
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

func (h *EmployeeHandler) GetInboundOrderReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employeeID := 0
		param := r.URL.Query().Get("id")
		if param != "" {
			id, err := strconv.Atoi(param)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid employee_id")
				return
			}
			employeeID = id
		}

		data, err := h.sv.ReportInboundOrders(employeeID)
		if err != nil {
			statusCode, msg := httperrors.GetErrorData(err)
			response.Error(w, statusCode, msg)
			return
		}
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"data": data,
		})
	}
}
