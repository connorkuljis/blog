package model

type Category struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Entries     []Entry
}

func NewCategory(title, description string) *Category {
	return &Category{Title: title, Description: description}
}
