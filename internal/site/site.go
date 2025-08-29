package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/internal/templates"
	"github.com/connorkuljis/blog/pkg/site"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type MySite struct {
	EnableDrafts bool

	Config        Config
	CreatedAt     time.Time
	Renderer      *templates.Renderer
	Categories    []*model.Category
	CategoriesMap map[string]*model.Category
	Tags          []*model.Tag
	TagsMap       map[string]*model.Tag
	NerdStats     *model.NerdStats
}

type Config struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Domain    string `json:"domain"`
	DirBuild  string `json:"dir_build"`
	DirAssets string `json:"dir_assets"`
	EmailList string `json:"email_list"`
}

func NewSite(
	config Config,
	enableDrafts bool,
	createdAt time.Time,
	db *sqlx.DB,
	nerdStats *model.NerdStats,
	markdown goldmark.Markdown,
	renderer *templates.Renderer,
) *MySite {
	site := &MySite{
		Config:    config,
		CreatedAt: createdAt,
		Renderer:  renderer,
		NerdStats: nerdStats,
	}

	categoryRepo := store.NewCategoryRepo(db)
	entryRepo := store.NewEntryRepo(db)
	tagRepo := store.NewTagRepo(db)

	if err := site.loadTags(tagRepo); err != nil {
		log.Fatal(err)
	}

	if err := site.loadCategories(categoryRepo, entryRepo, tagRepo, enableDrafts, markdown); err != nil {
		log.Fatal(err)
	}

	return site
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

	// Generate sitemap.xml after building pages
	if err := s.generateSiteMap(pages); err != nil {
		log.Printf("failed to generate sitemap: %v", err)
	}

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

func (s *MySite) loadTags(tagRepo *store.TagRepo) error {
	s.TagsMap = make(map[string]*model.Tag)
	tags, err := tagRepo.GetTags()
	if err != nil {
		return err
	}

	for _, t := range tags {
		tag := model.NewTag(t)
		s.Tags = append(s.Tags, tag)
		s.TagsMap[t.Name] = tag
	}
	return nil
}

// generateSiteMap creates a sitemap.xml file in the build directory
func (s *MySite) generateSiteMap(pages []site.Page) error {
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

func (s *MySite) loadCategories(
	categoryRepo *store.CategoryRepo,
	entryRepo *store.EntryRepo,
	tagRepo *store.TagRepo,
	enableDrafts bool,
	markdown goldmark.Markdown,
) error {
	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}

	for _, c := range categories {
		mCategory := model.NewCategory(c)
		s.Categories = append(s.Categories, mCategory)

		entries, err := entryRepo.ReadAllByCategoryID(c.ID, enableDrafts)
		if err != nil {
			return err
		}

		for _, e := range entries {
			mEntry := model.NewEntry(e, mCategory, markdown)

			tags, err := tagRepo.GetTagsForEntry(e.ID)
			if err != nil {
				return err
			}

			for _, t := range tags {
				tag := model.NewTag(t)
				mEntry.Tags = append(mEntry.Tags, tag)
				s.TagsMap[t.Name].Entries = append(s.TagsMap[t.Name].Entries, mEntry)
			}

			mCategory.AddEntry(mEntry)
		}
	}
	return nil
}
