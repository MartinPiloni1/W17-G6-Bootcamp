package repository_test

import (
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

func TestWarehouseRepositoryDB_Create(t *testing.T) {
	t.Run("Duplicate WarehouseCode should return a Conflict error 'the WarehouseCode already exists'", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		duplicateErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}

		mock.ExpectExec("INSERT INTO warehouses").
			WithArgs("WH001", "Test Address", "123456789", 1000, 5.5).
			WillReturnError(duplicateErr)

		repo := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repo.Create(models.WarehouseAttributes{
			WarehouseCode:      "WH001",
			Address:            "Test Address",
			Telephone:          "123456789",
			MinimunCapacity:    1000,
			MinimunTemperature: 5.5,
		})

		// assert
		assert.Error(t, err)
		assert.Equal(t, "the WarehouseCode already exists", err.Error())
		assert.ErrorAs(t, err, &httperrors.ConflictError{})

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Creates successfully a warehouse returning the instance with id 10", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		mock.ExpectExec("INSERT INTO warehouses").
			WithArgs("WH001", "Test Address", "123456789", 1000, 5.5).
			WillReturnResult(sqlmock.NewResult(10, 1))

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		warehouseAttributes := models.WarehouseAttributes{
			WarehouseCode:      "WH001",
			Address:            "Test Address",
			Telephone:          "123456789",
			MinimunCapacity:    1000,
			MinimunTemperature: 5.5,
		}

		// act
		warehouse, err := repoDB.Create(warehouseAttributes)

		// assert
		require.NoError(t, err)
		assert.Equal(t, 10, warehouse.Id)
		assert.Equal(t, "WH001", warehouse.WarehouseCode)
		assert.Equal(t, "Test Address", warehouse.Address)
		assert.Equal(t, "123456789", warehouse.Telephone)
		assert.Equal(t, 1000, warehouse.MinimunCapacity)
		assert.Equal(t, 5.5, warehouse.MinimunTemperature)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Sql exec returns a mysql error code that is not contemplated, Creates returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		notHandledSqlErrCode := &mysql.MySQLError{Number: 1037, Message: "Out of memory"} // 1037 out of memory

		mock.ExpectExec("INSERT INTO warehouses").
			WithArgs("WH001", "Test Address", "123456789", 1000, 5.5).
			WillReturnError(notHandledSqlErrCode)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		warehouseAttributes := models.WarehouseAttributes{
			WarehouseCode:      "WH001",
			Address:            "Test Address",
			Telephone:          "123456789",
			MinimunCapacity:    1000,
			MinimunTemperature: 5.5,
		}

		// act
		_, err = repoDB.Create(warehouseAttributes)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error creating warehouse", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("When LastInsertId fails, it returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		execResult := sqlmock.NewErrorResult(fmt.Errorf("cannot get last insert id"))

		mock.ExpectExec("INSERT INTO warehouses").
			WithArgs("WH001", "Test Address", "123456789", 1000, 5.5).
			WillReturnResult(execResult) // execResult has .LastInsertId() with the error

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		warehouseAttributes := models.WarehouseAttributes{
			WarehouseCode:      "WH001",
			Address:            "Test Address",
			Telephone:          "123456789",
			MinimunCapacity:    1000,
			MinimunTemperature: 5.5,
		}

		// act
		_, err = repoDB.Create(warehouseAttributes)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error obtaining last insert ID", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Creates warehouse with zero values successfully", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		mock.ExpectExec("INSERT INTO warehouses").
			WithArgs("", "", "", 0, 0.0).
			WillReturnResult(sqlmock.NewResult(15, 1))

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		warehouseAttributes := models.WarehouseAttributes{
			WarehouseCode:      "",
			Address:            "",
			Telephone:          "",
			MinimunCapacity:    0,
			MinimunTemperature: 0.0,
		}

		// act
		warehouse, err := repoDB.Create(warehouseAttributes)

		// assert
		require.NoError(t, err)
		assert.Equal(t, 15, warehouse.Id)
		assert.Equal(t, "", warehouse.WarehouseCode)
		assert.Equal(t, "", warehouse.Address)
		assert.Equal(t, "", warehouse.Telephone)
		assert.Equal(t, 0, warehouse.MinimunCapacity)
		assert.Equal(t, 0.0, warehouse.MinimunTemperature)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseRepositoryDB_GetAll(t *testing.T) {
	t.Run("QueryContext call returns an error, returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		queryError := &mysql.MySQLError{Number: 1037, Message: "Out of memory"}

		mock.ExpectQuery("FROM warehouses").WillReturnError(queryError)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.GetAll()

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error obtaining warehouses", err.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("in the Scan of rows, one of the rows has nil/incompatible values and breaks the scan", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		rowsWithInvalidFields := sqlmock.NewRows(
			[]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow("idNotAnInteger", "WH001", "Test Address", "123456789", 1000, 5.5)

		mock.ExpectQuery("FROM warehouses").
			WillReturnRows(rowsWithInvalidFields)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.GetAll()

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error reading warehouse data", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully returns 2 rows as a slice of warehouses", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		validRowsFromDB := sqlmock.NewRows(
			[]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow(1, "WH001", "Address 1", "123456789", 1000, 5.5).
			AddRow(2, "WH002", "Address 2", "987654321", 2000, -2.0)

		mock.ExpectQuery("FROM warehouses").WillReturnRows(validRowsFromDB)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		got, err := repoDB.GetAll()

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(got))
		assert.Equal(t, 1, got[0].Id)
		assert.Equal(t, "WH001", got[0].WarehouseCode)
		assert.Equal(t, "Address 1", got[0].Address)
		assert.Equal(t, "123456789", got[0].Telephone)
		assert.Equal(t, 1000, got[0].MinimunCapacity)
		assert.Equal(t, 5.5, got[0].MinimunTemperature)
		assert.Equal(t, 2, got[1].Id)
		assert.Equal(t, "WH002", got[1].WarehouseCode)
		assert.Equal(t, "Address 2", got[1].Address)
		assert.Equal(t, "987654321", got[1].Telephone)
		assert.Equal(t, 2000, got[1].MinimunCapacity)
		assert.Equal(t, -2.0, got[1].MinimunTemperature)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully returns empty slice when no warehouses exist", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		emptyRows := sqlmock.NewRows(
			[]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"})

		mock.ExpectQuery("FROM warehouses").WillReturnRows(emptyRows)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		got, err := repoDB.GetAll()

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(got))
		assert.Empty(t, got)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseRepositoryDB_GetByID(t *testing.T) {
	t.Run("the database throws an error and the GetByID returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		queryRowError := &mysql.MySQLError{Number: 1037, Message: "Out of memory"}
		mock.ExpectQuery("FROM warehouses").WillReturnError(queryRowError)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.GetByID(1)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error obtaining warehouse by ID", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("from the database the row has invalid fields, fails the scan of the warehouse", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		rowWithInvalidFields := mock.NewRows(
			[]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow(1, nil, nil, nil, nil, nil)

		mock.ExpectQuery("FROM warehouses").WillReturnRows(rowWithInvalidFields)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		got, err := repoDB.GetByID(1)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error obtaining warehouse by ID", err.Error())
		assert.Equal(t, models.Warehouse{}, got)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows will be fetched from the db returning NotFoundError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		emptyRow := mock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"})

		mock.ExpectQuery("FROM warehouses").WillReturnRows(emptyRow)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.GetByID(1)

		// assert
		expectedError := &httperrors.NotFoundError{Message: "warehouse not found"}
		assert.Error(t, err)
		assert.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully fetch the warehouse with the id 1", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		validRowsFromDB := mock.NewRows(
			[]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow(1, "WH001", "Test Address", "123456789", 1000, 5.5)

		mock.ExpectQuery("FROM warehouses").WillReturnRows(validRowsFromDB)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		expectedWarehouse := models.Warehouse{
			Id: 1,
			WarehouseAttributes: models.WarehouseAttributes{
				WarehouseCode:      "WH001",
				Address:            "Test Address",
				Telephone:          "123456789",
				MinimunCapacity:    1000,
				MinimunTemperature: 5.5,
			},
		}

		// act
		got, err := repoDB.GetByID(expectedWarehouse.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse.Id, got.Id)
		assert.Equal(t, expectedWarehouse.WarehouseCode, got.WarehouseCode)
		assert.Equal(t, expectedWarehouse.Address, got.Address)
		assert.Equal(t, expectedWarehouse.Telephone, got.Telephone)
		assert.Equal(t, expectedWarehouse.MinimunCapacity, got.MinimunCapacity)
		assert.Equal(t, expectedWarehouse.MinimunTemperature, got.MinimunTemperature)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseRepositoryDB_Update(t *testing.T) {
	warehouseToUpdate := models.Warehouse{
		Id: 1,
		WarehouseAttributes: models.WarehouseAttributes{
			WarehouseCode:      "WH001",
			Address:            "Updated Address",
			Telephone:          "987654321",
			MinimunCapacity:    1500,
			MinimunTemperature: -5.0,
		},
	}

	t.Run("update successfully the warehouse returning the warehouse data updated", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlResultUpdated := sqlmock.NewResult(0, 1)
		mock.ExpectExec("UPDATE warehouses").
			WithArgs(
				warehouseToUpdate.WarehouseCode,
				warehouseToUpdate.Address,
				warehouseToUpdate.Telephone,
				warehouseToUpdate.MinimunCapacity,
				warehouseToUpdate.MinimunTemperature,
				warehouseToUpdate.Id).
			WillReturnResult(sqlResultUpdated)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		got, err := repoDB.Update(warehouseToUpdate.Id, warehouseToUpdate.WarehouseAttributes)

		// assert
		require.NoError(t, err)
		assert.Equal(t, warehouseToUpdate, got)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ExecContext query returns a generic error returning InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlErrorTimeout := &mysql.MySQLError{Message: "Timeout", Number: 3024}
		mock.ExpectExec("UPDATE warehouses").
			WithArgs(
				warehouseToUpdate.WarehouseCode,
				warehouseToUpdate.Address,
				warehouseToUpdate.Telephone,
				warehouseToUpdate.MinimunCapacity,
				warehouseToUpdate.MinimunTemperature,
				warehouseToUpdate.Id).
			WillReturnError(sqlErrorTimeout)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.Update(warehouseToUpdate.Id, warehouseToUpdate.WarehouseAttributes)

		// assert
		require.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error updating warehouse", err.Error())

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RowsAffected returns an error, returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		errorRowsAffected := fmt.Errorf("RowAffected message err")
		returnResult := sqlmock.NewErrorResult(errorRowsAffected)

		mock.ExpectExec("UPDATE warehouses").
			WithArgs(
				warehouseToUpdate.WarehouseCode,
				warehouseToUpdate.Address,
				warehouseToUpdate.Telephone,
				warehouseToUpdate.MinimunCapacity,
				warehouseToUpdate.MinimunTemperature,
				warehouseToUpdate.Id).
			WillReturnResult(returnResult)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.Update(warehouseToUpdate.Id, warehouseToUpdate.WarehouseAttributes)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error checking update", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows where affected in the update returning NotFoundError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		argID := 1
		resultRowsAffected := sqlmock.NewResult(0, 0) // makes methods lastInsertedID = 0 and RowsAffected = 0
		mock.ExpectExec("UPDATE warehouses").
			WithArgs(
				warehouseToUpdate.WarehouseCode,
				warehouseToUpdate.Address,
				warehouseToUpdate.Telephone,
				warehouseToUpdate.MinimunCapacity,
				warehouseToUpdate.MinimunTemperature,
				argID).
			WillReturnResult(resultRowsAffected)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		_, err = repoDB.Update(argID, warehouseToUpdate.WarehouseAttributes)

		// assert
		expectedError := &httperrors.NotFoundError{Message: "warehouse not found or no changes made"}
		assert.Error(t, err)
		assert.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseRepositoryDB_Delete(t *testing.T) {
	t.Run("the db returns an error and cannot complete the delete, return InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		sqlErrorTimeout := &mysql.MySQLError{Message: "Timeout", Number: 3024}
		mock.ExpectExec("DELETE FROM warehouses").WithArgs(1).WillReturnError(sqlErrorTimeout)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		err = repoDB.Delete(1)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error deleting warehouse", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("foreign key constraint error returns InternalServerError with specific message", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		foreignKeyErr := &mysql.MySQLError{Number: 1451, Message: "Cannot delete or update a parent row"}
		mock.ExpectExec("DELETE FROM warehouses").WithArgs(1).WillReturnError(foreignKeyErr)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		err = repoDB.Delete(1)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "This warehouse cannot be deleted because it is associated with other entities", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("validating the RowsAffected() return error, return InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		errorRowsAffected := fmt.Errorf("RowAffected message err")
		returnResult := sqlmock.NewErrorResult(errorRowsAffected)

		mock.ExpectExec("DELETE FROM warehouses").WithArgs(1).WillReturnResult(returnResult)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		err = repoDB.Delete(1)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error checking delete", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no rows where affected in the deletion returning NotFoundError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		argID := 1
		resultRowsAffected := sqlmock.NewResult(0, 0) // rows affected = 0 meaning it was not found
		mock.ExpectExec("DELETE FROM warehouses").
			WithArgs(argID).
			WillReturnResult(resultRowsAffected)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		err = repoDB.Delete(argID)

		// assert
		expectedError := &httperrors.NotFoundError{Message: "warehouse not found"}
		assert.Error(t, err)
		assert.ErrorAs(t, err, expectedError)
		assert.Equal(t, err.Error(), expectedError.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully deletes a warehouse returning nil err", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		argID := 1
		resultRowsAffected := sqlmock.NewResult(0, 1) // rows affected = 1 meaning it was deleted
		mock.ExpectExec("DELETE FROM warehouses").
			WithArgs(argID).
			WillReturnResult(resultRowsAffected)

		repoDB := repository.NewWarehouseRepositoryDb(dbMocked)

		// act
		err = repoDB.Delete(argID)

		// assert
		assert.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
