package repository

import (
	"database/sql"
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

/*
ProductRepositoryDB is a SQL implementation of ProductRepository.
It uses the provided *sql.DB connection to perform CRUD operations
against the products table in the database.
*/
type ProductRepositoryDB struct {
	db *sql.DB
}

/*
NewProductRepositoryDB constructs a ProductRepositoryDB that uses
the given *sql.DB for all data operations.
*/
func NewProductRepositoryDB(db *sql.DB) ProductRepository {
	return &ProductRepositoryDB{
		db: db,
	}
}

/*
Create inserts a new product record into the database using the given
ProductAttributes, then fetches and returns the complete Product
(including its auto-generated ID). If any database operation fails,
it returns an InternalServerError.
*/
func (repository *ProductRepositoryDB) Create(productAttribbutes models.ProductAttributes) (models.Product, error) {
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

	result, err := repository.db.Exec(
		query,
		productAttribbutes.Description,
		productAttribbutes.ExpirationRate,
		productAttribbutes.FreezingRate,
		productAttribbutes.Height,
		productAttribbutes.Length,
		productAttribbutes.Width,
		productAttribbutes.NetWeight,
		productAttribbutes.ProductCode,
		productAttribbutes.RecommendedFreezingTemperature,
		productAttribbutes.ProductTypeID,
		productAttribbutes.SellerID,
	)
	if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Database error"}
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Database error"}
	}

	newProduct, err := repository.GetByID(int(lastId))
	if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Database error"}
	}
	return newProduct, nil
}

/*
GetAll retrieves all product records from the database.
It scans each row into a models.Product and returns the slice.
On any database error, it returns an InternalServerError.
*/
func (repository *ProductRepositoryDB) GetAll() ([]models.Product, error) {
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

	rows, err := repository.db.Query(query)
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

/*
GetByID retrieves a single Product by its integer ID.
If no matching row exists, returns a NotFoundError.
On other database errors, returns an InternalServerError.
*/
func (repository *ProductRepositoryDB) GetByID(id int) (models.Product, error) {
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

	row := repository.db.QueryRow(query, id)
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

/*
Update modifies the product record with the given ID in the database,
setting each column to the corresponding field in the provided models.Product.
After the UPDATE statement, it reloads and returns the updated Product.
On any database failure, it returns an InternalServerError.
*/
func (repository *ProductRepositoryDB) Update(id int, product models.Product) (models.Product, error) {
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

	_, err := repository.db.Exec(
		query,
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
		id,
	)
	if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Database error"}
	}

	updatedProduct, err := repository.GetByID(product.ID)
	if err != nil {
		return models.Product{},
			httperrors.InternalServerError{Message: "Database error"}
	}
	return updatedProduct, nil
}

/*
Delete removes the product record with the specified ID from the database.
If the deletion fails due to a DB error, it returns an InternalServerError.
If the product does not exist it returns a NotFoundError.
*/
func (repository *ProductRepositoryDB) Delete(id int) error {
	const query = `
		DELETE FROM products
		WHERE id=?	
    `

	res, err := repository.db.Exec(query, id)
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
