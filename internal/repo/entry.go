package repo

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
	res, err := r.db.Exec("INSERT INTO entries (category, title) VALUES (?, ?)", entry.Category, entry.Title)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	entry.Id = id

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

func (r *EntryRepository) ReadEntryByID(id int64) (*model.Entry, error) {
	var entry model.Entry
	err := r.db.Get(&entry, "SELECT * FROM entries WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting entry by id `%s`: %w", id, err)
	}

	return &entry, nil
}

func (r *EntryRepository) UpdateEntry(entry *model.Entry) error {
	q := "UPDATE entries SET category = ?, title = ?, content = ?, updated_at = ?, publish = ? WHERE id = ?"
	_, err := r.db.Exec(q, entry.Category, entry.Title, entry.Content, entry.UpdatedAt, entry.Publish, entry.Id)
	if err != nil {
		return err
	}

	return nil
}
