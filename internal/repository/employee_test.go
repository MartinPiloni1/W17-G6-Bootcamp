package repository_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/stretchr/testify/require"
)

func TestEmployeeRepositoryDB_Create(t *testing.T) {
	query := regexp.QuoteMeta(`
		INSERT INTO employees (
			card_number_id,
			first_name,
			last_name,
			warehouse_id
		) VALUES (?, ?, ?, ?)
	`)

	input := models.Employee{
		EmployeeAttributes: models.EmployeeAttributes{
			CardNumberID: "EMP123",
			FirstName:    "Juan",
			LastName:     "Pérez",
			WarehouseID:  4,
		},
	}
	expected := models.Employee{
		Id:                 1,
		EmployeeAttributes: input.EmployeeAttributes,
	}

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Employee
		expectedError error
	}{
		{
			testName: "Success: Employee is created successfully",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(input.EmployeeAttributes.CardNumberID, input.EmployeeAttributes.FirstName, input.EmployeeAttributes.LastName, input.EmployeeAttributes.WarehouseID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp:  expected,
			expectedError: nil,
		},
		{
			testName: "Error: Exec fails on insert statement",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(input.EmployeeAttributes.CardNumberID, input.EmployeeAttributes.FirstName, input.EmployeeAttributes.LastName, input.EmployeeAttributes.WarehouseID).
					WillReturnError(errors.New("db error"))
			},
			expectedResp:  models.Employee{},
			expectedError: errors.New("db error"),
		},
		{
			testName: "Error: LastInsertId fails",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(input.EmployeeAttributes.CardNumberID, input.EmployeeAttributes.FirstName, input.EmployeeAttributes.LastName, input.EmployeeAttributes.WarehouseID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("lastinsertid error")))
			},
			expectedResp:  models.Employee{},
			expectedError: errors.New("lastinsertid error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewEmployeeRepository(db)
			tc.mockSetup(mock)

			result, err := repo.Create(input)
			require.Equal(t, tc.expectedResp, result)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepositoryDB_GetAll(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT
			id,
			card_number_id,
			first_name,
			last_name,
			warehouse_id
		FROM employees
	`)

	employees := []models.Employee{
		{
			Id: 1,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberID: "AAA11",
				FirstName:    "Ana",
				LastName:     "López",
				WarehouseID:  1,
			},
		},
		{
			Id: 2,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberID: "BBB22",
				FirstName:    "Pedro",
				LastName:     "Paz",
				WarehouseID:  2,
			},
		},
	}

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  []models.Employee
		expectedError error
	}{
		{
			testName: "Success: Should retrieve all employees",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
					AddRow(1, "AAA11", "Ana", "López", 1).
					AddRow(2, "BBB22", "Pedro", "Paz", 2)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedResp:  employees,
			expectedError: nil,
		},
		{
			testName: "Error: Query execution fails",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WillReturnError(errors.New("query error"))
			},
			expectedResp:  nil,
			expectedError: errors.New("query error"),
		},
		{
			testName: "Error: Scan fails due to invalid type",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
					AddRow("invalid", "AAA11", "Ana", "López", 1)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedResp:  nil,
			expectedError: errors.New("sql: Scan error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewEmployeeRepository(db)
			tc.mockSetup(mock)

			result, err := repo.GetAll()

			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResp, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepositoryDB_GetByID(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT
			id,
			card_number_id,
			first_name,
			last_name,
			warehouse_id
		FROM employees
		WHERE id = ?
	`)
	employee := models.Employee{
		Id: 15,
		EmployeeAttributes: models.EmployeeAttributes{
			CardNumberID: "X1122",
			FirstName:    "María",
			LastName:     "Gómez",
			WarehouseID:  8,
		},
	}

	tests := []struct {
		testName      string
		inputID       int
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Employee
		expectedError error
	}{
		{
			testName: "Success: Employee found by ID",
			inputID:  15,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
					AddRow(employee.Id, employee.EmployeeAttributes.CardNumberID, employee.EmployeeAttributes.FirstName, employee.EmployeeAttributes.LastName, employee.EmployeeAttributes.WarehouseID)
				mock.ExpectQuery(query).WithArgs(15).WillReturnRows(rows)
			},
			expectedResp:  employee,
			expectedError: nil,
		},
		{
			testName: "Error: Employee not found (0 rows)",
			inputID:  123,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
				mock.ExpectQuery(query).WithArgs(123).WillReturnRows(rows)
			},
			expectedResp:  models.Employee{},
			expectedError: httperrors.NotFoundError{Message: "employee not found"},
		},
		{
			testName: "Error: Query execution fails",
			inputID:  333,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WithArgs(333).WillReturnError(errors.New("query error"))
			},
			expectedResp:  models.Employee{},
			expectedError: errors.New("query error"),
		},
		{
			testName: "Error: Scan fails due to invalid type",
			inputID:  888,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
					AddRow("invalid", "XYZ", "One", "Two", 1)
				mock.ExpectQuery(query).WithArgs(888).WillReturnRows(rows)
			},
			expectedResp:  models.Employee{},
			expectedError: errors.New("sql: Scan error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewEmployeeRepository(db)
			tc.mockSetup(mock)

			result, err := repo.GetByID(tc.inputID)
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedResp, result)
				// Use IsType for custom error types like NotFoundError
				if nf, ok := tc.expectedError.(httperrors.NotFoundError); ok {
					require.IsType(t, nf, err)
				} else {
					require.Error(t, err)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResp, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepositoryDB_Update(t *testing.T) {
	query := regexp.QuoteMeta(`
		UPDATE employees
		SET
			card_number_id = ?,
			first_name = ?,
			last_name = ?,
			warehouse_id = ?
		WHERE id = ?
	`)
	newAttrs := models.EmployeeAttributes{
		CardNumberID: "CC33",
		FirstName:    "Carlos",
		LastName:     "Sosa",
		WarehouseID:  10,
	}
	updated := models.Employee{Id: 21, EmployeeAttributes: newAttrs}

	tests := []struct {
		testName      string
		inputID       int
		inputAttrs    models.EmployeeAttributes
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Employee
		expectedError error
	}{
		{
			testName:   "Success: Employee updated successfully",
			inputID:    21,
			inputAttrs: newAttrs,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(newAttrs.CardNumberID, newAttrs.FirstName, newAttrs.LastName, newAttrs.WarehouseID, 21).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedResp:  updated,
			expectedError: nil,
		},
		{
			testName:   "Error: Employee not found/0 rows affected",
			inputID:    80,
			inputAttrs: newAttrs,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(newAttrs.CardNumberID, newAttrs.FirstName, newAttrs.LastName, newAttrs.WarehouseID, 80).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedResp:  models.Employee{},
			expectedError: httperrors.NotFoundError{Message: "employee not found"},
		},
		{
			testName:   "Error: Exec fails during update",
			inputID:    50,
			inputAttrs: newAttrs,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(newAttrs.CardNumberID, newAttrs.FirstName, newAttrs.LastName, newAttrs.WarehouseID, 50).
					WillReturnError(errors.New("update error"))
			},
			expectedResp:  models.Employee{},
			expectedError: errors.New("update error"),
		},
		{
			testName:   "Error: RowsAffected returns an error",
			inputID:    22,
			inputAttrs: newAttrs,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(newAttrs.CardNumberID, newAttrs.FirstName, newAttrs.LastName, newAttrs.WarehouseID, 22).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsaffected error")))
			},
			expectedResp:  models.Employee{},
			expectedError: errors.New("rowsaffected error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewEmployeeRepository(db)
			tc.mockSetup(mock)

			in := models.Employee{EmployeeAttributes: tc.inputAttrs}
			result, err := repo.Update(tc.inputID, in)
			if tc.expectedError != nil {
				if nf, ok := tc.expectedError.(httperrors.NotFoundError); ok {
					require.IsType(t, nf, err)
				} else {
					require.Error(t, err)
				}
				require.Equal(t, tc.expectedResp, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResp, result)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepositoryDB_Delete(t *testing.T) {
	query := regexp.QuoteMeta(`
		DELETE FROM employees
		WHERE id = ?
	`)

	tests := []struct {
		testName      string
		inputID       int
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			testName: "Success: Employee deleted",
			inputID:  99,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(99).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			testName: "Error: Employee not found/0 rows affected",
			inputID:  88,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(88).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: httperrors.NotFoundError{Message: "employee not found"},
		},
		{
			testName: "Error: Exec fails during delete",
			inputID:  77,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(77).
					WillReturnError(errors.New("delete error"))
			},
			expectedError: errors.New("delete error"),
		},
		{
			testName: "Error: RowsAffected returns an error",
			inputID:  66,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(query).
					WithArgs(66).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rowsaffected error")))
			},
			expectedError: errors.New("rowsaffected error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			repo := repository.NewEmployeeRepository(db)
			tc.mockSetup(mock)

			err = repo.Delete(tc.inputID)
			if tc.expectedError != nil {
				if nf, ok := tc.expectedError.(httperrors.NotFoundError); ok {
					require.IsType(t, nf, err)
				} else {
					require.Error(t, err)
				}
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
