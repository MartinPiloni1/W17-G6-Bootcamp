package repository

import (
	"database/sql"
	"errors"
	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

type CarryRepositoryDB struct {
	TableName string
	db        *sql.DB
}

// NewCarryRepositoryDb creates a new CarryRepositoryDB instance.
func NewCarryRepositoryDb(db *sql.DB) *CarryRepositoryDB {
	return &CarryRepositoryDB{TableName: "carries", db: db}
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
