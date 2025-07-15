package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/models"
	"github.com/aaguero_meli/W17-G6-Bootcamp/pkg/httperrors"
	"github.com/go-sql-driver/mysql"
)

type SellerRepositoryDB struct {
	db *sql.DB
}

func NewSellerRepository(db *sql.DB) SellerRepository {
	return &SellerRepositoryDB{db: db}
}

// This function gets a seller by its ID from the database.
// If the seller is not found, it returns a NotFoundError.
func (r SellerRepositoryDB) GetByID(id int) (models.Seller, error) {
	var seller models.Seller
	const query = `
		SELECT id, 
			cid, 
			company_name, 
			address, 
			telephone, 
			locality_id 
		FROM sellers 
		WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(&seller.ID, &seller.CID, &seller.CompanyName, &seller.Address, &seller.Telephone, &seller.LocalityID)
	if err == sql.ErrNoRows {
		return models.Seller{}, httperrors.NotFoundError{Message: "Seller not found"}
	}
	return seller, err
}

// This function retrieves all sellers from the database.
// It returns a slice of Seller models or an error if the query fails.
func (r SellerRepositoryDB) GetAll() ([]models.Seller, error) {
	const query = `
		SELECT id, cid, company_name, address, telephone, locality_id 
		FROM sellers
		`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var sellers []models.Seller
	for rows.Next() {
		var s models.Seller
		if err := rows.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.LocalityID); err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}
	return sellers, nil
}

// This function creates a new seller in the database.
// It takes a SellerAttributes struct as input and returns the created Seller or an error.
// It checks for unique constraints on the CID and foreign key constraints on the locality_id.
func (r *SellerRepositoryDB) Create(att models.SellerAttributes) (models.Seller, error) {
	const queryInsertSeller = `
        INSERT INTO sellers (cid, company_name, address, telephone, locality_id) 
        VALUES (?, ?, ?, ?, ?)
    `
	res, err := r.db.Exec(queryInsertSeller, att.CID, att.CompanyName, att.Address, att.Telephone, att.LocalityID)
	if err != nil {
		var sqlErr *mysql.MySQLError
		if errors.As(err, &sqlErr) {
			switch sqlErr.Number {
			case 1062:
				return models.Seller{}, httperrors.ConflictError{Message: "Seller CID already exists"}
			case 1452:
				//No existe la FK (locality_id)
				return models.Seller{}, httperrors.ConflictError{Message: "Locality ID does not exist"}
			}
		}
		return models.Seller{}, httperrors.InternalServerError{Message: "Error creating seller"}
	}
	id, _ := res.LastInsertId()
	return models.Seller{
		ID:               int(id),
		SellerAttributes: att,
	}, nil
}

// This function deletes a seller by its ID from the database.
// If the seller does not exist, it returns a NotFoundError.
func (r *SellerRepositoryDB) Delete(id int) error {
	const queryDelete = `DELETE FROM sellers WHERE id = ?`
	res, err := r.db.Exec(queryDelete, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return httperrors.NotFoundError{Message: "Seller not found"}
	}
	return nil
}

// This function updates an existing seller's attributes in the database.
// It first retrieves the seller by ID, then updates the fields that are provided in the attributes
func (r *SellerRepositoryDB) Update(id int, att *models.SellerAttributes) (models.Seller, error) {
	actual, err := r.GetByID(id)
	if err != nil {
		return models.Seller{}, err
	}
	if att.CID != 0 {
		actual.CID = att.CID
	}
	if att.CompanyName != "" {
		actual.CompanyName = att.CompanyName
	}
	if att.Address != "" {
		actual.Address = att.Address
	}
	if att.Telephone != "" {
		actual.Telephone = att.Telephone
	}
	if att.LocalityID != "" {
		actual.LocalityID = att.LocalityID
	}

	const queryUpdate = `
        UPDATE sellers SET
            cid = ?,
            company_name = ?,
            address = ?,
            telephone = ?,
            locality_id = ?
        WHERE id = ?`
	_, err = r.db.Exec(queryUpdate, actual.CID, actual.CompanyName, actual.Address, actual.Telephone, actual.LocalityID, id)
	if err != nil {
		return models.Seller{}, err
	}
	return r.GetByID(id)
}
