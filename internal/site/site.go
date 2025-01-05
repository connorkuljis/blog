package site

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
	Title      string
	Categories []model.Category
	Entries    []model.Entry

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
	categoryRepo := store.NewCategoryRepository(s.DB)
	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}
	s.Categories = categories

	entryRepo := store.NewEntryRepository(s.DB)
	siteEntries, err := entryRepo.ReadAllEntries()
	if err != nil {
		return err
	}
	s.Entries = siteEntries

	if s.Publish {
		s.Entries = filterPublishedEntries(siteEntries)
	}

	pages := []Page{
		HomePage{
			Site: s,
		},
	}

	for _, category := range categories {
		entries, err := entryRepo.ReadAllByCategoryID(category.ID)
		if err != nil {
			return err
		}

		if s.Publish {
			entries = filterPublishedEntries(siteEntries)
		}

		pages = append(pages, CategoryPage{
			Site:     s,
			Category: category,
			Entries:  entries,
		})

		// build the entry pages and parse md content to templatable safe HTML
		for _, entry := range entries {
			// parse markdown to html
			var buf bytes.Buffer
			err := s.MarkdownParser.Convert([]byte(entry.Content), &buf)
			if err != nil {
				return err
			}
			content := template.HTML(buf.String())

			pages = append(pages, EntryPage{
				Site:         s,
				Category:     category,
				CurrentEntry: entry,
				Content:      content,
			})
		}
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

	err = tpls.ExecuteTemplate(f, page.TemplateName(), page) // note: render page struct directly into the template data.
	if err != nil {
		return err
	}

	return nil
}
