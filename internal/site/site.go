package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
)

type Site struct {
	Title      string
	base       map[string]string
	components map[string]string
	layouts    map[string]string
	views      map[string]string
}

func Render() error {
	site := Site{
		Title: "Connor's Blog",
	}

	err := site.Init()
	if err != nil {
		return err
	}

	db, err := database.Connect()
	if err != nil {
		return err
	}

	er := model.NewEntryRepository(db)
	entries, err := er.ReadAllJoinCategories()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err := site.RenderEntry("public/posts", entry)
		if err != nil {
			return err
		}
	}

	cr := model.NewCategoryRepository(db)
	categories, err := cr.ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}
	for _, category := range categories {
		err := site.RenderCategory("public/categories", category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) Init() error {
	err := os.RemoveAll("public")
	if err != nil {
		return err
	}
	os.MkdirAll("public", os.ModePerm)

	s.base, err = loadTemplates("templates")
	if err != nil {
		return err
	}
	s.layouts, err = loadTemplates("templates/layouts")
	if err != nil {
		return err
	}
	s.components, err = loadTemplates("templates/components")
	if err != nil {
		return err
	}
	s.views, err = loadTemplates("templates/views")
	if err != nil {
		return err
	}

	return nil
}

func (s *Site) RenderEntry(dir string, entry model.Entry) error {
	templateStrings := []string{
		s.base["base.html"],
		s.base["head.html"],
		s.base["layout.html"],

		s.components["header.html"],

		s.views["entry.html"],
	}

	t, err := template.New("").Parse(strings.Join(templateStrings, " "))
	if err != nil {
		return err
	}

	os.MkdirAll(dir, os.ModePerm)

	filename := filepath.Join(dir, fmt.Sprintf("%d.html", entry.ID))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	data := map[string]any{
		"Site":  s,
		"Entry": entry,
	}

	err = t.ExecuteTemplate(f, "base", data)
	if err != nil {
		return err
	}

	log.Println(filename)

	return nil
}

func (s *Site) RenderCategory(dir string, category model.Category) error {
	templateStrings := []string{
		s.base["base.html"],
		s.base["head.html"],
		s.base["layout.html"],

		s.components["header.html"],

		s.views["category.html"],
	}

	t, err := template.New("").Parse(strings.Join(templateStrings, " "))
	if err != nil {
		return err
	}

	os.MkdirAll(dir, os.ModePerm)

	filename := filepath.Join(dir, fmt.Sprintf("%s.html", category.Title))
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	data := map[string]any{
		"Site":     s,
		"Category": category,
	}

	err = t.ExecuteTemplate(f, "base", data)
	if err != nil {
		return err
	}

	log.Println(filename)

	return nil
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
