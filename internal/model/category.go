package model

import (
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
)

type Category struct {
	ID          int64
	Title       string
	Description string
	Permalink   string
	Entries     []*Entry
}

func NewCategory(c *store.Category) *Category {
	return &Category{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Entries:     []*Entry{},
		Permalink:   "/" + util.Slugify(c.Title),
	}
}

// ToStoreCategory maps this model.Category to a store.Category.
func (c *Category) ToStoreCategory() *store.Category {
	return &store.Category{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
	}
}

func (c *Category) AddEntry(e ...*Entry) {
	c.Entries = append(c.Entries, e...)
}
