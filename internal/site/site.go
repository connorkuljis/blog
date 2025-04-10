package site

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/util"
)

type MySite struct {
	Title      string
	Categories []model.Category
	NerdStats  *model.NerdStats
}

func (s *MySite) Init() error {
	s.NerdStats.StartTime = time.Now()

	err := os.RemoveAll(s.RootDir())
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.RootDir(), os.ModePerm)
	if err != nil {
		return err
	}

	staticAssets := os.DirFS("static")
	err = os.CopyFS(s.RootDir(), staticAssets)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySite) Build() []model.Page {
	pages := []model.Page{
		HomePage{Site: s},
	}
	for _, category := range s.Categories {
		entries := category.Entries
		pages = append(pages, CategoryPage{
			Site:           s,
			Category:       category,
			GroupedEntries: model.GroupByYear(entries),
		})
		for _, entry := range entries {
			pages = append(pages, EntryPage{
				Site:     s,
				Category: category,
				Entry:    entry,
			})
		}
	}

	s.NerdStats.PageCount = len(pages)
	s.NerdStats.FinishTime = time.Now()

	return pages
}

func (s *MySite) Render(pages []model.Page) error {
	for _, page := range pages {
		dir := filepath.Dir(page.Filepath())
		os.MkdirAll(dir, os.ModePerm)

		f, err := os.Create(page.Filepath())
		if err != nil {
			return err
		}
		defer f.Close()

		err = s.Template().ExecuteTemplate(f, page.TemplateName(), page)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MySite) RootDir() string {
	return "public"
}

func (s *MySite) Template() *template.Template {
	var funcMap = template.FuncMap{
		"slugify":  util.Slugify,
		"truncate": util.Truncate,
		"sub": func(a, b int) int {
			return a - b
		},
	}

	return template.Must(template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html"))
}
