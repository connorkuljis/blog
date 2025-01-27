package ssg

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/store"
	"github.com/connorkuljis/content/internal/util"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

var funcMap = template.FuncMap{
	"slugify":  util.Slugify,
	"truncate": util.Truncate,
}

type Site struct {
	Title    string
	NavItems []string

	DB             *sqlx.DB
	MarkdownParser goldmark.Markdown

	Publish bool
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

func filterPublishedEntries(entries []model.Entry) []model.Entry {
	var publishedEntries []model.Entry
	for _, entry := range entries {
		if entry.Publish == 1 {
			publishedEntries = append(publishedEntries, entry)
		}

	}
	return publishedEntries
}

func (s *Site) Render() error {
	// Setup repositories
	categoryRepo := store.NewCategoryRepository(s.DB)
	entryRepo := store.NewEntryRepository(s.DB)

	// 1. Read all categories
	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}

	// 2. Read all entries by category
	for i, category := range categories {
		entries, err := entryRepo.ReadAllByCategoryID(category.ID)
		if err != nil {
			return err
		}
		categories[i].Entries = entries
	}

	var pages []Page
	for _, category := range categories {
		var cp CategoryPage

		cp = CategoryPage{
			Site:     s,
			Category: category,
		}

		pages = append(pages, cp)

		for _, entry := range category.Entries {
			var ep EntryPage

			var buf bytes.Buffer
			err := s.MarkdownParser.Convert([]byte(entry.Content), &buf)
			if err != nil {
				return err
			}
			content := template.HTML(buf.String())

			ep = EntryPage{
				Site:         s,
				Category:     category,
				CurrentEntry: entry,
				Content:      content,
			}

			pages = append(pages, ep)
		}
	}

	var navItems []string
	err = s.DB.Select(&navItems, "SELECT title FROM categories ORDER BY title")
	if err != nil {
		return err
	}
	s.NavItems = navItems

	pages = append(pages, HomePage{
		Site:       s,
		Categories: categories,
	})

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
	defer f.Close()

	err = tpls.ExecuteTemplate(f, page.TemplateName(), page) // note: render page struct directly into the template data.
	if err != nil {
		return err
	}

	return nil
}
