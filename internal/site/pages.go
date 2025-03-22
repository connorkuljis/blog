package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type Page interface {
	Filepath() string
	TemplateName() string
}

type HomePage struct {
	Site *Site
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
}

func (p CategoryPage) Filepath() string {
	return filepath.Join("public", p.Category.Slug, "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "_category.html"
}

type EntryPage struct {
	Site  *Site
	Entry model.Entry
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", p.Entry.Slug, "index.html")
}

func (p EntryPage) TemplateName() string {
	return "_entry.html"
}
