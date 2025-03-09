package store

import (
	"fmt"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/jmoiron/sqlx"
)

type EntryRepository struct {
	db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *EntryRepository {
	return &EntryRepository{db: db}
}

func (r *EntryRepository) CreateEntry(entry *model.Entry) error {
	res, err := r.db.Exec("INSERT INTO entries (category_id, title, created_at, updated_at) VALUES (?, ?, ?, ?)", entry.CategoryID, entry.Title, entry.CreatedAt.Format(time.RFC3339), entry.UpdatedAt.Format(time.RFC3339))
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
	err := r.db.Select(&entries, "SELECT * FROM entries ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadRecentEntries(limit int) ([]model.Entry, error) {
	var entries []model.Entry
	err := r.db.Select(&entries, "SELECT * FROM entries ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadRecentPublishedEntries(limit int) ([]model.Entry, error) {
	var entries []model.Entry
	err := r.db.Select(&entries, "SELECT * FROM entries WHERE is_draft = 0 ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadAllByCategoryID(categoryID int64) ([]model.Entry, error) {
	var entries []model.Entry
	err := r.db.Select(&entries, "SELECT * FROM entries WHERE category_id = ? ORDER BY created_at DESC", categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadAllPublishedByCategoryID(categoryID int64) ([]model.Entry, error) {
	var entries []model.Entry
	err := r.db.Select(&entries, "SELECT * FROM entries WHERE category_id = ? AND is_draft = 0 ORDER BY created_at DESC", categoryID)
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
	q := "UPDATE entries SET category_id = ?, title = ?, content = ?, description = ?, featured_image_url = ?, updated_at = ?, is_draft = ? WHERE id = ?"
	_, err := r.db.Exec(q, entry.CategoryID, entry.Title, entry.Content, entry.Description, entry.FeaturedImageURL, entry.UpdatedAt, entry.IsDraft, entry.ID)
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
