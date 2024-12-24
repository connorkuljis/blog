package site

import (
	"bytes"
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
	categories, err := store.NewCategoryRepository(s.DB).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	pages := []Page{
		HomePage{
			Site:       s,
			Categories: categories,
		},
	}

	for _, category := range categories {
		pages = append(pages, CategoryPage{
			Site:     s,
			Category: category,
		})

	}

	entries, err := store.NewEntryRepository(s.DB).ReadAllJoinCategories()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		var buf bytes.Buffer
		err := s.MarkdownParser.Convert([]byte(entry.Content), &buf)
		if err != nil {
			return err
		}

		pages = append(pages, EntryPage{
			Site:    s,
			Entry:   entry,
			Content: template.HTML(buf.String()),
		})
	}

	funcMap := template.FuncMap{
		"slugify": slugify,
	}

	tpls, err := template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html")
	if err != nil {
		return err
	}

	for _, page := range pages {
		err := renderPage(tpls, page)
		if err != nil {
			return err
		}
	}

	return nil
}

func renderPage(tpls *template.Template, page Page) error {
	dir := filepath.Dir(page.Filepath())
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
	Site       *Site
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
	return "view-category.html"
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
	return "view-entry.html"
}

func (p EntryPage) Data() map[string]any {
	return map[string]any{
		"Site":    p.Site,
		"Entry":   p.Entry,
		"Content": p.Content,
	}
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
