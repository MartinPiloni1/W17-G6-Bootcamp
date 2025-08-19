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

func TestCarryRepositoryDB_Create(t *testing.T) {
	t.Run("Duplicate Cid should return a Conflict error 'the Cid already exists'", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		duplicateErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("CID123", "Test Company", "Test Address", "123456789", "LOC001").
			WillReturnError(duplicateErr)

		repo := repository.NewCarryRepositoryDb(dbMocked)

		// act
		_, err = repo.Create(models.CarryAttributes{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  "LOC001",
		})

		// assert
		assert.Error(t, err)
		assert.Equal(t, "the Cid already exists", err.Error())
		assert.ErrorAs(t, err, &httperrors.ConflictError{})

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid LocalityId should return a Conflict error 'the LocalityId does not exist'", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		foreignKeyErr := &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"}

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("CID123", "Test Company", "Test Address", "123456789", "INVALID_LOC").
			WillReturnError(foreignKeyErr)

		repo := repository.NewCarryRepositoryDb(dbMocked)

		// act
		_, err = repo.Create(models.CarryAttributes{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  "INVALID_LOC",
		})

		// assert
		assert.Error(t, err)
		assert.Equal(t, "the LocalityId does not exist", err.Error())
		assert.ErrorAs(t, err, &httperrors.ConflictError{})

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Creates successfully a carry returning the instance with id 10", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("CID123", "Test Company", "Test Address", "123456789", "LOC001").
			WillReturnResult(sqlmock.NewResult(10, 1))

		repoDB := repository.NewCarryRepositoryDb(dbMocked)

		carryAttributes := models.CarryAttributes{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  "LOC001",
		}

		// act
		carry, err := repoDB.Create(carryAttributes)

		// assert
		require.NoError(t, err)
		assert.Equal(t, 10, carry.Id)
		assert.Equal(t, "CID123", carry.Cid)
		assert.Equal(t, "Test Company", carry.CompanyName)
		assert.Equal(t, "Test Address", carry.Address)
		assert.Equal(t, "123456789", carry.Telephone)
		assert.Equal(t, "LOC001", carry.LocalityId)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Sql exec returns a mysql error code that is not contemplated, Creates returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		notHandledSqlErrCode := &mysql.MySQLError{Number: 1037, Message: "Out of memory"} // 1037 out of memory

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("CID123", "Test Company", "Test Address", "123456789", "LOC001").
			WillReturnError(notHandledSqlErrCode)

		repoDB := repository.NewCarryRepositoryDb(dbMocked)

		carryAttributes := models.CarryAttributes{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  "LOC001",
		}

		// act
		_, err = repoDB.Create(carryAttributes)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error creating carry", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("When LastInsertId fails, it returns InternalServerError", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		execResult := sqlmock.NewErrorResult(fmt.Errorf("cannot get last insert id"))

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("CID123", "Test Company", "Test Address", "123456789", "LOC001").
			WillReturnResult(execResult) // execResult has .LastInsertId() with the error

		repoDB := repository.NewCarryRepositoryDb(dbMocked)

		carryAttributes := models.CarryAttributes{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  "LOC001",
		}

		// act
		_, err = repoDB.Create(carryAttributes)

		// assert
		assert.Error(t, err)
		assert.ErrorAs(t, err, &httperrors.InternalServerError{})
		assert.Equal(t, "error obtaining last insert ID", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Creates carry with empty strings successfully", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("", "", "", "", "").
			WillReturnResult(sqlmock.NewResult(15, 1))

		repoDB := repository.NewCarryRepositoryDb(dbMocked)

		carryAttributes := models.CarryAttributes{
			Cid:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityId:  "",
		}

		// act
		carry, err := repoDB.Create(carryAttributes)

		// assert
		require.NoError(t, err)
		assert.Equal(t, 15, carry.Id)
		assert.Equal(t, "", carry.Cid)
		assert.Equal(t, "", carry.CompanyName)
		assert.Equal(t, "", carry.Address)
		assert.Equal(t, "", carry.Telephone)
		assert.Equal(t, "", carry.LocalityId)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Creates carry with special characters successfully", func(t *testing.T) {
		// arrange
		dbMocked, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMocked.Close()

		mock.ExpectExec("INSERT INTO carries").
			WithArgs("CID-123", "Test & Company", "123 Main St., Apt #4", "+1-555-123-4567", "LOC-001").
			WillReturnResult(sqlmock.NewResult(20, 1))

		repoDB := repository.NewCarryRepositoryDb(dbMocked)

		carryAttributes := models.CarryAttributes{
			Cid:         "CID-123",
			CompanyName: "Test & Company",
			Address:     "123 Main St., Apt #4",
			Telephone:   "+1-555-123-4567",
			LocalityId:  "LOC-001",
		}

		// act
		carry, err := repoDB.Create(carryAttributes)

		// assert
		require.NoError(t, err)
		assert.Equal(t, 20, carry.Id)
		assert.Equal(t, "CID-123", carry.Cid)
		assert.Equal(t, "Test & Company", carry.CompanyName)
		assert.Equal(t, "123 Main St., Apt #4", carry.Address)
		assert.Equal(t, "+1-555-123-4567", carry.Telephone)
		assert.Equal(t, "LOC-001", carry.LocalityId)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
