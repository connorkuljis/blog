package repo

import (
	"github.com/connorkuljis/content/internal/model"
	"github.com/jmoiron/sqlx"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) CreateCategory(category *model.Category) error {
	_, err := r.db.Exec("INSERT INTO categories (title, description) VALUES (?, ?)", category.Title, category.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) ReadAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		return categories, err
	}

	return categories, nil
}
