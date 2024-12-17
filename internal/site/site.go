package site

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
	"github.com/jmoiron/sqlx"
)

type Site struct {
	Title      string
	Base       map[string]string // maps filename -> html string
	Components map[string]string // eg: header.html -> <h1>...</h1>
	Layouts    map[string]string
	Views      map[string]string

	DB *sqlx.DB
}

func (s *Site) AllComponents() []string {
	var cmps []string
	for _, component := range s.Components {
		cmps = append(cmps, component)
	}
	return cmps
}

func loadTemplates(src string) (map[string]string, error) {
	t := make(map[string]string)

	files, err := os.ReadDir(src)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" {
			filename := filepath.Join(src, file.Name())
			b, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}

			t[file.Name()] = string(b)
		}
	}

	return t, nil
}

func (s *Site) Init() error {
	// db connection
	db, err := database.Connect()
	if err != nil {
		return err
	}
	s.DB = db

	// clean public
	err = os.RemoveAll("public")
	if err != nil {
		return err
	}
	os.MkdirAll("public", os.ModePerm)

	// load template strings into maps
	s.Base, err = loadTemplates("templates")
	if err != nil {
		return err
	}
	s.Layouts, err = loadTemplates("templates/layouts")
	if err != nil {
		return err
	}
	s.Components, err = loadTemplates("templates/components")
	if err != nil {
		return err
	}
	s.Views, err = loadTemplates("templates/views")
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) Render() error {
	err := s.Init()
	if err != nil {
		return err
	}

	err = s.RenderEntry("public/entries")
	if err != nil {
		return err
	}

	err = s.RenderCategory("public/categories")
	if err != nil {
		return err
	}

	err = s.RenderIndex("public")

	return nil
}

func (s *Site) RenderEntry(dir string) error {
	entries, err := model.NewEntryRepository(s.DB).ReadAllJoinCategories()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		page := Page{
			Title:              entry.Title,
			Filepath:           filepath.Join(dir, fmt.Sprintf("%d.html", entry.ID)),
			BaseTemplate:       s.Base["base.html"],
			LayoutTemplate:     s.Base["layout.html"],
			HeadTemplate:       s.Base["head.html"],
			ViewTemplate:       s.Views["entry.html"],
			ComponentTemplates: s.AllComponents(),
			Data: map[string]any{
				"Site":  s,
				"Entry": entry,
			},
		}

		err := page.Render()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) RenderCategory(dir string) error {
	categories, err := model.NewCategoryRepository(s.DB).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}
	for _, category := range categories {
		page := Page{
			Title:              category.Title,
			Filepath:           filepath.Join(dir, fmt.Sprintf("%d.html", category.ID)),
			BaseTemplate:       s.Base["base.html"],
			LayoutTemplate:     s.Base["layout.html"],
			HeadTemplate:       s.Base["head.html"],
			ViewTemplate:       s.Views["category.html"],
			ComponentTemplates: s.AllComponents(),
			Data: map[string]any{
				"Site":     s,
				"Category": category,
			},
		}

		err := page.Render()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) RenderIndex(dir string) error {
	categories, err := model.NewCategoryRepository(s.DB).ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	page := Page{
		Title:              "index.html",
		Filepath:           filepath.Join(dir, "index.html"),
		BaseTemplate:       s.Base["base.html"],
		LayoutTemplate:     s.Base["layout.html"],
		HeadTemplate:       s.Base["head.html"],
		ViewTemplate:       s.Views["entry.html"],
		ComponentTemplates: s.AllComponents(),
		Data: map[string]any{
			"Site":       s,
			"Categories": categories,
		},
	}

	err = page.Render()
	if err != nil {
		return err
	}

	return nil
}
