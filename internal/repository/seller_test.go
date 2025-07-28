package repository_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSellerRepositoryDB_GetByID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
		AddRow(1, 1001, "Alkemy", "Calle 1", "12345", "1001")

	mock.ExpectQuery("SELECT id,").
		WithArgs(1).
		WillReturnRows(rows)

	seller, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, seller.ID)
	assert.Equal(t, "Alkemy", seller.CompanyName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSellerRepositoryDB_GetByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	mock.ExpectQuery("SELECT id,").
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(2)
	assert.Error(t, err)
	assert.IsType(t, httperrors.NotFoundError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestSellerRepositoryDB_Create_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	attr := models.SellerAttributes{
		CID:         2002,
		CompanyName: "TestCo",
		Address:     "ABC 22",
		Telephone:   "55555",
		LocalityID:  "9999",
	}

	mock.ExpectExec("INSERT INTO sellers").
		WithArgs(attr.CID, attr.CompanyName, attr.Address, attr.Telephone, attr.LocalityID).
		WillReturnResult(sqlmock.NewResult(5, 1)) // LastInsertID=5

	seller, err := repo.(*repository.SellerRepositoryDB).Create(attr)
	assert.NoError(t, err)
	assert.Equal(t, 5, seller.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSellerRepositoryDB_Create_DuplicateCID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	attr := models.SellerAttributes{CID: 1, CompanyName: "XYZ", Address: "YYY", Telephone: "333", LocalityID: "1001"}

	// Simulamos un error de clave duplicada de MySQL (1062)
	mysqlErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
	mock.ExpectExec("INSERT INTO sellers").
		WithArgs(attr.CID, attr.CompanyName, attr.Address, attr.Telephone, attr.LocalityID).
		WillReturnError(mysqlErr)

	_, err := repo.(*repository.SellerRepositoryDB).Create(attr)
	assert.Error(t, err)
	assert.IsType(t, httperrors.ConflictError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestSellerRepositoryDB_Delete_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	mock.ExpectExec("DELETE FROM sellers WHERE id = ?").
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := repo.(*repository.SellerRepositoryDB).Delete(3)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSellerRepositoryDB_Delete_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	mock.ExpectExec("DELETE FROM sellers WHERE id = ?").
		WithArgs(444).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err := repo.(*repository.SellerRepositoryDB).Delete(444)
	assert.Error(t, err)
	assert.IsType(t, httperrors.NotFoundError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestSellerRepositoryDB_GetAll_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
		AddRow(1, 1001, "A", "B", "C", "1001").
		AddRow(2, 2002, "B", "C", "D", "2002")

	mock.ExpectQuery("SELECT id, cid,").
		WillReturnRows(rows)

	result, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSellerRepositoryDB_Update_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	// El seller actual (de la base)
	actual := models.Seller{
		ID: 10,
		SellerAttributes: models.SellerAttributes{
			CID:         123,
			CompanyName: "Test",
			Address:     "Calle",
			Telephone:   "5555",
			LocalityID:  "1001",
		},
	}
	// Datos que se actualizan
	update := &models.SellerAttributes{
		CID:         200,
		CompanyName: "Nuevo",
		Address:     "Nueva Direccion",
		Telephone:   "22",
		LocalityID:  "2002",
	}

	// SELECT para obtener "actual"
	mock.ExpectQuery("SELECT id,").
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(actual.ID, actual.CID, actual.CompanyName, actual.Address, actual.Telephone, actual.LocalityID),
		)

	// UPDATE sellers SET ...
	mock.ExpectExec("UPDATE sellers SET").
		WithArgs(update.CID, update.CompanyName, update.Address, update.Telephone, update.LocalityID, 10).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// SELECT para devolver el actualizado
	mock.ExpectQuery("SELECT id,").
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(10, update.CID, update.CompanyName, update.Address, update.Telephone, update.LocalityID),
		)

	result, err := repo.(*repository.SellerRepositoryDB).Update(10, update)
	assert.NoError(t, err)
	assert.Equal(t, update.CID, result.CID)
	assert.Equal(t, update.CompanyName, result.CompanyName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSellerRepositoryDB_Update_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewSellerRepository(db)

	update := &models.SellerAttributes{
		CID: 111, CompanyName: "NoImporta", Address: "X", Telephone: "Y", LocalityID: "Z",
	}

	// SELECT (GetByID) da ErrNoRows
	mock.ExpectQuery("SELECT id,").
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.(*repository.SellerRepositoryDB).Update(99, update)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Seller not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}
