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
// ProductAttributes, then returns the complete Product
// (including its auto-generated ID). If any database operation fails,
// it returns an InternalServerError.
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
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			switch me.Number {
			case 1062:
				return models.Product{},
					httperrors.ConflictError{Message: "A product with the given product code already exists"}
			case 1452:
				return models.Product{},
					httperrors.ConflictError{Message: "The given seller id does not exists"}
			}
		}
		return models.Product{}, httperrors.InternalServerError{Message: "Error creating product"}
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Error creating product"}
	}

	newProduct := models.Product{
		ID:                int(lastId),
		ProductAttributes: productAttributes,
	}
	return newProduct, nil
}

// GetAll retrieves all product records from the database.
// It scans each row into a models.Product and returns the slice.
// On any database error, it returns an InternalServerError.
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
		return nil,
			httperrors.InternalServerError{Message: "Database error"}
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
			return nil,
				httperrors.InternalServerError{Message: "Database error"}
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil,
			httperrors.InternalServerError{Message: "Database error"}
	}
	return products, nil
}

// GetByID retrieves a single Product by its integer ID.
// If no matching row exists, returns a NotFoundError.
// On other database errors, returns an InternalServerError.
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
		return models.Product{},
			httperrors.NotFoundError{Message: "Product not found"}
	} else if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Database error"}
	}

	return product, nil
}

// Update modifies the product record with the given ID in the database,
// setting each column to the corresponding field in the provided models.Product.
// After the UPDATE statement, it returns the updated Product.
// On any database failure, it returns an InternalServerError.
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
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			switch me.Number {
			case 1062:
				return models.Product{},
					httperrors.ConflictError{Message: "A product with the given product code already exists"}
			case 1452:
				return models.Product{},
					httperrors.NotFoundError{Message: "The given seller id does not exists"}
			}
		}
		return models.Product{}, httperrors.InternalServerError{Message: "Error creating product"}
	}

	return updatedProduct, nil
}

// Delete removes the product record with the specified ID from the database.
// If the deletion fails due to a DB error, it returns an InternalServerError.
// If the product does not exist it returns a NotFoundError.
func (r *ProductRepositoryDB) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM products
		WHERE id=?	
    `

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return httperrors.InternalServerError{Message: "Database error"}
	}

	count, err := res.RowsAffected()
	if err != nil {
		return httperrors.InternalServerError{Message: "Database error"}
	} else if count == 0 {
		return httperrors.NotFoundError{Message: "Product not found"}
	}
	return nil
}
