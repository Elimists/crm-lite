package db

import "database/sql"

func SetupTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL,
			fname TEXT,
			lname TEXT,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL DEFAULT (lower(hex(randomblob(8)))),
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			phone TEXT,
			message TEXT NOT NULL,
			source_domain TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'new',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, stmt := range tables {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
