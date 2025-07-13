package model

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
	"github.com/yuin/goldmark"
)

type Entry struct {
	*store.Entry

	Category  *Category
	Tags      []*Tag
	Markdown  template.HTML
	Permalink string
	WordCount int
}

func NewEntry(e *store.Entry, c *Category, parser goldmark.Markdown) *Entry {
	entry := Entry{
		Entry:    e,
		Category: c,
	}

	// markdown
	var buf bytes.Buffer
	err := parser.Convert([]byte(e.Content.String), &buf)
	if err != nil {
		panic(err)
	}
	entry.Markdown = template.HTML(buf.String())

	// permalink
	timestamp := e.CreatedAt.Format("2006-01-02")
	title := util.Slugify(e.Title)
	slug := fmt.Sprintf("%s-%s", timestamp, title)
	entry.Permalink = filepath.Join(c.Permalink, slug)

	// word count
	entry.WordCount = len(strings.Split(e.Content.String, " "))

	return &entry
}
