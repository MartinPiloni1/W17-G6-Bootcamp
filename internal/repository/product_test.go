package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verifies the behavior of the repository layer responsible for creating a new Product. It covers:
//   - Successful creation of a new Product.
//   - Conflict when the given product_code is duplicated (MySQL error 1062).
//   - Conflict when the given seller_id does not exist (MySQL error 1452).
//   - InternalServerError on any other MySQL error code.
//   - InternalServerError when LastInsertId itself fails.
func TestProductRepository_Create(t *testing.T) {
	// Define the query and product used by the tests cases
	query := regexp.QuoteMeta(`
		INSERT INTO products (
			description,
			expiration_rate,
			freezing_rate,
			height,
			length,
			width,
			netweight,
			product_code,
			recommended_freezing_temperature,
			product_type_id,
			seller_id
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`)

	newProductAttributes := models.ProductAttributes{
		Description:                    "Yogurt helado",
		ExpirationRate:                 7,
		FreezingRate:                   2,
		Height:                         10.5,
		Length:                         20.0,
		Width:                          15.0,
		NetWeight:                      1.2,
		ProductCode:                    "YOG01",
		RecommendedFreezingTemperature: -5.0,
		ProductTypeID:                  3,
		SellerID:                       utils.Ptr(1),
	}

	newProduct := models.Product{
		ID:                1,
		ProductAttributes: newProductAttributes,
	}

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   productAttributes – the input attributes passed to repo.Create()
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedResp      – the Product value we expect Create() to return
	//   expectedError     – the error we expect Create() to return
	tests := []struct {
		testName          string
		productAttributes models.ProductAttributes
		mockSetup         func(mock sqlmock.Sqlmock)
		expectedResp      models.Product
		expectedError     error
	}{
		{
			testName:          "Success: Should create product correctly",
			productAttributes: newProductAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(
						newProductAttributes.Description,
						newProductAttributes.ExpirationRate,
						newProductAttributes.FreezingRate,
						newProductAttributes.Height,
						newProductAttributes.Length,
						newProductAttributes.Width,
						newProductAttributes.NetWeight,
						newProductAttributes.ProductCode,
						newProductAttributes.RecommendedFreezingTemperature,
						newProductAttributes.ProductTypeID,
						newProductAttributes.SellerID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp:  newProduct,
			expectedError: nil,
		},
		{
			testName:          "Error case: Duplicated product code (1062)",
			productAttributes: newProductAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1062})
			},
			expectedResp: models.Product{},
			expectedError: httperrors.ConflictError{
				Message: "A product with the given product code already exists",
			},
		},
		{
			testName:          "Error case: Non-existent seller ID (1452)",
			productAttributes: newProductAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			expectedResp: models.Product{},
			expectedError: httperrors.ConflictError{
				Message: "The given seller id does not exists",
			},
		},
		{
			testName:          "Error case: Internal Server Error on other MySQL error",
			productAttributes: newProductAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1146})
			},
			expectedResp:  models.Product{},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName:          "Error case: Internal Server Error when LastInsertId fails",
			productAttributes: newProductAttributes,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnResult(
						sqlmock.NewErrorResult(errors.New("last insert id error")),
					)
			},
			expectedResp:  models.Product{},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewProductRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			result, err := repo.Create(context.Background(), newProductAttributes)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

// Verifies the behavior of the repository layer responsible for fetching all Products. It covers:
//   - Successful retrieval of all rows and mapping into []Product.
//   - InternalServerError when the SQL query itself fails.
//   - InternalServerError when scanning a row into the Product struct fails.
//   - InternalServerError when rows iteration (rows.Err()) fails.
func TestProductRepository_GetAll(t *testing.T) {
	// Define the query and product used by the tests cases
	query := regexp.QuoteMeta(`
		SELECT
			id,
			description,
			expiration_rate,
			freezing_rate,
			height,
			length,
			width,
			netweight,
			product_code,
			recommended_freezing_temperature,
			product_type_id,
			seller_id
		FROM products 
	`)

	columns := []string{
		"id",
		"description",
		"expiration_rate",
		"freezing_rate",
		"height",
		"length",
		"width",
		"netweight",
		"product_code",
		"recommended_freezing_temperature",
		"product_type_id",
		"seller_id",
	}

	products := []models.Product{
		{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				Description:                    "Yogurt helado",
				ExpirationRate:                 7,
				FreezingRate:                   2,
				Height:                         10.5,
				Length:                         20.0,
				Width:                          15.0,
				NetWeight:                      1.2,
				ProductCode:                    "YOG01",
				RecommendedFreezingTemperature: -5.0,
				ProductTypeID:                  3,
				SellerID:                       utils.Ptr(1),
			},
		},
		{
			ID: 2,
			ProductAttributes: models.ProductAttributes{
				Description:                    "Pechuga de pollo",
				ExpirationRate:                 3,
				FreezingRate:                   1,
				Height:                         5.0,
				Length:                         25.0,
				Width:                          12.5,
				NetWeight:                      0.8,
				ProductCode:                    "POLLO01",
				RecommendedFreezingTemperature: 0.0,
				ProductTypeID:                  7,
				SellerID:                       utils.Ptr(2),
			},
		},
	}

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedResp      – the Products slice we expect GetAll() to return
	//   expectedError     – the error we expect GetAll() to return
	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  []models.Product
		expectedError error
	}{
		{
			testName: "Success: Should return all products",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						products[0].ID,
						products[0].Description,
						products[0].ExpirationRate,
						products[0].FreezingRate,
						products[0].Height,
						products[0].Length,
						products[0].Width,
						products[0].NetWeight,
						products[0].ProductCode,
						products[0].RecommendedFreezingTemperature,
						products[0].ProductTypeID,
						products[0].SellerID,
					).
					AddRow(
						products[1].ID,
						products[1].Description,
						products[1].ExpirationRate,
						products[1].FreezingRate,
						products[1].Height,
						products[1].Length,
						products[1].Width,
						products[1].NetWeight,
						products[1].ProductCode,
						products[1].RecommendedFreezingTemperature,
						products[1].ProductTypeID,
						products[1].SellerID,
					)
				mock.
					ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectedResp:  products,
			expectedError: nil,
		},
		{
			testName: "Error case: Internal Server Error on query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(query).
					WillReturnError(errors.New("db query error"))
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Error case: Internal Server Error on scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						"invalid_id",
						"Pechuga de pollo",
						3,
						1,
						5.0,
						25.0,
						12.5,
						0.8,
						"POLLO01",
						0.0,
						7,
						2,
					)

				mock.
					ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Error case: Internal Server Error on rows error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						products[0].ID,
						products[0].Description,
						products[0].ExpirationRate,
						products[0].FreezingRate,
						products[0].Height,
						products[0].Length,
						products[0].Width,
						products[0].NetWeight,
						products[0].ProductCode,
						products[0].RecommendedFreezingTemperature,
						products[0].ProductTypeID,
						products[0].SellerID,
					)

				rows.CloseError(errors.New("rows iteration error"))
				mock.
					ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewProductRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			result, err := repo.GetAll(context.Background())

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

// Verifies the behavior of the repository layer responsible for fetching a Product by its ID. It covers:
//   - Successful retrieval of an existing Product.
//   - NotFoundError when no row exists for the given ID.
//   - InternalServerError when the SQL query itself fails.
func TestProductRepository_GetByID(t *testing.T) {
	// Define the query, product and ID used by the tests cases
	query := regexp.QuoteMeta(`
		SELECT
			id,
			description,
			expiration_rate,
			freezing_rate,
			height,
			length,
			width,
			netweight,
			product_code,
			recommended_freezing_temperature,
			product_type_id,
			seller_id
		FROM products 
		WHERE id = ?
	`)

	columns := []string{
		"id",
		"description",
		"expiration_rate",
		"freezing_rate",
		"height",
		"length",
		"width",
		"netweight",
		"product_code",
		"recommended_freezing_temperature",
		"product_type_id",
		"seller_id",
	}

	inputID := 1

	product := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Yogurt helado",
			ExpirationRate:                 7,
			FreezingRate:                   2,
			Height:                         10.5,
			Length:                         20.0,
			Width:                          15.0,
			NetWeight:                      1.2,
			ProductCode:                    "YOG01",
			RecommendedFreezingTemperature: -5.0,
			ProductTypeID:                  3,
			SellerID:                       utils.Ptr(1),
		},
	}

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedResp      – the Product we expect GetByID() to return
	//   expectedError     – the error we expect GetByID() to return
	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  models.Product
		expectedError error
	}{
		{
			testName: "Success: Should return a product by ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						product.ID,
						product.Description,
						product.ExpirationRate,
						product.FreezingRate,
						product.Height,
						product.Length,
						product.Width,
						product.NetWeight,
						product.ProductCode,
						product.RecommendedFreezingTemperature,
						product.ProductTypeID,
						product.SellerID,
					)

				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnRows(rows)
			},
			expectedResp:  product,
			expectedError: nil,
		},
		{
			testName: "Error case: Product not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns)
				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnRows(rows)
			},
			expectedResp:  models.Product{},
			expectedError: httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName: "Error case: Internal Server Error on query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnError(errors.New("db query error"))
			},
			expectedResp:  models.Product{},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewProductRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			result, err := repo.GetByID(context.Background(), inputID)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

// Verifies the behavior of the repository layer responsible for retrieving record counts per Product.
// It covers:
//   - Successful retrieval of record count for a single Product by its ID.
//   - Successful retrieval of record counts for all Products.
//   - NotFoundError when querying a Product ID yields no rows.
//   - InternalServerError when scanning a row into the target struct fails.
//   - InternalServerError when the SQL query itself fails.
func TestProductRepository_GetRecordPerProduct(t *testing.T) {
	// Define the query, record and ID used by the tests cases
	query := regexp.QuoteMeta(`
			SELECT 
				p.id, 
				p.description,
				COUNT(pr.id) AS records_count
			FROM products p
			LEFT JOIN product_records pr ON pr.product_id = p.id
	`)

	inputID := utils.Ptr(1)

	columns := []string{
		"product_id",
		"description",
		"records_count",
	}

	expectedMultiRecord := []models.ProductRecordCount{
		{
			ProductID:    1,
			Description:  "Yogurt helado",
			RecordsCount: 3,
		},
		{
			ProductID:    2,
			Description:  "Pechuga de pollo",
			RecordsCount: 1,
		},
	}

	expectedSingleRecord := []models.ProductRecordCount{
		{
			ProductID:    2,
			Description:  "Pechuga de pollo",
			RecordsCount: 1,
		},
	}

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedResp      – the records per product count we expect GetRecordPerProduct() to return
	//   expectedError     – the error we expect GetRecordPerProduct() to return
	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedResp  []models.ProductRecordCount
		expectedError error
	}{
		{
			testName: "Success: Should return a product record count by product ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						expectedSingleRecord[0].ProductID,
						expectedSingleRecord[0].Description,
						expectedSingleRecord[0].RecordsCount,
					)

				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnRows(rows)
			},
			expectedResp:  expectedSingleRecord,
			expectedError: nil,
		},
		{
			testName: "Success: Should return all product reports",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						expectedMultiRecord[0].ProductID,
						expectedMultiRecord[0].Description,
						expectedMultiRecord[0].RecordsCount,
					).
					AddRow(
						expectedMultiRecord[1].ProductID,
						expectedMultiRecord[1].Description,
						expectedMultiRecord[1].RecordsCount,
					)

				mock.
					ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectedResp:  expectedMultiRecord,
			expectedError: nil,
		},
		{
			testName: "Error case: Product not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(columns)

				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnRows(rows)
			},
			expectedResp:  nil,
			expectedError: httperrors.NotFoundError{Message: "Product not found"},
		},
		{
			testName: "Error case: Internal Server Error on scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows(columns).
					AddRow(
						"anInvalidID",
						expectedSingleRecord[0].Description,
						expectedSingleRecord[0].RecordsCount,
					)

				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnRows(rows)
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Error case: Internal Server Error on query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(query).
					WithArgs(inputID).
					WillReturnError(errors.New("query error"))
			},
			expectedResp:  nil,
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewProductRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			result, err := repo.GetRecordsPerProduct(context.Background(), inputID)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

// Verifies the behavior of the repository layer responsible for updating an existing Product. It covers:
//   - Successful update of a Product.
//   - Conflict when the given product_code is duplicated (MySQL error 1062).
//   - Conflict when the given seller_id does not exist (MySQL error 1452).
//   - InternalServerError on any other database error.
func TestProductRepository_Update(t *testing.T) {
	// Define the query, product and ID used by the tests cases
	query := regexp.QuoteMeta(`
		UPDATE products
		SET
			description                         = ?,
			expiration_rate                     = ?,
			freezing_rate                       = ?,
			height                              = ?,
			length                              = ?,
			width                               = ?,
			netweight                           = ?,
			product_code                        = ?,
			recommended_freezing_temperature    = ?,
			product_type_id                     = ?,
			seller_id                           = ?
		WHERE id = ?
    `)

	updatedProduct := models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			Description:                    "Pechuga de pollo",
			ExpirationRate:                 3,
			FreezingRate:                   1,
			Height:                         5.0,
			Length:                         25.0,
			Width:                          12.5,
			NetWeight:                      0.8,
			ProductCode:                    "POLLO01",
			RecommendedFreezingTemperature: 0.0,
			ProductTypeID:                  7,
			SellerID:                       utils.Ptr(1),
		},
	}

	inputID := 1

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedResp      – the updated product we expect Update() to return
	//   expectedError     – the error we expect Update() to return
	tests := []struct {
		testName       string
		id             int
		updatedProduct models.Product
		mockSetup      func(mock sqlmock.Sqlmock)
		expectedResp   models.Product
		expectedError  error
	}{
		{
			testName:       "Success: Should update product correctly",
			id:             inputID,
			updatedProduct: updatedProduct,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(
						updatedProduct.Description,
						updatedProduct.ExpirationRate,
						updatedProduct.FreezingRate,
						updatedProduct.Height,
						updatedProduct.Length,
						updatedProduct.Width,
						updatedProduct.NetWeight,
						updatedProduct.ProductCode,
						updatedProduct.RecommendedFreezingTemperature,
						updatedProduct.ProductTypeID,
						updatedProduct.SellerID,
						inputID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedResp:  updatedProduct,
			expectedError: nil,
		},
		{
			testName:       "Fail: Duplicated product code (1062)",
			id:             inputID,
			updatedProduct: updatedProduct,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1062})
			},
			expectedResp: models.Product{},
			expectedError: httperrors.ConflictError{
				Message: "A product with the given product code already exists",
			},
		},
		{
			testName:       "Fail: Non-existent seller (1452)",
			id:             inputID,
			updatedProduct: updatedProduct,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(&mysql.MySQLError{Number: 1452})
			},
			expectedResp: models.Product{},
			expectedError: httperrors.ConflictError{
				Message: "The given seller id does not exists",
			},
		},
		{
			testName:       "Fail: Internal Server Error on other error",
			id:             inputID,
			updatedProduct: updatedProduct,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WillReturnError(errors.New("any database error"))
			},
			expectedResp:  models.Product{},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewProductRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			result, err := repo.Update(context.Background(), tc.id, tc.updatedProduct)

			// Assert
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedResp, result)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

// Verifies the behavior of the repository layer responsible for deleting a Product by its ID. It covers:
//   - Successful deletion of an existing Product.
//   - Conflict when the Product is still referenced by some product_records (MySQL error 1451).
//   - InternalServerError when the Exec call itself fails.
//   - NotFoundError when no rows are affected (Product not found).
//   - InternalServerError when retrieving RowsAffected fails.
func TestProductRepository_Delete(t *testing.T) {
	// Define the query and ID used by the tests cases
	query := regexp.QuoteMeta(`
		DELETE FROM products
		WHERE id=?	
	`)

	inputID := 1

	// Each test case is constructed by:
	//   testName          – a human‐readable description
	//   mockSetup         – sets up sqlmock expectations and returned results/errors
	//   expectedError     – the error we expect Delete() to return
	tests := []struct {
		testName      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			testName: "Success: Should delete product correctly",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(inputID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			testName: "Error case: Deleting referenced product",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(inputID).
					WillReturnError(&mysql.MySQLError{Number: 1451})
			},
			expectedError: httperrors.ConflictError{
				Message: "The product to delete is still referenced by some product records",
			},
		},
		{
			testName: "Error case: Internal Server Error on exec error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(inputID).
					WillReturnError(errors.New("db error on exec"))
			},
			expectedError: httperrors.InternalServerError{},
		},
		{
			testName: "Error case: Product not found",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(inputID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: httperrors.NotFoundError{Message: "Product not found"},
		},

		{
			testName: "Error case: Internal Server Error on RowsAffected error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(query).
					WithArgs(inputID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			expectedError: httperrors.InternalServerError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewProductRepositoryDB(db)
			tc.mockSetup(mock)

			// Act
			err = repo.Delete(context.Background(), inputID)

			// Assert
			require.Equal(t, tc.expectedError, err)
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
