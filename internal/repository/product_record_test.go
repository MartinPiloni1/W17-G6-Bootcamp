package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// Verifies the behavior of the repository layer responsible for creating a new ProductRecord
func TestProductRecordRepository_Create(t *testing.T) {
	// Define the query and product record used by the tests cases
	query := regexp.QuoteMeta(`
		INSERT INTO product_records (
		 	last_update_date,
			purchase_price,
			sale_price,
			product_id
		) VALUES (
			?, ?, ?, ?
		)
	`)

	newProductRecordAttributes := models.ProductRecordAttributes{
		LastUpdateDate: "2024-06-01",
		PurchasePrice:  150,
		SalePrice:      200,
		ProductID:      1,
	}

	newProductRecord := models.ProductRecord{
		ID:                      1,
		ProductRecordAttributes: newProductRecordAttributes,
	}

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   productAttributes – the input attributes passed to repo.Create()
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedResp      – the ProductRecord value we expect Create() to return
	//   expectedError     – the error we expect Create() to return
	tests := []struct {
		testName                string
		productRecordAttributes models.ProductRecordAttributes
		mockSetup               func(mock sqlmock.Sqlmock)
		expectedResp            models.ProductRecord
		expectedError           error
	}{
		{
			testName:                "Success: Should create product record correctly",
			productRecordAttributes: newProductRecordAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(
						newProductRecordAttributes.LastUpdateDate,
						newProductRecordAttributes.PurchasePrice,
						newProductRecordAttributes.SalePrice,
						newProductRecordAttributes.ProductID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp:  newProductRecord,
			expectedError: nil,
		},
		{
			testName:                "Error case: Non-existent product ID (1452)",
			productRecordAttributes: newProductRecordAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			expectedResp: models.ProductRecord{},
			expectedError: httperrors.ConflictError{
				Message: "a product with the given id does not exist",
			},
		},
		{
			testName:                "Error case: Internal Server Error on other MySQL error",
			productRecordAttributes: newProductRecordAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1146})
			},
			expectedResp:  models.ProductRecord{},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName:                "Error case: Internal Server Error when LastInsertId fails",
			productRecordAttributes: newProductRecordAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnResult(
						sqlmock.NewErrorResult(errors.New("last insert id error")),
					)
			},
			expectedResp:  models.ProductRecord{},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewProductRecordRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			result, err := repo.Create(context.Background(), newProductRecordAttributes)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
