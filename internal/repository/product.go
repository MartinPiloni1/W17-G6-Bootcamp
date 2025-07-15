package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

// ProductRepositoryDB is a SQL implementation of ProductRepository.
// It uses the provided *sql.DB connection to perform CRUD operations
// against the products table in the database.
type ProductRepositoryDB struct {
	db *sql.DB
}

// NewProductRepositoryDB constructs a ProductRepositoryDB that uses
// the given *sql.DB for all data operations.
func NewProductRepositoryDB(db *sql.DB) ProductRepository {
	return &ProductRepositoryDB{
		db: db,
	}
}

// Create inserts a new product into the database using the given
// ProductAttributes, then returns the complete Product (including its auto-generated ID).
//
// Behavior & Error Handling:
//   - On success, returns a models.Product populated with the generated ID
//     and the same attributes you passed in.
//   - If a product with the same product_code already exists (MySQL error #1062),
//     returns httperrors.ConflictError with a message about duplicate product code.
//   - If the provided seller_id does not exist (MySQL error #1452),
//     returns httperrors.ConflictError about the missing seller.
//   - Any other database error is returned.
func (r *ProductRepositoryDB) Create(ctx context.Context, productAttributes models.ProductAttributes) (models.Product, error) {
	const query = `
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
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		productAttributes.Description,
		productAttributes.ExpirationRate,
		productAttributes.FreezingRate,
		productAttributes.Height,
		productAttributes.Length,
		productAttributes.Width,
		productAttributes.NetWeight,
		productAttributes.ProductCode,
		productAttributes.RecommendedFreezingTemperature,
		productAttributes.ProductTypeID,
		productAttributes.SellerID,
	)
	if err != nil {
		var sqlError *mysql.MySQLError
		if errors.As(err, &sqlError) {
			switch sqlError.Number {
			case 1062:
				return models.Product{},
					httperrors.ConflictError{Message: "A product with the given product code already exists"}
			case 1452:
				return models.Product{},
					httperrors.ConflictError{Message: "The given seller id does not exists"}
			}
		}
		return models.Product{}, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Product{}, err
	}

	newProduct := models.Product{
		ID:                int(lastId),
		ProductAttributes: productAttributes,
	}
	return newProduct, nil
}

// GetAll fetches every product record from the database and returns them
// as a slice of models.Product. If any database error occurs it returns
// the error.
func (r *ProductRepositoryDB) GetAll(ctx context.Context) ([]models.Product, error) {
	const query = `
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
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(
			&product.ID,
			&product.Description,
			&product.ExpirationRate,
			&product.FreezingRate,
			&product.Height,
			&product.Length,
			&product.Width,
			&product.NetWeight,
			&product.ProductCode,
			&product.RecommendedFreezingTemperature,
			&product.ProductTypeID,
			&product.SellerID,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// GetByID retrieves a single product from the database by its unique ID.
// It returns the populated models.Product or if no matching row exists, returns a NotFoundError.
// On other database errors, returns the error.
func (r *ProductRepositoryDB) GetByID(ctx context.Context, id int) (models.Product, error) {
	const query = `
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
	`

	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return models.Product{}, err
	}

	var product models.Product
	err := row.Scan(
		&product.ID,
		&product.Description,
		&product.ExpirationRate,
		&product.FreezingRate,
		&product.Height,
		&product.Length,
		&product.Width,
		&product.NetWeight,
		&product.ProductCode,
		&product.RecommendedFreezingTemperature,
		&product.ProductTypeID,
		&product.SellerID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Product{}, err
	} else if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// GetRecordsPerProduct returns the count of product_records for each product.
// If id != nil, it filters by that product ID; otherwise it returns all products.
// It uses a LEFT JOIN so products with zero records appear with count = 0.
func (r *ProductRepositoryDB) GetRecordsPerProduct(ctx context.Context, id *int) ([]models.ProductRecordCount, error) {
	query := `
			SELECT 
				p.id, 
				p.description,
				COUNT(pr.id) AS records_count
			FROM products p
			LEFT JOIN product_records pr ON pr.product_id = p.id
		`

	var args []interface{}
	if id != nil {
		query += " WHERE p.id = ?"
		args = append(args, *id)
	}
	query += " GROUP BY p.id"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ProductsRecordsCount []models.ProductRecordCount
	for rows.Next() {
		var productRecordCount models.ProductRecordCount
		err := rows.Scan(
			&productRecordCount.ProductID,
			&productRecordCount.Description,
			&productRecordCount.RecordsCount,
		)
		if err != nil {
			return nil, err
		}
		ProductsRecordsCount = append(ProductsRecordsCount, productRecordCount)
	}
	return ProductsRecordsCount, nil
}

// Update modifies an existing product record in the database. It applies all
// fields from updatedProduct to the row identified by id.
//
// Behavior & Error Handling:
//   - On success, returns the updated product.
//   - If a product with the same product_code already exists (MySQL error #1062),
//     returns httperrors.ConflictError with a message about duplicate product code.
//   - If the provided seller_id does not exist (MySQL error #1452),
//     returns httperrors.ConflictError about the missing seller.
//   - Any other database error is returned.
func (r *ProductRepositoryDB) Update(ctx context.Context, id int, updatedProduct models.Product) (models.Product, error) {
	const query = `
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
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
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
		id,
	)
	if err != nil {
		var sqlError *mysql.MySQLError
		if errors.As(err, &sqlError) {
			switch sqlError.Number {
			case 1062:
				return models.Product{},
					httperrors.ConflictError{Message: "A product with the given product code already exists"}
			case 1452:
				return models.Product{},
					httperrors.ConflictError{Message: "The given seller id does not exists"}
			}
		}
		return models.Product{}, err
	}

	return updatedProduct, nil
}

// Delete removes the product record with the specified ID from the database.
// If the deletion fails due to a DB error, it returns the error.
// If the product does not exist it returns a NotFoundError.
func (r *ProductRepositoryDB) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM products
		WHERE id=?	
    `

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return httperrors.NotFoundError{Message: "Product not found"}
	}
	return nil
}
