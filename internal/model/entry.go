package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"html/template"

	"github.com/yuin/goldmark"
)

type Entry struct {
	ID               int64          `db:"id"`
	CategoryID       int64          `db:"category_id"`
	Title            string         `db:"title"`
	Content          string         `db:"content"`
	Description      sql.NullString `db:"description"`
	FeaturedImageURL sql.NullString `db:"featured_image_url"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	IsDraft          int            `db:"is_draft"`

	CategoryTitle string
	Markdown      template.HTML
	Slug          string
	WordCount     int
}

func NewEntry(categoryID int64, title string) *Entry {
	now := time.Now()
	return &Entry{
		CategoryID: categoryID,
		Title:      title,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("category_id: %d\n", e.CategoryID))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", e.CreatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", e.UpdatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("publish: %d\n", e.IsDraft))
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("%s", e.Content))

	return sb.String()
}

func (e *Entry) ToHTML(parser goldmark.Markdown) error {
	var buf bytes.Buffer
	err := parser.Convert([]byte(e.Content), &buf)
	if err != nil {
		return err
	}

	e.Markdown = template.HTML(buf.String())

	return nil
}
