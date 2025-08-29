package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/connorkuljis/blog/internal/markdown"
	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/templates"
	"github.com/connorkuljis/blog/pkg/site"
)

type MySite struct {
	Config    Config
	CreatedAt time.Time

	TemplateRenderer *templates.Renderer
	MarkdownRenderer *markdown.Renderer

	NerdStats *model.NerdStats

	CategoryRepo *store.CategoryRepo
	EntryRepo    *store.EntryRepo
	TagRepo      *store.TagRepo

	Categories    []*model.Category
	Entries       []*model.Entry
	RecentEntries []*model.Entry
	Tags          []*model.Tag
	TagsMap       map[string]*model.Tag
}

func NewSite(
	config Config,
	createdAt time.Time,
	renderer *templates.Renderer,
	markdown *markdown.Renderer,
	nerdStats *model.NerdStats,
	categoryRepo *store.CategoryRepo,
	entryRepo *store.EntryRepo,
	tagRepo *store.TagRepo,
) (*MySite, error) {
	site := &MySite{
		Config:           config,
		CreatedAt:        createdAt,
		TemplateRenderer: renderer,
		MarkdownRenderer: markdown,
		NerdStats:        nerdStats,
		CategoryRepo:     categoryRepo,
		EntryRepo:        entryRepo,
		TagRepo:          tagRepo,
		Categories:       make([]*model.Category, 0),
		Entries:          make([]*model.Entry, 0),
		RecentEntries:    make([]*model.Entry, 0),
		Tags:             make([]*model.Tag, 0),
		TagsMap:          make(map[string]*model.Tag),
	}

	return site, nil
}

func (s *MySite) Init() error {
	log.Println("[INIT]", s.Config.DirBuild, "removing")
	if err := os.RemoveAll(s.Config.DirBuild); err != nil {
		return err
	}

	log.Println("[INIT]", s.Config.DirBuild, "creating")
	if err := os.MkdirAll(s.Config.DirBuild, os.ModePerm); err != nil {
		return err
	}

	log.Println("[INIT]", s.Config.DirAssets, "->", s.Config.DirBuild)
	staticAssets := os.DirFS(s.Config.DirAssets)
	if err := os.CopyFS(s.Config.DirBuild, staticAssets); err != nil {
		return err
	}

	log.Println("[INIT]", "loading tags")
	tags, err := s.TagRepo.GetTags()
	if err != nil {
		return err
	}

	for _, t := range tags {
		mTag := model.NewTag(t)
		s.Tags = append(s.Tags, mTag)
		s.TagsMap[t.Name] = mTag
	}

	log.Println("[INIT]", "loading categories")
	categories, err := s.CategoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}
	for _, c := range categories {
		mCategory := model.NewCategory(c)
		s.Categories = append(s.Categories, mCategory)

		log.Println("[INIT]", "loading entries for", mCategory.Title)
		entries, err := s.EntryRepo.ReadAllByCategoryID(c.ID, s.Config.EnableDrafts)
		if err != nil {
			return err
		}

		for _, e := range entries {
			tags, err := s.TagRepo.GetTagsForEntry(e.ID)
			if err != nil {
				return err
			}
			var mTags []*model.Tag
			for _, t := range tags {
				mTags = append(mTags, model.NewTag(t))
			}

			mEntry := model.NewEntry(e, mCategory, mTags)
			mEntry.Markdown = s.MarkdownRenderer.RenderHTML(mEntry.Content)
			mCategory.AddEntry(mEntry)

			s.Entries = append(s.Entries, mEntry)
			for _, t := range mEntry.Tags {
				s.TagsMap[t.Name].Entries = append(s.TagsMap[t.Name].Entries, mEntry)
			}
		}
	}

	log.Println("[INIT]", "sorting entries by date")
	sort.Slice(s.Entries, func(i, j int) bool {
		return s.Entries[i].CreatedAt.After(s.Entries[j].CreatedAt)
	})

	log.Println("[INIT]", "collecting recent entries")
	s.RecentEntries = s.Entries
	if len(s.RecentEntries) > 5 {
		s.RecentEntries = s.Entries[:5]
	}

	return nil
}

func (s *MySite) Build() []site.Page {
	var pages = []site.Page{}

	log.Println("[BUILD]", "constructing Home page")
	pages = append(pages, NewHomePage(s))
	log.Println("[BUILD]", "constructing About page")
	pages = append(pages, NewAboutPage(s))

	for _, category := range s.Categories {
		log.Println("[BUILD]", "constructing", category.Title, "page")
		pages = append(pages, NewCategoryPage(s, category))

		log.Printf("[BUILD] constructing (%d) Entry pages for %s", len(category.Entries), category.Title)
		for i, entry := range category.Entries {
			var next *model.Entry
			var prev *model.Entry

			// Get previous entry if exists
			if i > 0 {
				prev = category.Entries[i-1]
			}

			// Get next entry if exists
			if i < len(category.Entries)-1 {
				next = category.Entries[i+1]
			}

			pages = append(pages, NewEntryPage(s, category, entry, next, prev))
		}
	}

	log.Println("[BUILD]", "constructing Tags page")
	pages = append(pages, NewTagsPage(s, s.Tags))

	for _, tag := range s.Tags {
		log.Println("[BUILD]", "constructing Tag page for", tag.Name)
		pages = append(pages, NewTagPage(s, tag))
	}

	log.Println("[BUILD]", "constructing Archive page")
	pages = append(pages, NewArchivePage(s))

	s.NerdStats.PageCount = len(pages)

	return pages
}

func (s *MySite) Render(pages []site.Page) error {
	log.Printf("[RENDER] rendering (%d) pages to %s", len(pages), s.Config.DirBuild)
	for _, page := range pages {
		filename := filepath.Join(s.Config.DirBuild, page.FileName())
		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		tmpl, ok := s.TemplateRenderer.TSet[page.TemplateName()]
		if !ok {
			return fmt.Errorf("template not found in TSet: %s", page.TemplateName())
		}
		if err := tmpl.ExecuteTemplate(f, "root", page); err != nil {
			return err
		}
	}
	log.Println("[RENDER] OK")
	return nil
}

// generateSiteMap creates a sitemap.xml file in the build directory
func (s *MySite) GenerateSiteMap(pages []site.Page) error {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	for _, page := range pages {
		url := strings.TrimRight(s.Config.Domain, "/") + "/" + strings.TrimLeft(page.FileName(), "/")
		url = strings.ReplaceAll(url, "index.html", "")
		b.WriteString(fmt.Sprintf("  <url>\n    <loc>%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>weekly</changefreq>\n    <priority>0.8</priority>\n  </url>\n", url, s.CreatedAt.Format("2006-01-02")))
	}

	b.WriteString("</urlset>\n")

	filename := filepath.Join(s.Config.DirBuild, "sitemap.xml")
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return fmt.Errorf("failed to write sitemap.xml: %w", err)
	}

	log.Println("[SITEMAP] generated", filename)
	return nil
}
