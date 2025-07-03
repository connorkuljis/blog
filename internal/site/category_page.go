package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/dto"
)

type CategoryPage struct {
	Site     *MySite
	Category *dto.Category
}

func NewCategoryPage(site *MySite, category *dto.Category) CategoryPage {
	p := CategoryPage{
		Site:     site,
		Category: category,
	}
	return p
}

func (p CategoryPage) Title() string {
	return p.Category.Title
}

func (p CategoryPage) FileName() string {
	return filepath.Join(p.Category.Permalink(), "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "_category.html"
}
