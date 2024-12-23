package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	conn = "file:database/content.db?_foreign_keys=true"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
