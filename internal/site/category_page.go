package site

import (
	"path/filepath"
	"sort"

	"github.com/connorkuljis/blog/internal/model"
)

type CategoryPage struct {
	Site           *MySite
	Category       *model.Category
	GroupedEntries [][]*model.Entry
}

func NewCategoryPage(site *MySite, category *model.Category) CategoryPage {
	p := CategoryPage{
		Site:     site,
		Category: category,
	}
	p.GroupEntriesByYear()
	return p
}

func (p CategoryPage) Filepath() string {
	return filepath.Join("public", p.Category.Permalink(), "index.html")
}

func (p CategoryPage) TemplateName() string {
	return "_category.html"
}

func (p *CategoryPage) GroupEntriesByYear() {
	entries := p.Category.Entries
	if len(entries) == 0 {
		return
	}

	// Group entries by year
	yearGroups := make(map[int][]*model.Entry)
	for _, entry := range entries {
		year := entry.CreatedAt.Year()
		yearGroups[year] = append(yearGroups[year], entry)
	}

	// Sort years in descending order
	years := make([]int, 0, len(yearGroups))
	for year := range yearGroups {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	// Build the final result slice
	p.GroupedEntries = make([][]*model.Entry, 0, len(years))
	for _, year := range years {
		p.GroupedEntries = append(p.GroupedEntries, yearGroups[year])
	}
}
