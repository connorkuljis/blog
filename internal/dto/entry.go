package dto

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/util"
	"github.com/yuin/goldmark"
)

type Entry struct {
	model.Entry
	Category *Category
	Markdown template.HTML
}

func NewEntryDTO(m model.Entry) *Entry {
	return &Entry{
		Entry: m,
	}
}

func (e *Entry) AddCategory(category *Category) {
	e.Category = category
}

func (e *Entry) ToHTML(parser goldmark.Markdown) error {
	var buf bytes.Buffer
	err := parser.Convert([]byte(e.Content.String), &buf)
	if err != nil {
		return err
	}

	e.Markdown = template.HTML(buf.String())

	return nil
}

func (e Entry) Permalink() string {
	if e.Category == nil || e.Category.Permalink() == "" {
		return "/"
	}

	return e.Category.Permalink() + "/" + e.CreatedAt.Format("2006-01-02") + "-" + util.Slugify(e.Title)
}

func (e *Entry) WordCount() int {
	return len(strings.Split(e.Content.String, " "))
}
