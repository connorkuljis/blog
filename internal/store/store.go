package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(conn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
