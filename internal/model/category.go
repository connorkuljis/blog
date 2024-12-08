package model

type Category struct {
	Title       string `db:"title"`
	Description string `db:"description"`
	Entries     []Entry
}

func NewCategory(title, description string) *Category {
	return &Category{Title: title, Description: description}
}
