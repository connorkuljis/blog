package model

import (
	"database/sql"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
)

type Entry struct {
	ID               int64
	Title            string
	Content          string
	Description      string
	FeaturedImageURL string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Category *Category
	Tags     []*Tag

	Markdown  template.HTML
	Permalink string
	WordCount int
}

func NewEntry(e *store.Entry, c *Category, t []*Tag) *Entry {
	entry := Entry{
		ID:               e.ID,
		Title:            e.Title,
		Content:          e.Content.String,
		Description:      e.Description.String,
		FeaturedImageURL: e.FeaturedImageURL.String,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
		Category:         c,
		Tags:             t,
	}

	// permalink
	timestamp := e.CreatedAt.Format("2006-01-02")
	title := util.Slugify(e.Title)
	slug := fmt.Sprintf("%s-%s", timestamp, title)
	entry.Permalink = filepath.Join(c.Permalink, slug)

	// word count
	entry.WordCount = len(strings.Split(e.Content.String, " "))

	return &entry
}

// ToStoreEntry maps this model.Entry to a store.Entry.
func (m *Entry) ToStoreEntry() *store.Entry {
	nullable := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{}
		}
		return sql.NullString{String: s, Valid: true}
	}

	var categoryID int64
	if m.Category != nil {
		categoryID = m.Category.ID
	}

	return &store.Entry{
		ID:               m.ID,
		CategoryID:       categoryID,
		Title:            m.Title,
		Content:          nullable(m.Content),
		Description:      nullable(m.Description),
		FeaturedImageURL: nullable(m.FeaturedImageURL),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
