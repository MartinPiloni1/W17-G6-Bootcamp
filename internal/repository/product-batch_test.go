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

func TestProductBatchRepository_Create(t *testing.T) {
	inputBatch := models.ProductBatchAttibutes{
		BatchNumber:        1243,
		CurrentQuantity:    20,
		CurrentTemperature: 21.5,
		DueDate:            "2024-06-11",
		InitialQuantity:    20,
		ManufacturingDate:  "2024-05-21",
		ManufacturingHour:  10,
		MinimumTemperature: 7.5,
		ProductID:          1,
		SectionID:          2,
	}
	expectedBatch := models.ProductBatch{
		ID:                   1,
		ProductBatchAttibutes: inputBatch,
	}

	expectedQuery := regexp.QuoteMeta(`
        INSERT INTO product_batches (
            batch_number, current_quantity, current_temperature, due_date,
            initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature,
            product_id, section_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.ProductBatch
		expectedError error
	}{
		{
			testName: "Success: Should create product batch correctly",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(
						inputBatch.BatchNumber, inputBatch.CurrentQuantity, inputBatch.CurrentTemperature, inputBatch.DueDate,
						inputBatch.InitialQuantity, inputBatch.ManufacturingDate, inputBatch.ManufacturingHour, inputBatch.MinimumTemperature,
						inputBatch.ProductID, inputBatch.SectionID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp:  expectedBatch,
			expectedError: nil,
		},
		{
			testName: "Fail: Conflict on duplicate batch number (1062)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WillReturnError(&mysql.MySQLError{Number: 1062})
			},
			expectedResp:  models.ProductBatch{},
			expectedError: httperrors.ConflictError{Message: "Batch number already exists."},
		},
		{
			testName: "Fail: Conflict on non-existent product/section (1452)",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			expectedResp:  models.ProductBatch{},
			expectedError: httperrors.ConflictError{Message: "Product or section does not exist."},
		},
		{
			testName: "Fail: Internal Server Error on other mysql error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WillReturnError(&mysql.MySQLError{Number: 1146})
			},
			expectedResp:  models.ProductBatch{},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Fail: Should return internal server error when LastInsertId fails",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))
			},
			expectedResp:  models.ProductBatch{},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Fail: Internal Server Error on generic db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WillReturnError(errors.New("a generic database error"))
			},
			expectedResp:  models.ProductBatch{},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewProductBatchRepositoryDB(db)
			tt.mockSetup(mock)

			result, err := repo.Create(context.Background(), inputBatch)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}