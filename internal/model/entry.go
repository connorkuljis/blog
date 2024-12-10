package model

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/jmoiron/sqlx"
)

type Entry struct {
	Id        int64     `db:"id"`
	Category  string    `db:"category"`
	Title     string    `db:"title" yaml:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `db:"updated_at" yaml:"updated_at"`
	Publish   int       `db:"publish" yaml:"publish"`
}

func NewEntry(category string, title string) *Entry {
	return &Entry{Title: title, Category: category}
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("category: %s\n", e.Category))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", e.CreatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", e.UpdatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("publish: %d\n", e.Publish))
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("%s", e.Content))

	return sb.String()
}

func (e *Entry) LoadFromContentString(r io.Reader) error {
	b, err := frontmatter.Parse(r, e)
	if err != nil {
		return err
	}

	e.Content = string(b)

	return nil
}

type EntryRepository struct {
	db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *EntryRepository {
	return &EntryRepository{db: db}
}

func (r *EntryRepository) CreateEntry(entry *Entry) error {
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

func (r *EntryRepository) ReadAllEntries() ([]Entry, error) {
	var entries []Entry
	err := r.db.Select(&entries, "SELECT * FROM entries")
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadEntryByID(id int64) (*Entry, error) {
	var entry Entry
	err := r.db.Get(&entry, "SELECT * FROM entries WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting entry by id `%d`: %w", id, err)
	}

	return &entry, nil
}

func (r *EntryRepository) UpdateEntry(entry *Entry) error {
	q := "UPDATE entries SET category = ?, title = ?, content = ?, updated_at = ?, publish = ? WHERE id = ?"
	_, err := r.db.Exec(q, entry.Category, entry.Title, entry.Content, entry.UpdatedAt, entry.Publish, entry.Id)
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
