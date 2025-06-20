package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/site"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	Title     = "kuljis.xyz"
	Author    = "Connor Kuljis"
	DirBuild  = "dist"
	DirAssets = "assets"
)

var (
	funcMap = template.FuncMap{
		"runeToString": func(r rune) string {
			return string(r)
		},
	}
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	t = template.Must(template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html"))
)

func main() {
	enableDrafts := flag.Bool("d", false, "enable drafts")

	flag.Parse()

	log.Println("Enable drafts:", *enableDrafts)

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	mySite := initialiseKuljisSite(*enableDrafts, md, db, t)

	err = mySite.Init()
	if err != nil {
		log.Fatal(err)
	}

	pages := mySite.Build()

	start := time.Now()

	err = mySite.Render(pages)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(time.Since(start))
}

func initialiseKuljisSite(enableDrafts bool, md goldmark.Markdown, db *sqlx.DB, t *template.Template) *site.MySite {
	categoryRepo := store.NewCategoryRepo(db)
	entryRepo := store.NewEntryRepo(db)

	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		log.Fatal(fmt.Errorf("error initialising site: %w", err))
	}

	for _, c := range categories {
		categoryEntries, err := entryRepo.ReadAllByCategoryID(c.ID, enableDrafts)
		if err != nil {
			log.Fatal(fmt.Errorf("error initialising site %w", c.Title, err))

		}

		c.AddEntry(categoryEntries...)

		for _, e := range categoryEntries {
			err := e.ToHTML(md)
			if err != nil {
				log.Fatal(err)
			}

			e.AddCategory(c)
		}
	}

	categoriesMap := make(map[string]*model.Category, len(categories))

	for _, category := range categories {
		categoriesMap[category.Title] = category
	}

	nerdStats := model.NewNerdStats(time.Now())

	return site.NewSite(Title, Author, DirBuild, DirAssets, time.Now(), t, categories, categoriesMap, nerdStats)
}
