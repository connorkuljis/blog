package site

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/pkg/site"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type MySite struct {
	Title        string
	Author       string
	DirBuild     string
	DirAssets    string
	EnableDrafts bool

	CreatedAt     time.Time
	T             *template.Template
	Categories    []*model.Category
	CategoriesMap map[string]*model.Category
	NerdStats     *model.NerdStats
}

// NewSite creates a new MySite instance.
// It takes the raw database models and transforms them into DTOs (Data Transfer Objects)
// that are suitable for rendering the website.
// This function is the main entry point for building the site's data structure.
func NewSite(
	title string,
	author string,
	dirBuild string,
	dirAssets string,
	enableDrafts bool,
	createdAt time.Time,
	t *template.Template,
	db *sqlx.DB,
	nerdStats *model.NerdStats,
	markdown goldmark.Markdown,
) *MySite {
	site := &MySite{
		Title:     title,
		Author:    author,
		DirBuild:  dirBuild,
		DirAssets: dirAssets,
		CreatedAt: createdAt,
		T:         t,
		NerdStats: nerdStats,
	}

	categories, err := store.NewCategoryRepo(db).ReadAllCategories()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range categories {
		mCategory := model.NewCategory(c)
		site.Categories = append(site.Categories, mCategory)

		entries, err := store.NewEntryRepo(db).ReadAllByCategoryID(c.ID, enableDrafts)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			mEntry := model.NewEntry(e, mCategory, markdown)
			mCategory.AddEntry(mEntry)
		}
	}

	return site
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

	pages = append(pages, NewHomePage(s))
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

	s.NerdStats.PageCount = len(pages)

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
