package kuljis

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type CategoryPage struct {
	Title string

	Site     *MySite
	Category *model.Category
}

func NewCategoryPage(site *MySite, category *model.Category, title string) CategoryPage {
	p := CategoryPage{
		Site:     site,
		Category: category,
		Title:    title,
	}
	return p
}

func (p CategoryPage) FileName() string {
	return filepath.Join(p.Category.Permalink(), "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "_category.html"
}
