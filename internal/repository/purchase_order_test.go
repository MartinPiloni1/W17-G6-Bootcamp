package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPurchaseOrderRepositoryDB_Create(t *testing.T) {
	newPurchaseOrderTime := time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC)

	t.Run("successfully creates a purchaseOrder", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		newPurchaseOrder := models.PurchaseOrderAttributes{
			OrderNumber:     "ON-001",
			OrderDate:       newPurchaseOrderTime,
			TrackingCode:    "TR-12345",
			BuyerId:         1,
			ProductRecordId: 100,
		}
		lastID := int64(22)

		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs(
				newPurchaseOrder.OrderNumber,
				newPurchaseOrder.OrderDate,
				newPurchaseOrder.TrackingCode,
				newPurchaseOrder.BuyerId,
				newPurchaseOrder.ProductRecordId).
			WillReturnResult(sqlmock.NewResult(lastID, 1))

		repo := repository.NewPurchaseOrderRepositoryDB(db)
		got, err := repo.Create(context.Background(), newPurchaseOrder)

		require.NoError(t, err)
		assert.Equal(t, int(lastID), got.Id)
		assert.Equal(t, newPurchaseOrder, got.PurchaseOrderAttributes)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("OrderNumber unique conflict (1062) error from db, return ConflictError", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		newPurchaseOrder := models.PurchaseOrderAttributes{
			OrderNumber:     "DUP-01",
			OrderDate:       newPurchaseOrderTime,
			TrackingCode:    "TR-02",
			BuyerId:         2,
			ProductRecordId: 33,
		}
		confErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}

		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs(
				newPurchaseOrder.OrderNumber,
				newPurchaseOrder.OrderDate,
				newPurchaseOrder.TrackingCode,
				newPurchaseOrder.BuyerId,
				newPurchaseOrder.ProductRecordId).
			WillReturnError(confErr)

		repo := repository.NewPurchaseOrderRepositoryDB(db)
		ctx := context.TODO()

		// act
		_, err = repo.Create(ctx, newPurchaseOrder)

		// assert
		expectedErr := &httperrors.ConflictError{Message: "OrderNumber already in use"}
		assert.Error(t, err)
		assert.ErrorAs(t, err, expectedErr)
		assert.Equal(t, err.Error(), expectedErr.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("dbError FK constraint returns ConflictError", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		newPurchaseOrder := models.PurchaseOrderAttributes{
			OrderNumber:     "DUP-01",
			OrderDate:       newPurchaseOrderTime,
			TrackingCode:    "TR-02",
			BuyerId:         2,
			ProductRecordId: 33,
		}
		confErr := &mysql.MySQLError{Number: 1452, Message: "foreign key constraint fails"}

		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs(
				newPurchaseOrder.OrderNumber,
				newPurchaseOrder.OrderDate,
				newPurchaseOrder.TrackingCode,
				newPurchaseOrder.BuyerId,
				newPurchaseOrder.ProductRecordId).
			WillReturnError(confErr)

		repo := repository.NewPurchaseOrderRepositoryDB(db)
		ctx := context.TODO()
		// act
		_, err = repo.Create(ctx, newPurchaseOrder)

		// assert
		expectedError := &httperrors.ConflictError{Message: "ProductRecordId and/or BuyerId does not exist"}
		require.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error from the db return that error as it is", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		newPurchaseOrder := models.PurchaseOrderAttributes{
			OrderNumber:     "DUP-01",
			OrderDate:       newPurchaseOrderTime,
			TrackingCode:    "TR-02",
			BuyerId:         2,
			ProductRecordId: 33,
		}
		queryRowError := &mysql.MySQLError{Number: 1037, Message: "Out of memory"}

		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs(
				newPurchaseOrder.OrderNumber,
				newPurchaseOrder.OrderDate,
				newPurchaseOrder.TrackingCode,
				newPurchaseOrder.BuyerId,
				newPurchaseOrder.ProductRecordId).
			WillReturnError(queryRowError)

		repo := repository.NewPurchaseOrderRepositoryDB(db)
		ctx := context.TODO()

		// act
		_, err = repo.Create(ctx, newPurchaseOrder)

		// assert
		assert.Error(t, err)
		assert.ErrorIs(t, queryRowError, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error at LastInsertId() return the error as it is", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		newPurchaseOrder := models.PurchaseOrderAttributes{
			OrderNumber:     "DUP-01",
			OrderDate:       newPurchaseOrderTime,
			TrackingCode:    "TR-02",
			BuyerId:         2,
			ProductRecordId: 33,
		}

		errLastInsertId := errors.New("cannot get last id")

		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs(
				newPurchaseOrder.OrderNumber,
				newPurchaseOrder.OrderDate,
				newPurchaseOrder.TrackingCode,
				newPurchaseOrder.BuyerId,
				newPurchaseOrder.ProductRecordId).
			WillReturnResult(sqlmock.NewErrorResult(errLastInsertId))

		repo := repository.NewPurchaseOrderRepositoryDB(db)
		ctx := context.TODO()

		// act
		_, err = repo.Create(ctx, newPurchaseOrder)

		// assert
		assert.Error(t, err)
		assert.Equal(t, err.Error(), errLastInsertId.Error())
		assert.ErrorIs(t, err, errLastInsertId)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
