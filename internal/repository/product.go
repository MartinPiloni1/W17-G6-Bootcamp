package repository

import (
	"database/sql"
	"fmt"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
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
	panic("implement")
}

func (repository *ProductRepositoryDB) GetAll() ([]models.Product, error) {
	query := `
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
	panic("implement")
}

func (repository *ProductRepositoryDB) Update(id int, product models.Product) (models.Product, error) {
	panic("implement")
}

func (repository *ProductRepositoryDB) Delete(id int) error {
	panic("implement")
}
