package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuyerRepositoryDB_Create(t *testing.T) {
	t.Run("Duplicate CardNumberID should return a Conflict error CardNumberId already in use", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		duplicateErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}

		mock.ExpectExec("INSERT INTO buyers").
			WithArgs(44728397, "Franco", "Colapinto").
			WillReturnError(duplicateErr)

		repo := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		_, err = repo.Create(ctx, models.BuyerAttributes{
			CardNumberId: 44728397,
			FirstName:    "Franco",
			LastName:     "Colapinto",
		})

		// assert
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "CardNumberId already in use")
		assert.ErrorAs(t, err, &httperrors.ConflictError{})

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Creates successfully a buyer returning the instance with id 10", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		mock.ExpectExec("INSERT INTO buyers").
			WithArgs(44728397, "Franco", "Colapinto").
			WillReturnResult(sqlmock.NewResult(10, 1))

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.Background()

		// act
		buyer, err := repoDB.Create(ctx, models.BuyerAttributes{
			CardNumberId: 44728397,
			FirstName:    "Franco",
			LastName:     "Colapinto",
		})

		// arrange
		require.NoError(t, err)
		assert.Equal(t, 10, buyer.Id)
		assert.Equal(t, 44728397, buyer.CardNumberId)
		assert.Equal(t, "Franco", buyer.FirstName)
		assert.Equal(t, "Colapinto", buyer.LastName)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Sql exec returns a mysql error code that is not contemplated, Creates returns the error", func(t *testing.T) {
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		notHandledSqlErrCode := &mysql.MySQLError{Number: 1037, Message: "Out of memory"} // 1037 out of memory

		mock.ExpectExec("INSERT INTO buyers").
			WillReturnError(notHandledSqlErrCode)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)

		ctx := context.TODO()
		_, err = repoDB.Create(
			ctx,
			models.BuyerAttributes{}, // pass a dummy BuyerAttributes to comply
		)

		assert.Error(t, err)
		assert.ErrorIs(t, err, notHandledSqlErrCode)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("When LastInsertId fails, it returns error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		execResult := sqlmock.NewErrorResult(fmt.Errorf("cannot get last insert id"))

		mock.ExpectExec("INSERT INTO buyers").
			WillReturnResult(execResult) // execResult has .LastInsertId() with the error

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO() // dummy ctx

		// act
		_, err = repoDB.Create(ctx, models.BuyerAttributes{})

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get last insert id")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositoryDB_GetAll(t *testing.T) {
	t.Run("QueryRowCOntext call returns an error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		queryRowError := &mysql.MySQLError{Number: 1037, Message: "Out of memory"}

		mock.ExpectQuery("FROM buyers").WillReturnError(queryRowError)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO() // dummy ctx

		// act
		_, err = repoDB.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, queryRowError)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("in the Scan of rows, one of the rows has nil/incompatible values and brakes the scan", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		rowsWithInvalidFields := sqlmock.NewRows(
			[]string{"id", "first_name", "last_name", "card_number_id"}).
			AddRow("idNotAnInteger", "Juan", nil, 1234)

		mock.ExpectQuery("FROM buyers").
			WillReturnRows(rowsWithInvalidFields)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO() // dummy ctx

		// act
		_, err = repoDB.GetAll(ctx)

		// assert
		assert.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully returns 2 rows as a slice of buyers", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)

		validRowsFromDB := sqlmock.NewRows(
			[]string{"id", "first_name", "last_name", "card_number_id"}).
			AddRow(1, "Juan", "Perez", 34567890).
			AddRow(2, "Juana", "Gonzalez", 23456789)

		mock.ExpectQuery("FROM buyers").WillReturnRows(validRowsFromDB)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO() // dummy context

		// act
		got, err := repoDB.GetAll(ctx)

		// assert
		expectedSecondBuyer := models.Buyer{
			Id: 2,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberId: 23456789,
				LastName:     "Gonzalez",
				FirstName:    "Juana",
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, 2, len(got))
		assert.Equal(t, expectedSecondBuyer.CardNumberId, got[1].CardNumberId)
		assert.Equal(t, expectedSecondBuyer.Id, got[1].Id)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error after the loop rows.Next() return error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlErrorRowsAfterLoopRows := errors.New("sql error after rows.Next()")

		rowWithErrorReturn := mock.NewRows(
			[]string{"id", "first_name", "last_name", "card_number_id"}).
			AddRow(1, "Juan", "Perez", 34567890).
			RowError(0, sqlErrorRowsAfterLoopRows) // the first row will return an error

		mock.ExpectQuery("FROM buyers").WillReturnRows(rowWithErrorReturn)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		_, err = repoDB.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, sqlErrorRowsAfterLoopRows)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositoryDB_GetByID(t *testing.T) {
	t.Run("the database throws an error and the create return that error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		queryRowError := &mysql.MySQLError{Number: 1037, Message: "Out of memory"}
		mock.ExpectQuery("FROM buyers").WillReturnError(queryRowError)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO() // dummy ctx

		// act
		_, err = repoDB.GetByID(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, queryRowError)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("from the database the row has invalidFIelds, fails the scan of the buyer ", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		rowWithInvalidFields := mock.NewRows(
			[]string{"id", "first_name", "last_name", "card_number_Id"}).
			AddRow(1, nil, nil, 12345678)

		mock.ExpectQuery("FROM buyers").WillReturnRows(rowWithInvalidFields)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetByID(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.Equal(t, models.Buyer{}, got)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows will be fetched from the db returning NotFoundError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		emptyRow := mock.NewRows([]string{"id", "first_name", "last_name", "card_number_Id"})

		mock.ExpectQuery("FROM buyers").WillReturnRows(emptyRow)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		_, err = repoDB.GetByID(ctx, 1)

		// assert
		expectedError := &httperrors.NotFoundError{Message: "Buyer not found"}
		assert.Error(t, err)
		assert.ErrorAs(t, err, expectedError)
		assert.Contains(t, err.Error(), expectedError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully fetch the the buyer with the id 1", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		validRowsFromDB := mock.NewRows(
			[]string{"id", "first_name", "last_name", "card_number_id"}).
			AddRow(1, "Juan", "Perez", 3456789)

		mock.ExpectQuery("FROM buyers").WillReturnRows(validRowsFromDB)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		expectedBuyer := models.Buyer{
			Id: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberId: 3456789,
				LastName:     "Perez",
				FirstName:    "Juan",
			},
		}

		// act
		got, err := repoDB.GetByID(ctx, expectedBuyer.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expectedBuyer.Id, got.Id)
		assert.Equal(t, expectedBuyer.CardNumberId, got.CardNumberId)
		assert.Equal(t, expectedBuyer.LastName, got.LastName)
		assert.Equal(t, expectedBuyer.FirstName, got.FirstName)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositoryDB_Delete(t *testing.T) {
	t.Run("the db returns an error and cannot complete the delete, return the db err", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlErrorTimeout := &mysql.MySQLError{Message: "Timeout", Number: 3024}
		mock.ExpectExec("DELETE FROM buyers").WithArgs(1).WillReturnError(sqlErrorTimeout)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		err = repoDB.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, sqlErrorTimeout)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("validating the RowsAffected() return error, return that error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		errorRowsAffected := errors.New("RowAffected message err")
		returnResult := sqlmock.NewErrorResult(errorRowsAffected)

		mock.ExpectExec("DELETE FROM buyers").WithArgs(1).WillReturnResult(returnResult)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		err = repoDB.Delete(ctx, 1)

		// assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorRowsAffected)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows where affected in the deletion returning NotFoundError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		argID := 1
		resultRowsAffected := sqlmock.NewResult(0, 0) // makes methods lastInsertedID = 0 and RowsAffected = 0
		mock.ExpectExec("DELETE FROM buyers").
			WithArgs(argID).
			WillReturnResult(resultRowsAffected)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()
		// act
		err = repoDB.Delete(ctx, argID)

		// assert
		expectedError := &httperrors.NotFoundError{Message: "Buyer not found"}
		assert.Error(t, err)
		assert.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully deletes a buyer returning nil err", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		argID := 1
		resultRowsAffected := sqlmock.NewResult(0, 1) // rows affected = 1 meaning it was deleted
		mock.ExpectExec("DELETE FROM buyers").
			WithArgs(argID).
			WillReturnResult(resultRowsAffected)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		err = repoDB.Delete(ctx, argID)

		// assert
		assert.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositoryDB_Update(t *testing.T) {
	buyerToUpdate := models.Buyer{
		Id: 1,
		BuyerAttributes: models.BuyerAttributes{
			CardNumberId: 12345678,
			FirstName:    "BuyerToUpdate",
			LastName:     "Bro",
		},
	}
	t.Run("update successfully the buyer returning the buyer data updated", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlResultUpdated := sqlmock.NewResult(0, 1)
		mock.ExpectExec("UPDATE buyers").
			WithArgs(
				buyerToUpdate.CardNumberId,
				buyerToUpdate.FirstName,
				buyerToUpdate.LastName,
				buyerToUpdate.Id).
			WillReturnResult(sqlResultUpdated)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.Update(ctx, buyerToUpdate.Id, buyerToUpdate)

		// assert
		require.NoError(t, err)
		assert.Equal(t, buyerToUpdate, got)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ExecContext query returns a generic error returning that error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlErrorTimeout := &mysql.MySQLError{Message: "Timeout", Number: 3024}
		mock.ExpectExec("UPDATE buyers").
			WithArgs(
				buyerToUpdate.CardNumberId,
				buyerToUpdate.FirstName,
				buyerToUpdate.LastName,
				buyerToUpdate.Id).
			WillReturnError(sqlErrorTimeout)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		_, err = repoDB.Update(ctx, buyerToUpdate.Id, buyerToUpdate)

		// assert
		require.Error(t, err)
		assert.Equal(t, sqlErrorTimeout, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("conflict trying to update to a cardNumberId that already exist", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		conflictErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
		mock.ExpectExec("UPDATE buyers").
			WithArgs(
				buyerToUpdate.CardNumberId,
				buyerToUpdate.FirstName,
				buyerToUpdate.LastName,
				buyerToUpdate.Id).
			WillReturnError(conflictErr)

		repo := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		_, err = repo.Update(ctx, buyerToUpdate.Id, buyerToUpdate)

		// assert
		expectedError := &httperrors.ConflictError{Message: "CardNumberId already in use"}
		require.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositoryDB_GetWithPurchaseOrdersCount(t *testing.T) {
	t.Run("method QueryContext fails with a random mysql error from the database", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlErrorTimeout := &mysql.MySQLError{Message: "Timeout", Number: 3024} // random error
		mock.ExpectQuery("FROM buyers b").WillReturnError(sqlErrorTimeout)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetWithPurchaseOrdersCount(ctx, nil) // pass nil to avoid WHERE and fetch all

		// arrange
		assert.Error(t, err)
		assert.Equal(t, err, sqlErrorTimeout)
		assert.Nil(t, got)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan of one of the rows fails because it has invalid fields, return the error", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		idArg := 42
		invalidRows := mock.NewRows(
			[]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow("invalidId", "thisShouldBeAnInt", "Juan", "Perez", 90)

		mock.ExpectQuery("FROM buyers b").
			WithArgs(idArg).
			WillReturnRows(invalidRows)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetWithPurchaseOrdersCount(ctx, &idArg)

		// assert
		assert.Error(t, err)
		assert.Nil(t, got)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err() returns an error after the loop of scan returning the error as it is", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		rowErr := errors.New("row error")
		idArg := 42 // id that will use the mock and the act section
		rowsThatHaveRowErr := mock.NewRows(
			[]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(idArg, 12345678, "Juan", "Perez", 90).
			RowError(0, rowErr)

		mock.ExpectQuery("FROM buyers b").WithArgs(idArg).WillReturnRows(rowsThatHaveRowErr)
		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetWithPurchaseOrdersCount(ctx, &idArg)

		// assert
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, rowErr.Error(), err.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("the id is provided but does not exist in the db returning NotFoundError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		idArg := 42 // id that will use the mock and the act section
		noRows := mock.NewRows(
			[]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"})

		mock.ExpectQuery("FROM buyers b").
			WithArgs(idArg).
			WillReturnRows(noRows)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetWithPurchaseOrdersCount(ctx, &idArg)

		// assert
		expectedError := &httperrors.NotFoundError{Message: "Buyer not found"}
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("id is not provided and successfully returns 2 BuyerWithPurchaseOrdersCount", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		buyerWithPurchaseOrdersCountOne := models.BuyerWithPurchaseOrdersCount{
			Buyer: models.Buyer{
				Id: 1,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberId: 12345678,
					FirstName:    "Juan",
					LastName:     "Perez",
				},
			},
			PurchaseOrdersCount: 90,
		}
		buyerWithPurchaseOrdersCountTwo := models.BuyerWithPurchaseOrdersCount{
			Buyer: models.Buyer{
				Id: 2,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberId: 12345678,
					FirstName:    "Jhon",
					LastName:     "Salch",
				},
			},
			PurchaseOrdersCount: 0,
		}

		twoRows := mock.NewRows(
			[]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(
				buyerWithPurchaseOrdersCountOne.Id,
				buyerWithPurchaseOrdersCountOne.CardNumberId,
				buyerWithPurchaseOrdersCountOne.FirstName,
				buyerWithPurchaseOrdersCountOne.LastName,
				buyerWithPurchaseOrdersCountOne.PurchaseOrdersCount).
			AddRow(
				buyerWithPurchaseOrdersCountTwo.Id,
				buyerWithPurchaseOrdersCountTwo.CardNumberId,
				buyerWithPurchaseOrdersCountTwo.FirstName,
				buyerWithPurchaseOrdersCountTwo.LastName,
				buyerWithPurchaseOrdersCountTwo.PurchaseOrdersCount)

		mock.ExpectQuery("FROM buyers b").WillReturnRows(twoRows)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetWithPurchaseOrdersCount(ctx, nil)

		// assert
		expectedResult := []models.BuyerWithPurchaseOrdersCount{
			buyerWithPurchaseOrdersCountOne,
			buyerWithPurchaseOrdersCountTwo,
		}

		assert.NoError(t, err)
		assert.Equal(t, len(expectedResult), len(got))
		assert.Equal(t, buyerWithPurchaseOrdersCountTwo.PurchaseOrdersCount, got[1].PurchaseOrdersCount)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("id is not provided and successfully returns 2 BuyerWithPurchaseOrdersCount", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		buyerWithPurchaseOrdersCount := models.BuyerWithPurchaseOrdersCount{
			Buyer: models.Buyer{
				Id: 1,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberId: 12345678,
					FirstName:    "Juan",
					LastName:     "Perez",
				},
			},
			PurchaseOrdersCount: 90,
		}

		twoRows := mock.NewRows(
			[]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(
				buyerWithPurchaseOrdersCount.Id,
				buyerWithPurchaseOrdersCount.CardNumberId,
				buyerWithPurchaseOrdersCount.FirstName,
				buyerWithPurchaseOrdersCount.LastName,
				buyerWithPurchaseOrdersCount.PurchaseOrdersCount)

		mock.ExpectQuery("FROM buyers b").WillReturnRows(twoRows)

		repoDB := repository.NewBuyerRepositoryDB(dbMocked)
		ctx := context.TODO()

		// act
		got, err := repoDB.GetWithPurchaseOrdersCount(ctx, nil)

		// assert

		assert.NoError(t, err)
		assert.Equal(t, 1, len(got))
		assert.Equal(t, buyerWithPurchaseOrdersCount.PurchaseOrdersCount, got[0].PurchaseOrdersCount)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
