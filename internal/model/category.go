package model

import (
	"github.com/jmoiron/sqlx"
)

type Category struct {
	Title       string `db:"title"`
	Description string `db:"description"`
	Entries     []Entry
}

func NewCategory(title, description string) *Category {
	return &Category{Title: title, Description: description}
}

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) CreateCategory(category *Category) error {
	_, err := r.db.Exec("INSERT INTO categories (title, description) VALUES (?, ?)", category.Title, category.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) ReadAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return categories, err
	}

	return categories, nil
}
