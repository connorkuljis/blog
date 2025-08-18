package site

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/pkg/site"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type MySite struct {
	EnableDrafts bool

	Config    Config
	CreatedAt time.Time
	// TRoot is the base template parsed once; we'll clone it per-page to build TSet entries.
	TRoot         *template.Template
	TSet          map[string]*template.Template
	Categories    []*model.Category
	CategoriesMap map[string]*model.Category
	Tags          []*model.Tag
	TagsMap       map[string]*model.Tag
	NerdStats     *model.NerdStats
}

type Config struct {
	Title     string `toml:"title"`
	Author    string `toml:"author"`
	Domain    string `toml:"domain"`
	DirBuild  string `toml:"dir_build"`
	DirAssets string `toml:"dir_assets"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func NewSite(
	config Config,
	enableDrafts bool,
	createdAt time.Time,
	t *template.Template,
	db *sqlx.DB,
	nerdStats *model.NerdStats,
	markdown goldmark.Markdown,
) *MySite {
	site := &MySite{
		Config:    config,
		CreatedAt: createdAt,
		TRoot:     t,
		TSet:      make(map[string]*template.Template),
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

	return pages
}

func (s *MySite) Render(pages []site.Page) error {
	// Ensure TSet is populated
	if len(s.TSet) == 0 {
		if err := s.buildTSet(); err != nil {
			return err
		}
	}

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

		// Lookup the prebuilt template for this page and execute it.
		tmplName := page.TemplateName()
		tmpl, ok := s.TSet[tmplName]
		if !ok {
			return fmt.Errorf("template not found in TSet: %s", tmplName)
		}
		// Execute the layout template inside this cloned template set so the layout can call the page view.
		if err := tmpl.ExecuteTemplate(f, "_layout.html", page); err != nil {
			return err
		}

		log.Println("created:", page.FileName())
	}

	return nil
}

func (s *MySite) buildTSet() error {
	// For each page template file, create a new template set that contains:
	// - the _layout.html as the top level template
	// - all component templates (head, header, footer)
	// - the single page view template (from templates/pages/<name>) parsed under its filename

	// Get list of page templates
	pages, err := filepath.Glob("templates/pages/*.html")
	if err != nil {
		return err
	}

	for _, p := range pages {
		name := filepath.Base(p)
		// Clone root template
		root, err := s.TRoot.Clone()
		if err != nil {
			return err
		}

		// Parse only the page file into the cloned template under its filename
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		// Ensure parsed as named template with name equal to filename
		if _, err := root.New(name).Parse(string(content)); err != nil {
			return err
		}

		// Store in TSet keyed by the filename used by pages (TemplateName())
		s.TSet[name] = root
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
