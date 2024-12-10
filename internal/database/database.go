package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = "sqlite/create_tables.sql"
	conn   = "sqlite/content.db"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conn)
	if err != nil {
		return nil, err
	}

	text, err := os.ReadFile(schema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(text))
	if err != nil {
		return nil, err
	}

	return db, nil
}
