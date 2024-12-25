package site

import (
	"html/template"
	"path/filepath"

	"github.com/connorkuljis/content/internal/model"
)

type Page interface {
	Filepath() string
	Data() map[string]any
	TemplateName() string
}

type HomePage struct {
	Site       *Site
	Categories []model.Category
}

func (p HomePage) Filepath() string {
	return "public/index.html"
}

func (p HomePage) TemplateName() string {
	return "index.html"
}

func (p HomePage) Data() map[string]any {
	return map[string]any{
		"Site":       p.Site,
		"Categories": p.Categories,
	}
}

type CategoryPage struct {
	Site     *Site
	Category model.Category
}

func (p CategoryPage) Filepath() string {
	return filepath.Join("public", slugify(p.Category.Title), "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "category.html"
}

func (p CategoryPage) Data() map[string]any {
	return map[string]any{
		"Site":     p.Site,
		"Category": p.Category,
	}
}

type EntryPage struct {
	Site    *Site
	Entry   model.Entry
	Content template.HTML
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", slugify(p.Entry.CategoryTitle), slugify(p.Entry.Title), "index.html")
}

func (p EntryPage) TemplateName() string {
	return "entry.html"
}

func (p EntryPage) Data() map[string]any {
	return map[string]any{
		"Site":    p.Site,
		"Entry":   p.Entry,
		"Content": p.Content,
	}
}
