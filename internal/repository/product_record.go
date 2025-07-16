package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

// ProductRecordRepositoryDB is a SQL implementation of ProductRecordRepository.
// It uses the provided *sql.DB connection to perform CRUD operations
// against the product_records table in the database.
type ProductRecordRepositoryDB struct {
	db *sql.DB
}

// NewProductRecordRepositoryDB constructs a ProductRecordRepositoryDB that uses
// the given *sql.DB for all data operations.
func NewProductRecordRepositoryDB(db *sql.DB) ProductRecordRepository {
	return &ProductRecordRepositoryDB{
		db: db,
	}
}

// Create inserts a new product record into the database using the given
// attributes, then returns the complete ProductRecord (including its auto-generated ID).
// If any database operation fails, it returns an InternalServerError.
func (r *ProductRecordRepositoryDB) Create(ctx context.Context, attributes models.ProductRecordAttributes) (models.ProductRecord, error) {
	const query = `
		INSERT INTO product_records (
		 	last_update_date,
			purchase_price,
			sale_price,
			product_id
		) VALUES (
			?, ?, ?, ?
		)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		attributes.LastUpdateDate,
		attributes.PurchasePrice,
		attributes.SalePrice,
		attributes.ProductID,
	)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1452 {
			return models.ProductRecord{},
				httperrors.ConflictError{Message: "a product with the given id does not exist"}
		}
		return models.ProductRecord{}, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.ProductRecord{}, err
	}

	newProductRecord := models.ProductRecord{
		ID:                      int(lastId),
		ProductRecordAttributes: attributes,
	}
	return newProductRecord, nil
}
