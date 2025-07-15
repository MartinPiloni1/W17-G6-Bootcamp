package repository

import (
	"database/sql"
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type EmployeeRepositoryDB struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &EmployeeRepositoryDB{db: db}
}

func (e *EmployeeRepositoryDB) Create(employee models.Employee) (models.Employee, error) {
	const query = `
		INSERT INTO employees (
			card_number_id,
			first_name,
			last_name,
			warehouse_id
		) VALUES (?, ?, ?, ?)
	`
	res, err := e.db.Exec(query, employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID)
	if err != nil {
		return models.Employee{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Employee{}, err
	}
	employee.Id = int(id)
	return employee, nil
}

func (e *EmployeeRepositoryDB) GetAll() ([]models.Employee, error) {
	const query = `
		SELECT
			id,
			card_number_id,
			first_name,
			last_name,
			warehouse_id
		FROM employees
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.Id, &emp.CardNumberID, &emp.FirstName, &emp.LastName, &emp.WarehouseID); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (e *EmployeeRepositoryDB) GetByID(id int) (models.Employee, error) {
	const query = `
		SELECT
			id,
			card_number_id,
			first_name,
			last_name,
			warehouse_id
		FROM employees
		WHERE id = ?
	`
	var emp models.Employee
	err := e.db.QueryRow(query, id).Scan(
		&emp.Id, &emp.CardNumberID, &emp.FirstName, &emp.LastName, &emp.WarehouseID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Employee{}, httperrors.NotFoundError{Message: "employee not found"}
	}
	if err != nil {
		return models.Employee{}, err
	}
	return emp, nil
}

func (e *EmployeeRepositoryDB) Update(id int, employee models.Employee) (models.Employee, error) {
	const query = `
		UPDATE employees
		SET
			card_number_id = ?,
			first_name = ?,
			last_name = ?,
			warehouse_id = ?
		WHERE id = ?
	`
	res, err := e.db.Exec(query, employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID, id)
	if err != nil {
		return models.Employee{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.Employee{}, err
	}
	if rows == 0 {
		return models.Employee{}, httperrors.NotFoundError{Message: "employee not found"}
	}
	employee.Id = id
	return employee, nil
}

func (e *EmployeeRepositoryDB) Delete(id int) error {
	const query = `
		DELETE FROM employees
		WHERE id = ?
	`
	res, err := e.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return httperrors.NotFoundError{Message: "employee not found"}
	}
	return nil
}
