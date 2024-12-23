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
	ID         int64     `db:"id"`
	CategoryID int64     `db:"category_id" yaml:"category_id"`
	Title      string    `db:"title" yaml:"title"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at" yaml:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" yaml:"updated_at"`
	Publish    int       `db:"publish" yaml:"publish"`

	CategoryTitle       string `db:"category_title"`
	CategoryDescription string `db:"category_description"`
}

func NewEntry(categoryID int64, title string) *Entry {
	return &Entry{CategoryID: categoryID, Title: title}
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("category_id: %d\n", e.CategoryID))
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

func (r *EntryRepository) ReadAllEntries() ([]Entry, error) {
	var entries []Entry
	err := r.db.Select(&entries, "SELECT * FROM entries")
	if err != nil {
		return nil, fmt.Errorf("Error getting all entries: %w", err)
	}

	return entries, nil
}

func (r *EntryRepository) ReadAllJoinCategories() ([]Entry, error) {
	var entries []Entry
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

func (r *EntryRepository) ReadEntryByID(id int64) (*Entry, error) {
	var entry Entry
	err := r.db.Get(&entry, "SELECT * FROM entries WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting entry by id `%d`: %w", id, err)
	}

	return &entry, nil
}

func (r *EntryRepository) UpdateEntry(entry *Entry) error {
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
