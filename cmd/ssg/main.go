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

const (
	title = "kuljis.xyz"
)

func main() {
	start := time.Now()

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
			html.WithUnsafe(),
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	site := site.Site{
		Title:          title,
		DB:             db,
		MarkdownParser: md,
	}

	err = site.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = site.Render()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("built site in: %d ms\n", time.Since(start).Milliseconds())
}
