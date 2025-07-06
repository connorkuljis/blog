package model

import (
	"fmt"
	"strings"

	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
)

type Category struct {
	*store.Category

	Permalink string
	Entries   []*Entry
}

func NewCategory(c *store.Category) *Category {
	return &Category{
		Category:  c,
		Entries:   []*Entry{},
		Permalink: "/" + util.Slugify(c.Title),
	}
}

func (c *Category) AddEntry(e ...*Entry) {
	c.Entries = append(c.Entries, e...)
}

func (c Category) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprint("Category:", c.Title))
	sb.WriteString(fmt.Sprint("ID:", c.ID))
	sb.WriteString(fmt.Sprint("Description:", c.Description))
	return sb.String()
}
