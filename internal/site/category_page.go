package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type CategoryPage struct {
	Site           *MySite
	Category       model.Category
	GroupedEntries [][]model.Entry
}

func (p CategoryPage) Filepath() string {
	return filepath.Join("public", p.Category.Slug(), "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "_category.html"
}
