package site

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/content/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type Site struct {
	Title string
}

func (s *Site) Init() error {
	// clean public
	err := os.RemoveAll("public")
	if err != nil {
		return err
	}
	os.MkdirAll("public", os.ModePerm)

	return nil
}

func (s *Site) Render(db *sqlx.DB, lib *HTMLTemplateLibrary, md goldmark.Markdown) error {
	err := s.Init()
	if err != nil {
		return err
	}

	entryPageView := PageView{
		Name:               "entry.html",
		BaseTemplate:       lib.Base["base.html"],
		LayoutTemplate:     lib.Base["layout.html"],
		HeadTemplate:       lib.Base["head.html"],
		ComponentTemplates: lib.AllComponents(),
		ViewTemplate:       lib.Views["entry.html"],
	}
	entryTemplate, err := entryPageView.Template()
	if err != nil {
		return err
	}

	err = s.RenderEntry("public", db, md, entryTemplate)
	if err != nil {
		return err
	}

	categoryPageView := PageView{
		Name:               "category.html",
		BaseTemplate:       lib.Base["base.html"],
		LayoutTemplate:     lib.Base["layout.html"],
		HeadTemplate:       lib.Base["head.html"],
		ComponentTemplates: lib.AllComponents(),
		ViewTemplate:       lib.Views["category.html"],
	}
	categoryTemplate, err := categoryPageView.Template()
	if err != nil {
		return err
	}

	err = s.RenderCategory("public", db, categoryTemplate)
	if err != nil {
		return err
	}

	indexPageView := PageView{
		Name:               "index.html",
		BaseTemplate:       lib.Base["base.html"],
		LayoutTemplate:     lib.Base["layout.html"],
		HeadTemplate:       lib.Base["head.html"],
		ComponentTemplates: lib.AllComponents(),
		ViewTemplate:       lib.Views["index.html"],
	}
	indexTemplate, err := indexPageView.Template()
	if err != nil {
		return err
	}

	err = s.RenderIndex("public", db, indexTemplate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) RenderEntry(dir string, db *sqlx.DB, md goldmark.Markdown, t *template.Template) error {
	entries, err := model.NewEntryRepository(db).ReadAllJoinCategories()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		filepath := filepath.Join(dir, slugify(entry.CategoryTitle), slugify(entry.Title), "index.html")

		page := Page{
			Title:    entry.Title,
			Filepath: filepath,
			T:        t,
		}

		// convert markdown to html
		var buf bytes.Buffer
		if err := md.Convert([]byte(entry.Content), &buf); err != nil {
			panic(err)
		}

		data := map[string]any{
			"Site":    s,
			"Entry":   entry,
			"Content": template.HTML(buf.String()),
		}

		err := page.Render(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) RenderCategory(dir string, db *sqlx.DB, t *template.Template) error {
	categories, err := model.NewCategoryRepository(db).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	for _, category := range categories {
		filepath := filepath.Join(dir, slugify(category.Title), "index.html")

		page := Page{
			Title:    category.Title,
			Filepath: filepath,
			T:        t,
		}

		data := map[string]any{
			"Site":     s,
			"Category": category,
		}

		err := page.Render(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) RenderIndex(dir string, db *sqlx.DB, t *template.Template) error {
	categories, err := model.NewCategoryRepository(db).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	page := Page{
		Title:    "index.html",
		Filepath: filepath.Join(dir, "index.html"),
		T:        t,
	}

	data := map[string]any{
		"Site":       s,
		"Categories": categories,
	}

	err = page.Render(data)
	if err != nil {
		return err
	}

	return nil
}
