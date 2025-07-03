package site

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/connorkuljis/blog/internal/dto"
	"github.com/connorkuljis/blog/internal/model"
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
	Categories    []*dto.Category
	CategoriesMap map[string]*dto.Category
	NerdStats     *dto.NerdStats
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
	categoriesModel []*model.Category,
	entriesModel []*model.Entry,
	nerdStats *dto.NerdStats,
	markdown goldmark.Markdown,
) *MySite {
	// Create maps for efficient lookups.
	// dtoCategoriesMap maps category IDs to category DTOs.
	// dtoCategoriesByTitle maps category titles to category DTOs.
	dtoCategoriesMap := make(map[int64]*dto.Category)
	dtoCategories := make([]*dto.Category, 0, len(categoriesModel))
	dtoCategoriesByTitle := make(map[string]*dto.Category)

	// Transform category models to category DTOs and populate the maps.
	for _, mCat := range categoriesModel {
		dtoCat := dto.NewCategoryDTO(*mCat)

		dtoCategories = append(dtoCategories, dtoCat)
		dtoCategoriesMap[dtoCat.ID] = dtoCat
		dtoCategoriesByTitle[dtoCat.Title] = dtoCat
	}

	// Transform entry models to entry DTOs.
	// It also associates each entry with its corresponding category DTO.
	for _, mEntry := range entriesModel {
		dtoEntry := dto.NewEntryDTO(*mEntry)

		// Find the category for the entry and add the entry to the category's list.
		if dtoCat, ok := dtoCategoriesMap[dtoEntry.CategoryID]; ok {
			dtoEntry.Category = dtoCat
			dtoCat.AddEntry(dtoEntry)
			// Convert the entry's markdown content to HTML.
			dtoEntry.ToHTML(markdown)
		}
	}

	return &MySite{
		Title:         title,
		Author:        author,
		DirBuild:      dirBuild,
		DirAssets:     dirAssets,
		CreatedAt:     createdAt,
		T:             t,
		Categories:    dtoCategories,
		CategoriesMap: dtoCategoriesByTitle,
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
