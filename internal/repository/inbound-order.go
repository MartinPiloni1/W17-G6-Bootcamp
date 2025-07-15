package repository

import (
	"database/sql"
	"errors"
	"strings"
	_ "time"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
)

type InboundOrderRepositoryDB struct {
	db *sql.DB
}

func NewInboundOrderRepository(db *sql.DB) InboundOrderRepository {
	return &InboundOrderRepositoryDB{db: db}
}

// Create inserts a new InboundOrder and returns the created record with its ID and OrderDate.
func (r *InboundOrderRepositoryDB) Create(order models.InboundOrder) (models.InboundOrder, error) {
	const query = `
		INSERT INTO inbound_orders (
			order_number,
			order_date,
			employee_id,
			warehouse_id,
			product_batch_id
		) VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(query,
		order.OrderNumber,
		order.OrderDate, // Debes asegurarte de, si llega vacía, poner time.Now() antes de llamar a Create
		order.EmployeeID,
		order.WarehouseID,
		order.ProductBatchID,
	)
	if err != nil {
		return models.InboundOrder{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.InboundOrder{}, err
	}
	order.ID = int(id)
	return order, nil
}

// GetByOrderNumber returns the InboundOrder with that order_number, or a zero-value if not found.
func (r *InboundOrderRepositoryDB) GetByOrderNumber(orderNumber string) (models.InboundOrder, error) {
	const query = `
		SELECT
			id,
			order_number,
			order_date,
			employee_id,
			warehouse_id,
			product_batch_id
		FROM inbound_orders
		WHERE order_number = ?
	`
	var order models.InboundOrder
	err := r.db.QueryRow(query, orderNumber).
		Scan(&order.ID, &order.OrderNumber, &order.OrderDate, &order.EmployeeID, &order.WarehouseID, &order.ProductBatchID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.InboundOrder{}, nil // Patron tuyo: not found = objeto vacío y err=nil
	}
	if err != nil {
		return models.InboundOrder{}, err
	}
	return order, nil
}

func (r *InboundOrderRepositoryDB) CountInboundOrdersForEmployees(employeeIDs []int) (map[int]int, error) {
	result := make(map[int]int)
	if len(employeeIDs) == 0 {
		return result, nil
	}
	// Utiliza IN (?,?,?)...
	query := "SELECT employee_id, COUNT(*) FROM inbound_orders WHERE employee_id IN (?" + strings.Repeat(",?", len(employeeIDs)-1) + ") GROUP BY employee_id"
	args := make([]interface{}, len(employeeIDs))
	for i, id := range employeeIDs {
		args[i] = id
	}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var eid, cnt int
		if err := rows.Scan(&eid, &cnt); err != nil {
			return nil, err
		}
		result[eid] = cnt
	}
	return result, nil
}

func (r *InboundOrderRepositoryDB) CountInboundOrdersForEmployee(employeeID int) (int, error) {
	const query = `SELECT COUNT(*) FROM inbound_orders WHERE employee_id = ?`
	var count int
	err := r.db.QueryRow(query, employeeID).Scan(&count)
	return count, err
}
