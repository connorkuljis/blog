package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"html/template"

	"github.com/connorkuljis/blog/internal/util"
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

	// below fields are computed values, that may or may not be populated.
	Category *Category
	Markdown template.HTML
}

func NewEntry(categoryID int64, title string) *Entry {
	now := time.Now()
	return &Entry{
		CategoryID: categoryID,
		Title:      title,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsDraft:    util.BoolToInt(true),
	}
}

func (e Entry) String() string {
	var sb strings.Builder

	// frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("category_id: %d\n", e.CategoryID))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", e.CreatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", e.UpdatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("is_draft: %d\n", e.IsDraft))
	sb.WriteString("---\n")
	// ./frontmatter

	sb.WriteString(fmt.Sprintf("%s", e.Content))

	return sb.String()
}

func (e *Entry) AddCategory(category *Category) {
	e.Category = category
}

func (e *Entry) ToHTML(parser goldmark.Markdown) error {
	var buf bytes.Buffer
	err := parser.Convert([]byte(e.Content), &buf)
	if err != nil {
		return err
	}

	e.Markdown = template.HTML(buf.String())

	e.WordCount()

	return nil
}

func (e Entry) Permalink() string {
	if e.Category.Permalink() == "" {
		return "/"
	}

	return e.Category.Permalink() + "/" + e.CreatedAt.Format("2006-01-02") + "-" + util.Slugify(e.Title)
}

func (e *Entry) WordCount() int {
	return len(strings.Split(e.Content, " "))
}
