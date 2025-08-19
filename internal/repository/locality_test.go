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

func TestLocalityRepositoryDB_Create_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	locality := models.Locality{
		ID:           "456",
		LocalityName: "Loc",
		ProvinceName: "BA",
		CountryName:  "Arg",
	}

	mock.ExpectExec("INSERT INTO localities").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnResult(sqlmock.NewResult(0, 1))

	res, err := repo.(*repository.LocalityRepositoryDB).Create(locality)
	assert.NoError(t, err)
	assert.Equal(t, locality, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_Create_Conflict(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	locality := models.Locality{ID: "X", LocalityName: "Z", ProvinceName: "P", CountryName: "C"}

	mysqlErr := &mysql.MySQLError{Number: 1062, Message: "duplicate"}
	mock.ExpectExec("INSERT INTO localities").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnError(mysqlErr)

	_, err := repo.(*repository.LocalityRepositoryDB).Create(locality)
	assert.Error(t, err)
	assert.IsType(t, httperrors.ConflictError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetByID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	locality := models.Locality{ID: "5", LocalityName: "A", ProvinceName: "B", CountryName: "C"}

	rows := sqlmock.NewRows([]string{"id", "locality_name", "province_name", "country_name"}).
		AddRow(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName)

	mock.ExpectQuery("SELECT id, locality_name").
		WithArgs("5").
		WillReturnRows(rows)

	res, err := repo.(*repository.LocalityRepositoryDB).GetByID("5")
	assert.NoError(t, err)
	assert.Equal(t, locality, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	mock.ExpectQuery("SELECT id, locality_name").
		WithArgs("7").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.(*repository.LocalityRepositoryDB).GetByID("7")
	assert.Error(t, err)
	assert.IsType(t, httperrors.NotFoundError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestLocalityRepositoryDB_GetSellerReport_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	localityID := "10"
	rows := sqlmock.NewRows([]string{"id", "locality_name", "sellers_count"}).
		AddRow(localityID, "LocName", 2).
		AddRow("11", "OtherLoc", 0)

	mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\)").
		WillReturnRows(rows)

	// Cuando localityID es nil --> reporte global
	res, err := repo.(*repository.LocalityRepositoryDB).GetSellerReport(nil)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "LocName", res[0].LocalityName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetSellerReport_ByID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	localityID := "10"
	rows := sqlmock.NewRows([]string{"id", "locality_name", "sellers_count"}).
		AddRow(localityID, "LocName", 2)

	mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\)").
		WithArgs(localityID).
		WillReturnRows(rows)

	res, err := repo.(*repository.LocalityRepositoryDB).GetSellerReport(&localityID)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "LocName", res[0].LocalityName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetSellerReport_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	localityID := "200"
	rows := sqlmock.NewRows([]string{"id", "locality_name", "sellers_count"})
	mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\)").
		WithArgs(localityID).
		WillReturnRows(rows)

	_, err := repo.(*repository.LocalityRepositoryDB).GetSellerReport(&localityID)
	assert.Error(t, err)
	assert.IsType(t, httperrors.NotFoundError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetReportByLocalityId_OneRow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	expectedReport := models.CarryReport{
		LocalityId:   "700",
		LocalityName: "Lujan",
		CarriesCount: 3,
	}

	rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
		AddRow(expectedReport.LocalityId, expectedReport.LocalityName, expectedReport.CarriesCount)

	mock.ExpectQuery("SELECT\\s+l.id AS locality_id").
		WithArgs("700").
		WillReturnRows(rows)

	got, err := repo.(*repository.LocalityRepositoryDB).GetReportByLocalityId("700")
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, expectedReport, got[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetReportByLocalityId_Global(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
		AddRow("101", "A", 2).
		AddRow("202", "B", 0)

	mock.ExpectQuery("SELECT\\s+l.id AS locality_id").
		WillReturnRows(rows)

	got, err := repo.(*repository.LocalityRepositoryDB).GetReportByLocalityId("")
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "A", got[0].LocalityName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetReportByLocalityId_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"})

	mock.ExpectQuery("SELECT\\s+l.id AS locality_id").
		WithArgs("000").
		WillReturnRows(rows)

	got, err := repo.(*repository.LocalityRepositoryDB).GetReportByLocalityId("000")
	assert.NoError(t, err)
	assert.Len(t, got, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetReportByLocalityId_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	mock.ExpectQuery("SELECT\\s+l.id AS locality_id").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.(*repository.LocalityRepositoryDB).GetReportByLocalityId("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error obtaining Report by LocalityId")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_Create_MySQLError_OtroNumero(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	locality := models.Locality{ID: "X", LocalityName: "Z", ProvinceName: "P", CountryName: "C"}

	mysqlErr := &mysql.MySQLError{Number: 999, Message: "otro"}
	mock.ExpectExec("INSERT INTO localities").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnError(mysqlErr)

	_, err := repo.(*repository.LocalityRepositoryDB).Create(locality)
	assert.Error(t, err)
	assert.Equal(t, mysqlErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_Create_NotFoundError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	locality := models.Locality{ID: "Z", LocalityName: "Y", ProvinceName: "X", CountryName: "W"}

	mysqlErr := &mysql.MySQLError{Number: 1452, Message: "foreign key fails"}
	mock.ExpectExec("INSERT INTO localities").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnError(mysqlErr)

	_, err := repo.(*repository.LocalityRepositoryDB).Create(locality)
	assert.Error(t, err)
	assert.IsType(t, httperrors.NotFoundError{}, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetByID_GenError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	mock.ExpectQuery("SELECT id, locality_name").
		WithArgs("8").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.(*repository.LocalityRepositoryDB).GetByID("8")
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetSellerReport_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	rows := sqlmock.NewRows([]string{"id", "locality_name", "sellers_count"}).
		AddRow(1, "Loc", "bad_int") // intentionally 'bad' type so Scan falle

	mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\)").
		WillReturnRows(rows)

	_, err := repo.(*repository.LocalityRepositoryDB).GetSellerReport(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error reading Seller data")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLocalityRepositoryDB_GetReportByLocalityId_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := repository.NewLocalityRepository(db)

	rows := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
		AddRow(1, "Loc", "bad_int")

	mock.ExpectQuery("SELECT\\s+l.id AS locality_id").
		WillReturnRows(rows)

	_, err := repo.(*repository.LocalityRepositoryDB).GetReportByLocalityId("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading Report by LocalityId data")
	assert.NoError(t, mock.ExpectationsWereMet())
}
