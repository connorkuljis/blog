package kuljis

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type MySite struct {
	Title         string
	DirWWWRoot    string
	DirAssets     string
	T             *template.Template
	Categories    []*model.Category
	CategoriesMap map[string]*model.Category
	NerdStats     *model.NerdStats
}

func NewSite(
	title string,
	dirWWWRoot string,
	dirAssets string,
	categories []*model.Category,
	categoriesMap map[string]*model.Category,
	nerdStats *model.NerdStats,
	t *template.Template,
) *MySite {
	return &MySite{
		Title:         title,
		DirWWWRoot:    dirWWWRoot,
		DirAssets:     dirAssets,
		Categories:    categories,
		CategoriesMap: categoriesMap,
		NerdStats:     nerdStats,
		T:             t,
	}
}

func (s *MySite) Init() error {
	err := os.RemoveAll(s.DirWWWRoot)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.DirWWWRoot, os.ModePerm)
	if err != nil {
		return err
	}

	staticAssets := os.DirFS(s.DirAssets)
	err = os.CopyFS(s.DirWWWRoot, staticAssets)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySite) Build() []model.Page {
	var pages = []model.Page{}

	pages = append(pages, NewHomePage(s, "Home"))

	for _, category := range s.Categories {
		pages = append(pages, NewCategoryPage(s, category, category.Title))

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

			pages = append(pages, NewEntryPage(s, category, entry, next, prev, entry.Title))
		}
	}

	s.NerdStats.SetPageCount(len(pages))

	return pages
}

func (s *MySite) Render(pages []model.Page) error {
	for _, page := range pages {
		filename := filepath.Join(s.DirWWWRoot, page.FileName())

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

		err = s.T.ExecuteTemplate(f, page.TemplateName(), page)
		if err != nil {
			return err
		}
		fmt.Println("-->", filename)
	}

	return nil
}
