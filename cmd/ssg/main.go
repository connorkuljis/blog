package main

import (
	"fmt"
	"log"
	"time"

	"github.com/connorkuljis/content/internal/ssg"
	"github.com/connorkuljis/content/internal/store"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	title   = "Connor's Blog"
	publish = false
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
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	site := ssg.Site{
		Title:          title,
		DB:             db,
		MarkdownParser: md,
		Publish:        publish,
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
