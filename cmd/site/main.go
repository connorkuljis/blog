package main

import (
	"fmt"
	"html/template"
	"log"
	"sort"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/site"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/util"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	title           string = "kuljis.xyz"
	limit           int    = 10
	enableDraftMode bool   = true
)

var funcMap = template.FuncMap{
	"slugify":  util.Slugify,
	"truncate": util.Truncate,
	"sub": func(a, b int) int {
		return a - b
	},
	"groupByYear": func(entries []model.Entry) [][]model.Entry {
		// Handle empty or nil input slice gracefully
		if len(entries) == 0 {
			return [][]model.Entry{} // Return an empty slice, not nil
		}

		// 1. Use a map to group entries by year.
		//    Key: Year (int)
		//    Value: Slice of entries for that year ([]model.Entry)
		groups := make(map[int][]model.Entry)

		for _, entry := range entries {
			year := entry.CreatedAt.Year()
			// Append the entry to the slice associated with its year.
			// If the key (year) doesn't exist yet, it will be created with a new slice.
			groups[year] = append(groups[year], entry)
		}

		// 2. Get the years (keys) from the map to sort them.
		years := make([]int, 0, len(groups))
		for year := range groups {
			years = append(years, year)
		}

		// 3. Sort the years chronologically.
		sort.Sort(sort.Reverse(sort.IntSlice(years)))

		// 4. Build the final result slice, ordered by the sorted years.
		//    Pre-allocate capacity for efficiency.
		result := make([][]model.Entry, 0, len(years))
		for _, year := range years {
			// Append the slice of entries for the current year (from the map)
			// to the result slice.
			result = append(result, groups[year])
		}

		return result
	},
}

func main() {
	start := time.Now()

	site := &site.Site{
		Title: title,
		Socials: []site.Social{
			{URL: "https://github.com/connorkuljis"},
			{URL: "https://linkedin.com/in/connor-kuljis"},
		},
	}

	err := site.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	md := markdownParser()

	site.Categories, err = getCategoriesAndEntries(db, md)
	if err != nil {
		log.Fatal(err)
	}

	site.RecentEntries, err = getRecentEntries(db, md, limit)
	if err != nil {
		log.Fatal(err)
	}

	pages := site.BuildPages()

	t, err := template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	err = site.RenderPages(t, pages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("built %d pages for %s in: %d ms\n", len(pages), site.Title, time.Since(start).Milliseconds())
}

// getCategoriesAndEntries retrieves all categories from the database and their associated entries.
// It populates the Entries field of each Category with the corresponding entries,
// generates slugs for categories and entries, and converts the Markdown content of entries to HTML.
func getCategoriesAndEntries(db *sqlx.DB, md goldmark.Markdown) ([]model.Category, error) {
	categories, err := store.NewCategoryRepository(db).ReadAllCategories()
	if err != nil {
		return categories, err
	}

	entryRepo := store.NewEntryRepository(db)

	for i := range categories {
		var entries []model.Entry
		var err error

		if enableDraftMode {
			entries, err = entryRepo.ReadAllByCategoryID(categories[i].ID)
		} else {
			entries, err = entryRepo.ReadAllPublishedByCategoryID(categories[i].ID)
		}

		if err != nil {
			return categories, err
		}

		categories[i].Slugify()

		for j := range entries {
			err := entries[j].ToHTML(md) // side-effect: updates markdown field.
			if err != nil {
				return categories, err
			}

			entries[j].Slugify(categories[i].Slug)
			entries[j].CalculateWordCount()
			entries[j].AddCategory(categories[i])
		}

		categories[i].AddEntries(entries)
	}

	return categories, nil
}

func getRecentEntries(db *sqlx.DB, md goldmark.Markdown, limit int) ([]model.Entry, error) {
	entryRepo := store.NewEntryRepository(db)
	categoryRepo := store.NewCategoryRepository(db)

	var entries []model.Entry
	var err error

	if enableDraftMode {
		entries, err = entryRepo.ReadRecentEntries(limit)
	} else {
		entries, err = entryRepo.ReadRecentPublishedEntries(limit)
	}

	if err != nil {
		return entries, err
	}

	for i := range entries {
		category, err := categoryRepo.ReadCategoryByID(entries[i].CategoryID)
		if err != nil {
			return entries, err
		}
		category.Slugify()

		entries[i].AddCategory(category)
		entries[i].Slugify(category.Slug)
		err = entries[i].ToHTML(md)
		if err != nil {
			return entries, err
		}
	}

	return entries, nil
}

func markdownParser() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			// html.WithUnsafe(),
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}
