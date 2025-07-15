package repository

import (
	"database/sql"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
)

type LocalityRepositoryDB struct {
	db *sql.DB
}

func NewLocalityRepository(db *sql.DB) LocalityRepository {
	return &LocalityRepositoryDB{db: db}
}

// This function creates a new locality in the database.
// It checks if all required fields are provided and if the locality ID already exists.
func (r *LocalityRepositoryDB) Create(locality models.Locality) (models.Locality, error) {
	var exists int
	const queryCheckExists = `
        SELECT COUNT(*) 
        FROM localities 
        WHERE id = ?
    `
	err := r.db.QueryRow(queryCheckExists, locality.ID).Scan(&exists)
	if err != nil {
		return models.Locality{}, err
	}
	if exists > 0 {
		return models.Locality{}, httperrors.ConflictError{Message: "Locality id already exists"}
	}

	const queryInsertLocality = `
        INSERT INTO localities (id, locality_name, province_name, country_name)
        VALUES (?, ?, ?, ?)
    `
	_, err = r.db.Exec(queryInsertLocality, locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName)
	if err != nil {
		return models.Locality{}, err
	}

	return locality, nil
}

// This function retrieves a locality by its ID.
// It returns an error if the locality is not found or if there is a database error.
func (r *LocalityRepositoryDB) GetByID(id string) (models.Locality, error) {
	var locality models.Locality
	const queryGetByID = `
		SELECT id, locality_name, province_name, country_name 
		FROM localities 
		WHERE id = ?
	`
	err := r.db.QueryRow(queryGetByID, id).Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceName, &locality.CountryName)

	if err == sql.ErrNoRows {
		return models.Locality{}, httperrors.NotFoundError{Message: "Locality not found"}
	}
	return locality, err
}

// This function retrieves a report of sellers by locality.
// If an ID is provided, it returns the report for that specific locality.
// If no ID is provided, it returns the report for all localities.
func (r *LocalityRepositoryDB) GetSellerReport(localityID *string) ([]models.SellerReport, error) {
	var rows *sql.Rows
	var err error

	if localityID != nil && *localityID != "" {
		rows, err = r.db.Query(`
            SELECT l.id, l.locality_name, COUNT(s.id)
            FROM localities l
            LEFT JOIN sellers s ON l.id = s.locality_id
            WHERE l.id = ?
            GROUP BY l.id, l.locality_name
        `, *localityID)
	} else {
		rows, err = r.db.Query(`
            SELECT l.id, l.locality_name, COUNT(s.id)
            FROM localities l
            LEFT JOIN sellers s ON l.id = s.locality_id
            GROUP BY l.id, l.locality_name
        `)
	}
	if err != nil {
		return nil, httperrors.InternalServerError{Message: "Error obtaining Report by LocalityId"}
	}
	defer rows.Close()

	var reports []models.SellerReport
	for rows.Next() {
		var report models.SellerReport
		if err := rows.Scan(
			&report.LocalityID,
			&report.LocalityName,
			&report.SellersCount,
		); err != nil {
			return nil, httperrors.InternalServerError{Message: "Error reading Seller data"}
		}
		reports = append(reports, report)
	}
	if localityID != nil && len(reports) == 0 {
		return nil, httperrors.NotFoundError{Message: "Locality not found"}
	}
	return reports, nil
}
