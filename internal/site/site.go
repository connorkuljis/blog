package site

import (
	"html/template"
	"os"
	"path/filepath"

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

func (s *Site) Render() error {
	categoryRepo := store.NewCategoryRepository(s.DB)
	entryRepo := store.NewEntryRepository(s.DB)

	var navItems []string
	err := s.DB.Select(&navItems, "SELECT title FROM categories ORDER BY title")
	if err != nil {
		return err
	}
	s.NavItems = navItems

	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}

	var pages []Page
	for i := range categories {
		entries, err := entryRepo.ReadAllByCategoryID(categories[i].ID)
		if err != nil {
			return err
		}

		categories[i].Entries = entries

		for j := range categories[i].Entries {
			err := categories[i].Entries[j].ContentMdToHTML(s.MarkdownParser)
			if err != nil {
				return err
			}

			pages = append(pages, EntryPage{
				Site:         s,
				Category:     categories[i],
				CurrentEntry: categories[i].Entries[j],
			})
		}

		pages = append(pages, CategoryPage{
			Site:     s,
			Category: categories[i],
		})

	}

	pages = append(pages, HomePage{Site: s, Categories: categories})

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
