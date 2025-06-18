package site

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/pkg/site"
)

type MySite struct {
	Title     string
	Author    string
	DirBuild  string
	DirAssets string

	CreatedAt     time.Time
	T             *template.Template
	Categories    []*model.Category
	CategoriesMap map[string]*model.Category
	NerdStats     *model.NerdStats
}

func NewSite(
	title string,
	author string,
	dirBuild string,
	dirAssets string,
	createdAt time.Time,
	t *template.Template,
	categories []*model.Category,
	categoriesMap map[string]*model.Category,
	nerdStats *model.NerdStats,
) *MySite {
	return &MySite{
		Title:         title,
		Author:        author,
		DirBuild:      dirBuild,
		DirAssets:     dirAssets,
		CreatedAt:     createdAt,
		T:             t,
		Categories:    categories,
		CategoriesMap: categoriesMap,
		NerdStats:     nerdStats,
	}
}

func (s *MySite) Init() error {
	err := os.RemoveAll(s.DirBuild)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.DirBuild, os.ModePerm)
	if err != nil {
		return err
	}

	staticAssets := os.DirFS(s.DirAssets)
	err = os.CopyFS(s.DirBuild, staticAssets)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySite) Build() []site.Page {
	var pages = []site.Page{}

	var latestEntry *model.Entry
	for _, c := range s.Categories {
		for _, e := range c.Entries {
			if latestEntry == nil || e.CreatedAt.After(latestEntry.CreatedAt) {
				latestEntry = e
			}
		}
	}

	pages = append(pages, NewHomePage(s, latestEntry))
	pages = append(pages, NewAboutPage(s))

	for _, category := range s.Categories {
		pages = append(pages, NewCategoryPage(s, category))

		for i, entry := range category.Entries {
			var next *model.Entry
			var prev *model.Entry

			// Get previous entry if exists
			if i > 0 {
				prev = category.Entries[i-1]
			}

			// Get next entry if exists
			if i < len(category.Entries)-1 {
				next = category.Entries[i+1]
			}

			pages = append(pages, NewEntryPage(s, category, entry, next, prev))
		}
	}

	s.NerdStats.SetPageCount(len(pages))

	return pages
}

func (s *MySite) Render(pages []site.Page) error {
	for _, page := range pages {
		filename := filepath.Join(s.DirBuild, page.FileName())

		dir := filepath.Dir(filename)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		err = s.T.ExecuteTemplate(f, page.TemplateName(), page) // note: assume each page is rendered as concrete type when accessing data in templates.
		if err != nil {
			return err
		}
		log.Println("created:", page.FileName())
	}

	return nil
}
