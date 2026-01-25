package db

import (
	"crm-lite/models"
	"database/sql"
	"fmt"
)

func InsertContact(db *sql.DB, c models.Contact) (sql.Result, error) {
	result, err := db.Exec(`
		INSERT INTO contacts(
			name, 
			email, 
			phone, 
			message, 
			source_domain
			) 
		VALUES(?, ?, ?, ?, ?)`,
		c.Name,
		c.Email,
		c.Phone,
		c.Message,
		c.SourceDomain,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func QueryContactBy(db *sql.DB, field string, value any) (*models.Contact, error) {

	allowedfields := map[string]bool{
		"id":    true,
		"uuid":  true,
		"email": true,
	}
	if !allowedfields[field] {
		return nil, fmt.Errorf("invalid field name: %s", field)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			uuid,
			name,
			email,
			phone,
			message,
			source_domain,
			status,
			created_at,
			updated_at
		FROM contacts
		WHERE %s = ?`, field) // safe because we validated field

	row := db.QueryRow(query, value)

	var c models.Contact
	err := row.Scan(
		&c.ID,
		&c.UUID,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Message,
		&c.SourceDomain,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
