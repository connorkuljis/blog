package model

import (
	"database/sql"
	"time"
)

type Entry struct {
	ID               int64          `db:"id"`
	CategoryID       int64          `db:"category_id"`
	Title            string         `db:"title"`
	Content          sql.NullString `db:"content"`
	Description      sql.NullString `db:"description"`
	FeaturedImageURL sql.NullString `db:"featured_image_url"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	IsDraft          int            `db:"is_draft"`
}
