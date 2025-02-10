package main

import (
	"fmt"
	"log"
	"time"

	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/site"
	"github.com/connorkuljis/content/internal/store"
	"github.com/connorkuljis/content/internal/util"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	title = "kuljis.xyz"
)

func main() {
	start := time.Now()

	site := &site.Site{
		Title: title,
	}

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	md := markdownParser()

	categories, err := getCategoriesAndEntries(db, md)
	if err != nil {
		log.Fatal(err)
	}
	site.Categories = categories

	recentEntries, err := getRecentEntries(db, md)
	if err != nil {
		log.Fatal(err)
	}
	site.RecentEntries = recentEntries

	err = site.Init()
	if err != nil {
		log.Fatal(err)
	}

	pages := site.BuildPages()

	err = site.RenderPages(pages)
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
		// Retrieve all entries associated with the current category's ID.
		entries, err := entryRepo.ReadAllByCategoryID(categories[i].ID)
		if err != nil {
			return categories, err
		}

		// Assign retrieved entries to corresponding category.
		// Establishes relationship between a category and its entries.
		categories[i].Entries = entries

		// Generate a slug for the category based on its title.
		categories[i].Slug = util.Slugify(categories[i].Title)

		for j := range entries {
			// Convert the Markdown content of the entry to HTML.
			err := entries[j].ContentMdToHTML(md)
			if err != nil {
				return categories, err
			}

			// Generate a slug for the entry, combining the category slug and the entry title.
			entries[j].Slug = categories[i].Slug + "/" + util.Slugify(entries[j].Title)
		}
	}

	return categories, nil
}

func getRecentEntries(db *sqlx.DB, md goldmark.Markdown) ([]model.Entry, error) {
	entries, err := store.NewEntryRepository(db).ReadRecentEntries(5)
	if err != nil {
		return entries, err
	}

	categoryRepo := store.NewCategoryRepository(db)
	for i := range entries {
		category, err := categoryRepo.ReadCategoryByID(entries[i].CategoryID)
		if err != nil {
			return entries, err
		}

		entries[i].Slug = fmt.Sprintf("%s/%s", util.Slugify(category.Title), util.Slugify(entries[i].Title))
		entries[i].ContentMdToHTML(md)
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
			html.WithUnsafe(),
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}
