package model

import "fmt"

type Category struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Entries     []Entry
}

func NewCategory(title, description string) *Category {
	return &Category{Title: title, Description: description}
}

func (c Category) String() string {
	return fmt.Sprintf("Category %q (ID: %d)", c.Title, c.ID)
}
