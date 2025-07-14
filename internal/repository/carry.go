package repository

import (
	"database/sql"
	"errors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
	"log"
)

type CarryRepositoryDB struct {
	db *sql.DB
}

// NewCarryRepositoryDb creates a new CarryRepositoryDB instance.
func NewCarryRepositoryDb(db *sql.DB) *CarryRepositoryDB {
	return &CarryRepositoryDB{db: db}
}

// Create inserts a new carry into the database.
func (p *CarryRepositoryDB) Create(carryAttributes models.CarryAttributes) (models.Carry, error) {

	query := `
		INSERT INTO carries (
			cid, 
			company_name, 
			address,
			telephone,
			locality_id
		) VALUES (?, ?, ?, ?, ?)
	`
	result, err := p.db.Exec(query,
		carryAttributes.Cid,
		carryAttributes.CompanyName,
		carryAttributes.Address,
		carryAttributes.Telephone,
		carryAttributes.LocalityId,
	)

	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			if me.Number == 1452 {
				return models.Carry{}, httperrors.NotFoundError{Message: "the LocalityId does not exist"}
			}
			if me.Number == 1062 {
				return models.Carry{}, httperrors.ConflictError{Message: "the Cid already exists"}
			}
			return models.Carry{}, httperrors.InternalServerError{Message: "error creating carry"}
		}
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return models.Carry{}, httperrors.InternalServerError{Message: "error obtaining last insert ID"}
	}

	newCarry := models.Carry{
		Id:              int(lastInsertId),
		CarryAttributes: carryAttributes,
	}

	return newCarry, nil
}

// GetReportByLocalityId retrieves a report of carries by locality ID.
func (p *CarryRepositoryDB) GetReportByLocalityId(localityId string) ([]models.CarryReport, error) {
	var (
		query string
		rows  *sql.Rows
		err   error
	)
	if localityId == "" {
		query = `
			SELECT
				c.locality_id,
				l.locality_name,
				COUNT(c.id) AS carries_count
			FROM carries c
			INNER JOIN localities l ON c.locality_id = l.id
			GROUP BY c.locality_id, l.locality_name
		`
		rows, err = p.db.Query(query)
	} else {
		query = `
			SELECT
				c.locality_id,
				l.locality_name,
				COUNT(c.id) AS carries_count
			FROM carries c
			INNER JOIN localities l ON c.locality_id = l.id
			WHERE c.locality_id = ?
			GROUP BY c.locality_id, l.locality_name
		`
		rows, err = p.db.Query(query, localityId)
	}
	if err != nil {
		return nil, httperrors.InternalServerError{Message: "error obtaining Report by LocalityId"}
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var reportCarries []models.CarryReport
	for rows.Next() {
		var reportCarry models.CarryReport
		if err := rows.Scan(
			&reportCarry.LocalityId,
			&reportCarry.LocalityName,
			&reportCarry.CarriesCount,
		); err != nil {
			return nil, httperrors.InternalServerError{Message: "error reading warehouse data"}
		}
		reportCarries = append(reportCarries, reportCarry)
	}
	return reportCarries, nil
}
