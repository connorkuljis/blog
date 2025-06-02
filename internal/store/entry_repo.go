package store

import (
	"fmt"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/jmoiron/sqlx"
)

type EntryRepo struct {
	db *sqlx.DB
}

func NewEntryRepo(db *sqlx.DB) *EntryRepo {
	return &EntryRepo{db: db}
}

func (r *EntryRepo) CreateEntry(entry *model.Entry) error {
	q := "INSERT INTO entries (category_id, title, created_at, updated_at) VALUES (?, ?, ?, ?)"

	res, err := r.db.Exec(q, entry.CategoryID, entry.Title,
		entry.CreatedAt.Format(time.RFC3339), entry.UpdatedAt.Format(time.RFC3339))
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

func (r *EntryRepo) ReadAllEntries(includeDrafts bool) ([]*model.Entry, error) {
	var entries []*model.Entry

	q := "SELECT * FROM entries"

	// if includeDrafts is false (release mode), get only non-draft (released) entries.
	// note: sqlite does not have bools, so we use integers.
	if !includeDrafts {
		q += " WHERE is_draft = 0"
	}

	q += " ORDER BY created_at DESC"

	// Execute the query
	err := r.db.Select(&entries, q)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) ReadRecentEntries(limit int) ([]model.Entry, error) {
	var entries []model.Entry
	err := r.db.Select(&entries, "SELECT * FROM entries ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) ReadAllByCategoryID(categoryID int64, enableDrafts bool) ([]*model.Entry, error) {
	var entries []*model.Entry

	q := "SELECT * FROM entries WHERE category_id = ?"

	if !enableDrafts {
		q += "AND is_draft = 0"
	}

	q += " ORDER BY created_at DESC"

	err := r.db.Select(&entries, q, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepo) ReadEntryByID(id int64) (*model.Entry, error) {
	var entry model.Entry
	err := r.db.Get(&entry, "SELECT * FROM entries WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting entry by id `%d`: %w", id, err)
	}

	return &entry, nil
}

func (r *EntryRepo) UpdateEntry(entry *model.Entry) error {
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
	_, err := r.db.Exec("DELETE FROM entries WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Error deleting entry by id `%d`: %w", id, err)
	}

	return nil
}
