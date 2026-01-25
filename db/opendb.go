package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(path string) (*sql.DB, error) {

	dbConn, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(1)

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	_, err = dbConn.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Println("Warning: failed to enable WAL mode:", err)
	}

	return dbConn, nil
}
