package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type MySite struct {
	Title         string
	RootDir       string
	Template      *template.Template
	Categories    []*model.Category
	RecentEntries []*model.Entry
	CategoriesMap map[string]*model.Category
	NerdStats     *model.NerdStats
}

func NewSite(
	title string,
	rootDir string,
	categories []*model.Category,
	recentEntries []*model.Entry,
	categoriesMap map[string]*model.Category,
	nerdStats *model.NerdStats,
	t *template.Template,
) *MySite {
	return &MySite{
		Title:         title,
		RootDir:       rootDir,
		Categories:    categories,
		RecentEntries: recentEntries,
		CategoriesMap: categoriesMap,
		NerdStats:     nerdStats,
		Template:      t,
	}
}

func (s *MySite) Init() error {
	err := os.RemoveAll(s.RootDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.RootDir, os.ModePerm)
	if err != nil {
		return err
	}

	staticAssets := os.DirFS("static")
	err = os.CopyFS(s.RootDir, staticAssets)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySite) Build() []model.Page {
	var pages = []model.Page{}

	pages = append(pages, NewHomePage(s))

	for _, category := range s.Categories {
		pages = append(pages, NewCategoryPage(s, category))

		for _, entry := range category.Entries {
			pages = append(pages, NewEntryPage(s, category, entry))
		}
	}

	s.NerdStats.SetPageCount(len(pages))

	return pages
}

func (s *MySite) Render(pages []model.Page) error {
	for _, page := range pages {
		dir := filepath.Dir(page.Filepath())
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(page.Filepath())
		if err != nil {
			return err
		}
		defer f.Close()

		err = s.Template.ExecuteTemplate(f, page.TemplateName(), page)
		if err != nil {
			return err
		}
		fmt.Println("-->", page.Filepath())
	}

	return nil
}
