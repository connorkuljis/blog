package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(conn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conn)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database '%s': %w", conn, err)
	}

	return db, nil
}

func Exec(db *sqlx.DB, sql string) (sql.Result, error) {
	res, err := db.Exec(sql)
	if err != nil {
		return res, err
	}

	return res, nil
}
