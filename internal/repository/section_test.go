package repository

import (
	"context"
	"database/sql"
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

func TestSectionRepository_Delete(t *testing.T) {
	inputID := 1
	expectedQuery := regexp.QuoteMeta("DELETE FROM sections WHERE id = ?")

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			testName: "Success: Should delete section correctly",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(inputID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			testName: "Fail: Section not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(inputID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: httperrors.NotFoundError{Message: "Section not found"},
		},
		{
			testName: "Fail: Should return internal server error on exec error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(inputID).
					WillReturnError(errors.New("db error on exec"))
			},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Fail: Should return internal server error on RowsAffected error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(inputID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
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

			err = repo.Delete(context.Background(), inputID)

			assert.Equal(t, tt.expectedError, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestSectionRepository_GetAll(t *testing.T) {
	expectedQuery := regexp.QuoteMeta(`
        SELECT
            id, section_number, current_temperature, minimum_temperature,
            current_capacity, minimum_capacity, maximum_capacity,
            warehouse_id, product_type_id
        FROM sections
    `)

	columns := []string{
		"id", "section_number", "current_temperature", "minimum_temperature",
		"current_capacity", "minimum_capacity", "maximum_capacity",
		"warehouse_id", "product_type_id",
	}

	expectedSections := []models.Section{
		{ID: 1, SectionNumber: "SEC-101", CurrentTemperature: 20, MinimumTemperature: 15, CurrentCapacity: 50, MinimumCapacity: 10, MaximumCapacity: 100, WarehouseID: 1, ProductTypeID: 1},
		{ID: 2, SectionNumber: "SEC-102", CurrentTemperature: 22, MinimumTemperature: 18, CurrentCapacity: 70, MinimumCapacity: 20, MaximumCapacity: 150, WarehouseID: 1, ProductTypeID: 2},
	}

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  []models.Section
		expectedError error
	}{
		{
			testName: "Success: Should return all sections",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow(expectedSections[0].ID, expectedSections[0].SectionNumber, expectedSections[0].CurrentTemperature, expectedSections[0].MinimumTemperature, expectedSections[0].CurrentCapacity, expectedSections[0].MinimumCapacity, expectedSections[0].MaximumCapacity, expectedSections[0].WarehouseID, expectedSections[0].ProductTypeID).
					AddRow(expectedSections[1].ID, expectedSections[1].SectionNumber, expectedSections[1].CurrentTemperature, expectedSections[1].MinimumTemperature, expectedSections[1].CurrentCapacity, expectedSections[1].MinimumCapacity, expectedSections[1].MaximumCapacity, expectedSections[1].WarehouseID, expectedSections[1].ProductTypeID)
				
				mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
			},
			expectedResp:  expectedSections,
			expectedError: nil,
		},
		{
			testName: "Fail: Should return internal server error on query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("db query error"))
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Fail: Should return internal server error on scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).AddRow("invalid_id", "SEC-101", 20, 15, 50, 10, 100, 1, 1)
				mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
        {
            testName: "Fail: Should return internal server error on rows error",
            mockSetup: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows(columns).
                    AddRow(expectedSections[0].ID, expectedSections[0].SectionNumber, expectedSections[0].CurrentTemperature, expectedSections[0].MinimumTemperature, expectedSections[0].CurrentCapacity, expectedSections[0].MinimumCapacity, expectedSections[0].MaximumCapacity, expectedSections[0].WarehouseID, expectedSections[0].ProductTypeID)
                
                rows.CloseError(errors.New("rows iteration error"))
                
                mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
            },
            expectedResp:  nil,
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

			result, err := repo.GetAll(context.Background())

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestSectionRepository_GetByID(t *testing.T) {
	inputID := 1
	expectedQuery := regexp.QuoteMeta(`
        SELECT
            id, section_number, current_temperature, minimum_temperature,
            current_capacity, minimum_capacity, maximum_capacity,
            warehouse_id, product_type_id
        FROM sections
        WHERE id = ?
    `)

	columns := []string{
		"id", "section_number", "current_temperature", "minimum_temperature",
		"current_capacity", "minimum_capacity", "maximum_capacity",
		"warehouse_id", "product_type_id",
	}

	expectedSection := models.Section{
		ID: 1, SectionNumber: "SEC-101", CurrentTemperature: 20, MinimumTemperature: 15,
		CurrentCapacity: 50, MinimumCapacity: 10, MaximumCapacity: 100,
		WarehouseID: 1, ProductTypeID: 1,
	}

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Section
		expectedError error
	}{
		{
			testName: "Success: Should return a section by ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow(expectedSection.ID, expectedSection.SectionNumber, expectedSection.CurrentTemperature, expectedSection.MinimumTemperature, expectedSection.CurrentCapacity, expectedSection.MinimumCapacity, expectedSection.MaximumCapacity, expectedSection.WarehouseID, expectedSection.ProductTypeID)
				
				mock.ExpectQuery(expectedQuery).WithArgs(inputID).WillReturnRows(rows)
			},
			expectedResp:  expectedSection,
			expectedError: nil,
		},
		{
			testName: "Fail: Section not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns)
				mock.ExpectQuery(expectedQuery).WithArgs(inputID).WillReturnRows(rows)
			},
			expectedResp:  models.Section{},
			expectedError: httperrors.NotFoundError{Message: "Section not found"},
		},
		{
			testName: "Fail: Should return internal server error on query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WithArgs(inputID).WillReturnError(errors.New("db query error"))
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

			result, err := repo.GetByID(context.Background(), inputID)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestSectionRepository_GetProductsReport(t *testing.T) {
	inputID := 1
	expectedQuery := regexp.QuoteMeta(`
        SELECT s.id, s.section_number, COUNT(pb.id) as products_count
        FROM sections s
        LEFT JOIN product_batches pb ON s.id = pb.section_id
        WHERE s.id = ?
        GROUP BY s.id, s.section_number;
    `)

	columns := []string{"id", "section_number", "products_count"}

	expectedReport := models.SectionProductsReport{
		SectionID:     1,
		SectionNumber: "SEC-101",
		ProductsCount: 50,
	}

	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.SectionProductsReport
		expectedError error
	}{
		{
			testName: "Success: Should return a product report by section ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns).
					AddRow(expectedReport.SectionID, expectedReport.SectionNumber, expectedReport.ProductsCount)
				mock.ExpectQuery(expectedQuery).WithArgs(inputID).WillReturnRows(rows)
			},
			expectedResp:  expectedReport,
			expectedError: nil,
		},
		{
			testName: "Fail: Section not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WithArgs(inputID).WillReturnError(sql.ErrNoRows)
			},
			expectedResp:  models.SectionProductsReport{},
			expectedError: httperrors.NotFoundError{Message: "Section not found"},
		},
		{
			testName: "Fail: Should return internal server error on scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).WithArgs(inputID).WillReturnError(errors.New("any other scan error"))
			},
			expectedResp:  models.SectionProductsReport{},
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

			result, err := repo.GetProductsReport(context.Background(), inputID)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, result)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}