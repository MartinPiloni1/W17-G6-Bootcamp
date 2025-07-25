package service_test

import (
	"errors"
	mocks "github.com/aaguero_meli/W17-G6-Bootcamp/internal/mocks/repository"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/service"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/require"
)

func getValidEmployeeAttributes() models.EmployeeAttributes {
	return models.EmployeeAttributes{CardNumberID: "10000000", FirstName: "Thomas", LastName: "Shelby", WarehouseID: 1}
}

func TestEmployeeServiceDefault_Create(t *testing.T) {
	attrs := getValidEmployeeAttributes()
	toCreateEmp := models.Employee{EmployeeAttributes: attrs}
	expectedEmp := models.Employee{Id: 1, EmployeeAttributes: attrs}

	t.Run("create_ok: success to create employee", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetAll").Return([]models.Employee{}, nil)
		mockRepo.On("Create", toCreateEmp).Return(expectedEmp, nil)
		emp, err := serviceEmp.Create(attrs)
		require.NoError(t, err)
		require.Equal(t, expectedEmp, emp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_conflict: duplicate card_number_id returns conflict", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetAll").Return([]models.Employee{expectedEmp}, nil)
		emp, err := serviceEmp.Create(attrs)
		require.Error(t, err)
		require.ErrorAs(t, err, &httperrors.ConflictError{})
		require.Equal(t, models.Employee{}, emp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_repo_error: returns error if repo fails on GetAll", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetAll").Return(nil, errors.New("db error"))
		emp, err := serviceEmp.Create(attrs)
		require.Error(t, err)
		require.Equal(t, models.Employee{}, emp)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeServiceDefault_GetAll(t *testing.T) {
	attrs := getValidEmployeeAttributes()
	expectedEmp := models.Employee{Id: 1, EmployeeAttributes: attrs}

	t.Run("find_all: should return all employees", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		employeeA := expectedEmp
		employeeB := models.Employee{Id: 2, EmployeeAttributes: models.EmployeeAttributes{
			CardNumberID: "10000001",
			FirstName:    "Arthur",
			LastName:     "Shelby",
			WarehouseID:  2,
		}}
		employees := []models.Employee{employeeA, employeeB}
		mockRepo.On("GetAll").Return(employees, nil)
		result, err := serviceEmp.GetAll()
		require.NoError(t, err)
		require.Equal(t, employees, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("find_all_error: should return error on repo failure", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetAll").Return(nil, errors.New("db error"))
		result, err := serviceEmp.GetAll()
		require.Error(t, err)
		require.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeServiceDefault_GetByID(t *testing.T) {
	attrs := getValidEmployeeAttributes()
	expectedEmp := models.Employee{Id: 1, EmployeeAttributes: attrs}

	t.Run("find_by_id_non_existent: should return not found", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetByID", 99).Return(models.Employee{}, errors.New("not found"))
		emp, err := serviceEmp.GetByID(99)
		require.Error(t, err)
		require.Equal(t, models.Employee{}, emp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("find_by_id_existent: should return single employee", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetByID", 1).Return(expectedEmp, nil)
		emp, err := serviceEmp.GetByID(1)
		require.NoError(t, err)
		require.Equal(t, expectedEmp, emp)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeServiceDefault_Update(t *testing.T) {
	attrs := getValidEmployeeAttributes()
	expectedEmp := models.Employee{Id: 1, EmployeeAttributes: attrs}

	t.Run("update_existent: successfully updates and returns employee", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetByID", 1).Return(expectedEmp, nil)
		mockRepo.On("GetAll").Return([]models.Employee{expectedEmp}, nil)
		updatedAttrs := models.EmployeeAttributes{
			CardNumberID: "10000002",
			FirstName:    "Tommy",
			LastName:     "Shelby",
			WarehouseID:  1,
		}
		expectedUpdated := models.Employee{Id: 1, EmployeeAttributes: updatedAttrs}
		mockRepo.On("Update", 1, expectedUpdated).Return(expectedUpdated, nil)
		emp, err := serviceEmp.Update(1, updatedAttrs)
		require.NoError(t, err)
		require.Equal(t, expectedUpdated, emp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update_non_existent: returns error if employee does not exist", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetByID", 99).Return(models.Employee{}, errors.New("not found"))
		emp, err := serviceEmp.Update(99, attrs)
		require.Error(t, err)
		require.Equal(t, models.Employee{}, emp)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeServiceDefault_Delete(t *testing.T) {
	t.Run("delete_non_existent: returns error if employee does not exist", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("Delete", 99).Return(errors.New("not found"))
		err := serviceEmp.Delete(99)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete_ok: successfully deletes employee", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("Delete", 1).Return(nil)
		err := serviceEmp.Delete(1)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeServiceDefault_ReportInboundOrders(t *testing.T) {

	t.Run("one_employee_not_found_by_error", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetByID", 123).Return(models.Employee{}, errors.New("db error"))
		res, err := serviceEmp.ReportInboundOrders(123)
		require.Error(t, err)
		require.Nil(t, res)
		require.ErrorAs(t, err, &httperrors.NotFoundError{})
		mockRepo.AssertExpectations(t)
	})

	t.Run("one_employee_not_found_by_zero_id", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetByID", 123).Return(models.Employee{}, nil)
		res, err := serviceEmp.ReportInboundOrders(123)
		require.Error(t, err)
		require.Nil(t, res)
		require.ErrorAs(t, err, &httperrors.NotFoundError{})
		mockRepo.AssertExpectations(t)
	})

	t.Run("count_inbound_orders_for_employee_error", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		emp := models.Employee{Id: 1, EmployeeAttributes: getValidEmployeeAttributes()}
		mockRepo.On("GetByID", 1).Return(emp, nil)
		mockInboundOrderRepo.On("CountInboundOrdersForEmployee", 1).Return(0, errors.New("count error"))
		res, err := serviceEmp.ReportInboundOrders(1)
		require.Error(t, err)
		require.Nil(t, res)
		mockRepo.AssertExpectations(t)
		mockInboundOrderRepo.AssertExpectations(t)
	})

	t.Run("get_all_error_when_all_employees", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		mockRepo.On("GetAll").Return([]models.Employee{}, errors.New("db error"))
		res, err := serviceEmp.ReportInboundOrders(0)
		require.Error(t, err)
		require.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("count_inbound_orders_for_employees_error", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		emps := []models.Employee{
			{Id: 1, EmployeeAttributes: getValidEmployeeAttributes()},
		}
		mockRepo.On("GetAll").Return(emps, nil)
		mockInboundOrderRepo.On("CountInboundOrdersForEmployees").Return(map[int]int{}, errors.New("count error"))
		res, err := serviceEmp.ReportInboundOrders(0)
		require.Error(t, err)
		require.Nil(t, res)
		mockRepo.AssertExpectations(t)
		mockInboundOrderRepo.AssertExpectations(t)
	})

	t.Run("one_employee_report_success", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		employeeID := 42
		emp := models.Employee{
			Id: employeeID,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberID: "abc", FirstName: "Rick", LastName: "Sánchez", WarehouseID: 99,
			},
		}
		count := 123
		mockRepo.On("GetByID", employeeID).Return(emp, nil)
		mockInboundOrderRepo.On("CountInboundOrdersForEmployee", employeeID).Return(count, nil)

		expected := []models.EmployeeWithInboundCount{
			{Id: emp.Id, CardNumberID: emp.CardNumberID, FirstName: emp.FirstName, LastName: emp.LastName, WarehouseID: emp.WarehouseID, InboundOrdersCount: count},
		}

		res, err := serviceEmp.ReportInboundOrders(employeeID)
		require.NoError(t, err)
		require.Equal(t, expected, res)
		mockRepo.AssertExpectations(t)
		mockInboundOrderRepo.AssertExpectations(t)
	})

	t.Run("all_employees_report_success", func(t *testing.T) {
		mockRepo := mocks.MockEmployeeRepository{}
		mockInboundOrderRepo := mocks.MockInboundOrderRepository{}
		serviceEmp := service.NewEmployeeService(&mockRepo, &mockInboundOrderRepo)
		emps := []models.Employee{
			{Id: 1, EmployeeAttributes: models.EmployeeAttributes{
				CardNumberID: "1", FirstName: "A", LastName: "AA", WarehouseID: 1,
			}},
			{Id: 2, EmployeeAttributes: models.EmployeeAttributes{
				CardNumberID: "2", FirstName: "B", LastName: "BB", WarehouseID: 2,
			}},
		}
		counts := map[int]int{1: 10, 2: 20}

		mockRepo.On("GetAll").Return(emps, nil)
		mockInboundOrderRepo.On("CountInboundOrdersForEmployees").Return(counts, nil)

		expected := []models.EmployeeWithInboundCount{
			{Id: 1, CardNumberID: "1", FirstName: "A", LastName: "AA", WarehouseID: 1, InboundOrdersCount: 10},
			{Id: 2, CardNumberID: "2", FirstName: "B", LastName: "BB", WarehouseID: 2, InboundOrdersCount: 20},
		}

		res, err := serviceEmp.ReportInboundOrders(0)
		require.NoError(t, err)
		require.Equal(t, expected, res)
		mockRepo.AssertExpectations(t)
		mockInboundOrderRepo.AssertExpectations(t)
	})
}
