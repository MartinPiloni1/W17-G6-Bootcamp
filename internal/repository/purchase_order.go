package repository

import (
	"context"
	"database/sql"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

type PurchaseOrderRepositoryDB struct {
	db *sql.DB
}

// Generates a repository with the db connection as parameter
func NewPurchaseOrderRepositoryDB(db *sql.DB) PurchaseOrderRepository {
	return &PurchaseOrderRepositoryDB{db: db}
}

func (r *PurchaseOrderRepositoryDB) Create(
	ctx context.Context,
	newPurchaseOrder models.PurchaseOrderAttributes) (models.PurchaseOrder, error) {
	const query = `
		INSERT INTO purchase_orders
			(order_number, order_date, tracking_code, buyer_id, product_record_id)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		newPurchaseOrder.OrderNumber,
		newPurchaseOrder.OrderDate,
		newPurchaseOrder.TrackingCode,
		newPurchaseOrder.BuyerId,
		newPurchaseOrder.ProductRecordId,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				// unique key OrderNumber duplicated
				err = httperrors.ConflictError{Message: "OrderNumber already in use"}
			case 1452:
				// FK violation constraint ProductRecordId or BuyerId
				err = httperrors.ConflictError{Message: "ProductRecordId and/or BuyerId does not exist"}
			}
		}
		return models.PurchaseOrder{}, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.PurchaseOrder{}, err
	}

	purchaseOrder := models.PurchaseOrder{
		Id:                      int(lastId),
		PurchaseOrderAttributes: newPurchaseOrder,
	}

	return purchaseOrder, nil
}
