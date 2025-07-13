package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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

type EntryRepo struct {
	db *sqlx.DB
}

func NewEntryRepo(db *sqlx.DB) *EntryRepo {
	return &EntryRepo{db: db}
}

func (r *EntryRepo) CreateEntry(entry *Entry) error {
	q := `
INSERT INTO 
entries 
	(category_id, title, author_id, created_at, updated_at) 
VALUES 
	($1, $2, $3, $4, $5)
`

	res, err := r.db.Exec(q,
		entry.CategoryID,
		entry.Title,
		entry.CreatedAt.Format(time.RFC3339),
		entry.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	entry.ID = id

	return nil
}

func (r *EntryRepo) ReadAllEntries(includeDrafts bool) ([]*Entry, error) {
	var entries []*Entry

	q := "SELECT * FROM entries"

	if !includeDrafts {
		q += " WHERE is_draft != 1"
	}

	q += " ORDER BY created_at DESC"

	err := r.db.Select(&entries, q)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) ReadRecentEntries(limit int) ([]Entry, error) {
	var entries []Entry
	err := r.db.Select(&entries, "SELECT * FROM entries ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) ReadAllByCategoryID(categoryID int64, enableDrafts bool) ([]*Entry, error) {
	var entries []*Entry

	q := "SELECT * FROM entries WHERE category_id = ?"

	if !enableDrafts {
		q += "AND is_draft != 1"
	}

	q += " ORDER BY created_at DESC"

	err := r.db.Select(&entries, q, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) ReadEntryByID(id int64) (*Entry, error) {
	var entry Entry
	err := r.db.Get(&entry, "SELECT * FROM entries WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting entry by id `%d`: %w", id, err)
	}

	return &entry, nil
}

func (r *EntryRepo) UpdateEntry(entry *Entry) error {
	q := `
UPDATE 
	entries 
SET 
	category_id = ?, 
	title = ?, 
	content = ?, 
	description = ?, 
	featured_image_url = ?, 
	updated_at = ?, 
	is_draft = ? 
WHERE 
	id = ?
`
	_, err := r.db.Exec(q,
		entry.CategoryID,
		entry.Title,
		entry.Content,
		entry.Description,
		entry.FeaturedImageURL,
		entry.UpdatedAt,
		entry.IsDraft,
		entry.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *EntryRepo) DeleteEntryByID(id int64) error {
	_, err := r.db.Exec("DELETE FROM entries WHERE id = ?", id)
	if err != nil {
		// TODO: better error handling
		return err
	}

	return nil
}

func (r *EntryRepo) ReadAllByTag(tag string, enableDrafts bool) ([]*Entry, error) {
	var entries []*Entry

	q := `
	SELECT e.*
	FROM entries e
	JOIN entry_tags et ON e.id = et.entry_id
	JOIN tags t ON et.tag_id = t.id
	WHERE t.name = ?
	`

	if !enableDrafts {
		q += "AND e.is_draft != 1"
	}

	q += " ORDER BY e.created_at DESC"

	err := r.db.Select(&entries, q, tag)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries by tag: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) AddTagToEntry(entryID int64, tagID int64) error {
	q := "INSERT INTO entry_tags (entry_id, tag_id) VALUES (?, ?)"
	_, err := r.db.Exec(q, entryID, tagID)
	return err
}

func (r *EntryRepo) RemoveTagFromEntry(entryID int64, tagID int64) error {
	q := "DELETE FROM entry_tags WHERE entry_id = ? AND tag_id = ?"
	_, err := r.db.Exec(q, entryID, tagID)
	return err
}
