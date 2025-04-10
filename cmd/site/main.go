package main

import (
	"fmt"
	"log"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/site"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	title = "kuljis.xyz"
)

func main() {
	md := goldmark.New(
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

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	categories, err := store.NewCategoryRepo(db).ReadAllCategories()
	if err != nil {
		log.Fatal(err)
	}

	entryRepo := store.NewEntryRepo(db)

	for i := range categories {
		entries, err := entryRepo.ReadAllByCategoryID(categories[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		for j := range entries {
			err := entries[j].ToHTML(md) // side-effect: updates markdown field.
			if err != nil {
				log.Fatal(err)
			}
			entries[j].AddCategory(categories[i])
			categories[i].AddEntry(entries[j])
		}
	}

	start := time.Now()

	site := &site.MySite{
		Title:      title,
		Categories: categories,
		NerdStats:  &model.NerdStats{},
	}

	err = site.Init()
	if err != nil {
		log.Fatal(err)
	}

	pages := site.Build()

	err = site.Render(pages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("built %d pages for %s in: %d ms\n", site.NerdStats.PageCount,
		site.Title, time.Since(start).Milliseconds())
}
