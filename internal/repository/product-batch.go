package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

// ProductBatchRepositoryDB implements ProductBatchRepository
type ProductBatchRepositoryDB struct {
	db *sql.DB
}

/*
NewProductBatchRepositoryDB constructs a ProductBatchRepositoryDB that uses
the given *sql.DB for all data operations.
*/
func NewProductBatchRepositoryDB(db *sql.DB) ProductBatchRepository {
	return &ProductBatchRepositoryDB{
		db: db,
	}
}


// Create creates a new product batch in the repository
func (repository *ProductBatchRepositoryDB) Create(ctx context.Context, productBatch models.ProductBatchAttibutes) (models.ProductBatch, error) {
	const query = `
        INSERT INTO product_batches (
            batch_number, current_quantity, current_temperature, due_date,
            initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature,
            product_id, section_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := repository.db.ExecContext(ctx, query,
		productBatch.BatchNumber, productBatch.CurrentQuantity, productBatch.CurrentTemperature, productBatch.DueDate,
		productBatch.InitialQuantity, productBatch.ManufacturingDate, productBatch.ManufacturingHour, productBatch.MinimumTemperature,
		productBatch.ProductID, productBatch.SectionID,
	)

	// Handle database errors
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return models.ProductBatch{}, httperrors.ConflictError{Message: "Batch number already exists."}
			case 1452:
				return models.ProductBatch{}, httperrors.ConflictError{Message: "Product or section does not exist."}
			default:
				return models.ProductBatch{}, httperrors.InternalServerError{}
			}
		}
		return models.ProductBatch{}, httperrors.InternalServerError{}
	}

	// Get the last inserted ID
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.ProductBatch{}, httperrors.InternalServerError{}
	}
	productCreated := models.ProductBatch{ID: int(lastId), ProductBatchAttibutes: productBatch}
	return productCreated, nil
}