package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/templates"
	"github.com/connorkuljis/blog/pkg/site"
	"github.com/yuin/goldmark"
)

type MySite struct {
	Config    Config
	CreatedAt time.Time

	Renderer *templates.Renderer
	Markdown goldmark.Markdown

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
	markdown goldmark.Markdown,
	nerdStats *model.NerdStats,
	categoryRepo *store.CategoryRepo,
	entryRepo *store.EntryRepo,
	tagRepo *store.TagRepo,
) (*MySite, error) {
	site := &MySite{
		Config:        config,
		CreatedAt:     createdAt,
		Renderer:      renderer,
		Markdown:      markdown,
		NerdStats:     nerdStats,
		CategoryRepo:  categoryRepo,
		EntryRepo:     entryRepo,
		TagRepo:       tagRepo,
		Categories:    make([]*model.Category, 0),
		Entries:       make([]*model.Entry, 0),
		RecentEntries: make([]*model.Entry, 0),
		Tags:          make([]*model.Tag, 0),
		TagsMap:       make(map[string]*model.Tag),
	}

	return site, nil
}

func (s *MySite) Init() error {
	if err := os.RemoveAll(s.Config.DirBuild); err != nil {
		return err
	}

	if err := os.MkdirAll(s.Config.DirBuild, os.ModePerm); err != nil {
		return err
	}

	staticAssets := os.DirFS(s.Config.DirAssets)
	if err := os.CopyFS(s.Config.DirBuild, staticAssets); err != nil {
		return err
	}

	tags, err := s.TagRepo.GetTags()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		mTag := model.NewTag(tag)
		s.Tags = append(s.Tags, mTag)
		s.TagsMap[tag.Name] = mTag
	}

	categories, err := s.CategoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}

	for _, category := range categories {
		mCategory := model.NewCategory(category)
		s.Categories = append(s.Categories, mCategory)

		entries, err := s.EntryRepo.ReadAllByCategoryID(category.ID, s.Config.EnableDrafts)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			mEntry := model.NewEntry(entry, mCategory, s.Markdown)
			s.Entries = append(s.Entries, mEntry)
			mCategory.AddEntry(mEntry)

			tags, err := s.TagRepo.GetTagsForEntry(mEntry.ID)
			if err != nil {
				return err
			}

			for _, t := range tags {
				tag := model.NewTag(t)
				mEntry.Tags = append(mEntry.Tags, tag)
				s.TagsMap[t.Name].Entries = append(s.TagsMap[t.Name].Entries, mEntry)
			}
		}
	}

	sort.Slice(s.Entries, func(i, j int) bool {
		return s.Entries[i].CreatedAt.After(s.Entries[j].CreatedAt)
	})

	s.RecentEntries = s.Entries
	if len(s.RecentEntries) > 5 {
		s.RecentEntries = s.Entries[:5]
	}

	return nil
}

func (s *MySite) Build() []site.Page {
	var pages = []site.Page{}

	pages = append(pages, NewHomePage(s))
	pages = append(pages, NewAboutPage(s))

	for _, category := range s.Categories {
		pages = append(pages, NewCategoryPage(s, category))

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

	pages = append(pages, NewTagsPage(s, s.Tags))

	for _, tag := range s.Tags {
		pages = append(pages, NewTagPage(s, tag))
	}

	s.NerdStats.PageCount = len(pages)

	return pages
}

func (s *MySite) Render(pages []site.Page) error {
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

		tmpl, ok := s.Renderer.TSet[page.TemplateName()]
		if !ok {
			return fmt.Errorf("template not found in TSet: %s", page.TemplateName())
		}
		if err := tmpl.ExecuteTemplate(f, "root", page); err != nil {
			return err
		}
	}
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

	log.Println("created: sitemap.xml")
	return nil
}
