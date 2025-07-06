package site

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/connorkuljis/blog/pkg/site"
	"github.com/jmoiron/sqlx"
	"github.com/yuin/goldmark"
)

type MySite struct {
	Title        string
	Author       string
	DirBuild     string
	DirAssets    string
	EnableDrafts bool

	CreatedAt     time.Time
	T             *template.Template
	Categories    []*model.Category
	CategoriesMap map[string]*model.Category
	Tags          []*model.Tag
	TagsMap       map[string]*model.Tag
	NerdStats     *model.NerdStats
}

func NewSite(
	title string,
	author string,
	dirBuild string,
	dirAssets string,
	enableDrafts bool,
	createdAt time.Time,
	t *template.Template,
	db *sqlx.DB,
	nerdStats *model.NerdStats,
	markdown goldmark.Markdown,
) *MySite {
	site := &MySite{
		Title:     title,
		Author:    author,
		DirBuild:  dirBuild,
		DirAssets: dirAssets,
		CreatedAt: createdAt,
		T:         t,
		NerdStats: nerdStats,
	}

	site.TagsMap = make(map[string]*model.Tag)
	tags, err := store.NewTagRepo(db).GetTags()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range tags {
		tag := model.NewTag(t)
		site.Tags = append(site.Tags, tag)
		site.TagsMap[t.Name] = tag
	}

	categories, err := store.NewCategoryRepo(db).ReadAllCategories()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range categories {
		mCategory := model.NewCategory(c)
		site.Categories = append(site.Categories, mCategory)

		entries, err := store.NewEntryRepo(db).ReadAllByCategoryID(c.ID, enableDrafts)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			mEntry := model.NewEntry(e, mCategory, markdown)

			tags, err := store.NewTagRepo(db).GetTagsForEntry(e.ID)
			if err != nil {
				log.Fatal(err)
			}

			for _, t := range tags {
				tag := model.NewTag(t)
				mEntry.Tags = append(mEntry.Tags, tag)
				site.TagsMap[t.Name].Entries = append(site.TagsMap[t.Name].Entries, mEntry)
			}

			mCategory.AddEntry(mEntry)
		}
	}

	return site
}

func (s *MySite) Init() error {
	err := os.RemoveAll(s.DirBuild)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.DirBuild, os.ModePerm)
	if err != nil {
		return err
	}

	staticAssets := os.DirFS(s.DirAssets)
	err = os.CopyFS(s.DirBuild, staticAssets)
	if err != nil {
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
	for _, page := range pages {
		filename := filepath.Join(s.DirBuild, page.FileName())

		dir := filepath.Dir(filename)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		err = s.T.ExecuteTemplate(f, page.TemplateName(), page) // note: assume each page is rendered as concrete type when accessing data in templates.
		if err != nil {
			return err
		}
		log.Println("created:", page.FileName())
	}

	return nil
}
