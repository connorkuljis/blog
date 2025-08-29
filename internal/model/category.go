package model

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
)

type Category struct {
	ID          int64
	Title       string
	Description string
	Entries     []*Entry
}

func NewCategory(c *store.Category) *Category {
	return &Category{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Entries:     []*Entry{},
	}
}

func (c *Category) Permalink() string {
	return "/" + filepath.Join("categories", util.Slugify(c.Title))
}

func (c *Category) AddEntry(e ...*Entry) {
	c.Entries = append(c.Entries, e...)
}

func (c *Category) ToStoreCategory() *store.Category {
	return &store.Category{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
	}
}
