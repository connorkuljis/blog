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
	"github.com/yuin/goldmark"
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
	markdown      goldmark.Markdown
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
	createdAt time.Time,
	t *template.Template,
	categoriesData []*store.Category,
	entriesData []*store.Entry,
	nerdStats *model.NerdStats,
	markdown goldmark.Markdown,
) *MySite {
	// Create maps for efficient lookups.
	allCategories := make([]*model.Category, 0, len(categoriesData))
	mapCategoryIDToCategory := make(map[int64]*model.Category)
	mapCategoryTitleToCategory := make(map[string]*model.Category)

	// Transform category models to category DTOs and populate the maps.
	for _, c := range categoriesData {
		category := model.NewCategory(*c)

		allCategories = append(allCategories, category)
		mapCategoryIDToCategory[category.ID] = category
		mapCategoryTitleToCategory[category.Title] = category
	}

	// Transform entry models to entry DTOs.
	// It also associates each entry with its corresponding category DTO.
	for _, e := range entriesData {
		entry := model.NewEntry(*e)

		// Find the category for the entry and add the entry to the category's list.
		if category, ok := mapCategoryIDToCategory[entry.CategoryID]; ok {
			entry.AddCategory(category)
			category.AddEntry(entry)
			// Convert the entry's markdown content to HTML.
			entry.ToHTML(markdown)
		}
	}

	return &MySite{
		Title:         title,
		Author:        author,
		DirBuild:      dirBuild,
		DirAssets:     dirAssets,
		CreatedAt:     createdAt,
		T:             t,
		Categories:    allCategories,
		CategoriesMap: mapCategoryTitleToCategory,
		NerdStats:     nerdStats,
		markdown:      markdown,
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

	var latestEntry *dto.Entry
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
			var next *dto.Entry
			var prev *dto.Entry

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
