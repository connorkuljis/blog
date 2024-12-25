package site

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/connorkuljis/content/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type Site struct {
	Title string

	DB             *sqlx.DB
	MarkdownParser goldmark.Markdown
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
	var pages []Page

	categories, err := store.NewCategoryRepository(s.DB).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	homePage := HomePage{
		Site:       s,
		Categories: categories,
	}
	pages = append(pages, homePage)

	for _, category := range categories {
		categoryPage := CategoryPage{
			Site:     s,
			Category: category,
		}
		pages = append(pages, categoryPage)
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

		entryPage := EntryPage{
			Site:    s,
			Entry:   entry,
			Content: template.HTML(buf.String()),
		}
		pages = append(pages, entryPage)
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

func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with a hyphen
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
