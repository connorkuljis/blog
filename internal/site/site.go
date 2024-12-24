package site

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type Site struct {
	Title string

	DB             *sqlx.DB
	MarkdownParser goldmark.Markdown
}

type Page interface {
	Filepath() string
	Data() map[string]any
	TemplateName() string
}

func (s *Site) Init() error {
	// clean public
	err := os.RemoveAll("public")
	if err != nil {
		return err
	}
	os.MkdirAll("public", os.ModePerm)

	// copy all static directory contents into public
	err = os.CopyFS("public", os.DirFS("static"))
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) Render() error {
	tpls, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return err
	}

	categories, err := store.NewCategoryRepository(s.DB).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	pages := []Page{
		HomePage{Categories: categories},
	}

	for _, category := range categories {
		pages = append(pages, CategoryPage{Category: category})
		for _, entry := range category.Entries {
			pages = append(pages, EntryPage{Category: category, Entry: entry})
		}
	}

	for _, page := range pages {
		err := Render(tpls, page)
		if err != nil {
			return err
		}
	}

	return nil
}

func Render(tpls *template.Template, page Page) error {
	dir, _ := filepath.Split(page.Filepath())
	os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(page.Filepath())
	if err != nil {
		return err
	}

	err = tpls.ExecuteTemplate(f, page.TemplateName(), page.Data())
	if err != nil {
		return err
	}

	return nil
}

type HomePage struct {
	Categories []model.Category
}

func (p HomePage) Filepath() string {
	return "public/index.html"
}

func (p HomePage) TemplateName() string {
	return "view-index.html"
}

func (p HomePage) Data() map[string]any {
	return map[string]any{
		"Categories": p.Categories,
	}
}

type CategoryPage struct {
	Category model.Category
}

func (p CategoryPage) Filepath() string {
	return filepath.Join("public", p.Category.Title, "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "view-category.html"
}

func (p CategoryPage) Data() map[string]any {
	return map[string]any{
		"Category": p.Category,
	}
}

type EntryPage struct {
	Category model.Category
	Entry    model.Entry
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", p.Category.Title, p.Entry.Title, "index.html")
}

func (p EntryPage) TemplateName() string {
	return "view-entry.html"
}

func (p EntryPage) Data() map[string]any {
	return map[string]any{
		"Entry": p.Entry,
	}
}

var funcMap = template.FuncMap{
	"slugify": slugify,
}

func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with a hyphen
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
