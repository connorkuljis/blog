package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/connorkuljis/blog/internal/markdown"
	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/site"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/templates"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

const (
	Title     = "kuljis.xyz"
	Author    = "Connor Kuljis"
	DirBuild  = "dist"
	DirAssets = "assets"
)

var (
	md = markdown.NewMarkdown()
)

func main() {
	enableDrafts := flag.Bool("d", false, "enable drafts")

	flag.Parse()

	// Setup file logging
	os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile("logs/site-build.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// Set up multi-writer to log to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("Enable drafts:", *enableDrafts)

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	t := templates.NewTemplate()

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
	nerdStats := model.NewNerdStats(time.Now())

	return site.NewSite(
		Title,
		Author,
		DirBuild,
		DirAssets,
		enableDrafts,
		time.Now(),
		t,
		db,
		nerdStats,
		md,
	)
}
