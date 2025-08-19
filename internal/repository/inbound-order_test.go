package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestInboundOrderRepositoryDB_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewInboundOrderRepository(db)

	order := models.InboundOrder{
		InboundOrderAttributes: models.InboundOrderAttributes{
			OrderNumber:    "A123",
			OrderDate:      time.Date(2024, 6, 13, 14, 0, 0, 0, time.UTC),
			EmployeeID:     1,
			WarehouseID:    2,
			ProductBatchID: 3,
		},
	}
	// Case: successful create
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO inbound_orders (
			order_number,
			order_date,
			employee_id,
			warehouse_id,
			product_batch_id
		) VALUES (?, ?, ?, ?, ?)
	`)).
		WithArgs(order.OrderNumber, order.OrderDate, order.EmployeeID, order.WarehouseID, order.ProductBatchID).
		WillReturnResult(sqlmock.NewResult(42, 1))

	created, err := repo.Create(order)
	require.NoError(t, err)
	require.Equal(t, 42, created.ID)
	require.Equal(t, order.OrderNumber, created.OrderNumber)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: Exec error
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO inbound_orders (
			order_number,
			order_date,
			employee_id,
			warehouse_id,
			product_batch_id
		) VALUES (?, ?, ?, ?, ?)
	`)).
		WithArgs(order.OrderNumber, order.OrderDate, order.EmployeeID, order.WarehouseID, order.ProductBatchID).
		WillReturnError(errors.New("db fail"))

	created, err = repo.Create(order)
	require.Error(t, err)
	require.Equal(t, models.InboundOrder{}, created)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: LastInsertId error (simulate)
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO inbound_orders (
			order_number,
			order_date,
			employee_id,
			warehouse_id,
			product_batch_id
		) VALUES (?, ?, ?, ?, ?)
	`)).
		WithArgs(order.OrderNumber, order.OrderDate, order.EmployeeID, order.WarehouseID, order.ProductBatchID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

	created, err = repo.Create(order)
	require.Error(t, err)
	require.Equal(t, models.InboundOrder{}, created)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInboundOrderRepositoryDB_GetByOrderNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.NewInboundOrderRepository(db)

	query := regexp.QuoteMeta(`
		SELECT
			id,
			order_number,
			order_date,
			employee_id,
			warehouse_id,
			product_batch_id
		FROM inbound_orders
		WHERE order_number = ?
	`)

	orderDate := time.Date(2024, 6, 13, 13, 0, 0, 0, time.UTC)
	order := models.InboundOrder{
		ID: 7,
		InboundOrderAttributes: models.InboundOrderAttributes{
			OrderNumber:    "A123",
			OrderDate:      orderDate,
			EmployeeID:     1,
			WarehouseID:    2,
			ProductBatchID: 3,
		},
	}

	// Case: found
	rows := sqlmock.NewRows([]string{
		"id", "order_number", "order_date", "employee_id", "warehouse_id", "product_batch_id",
	}).AddRow(order.ID, order.OrderNumber, order.OrderDate, order.EmployeeID, order.WarehouseID, order.ProductBatchID)
	mock.ExpectQuery(query).WithArgs("A123").WillReturnRows(rows)

	got, err := repo.GetByOrderNumber("A123")
	require.NoError(t, err)
	require.Equal(t, order, got)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: not found (returns zero-value)
	rows = sqlmock.NewRows([]string{
		"id", "order_number", "order_date", "employee_id", "warehouse_id", "product_batch_id",
	})
	mock.ExpectQuery(query).WithArgs("ZZZ").WillReturnRows(rows)

	got, err = repo.GetByOrderNumber("ZZZ")
	require.NoError(t, err)
	require.Equal(t, models.InboundOrder{}, got)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: error (other than no rows)
	mock.ExpectQuery(query).WithArgs("fail").WillReturnError(errors.New("db error"))
	got, err = repo.GetByOrderNumber("fail")
	require.Error(t, err)
	require.Equal(t, models.InboundOrder{}, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInboundOrderRepositoryDB_CountInboundOrdersForEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.NewInboundOrderRepository(db)

	query := regexp.QuoteMeta(`SELECT COUNT(*) FROM inbound_orders WHERE employee_id = ?`)

	// Case: ok
	mock.ExpectQuery(query).WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(8))
	count, err := repo.CountInboundOrdersForEmployee(5)
	require.NoError(t, err)
	require.Equal(t, 8, count)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: db error
	mock.ExpectQuery(query).WithArgs(99).WillReturnError(errors.New("db fail"))
	_, err = repo.CountInboundOrdersForEmployee(99)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: scan error
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow("invalid"))
	_, err = repo.CountInboundOrdersForEmployee(1)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInboundOrderRepositoryDB_CountInboundOrdersForEmployees(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := repository.NewInboundOrderRepository(db)

	query := regexp.QuoteMeta("SELECT employee_id, COUNT(*) FROM inbound_orders GROUP BY employee_id")

	// Case: more than 1 employee
	rows := sqlmock.NewRows([]string{"employee_id", "COUNT(*)"}).
		AddRow(1, 5).
		AddRow(2, 6)
	mock.ExpectQuery(query).WillReturnRows(rows)

	got, err := repo.CountInboundOrdersForEmployees()
	require.NoError(t, err)
	require.Equal(t, map[int]int{1: 5, 2: 6}, got)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: error on query
	mock.ExpectQuery(query).WillReturnError(errors.New("db error"))
	_, err = repo.CountInboundOrdersForEmployees()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

	// Case: scan error
	rows = sqlmock.NewRows([]string{"employee_id", "COUNT(*)"}).AddRow("bad", 9)
	mock.ExpectQuery(query).WillReturnRows(rows)
	_, err = repo.CountInboundOrdersForEmployees()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
