package site

import (
	"html/template"
	"path/filepath"

	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/util"
)

type Page interface {
	Filepath() string
	Data() map[string]any
	TemplateName() string
}

type HomePage struct {
	Site       *Site
	Categories []model.Category
	Entries    []model.Entry
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
		"Entries":    p.Entries,
	}
}

type CategoryPage struct {
	Site     *Site
	Category model.Category
	Entries  []model.Entry
}

func (p CategoryPage) Filepath() string {
	return filepath.Join("public", util.Slugify(p.Category.Title), "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "category.html"
}

func (p CategoryPage) Data() map[string]any {
	return map[string]any{
		"Site":     p.Site,
		"Category": p.Category,
		"Entries":  p.Entries,
	}
}

type EntryPage struct {
	Site     *Site
	Category model.Category
	Entry    model.Entry
	Content  template.HTML
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", util.Slugify(p.Category.Title), util.Slugify(p.Entry.Title), "index.html")
}

func (p EntryPage) TemplateName() string {
	return "entry.html"
}

func (p EntryPage) Data() map[string]any {
	return map[string]any{
		"Site":     p.Site,
		"Entry":    p.Entry,
		"Content":  p.Content,
		"Category": p.Category,
	}
}
