package main

import (
	"fmt"
	"log"
	"time"

	"github.com/connorkuljis/content/internal/site"
	"github.com/connorkuljis/content/internal/store"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const title = "Connor's Blog"

func main() {
	start := time.Now()
	err := generateSite()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("took: %d ms", time.Since(start).Milliseconds())
}

func generateSite() error {
	db, err := store.Connect()
	if err != nil {
		log.Fatal(err)
	}

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

	s := site.Site{
		Title:          title,
		DB:             db,
		MarkdownParser: md,
	}

	err = s.Init()
	if err != nil {
		return err
	}

	err = s.Render()
	if err != nil {
		return err
	}

	return nil
}
