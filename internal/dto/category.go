package dto

import (
	"fmt"
	"strings"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/util"
)

type Category struct {
	model.Category
	Entries []*Entry
}

func NewCategoryDTO(m model.Category) *Category {
	return &Category{
		Category: m,
		Entries:  []*Entry{},
	}
}

func (c *Category) AddEntry(e ...*Entry) {
	c.Entries = append(c.Entries, e...)
}

func (c Category) Permalink() string {
	return "/" + util.Slugify(c.Title)
}

func (c Category) String() string {
	var sb strings.Builder

	sb.WriteString("Category: ")
	sb.WriteString(c.Title)
	sb.WriteString("ID: ")
	sb.WriteString(fmt.Sprint(c.ID))
	sb.WriteString("Description: ")
	sb.WriteString(c.Description)

	return sb.String()
}
