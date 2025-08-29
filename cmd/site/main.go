package main

import (
	"flag"
	"log"
	"time"

	"github.com/connorkuljis/blog/internal/markdown"
	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/site"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/templates"
)

var (
	flagEnableDrafts   = flag.Bool("d", false, "enable drafts")
	flagConfigFilename = flag.String("config", "config.json", "path to config file")
)

func main() {
	flag.Parse()
	start := time.Now()

	cfg, err := site.LoadConfig(*flagConfigFilename)
	if err != nil {
		log.Fatal(err)
	}
	if *flagEnableDrafts {
		cfg.EnableDrafts = *flagEnableDrafts
	}

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}
	categoryRepo := store.NewCategoryRepo(db)
	entryRepo := store.NewEntryRepo(db)
	tagRepo := store.NewTagRepo(db)

	nerdStats := model.NewNerdStats(time.Now())
	markdown := markdown.NewMarkdown()
	renderer := templates.NewRenderer()

	mySite, err := site.NewSite(
		*cfg,
		start,
		renderer,
		markdown,
		nerdStats,
		categoryRepo,
		entryRepo,
		tagRepo,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = mySite.Init()
	if err != nil {
		log.Fatal(err)
	}
	pages := mySite.Build()
	err = mySite.GenerateSiteMap(pages)
	if err != nil {
		log.Fatal(err)
	}
	err = mySite.Render(pages)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(time.Since(start))
}
