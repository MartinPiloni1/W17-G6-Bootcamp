package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type ProductRepositoryDB struct {
	db *sql.DB
}

func NewProductRepositoryDB(db *sql.DB) ProductRepository {
	return &ProductRepositoryDB{
		db: db,
	}
}

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
		return models.Product{}, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Product{}, err
	}

	newProduct := models.Product{
		ID:                int(lastId),
		ProductAttributes: productAttribbutes,
	}
	return newProduct, nil
}

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
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID,
			&p.Description,
			&p.ExpirationRate,
			&p.FreezingRate,
			&p.Height,
			&p.Length,
			&p.Width,
			&p.NetWeight,
			&p.ProductCode,
			&p.RecommendedFreezingTemperature,
			&p.ProductTypeID,
			&p.SellerID,
		)
		if err != nil {
			return nil, fmt.Errorf("database scan error: %w", err)
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return products, nil
}

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
		WHEREE id = ?
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
		return models.Product{}, fmt.Errorf("database scan error: %w", err)
	}

	return product, nil
}

func (repository *ProductRepositoryDB) Update(id int, product models.Product) (models.Product, error) {
	panic("implement")
}

func (repository *ProductRepositoryDB) Delete(id int) error {
	const query = `
		DELETE FROM products
		WHERE id=?	
    `

	res, err := repository.db.Exec(query, id)
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
