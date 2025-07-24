package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSectionRepository_Create(t *testing.T) {
	inputSection := models.Section{
		SectionNumber:      "SEC-101",
		CurrentTemperature: 20,
		MinimumTemperature: 15,
		CurrentCapacity:    50,
		MinimumCapacity:    10,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	expectedSection := inputSection
	expectedSection.ID = 1

	expectedQuery := regexp.QuoteMeta(`
        INSERT INTO sections (
            section_number, current_temperature, minimum_temperature,
            current_capacity, minimum_capacity, maximum_capacity,
            warehouse_id, product_type_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Section
		expectedError error
	}{
		{
			testName: "Success: Should create section correctly",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(inputSection.SectionNumber, inputSection.CurrentTemperature, inputSection.MinimumTemperature, inputSection.CurrentCapacity, inputSection.MinimumCapacity, inputSection.MaximumCapacity, inputSection.WarehouseID, inputSection.ProductTypeID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp:  expectedSection,
			expectedError: nil,
		},
		{
			testName: "Fail: Conflict on duplicate section number (1062)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(&mysql.MySQLError{Number: 1062})
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.ConflictError{Message: "Section number already exists."},
		},
		{
			testName: "Fail: Conflict on non-existent warehouse (1452)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.ConflictError{Message: "Warehouse does not exist."},
		},
		{
			testName: "Fail: Internal Server Error on other mysql error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(&mysql.MySQLError{Number: 1146})
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Fail: Should return internal server error when LastInsertId fails",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.InternalServerError{},
		},
        {
            testName: "Fail: Internal Server Error on generic db error",
            mockSetup: func(mock sqlmock.Sqlmock) {
                mock.ExpectExec(expectedQuery).
                    WillReturnError(errors.New("a generic database error"))
            },
            expectedResp:  models.Section{},
            expectedError: httperrors.InternalServerError{},
        },
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewSectionRepositoryDB(db)
			
			tt.mockSetup(mock)

			result, err := repo.Create(context.Background(), inputSection)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)
			
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestSectionRepository_Update(t *testing.T) {
	inputID := 1
	sectionToUpdate := models.Section{
		ID:                 1,
		SectionNumber:      "SEC-102",
		CurrentTemperature: 25,
		MinimumTemperature: 18,
		CurrentCapacity:    75,
		MinimumCapacity:    20,
		MaximumCapacity:    150,
		WarehouseID:        2,
		ProductTypeID:      2,
	}

	expectedQuery := regexp.QuoteMeta(`
        UPDATE sections SET
            section_number = ?, current_temperature = ?, minimum_temperature = ?,
            current_capacity = ?, minimum_capacity = ?, maximum_capacity = ?,
            warehouse_id = ?, product_type_id = ?
        WHERE id = ?
    `)

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Section
		expectedError error
	}{
		{
			testName: "Success: Should update section correctly",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(sectionToUpdate.SectionNumber, sectionToUpdate.CurrentTemperature, sectionToUpdate.MinimumTemperature, sectionToUpdate.CurrentCapacity, sectionToUpdate.MinimumCapacity, sectionToUpdate.MaximumCapacity, sectionToUpdate.WarehouseID, sectionToUpdate.ProductTypeID, inputID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedResp:  sectionToUpdate,
			expectedError: nil,
		},
		{
			testName: "Fail: Conflict on duplicate section number (1062)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(&mysql.MySQLError{Number: 1062})
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.ConflictError{Message: "Section number already exists."},
		},
		{
			testName: "Fail: Conflict on non-existent warehouse (1452)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.ConflictError{Message: "Warehouse does not exist."},
		},
		{
			testName: "Fail: Internal Server Error on other mysql error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(&mysql.MySQLError{Number: 1146})
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Fail: Internal Server Error on generic db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).WillReturnError(errors.New("any database error"))
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewSectionRepositoryDB(db)
			
			tt.mockSetup(mock)

			result, err := repo.Update(context.Background(), inputID, sectionToUpdate)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)
			
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}