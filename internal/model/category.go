package model

import (
	"fmt"
	"strings"

	"github.com/connorkuljis/blog/internal/util"
)

type Category struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`

	Entries []Entry
	Slug    string
}

func NewCategory(title, description string) *Category {
	return &Category{Title: title, Description: description}
}

func (c *Category) AddEntries(entries []Entry) {
	c.Entries = append(c.Entries, entries...)
}

func (c *Category) Slugify() {
	c.Slug = util.Slugify(c.Title)
}

func (c Category) String() string {
	var sb strings.Builder

	sb.WriteString("Category: ")
	sb.WriteString(c.Title)
	sb.WriteString("\nID: ")
	sb.WriteString(fmt.Sprint(c.ID))
	sb.WriteString("\nDescription: ")
	sb.WriteString(c.Description)

	return sb.String()
}
