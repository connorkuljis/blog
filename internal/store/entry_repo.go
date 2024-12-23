package store

import (
	"fmt"

	"github.com/connorkuljis/content/internal/model"
	"github.com/jmoiron/sqlx"
)

type EntryRepository struct {
	db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *EntryRepository {
	return &EntryRepository{db: db}
}

func (r *EntryRepository) CreateEntry(entry *model.Entry) error {
	res, err := r.db.Exec("INSERT INTO entries (category_id, title) VALUES (?, ?)", entry.CategoryID, entry.Title)
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

func (r *EntryRepository) ReadAllEntries() ([]model.Entry, error) {
	var entries []model.Entry
	err := r.db.Select(&entries, "SELECT * FROM entries")
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadAllJoinCategories() ([]model.Entry, error) {
	var entries []model.Entry
	q := `
SELECT 
	e.id,
	e.title,
	e.category_id,
	e.content,
	e.created_at,
	e.updated_at,
	e.publish,
	c.title AS category_title,
	c.description AS category_description
FROM 
	entries AS e
INNER JOIN
	categories as c
ON
	c.id = e.category_id
ORDER BY 
	c.id
`

	err := r.db.Select(&entries, q)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadEntryByID(id int64) (*model.Entry, error) {
	var entry model.Entry
	err := r.db.Get(&entry, "SELECT * FROM entries WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting entry by id `%d`: %w", id, err)
	}

	return &entry, nil
}

func (r *EntryRepository) UpdateEntry(entry *model.Entry) error {
	q := "UPDATE entries SET category_id = ?, title = ?, content = ?, updated_at = ?, publish = ? WHERE id = ?"
	_, err := r.db.Exec(q, entry.CategoryID, entry.Title, entry.Content, entry.UpdatedAt, entry.Publish, entry.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *EntryRepository) DeleteEntryByID(id int64) error {
	_, err := r.db.Exec("DELETE FROM entries WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Error deleting entry by id `%d`: %w", id, err)
	}

	return nil
}
