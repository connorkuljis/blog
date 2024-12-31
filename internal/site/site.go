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

type Site struct {
	Title      string
	Categories []model.Category

	DB             *sqlx.DB
	MarkdownParser goldmark.Markdown
}

type NavItem struct {
	Name     string
	Resource string
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
	categories, err := store.NewCategoryRepository(s.DB).ReadAllCategories()
	if err != nil {
		return err
	}

	s.Categories = categories

	publishedEntries, err := store.NewEntryRepository(s.DB).ReadPublishedEntries()
	if err != nil {
		return err
	}

	recentEntries, err := store.NewEntryRepository(s.DB).ReadRecentlyPublishedEntries(10)
	if err != nil {
		return err
	}

	pages := []Page{
		HomePage{
			Site:          s,
			RecentEntries: recentEntries,
		},
	}

	for _, entry := range publishedEntries {
		// parse markdown to html
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

	for _, category := range categories {
		entries, err := store.NewEntryRepository(s.DB).ReadPublishedEntriesByCategoryID(category.ID)
		if err != nil {
			return err
		}

		pages = append(pages, CategoryPage{
			Site:     s,
			Category: category,
			Entries:  entries,
		})
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
