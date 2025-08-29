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
	md               = markdown.NewMarkdown()
	enableDraftsFlag = flag.Bool("d", false, "enable drafts")
	configFlag       = flag.String("config", "config.jsonc", "path to config file")
)

func main() {
	flag.Parse()

	start := time.Now()

	var filename string = *configFlag
	cfg, err := site.LoadConfig(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded config:", filename)

	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

	nerdStats := model.NewNerdStats(time.Now())

	renderer := templates.NewRenderer()

	mySite := site.NewSite(
		*cfg,
		*enableDraftsFlag,
		time.Now(),
		db,
		nerdStats,
		md,
		renderer,
	)

	err = mySite.Init()
	if err != nil {
		log.Fatal(err)
	}

	pages := mySite.Build()

	err = mySite.Render(pages)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(time.Since(start))
}
