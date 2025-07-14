package repository

import (
	"database/sql"
	"errors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"log"
)

type WarehouseRepositoryDB struct {
	db *sql.DB
}

// NewWarehouseRepositoryDb creates a new instance of WarehouseRepositoryDB.
func NewWarehouseRepositoryDb(db *sql.DB) *WarehouseRepositoryDB {
	return &WarehouseRepositoryDB{db: db}
}

// Create adds a new warehouse and returns the created warehouse.
func (p *WarehouseRepositoryDB) Create(warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	query := `
		INSERT INTO warehouse (
			warehouse_code, 
			address, 
			telephone,
			minimun_capacity,
			minimun_temperature
		) VALUES (?, ?, ?, ?, ?)
	`
	result, err := p.db.Exec(query,
		warehouseAttributes.WarehouseCode,
		warehouseAttributes.Address,
		warehouseAttributes.Telephone,
		warehouseAttributes.MinimunCapacity,
		warehouseAttributes.MinimunTemperature,
	)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			if me.Number == 1062 {
				return models.Warehouse{}, httperrors.ConflictError{Message: "the WarehouseCode already exists"}
			}
			return models.Warehouse{}, httperrors.InternalServerError{Message: "error creating warehouse"}
		}
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return models.Warehouse{}, httperrors.InternalServerError{Message: "error obtaining last insert ID"}
	}
	newWarehouse := models.Warehouse{
		Id:                  int(lastInsertId),
		WarehouseAttributes: warehouseAttributes,
	}
	return newWarehouse, nil
}

// GetAll returns all warehouses.
func (p *WarehouseRepositoryDB) GetAll() ([]models.Warehouse, error) {
	query := `
		SELECT 
			id,
			warehouse_code, 
			address, 
			telephone,
			minimun_capacity,
			minimun_temperature
		FROM warehouse
	`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, httperrors.InternalServerError{Message: "error obtaining warehouses"}
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()
	var warehouses []models.Warehouse
	for rows.Next() {
		var warehouse models.Warehouse
		if err := rows.Scan(
			&warehouse.Id,
			&warehouse.WarehouseCode,
			&warehouse.Address,
			&warehouse.Telephone,
			&warehouse.MinimunCapacity,
			&warehouse.MinimunTemperature,
		); err != nil {
			return nil, httperrors.InternalServerError{Message: "error reading warehouse data"}
		}
		warehouses = append(warehouses, warehouse)
	}
	return warehouses, nil
}

// GetByID returns a warehouse by its ID.
func (p *WarehouseRepositoryDB) GetByID(id int) (models.Warehouse, error) {
	query := `
		SELECT 
			id,
			warehouse_code, 
			address, 
			telephone,
			minimun_capacity,
			minimun_temperature
		FROM warehouse
		WHERE id = ?
	`
	row := p.db.QueryRow(query, id)
	var warehouse models.Warehouse
	if err := row.Scan(
		&warehouse.Id,
		&warehouse.WarehouseCode,
		&warehouse.Address,
		&warehouse.Telephone,
		&warehouse.MinimunCapacity,
		&warehouse.MinimunTemperature,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Warehouse{}, httperrors.NotFoundError{Message: "warehouse not found"}
		}
		return models.Warehouse{}, httperrors.InternalServerError{Message: "error obtaining warehouse by ID"}
	}
	return warehouse, nil
}

// Update modifies an existing warehouse and returns the updated warehouse.
func (p *WarehouseRepositoryDB) Update(id int, warehouseAttributes models.WarehouseAttributes) (models.Warehouse, error) {
	query := `
		UPDATE warehouse SET
			warehouse_code = ?,
			address = ?,
			telephone = ?,
			minimun_capacity = ?,
			minimun_temperature = ?
		WHERE id = ?
	`
	result, err := p.db.Exec(query,
		warehouseAttributes.WarehouseCode,
		warehouseAttributes.Address,
		warehouseAttributes.Telephone,
		warehouseAttributes.MinimunCapacity,
		warehouseAttributes.MinimunTemperature,
		id,
	)
	if err != nil {
		return models.Warehouse{}, httperrors.InternalServerError{Message: "error updating warehouse"}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Warehouse{}, httperrors.InternalServerError{Message: "error checking update"}
	}
	if rowsAffected == 0 {
		return models.Warehouse{}, httperrors.NotFoundError{Message: "warehouse not found"}
	}
	updatedWarehouse := models.Warehouse{
		Id:                  id,
		WarehouseAttributes: warehouseAttributes,
	}
	return updatedWarehouse, nil
}

// Delete removes a warehouse by its ID.
func (p *WarehouseRepositoryDB) Delete(id int) error {
	query := `DELETE FROM warehouse WHERE id = ?`
	result, err := p.db.Exec(query, id)
	if err != nil {
		return httperrors.InternalServerError{Message: "error deleting warehouse"}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return httperrors.InternalServerError{Message: "error checking delete"}
	}
	if rowsAffected == 0 {
		return httperrors.NotFoundError{Message: "warehouse not found"}
	}
	return nil
}
