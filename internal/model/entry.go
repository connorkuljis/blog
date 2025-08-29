package model

import (
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
