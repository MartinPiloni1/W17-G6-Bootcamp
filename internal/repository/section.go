package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

// SectionRepositoryDB implements SectionRepository
type SectionRepositoryDB struct {
	db *sql.DB
}

/*
NewSectionRepositoryDB constructs a SectionRepositoryDB that uses
the given *sql.DB for all data operations.
*/
func NewSectionRepositoryDB(db *sql.DB) SectionRepository {
	return &SectionRepositoryDB{
		db: db,
	}
}


// Create, creates a new section in the repository
func (repository *SectionRepositoryDB) Create(ctx context.Context, section models.Section) (models.Section, error) {
	const query = `
        INSERT INTO sections (
            section_number, current_temperature, minimum_temperature,
            current_capacity, minimum_capacity, maximum_capacity,
            warehouse_id, product_type_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := repository.db.ExecContext(ctx, query,
		section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature,
		section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity,
		section.WarehouseID, section.ProductTypeID,
	)

	// Handle database errors
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return models.Section{}, httperrors.ConflictError{Message: "Section number already exists."}
			case 1452:
				return models.Section{}, httperrors.ConflictError{Message: "Warehouse does not exist."}
			default:
				return models.Section{}, httperrors.InternalServerError{}
			}
		}
		return models.Section{}, httperrors.InternalServerError{}
	}

	// Get the last inserted ID
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Section{}, httperrors.InternalServerError{}
	}
	section.ID = int(lastId)
	return section, nil
}

// Update, updates a section in the repository
func (repository *SectionRepositoryDB) Update(ctx context.Context, id int, data models.Section) (models.Section, error) {
	const query = `
        UPDATE sections SET
            section_number = ?, current_temperature = ?, minimum_temperature = ?,
            current_capacity = ?, minimum_capacity = ?, maximum_capacity = ?,
            warehouse_id = ?, product_type_id = ?
        WHERE id = ?
    `
	_, err := repository.db.ExecContext(ctx, query,
		data.SectionNumber, data.CurrentTemperature, data.MinimumTemperature,
		data.CurrentCapacity, data.MinimumCapacity, data.MaximumCapacity,
		data.WarehouseID, data.ProductTypeID, id,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return models.Section{}, httperrors.ConflictError{Message: "Section number already exists."}
			case 1452:
				return models.Section{}, httperrors.ConflictError{Message: "Warehouse does not exist."}
			default:
				return models.Section{}, httperrors.InternalServerError{}
			}
		}
		return models.Section{}, httperrors.InternalServerError{}
	}

	return data, nil
}


// Delete, deletes a section from the repository
func (repository *SectionRepositoryDB) Delete(ctx context.Context, id int) error {
    const query = `
        DELETE FROM sections
        WHERE id = ?
    `

    // Execute the query
    result, err := repository.db.ExecContext(ctx, query, id)
    if err != nil {
        return httperrors.InternalServerError{}
    }

    // Get the number of rows affected
    count, err := result.RowsAffected()
    if err != nil {
        return httperrors.InternalServerError{}
    } else if count == 0 {
        return httperrors.NotFoundError{Message: "Section not found"}
    }
    return nil
}

// GetAll returns all sections in the repository
func (repository *SectionRepositoryDB) GetAll(ctx context.Context) ([]models.Section, error) {
    const query = `
        SELECT
            id,
            section_number,
            current_temperature,
            minimum_temperature,
            current_capacity,
            minimum_capacity,
            maximum_capacity,
            warehouse_id,
            product_type_id
        FROM sections
    `

    // Execute the query
    rows, err := repository.db.QueryContext(ctx, query)
    if err != nil {
        return nil, httperrors.InternalServerError{}
    }
    defer rows.Close()

    var sections []models.Section
    for rows.Next() {
        var section models.Section
        err = rows.Scan(
            &section.ID,
            &section.SectionNumber,
            &section.CurrentTemperature,
            &section.MinimumTemperature,
            &section.CurrentCapacity,
            &section.MinimumCapacity,
            &section.MaximumCapacity,
            &section.WarehouseID,
            &section.ProductTypeID,
        )
        if err != nil {
            return nil, httperrors.InternalServerError{}
        }

        sections = append(sections, section)
    }

    if err := rows.Err(); err != nil {
        return nil, httperrors.InternalServerError{}
    }
    return sections, nil
}

// GetByID returns a section by its ID
func (repository *SectionRepositoryDB) GetByID(ctx context.Context, id int) (models.Section, error) {
    const query = `
        SELECT
            id,
            section_number,
            current_temperature,
            minimum_temperature,
            current_capacity,
            minimum_capacity,
            maximum_capacity,
            warehouse_id,
            product_type_id
        FROM sections
        WHERE id = ?
    `

    row := repository.db.QueryRowContext(ctx, query, id)
    if err := row.Err(); err != nil {
        return models.Section{}, httperrors.InternalServerError{}
    }

    var section models.Section
    err := row.Scan(
        &section.ID,
        &section.SectionNumber,
        &section.CurrentTemperature,
        &section.MinimumTemperature,
        &section.CurrentCapacity,
        &section.MinimumCapacity,
        &section.MaximumCapacity,
        &section.WarehouseID,
        &section.ProductTypeID,
    )
    if err != nil {
        return models.Section{}, httperrors.NotFoundError{Message: "Section not found"}
    }

    return section, nil
}