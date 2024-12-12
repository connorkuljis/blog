package model

import (
	"github.com/jmoiron/sqlx"
)

type Category struct {
	ID          int64  `db:"id"`
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

func (r *CategoryRepo) ReadCategoryByID(id int64) (Category, error) {
	var category Category
	err := r.db.Get(&category, "SELECT * FROM categories WHERE id = ?", id)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *CategoryRepo) ReadAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (r *CategoryRepo) ReadAllCategoriesWithEntries() ([]Category, error) {
	categories, err := r.ReadAllCategories()
	if err != nil {
		return []Category{}, err
	}

	for i, c := range categories {
		var entries []Entry
		err := r.db.Select(&entries, "SELECT * FROM entries WHERE category_id = ?", c.ID)
		if err != nil {
			return []Category{}, err
		}
		categories[i].Entries = entries
	}

	return categories, nil
}

func (r *CategoryRepo) DeleteCategoryByTitle(id int64) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
