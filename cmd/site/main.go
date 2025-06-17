package main

import (
	"fmt"
	"html/template"
	"log"
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
	// Title     = "kuljis.xyz"
	Title     = "connorkuljis.com"
	DirBuild  = "dist"
	DirAssets = "assets"
)

var funcMap = template.FuncMap{
	"slugify":  util.Slugify,
	"truncate": util.Truncate,
	"runeToString": func(r rune) string {
		return string(r)
	},
	"sub": func(a, b int) int {
		return a - b
	},
}

func main() {
	md := goldmark.New(
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

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	t := template.Must(template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html"))

	enableDrafts := flag.Bool("d", false, "enable drafts")
	flag.Parse()

	if *enableDrafts {
		fmt.Println("draft mode is enabled")
	} else {
		fmt.Println("draft mode is disabled")
	}

	// inject dependencies and use interface type, rather than concrete type
	var site *site.MySite = initialiseKuljisSite(*enableDrafts, md, db, t)

	err = site.Init()
	if err != nil {
		log.Fatal(err)
	}

	pages := site.Build()

	start := time.Now()

	err = site.Render(pages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(time.Since(start))
}

func initialiseKuljisSite(enableDrafts bool, md goldmark.Markdown, db *sqlx.DB, t *template.Template) *site.MySite {
	var (
		categories = store.NewCategoryRepo(db)
		entries    = store.NewEntryRepo(db)
	)

	allCategories, err := categories.ReadAllCategories()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range allCategories {
		categoryEntries, err := entries.ReadAllByCategoryID(c.ID, enableDrafts)
		if err != nil {
			log.Fatal(err)
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

	categoriesMap := make(map[string]*model.Category, len(allCategories))
	for _, category := range allCategories {
		categoriesMap[category.Title] = category
	}

	nerdStats := model.NewNerdStats(time.Now())

	return site.NewSite(Title, DirBuild, DirAssets, time.Now(), t, allCategories, categoriesMap, nerdStats)
}
