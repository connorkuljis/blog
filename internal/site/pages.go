package site

import (
	"html/template"
	"path/filepath"

	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/util"
)

type Page interface {
	Filepath() string
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
	return "_index.html"
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
	return "_category.html"
}

type EntryPage struct {
	Site         *Site
	Category     model.Category
	CurrentEntry model.Entry
	Content      template.HTML
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", util.Slugify(p.Category.Title), util.Slugify(p.CurrentEntry.Title), "index.html")
}

func (p EntryPage) TemplateName() string {
	return "_entry.html"
}
