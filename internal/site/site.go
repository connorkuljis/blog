package site

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/content/internal/store"
	"github.com/connorkuljis/content/internal/util"
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
	categories, err := store.NewCategoryRepository(s.DB).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	allEntries, err := store.NewEntryRepository(s.DB).ReadAllEntriesWithCategoryData()
	if err != nil {
		return err
	}

	pages := []Page{
		HomePage{
			Site:       s,
			Categories: categories,
			Entries:    allEntries,
		},
	}

	for _, category := range categories {
		matchedEntries, err := store.NewEntryRepository(s.DB).ReadAllEntriesByCategoryID(category.ID)
		if err != nil {
			return err
		}

		pages = append(pages, CategoryPage{
			Site:     s,
			Category: category,
			Entries:  matchedEntries,
		})

		for _, entry := range matchedEntries {
			var buf bytes.Buffer
			err := s.MarkdownParser.Convert([]byte(entry.Content), &buf)
			if err != nil {
				return err
			}

			pages = append(pages, EntryPage{
				Site:     s,
				Category: category,
				Entry:    entry,
				Content:  template.HTML(buf.String()),
			})
		}
	}

	funcMap := template.FuncMap{
		"slugify":  util.Slugify,
		"truncate": util.Truncate,
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
