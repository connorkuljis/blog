package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/connorkuljis/blog/internal/dto"
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

	entries, err := entryRepo.ReadAllEntries(enableDrafts)
	if err != nil {
		log.Fatal(fmt.Errorf("error initialising site: %w", err))
	}

	nerdStats := dto.NewNerdStats(time.Now())

	return site.NewSite(Title, Author, DirBuild, DirAssets, time.Now(), t, categories, entries, nerdStats, md)
}
