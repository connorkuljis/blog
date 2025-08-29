package model

import (
	"fmt"
	"strings"

	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
)

type Category struct {
	ID          int64
	Title       string
	Description string

	Permalink string
	Entries   []*Entry
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
	if c == nil {
		return nil
	}
	return &store.Category{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
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
