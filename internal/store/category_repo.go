package store

import (
	"github.com/connorkuljis/blog/internal/model"
	"github.com/jmoiron/sqlx"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) CreateCategory(category *model.Category) error {
	_, err := r.db.Exec("INSERT INTO categories (title, description) VALUES (?, ?)", category.Title, category.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) ReadCategoryByID(id int64) (model.Category, error) {
	var category model.Category
	err := r.db.Get(&category, "SELECT * FROM categories WHERE id = ?", id)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *CategoryRepo) ReadCategoryByTitle(title string) (model.Category, error) {
	var category model.Category
	err := r.db.Get(&category, "SELECT * FROM categories WHERE title LIKE ?", title)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r *CategoryRepo) ReadAllCategories() ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return categories, err
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
