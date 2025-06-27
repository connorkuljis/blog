package site

import (
	"html/template"
	"os"
	"testing"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/pkg/site"
)

func TestNewSite(t *testing.T) {
	title := "Test Site"
	author := "Test Author"
	dirBuild := "build"
	dirAssets := "assets"
	createdAt := time.Now()
	tmpl := template.New("test")
	categories := []*model.Category{}
	categoriesMap := make(map[string]*model.Category)
	nerdStats := &model.NerdStats{}

	site := NewSite(title, author, dirBuild, dirAssets, createdAt, tmpl, categories, categoriesMap, nerdStats)

	if site.Title != title {
		t.Errorf("Expected title %s, got %s", title, site.Title)
	}
	if site.Author != author {
		t.Errorf("Expected author %s, got %s", author, site.Author)
	}
	if site.DirBuild != dirBuild {
		t.Errorf("Expected build directory %s, got %s", dirBuild, site.DirBuild)
	}
	if site.DirAssets != dirAssets {
		t.Errorf("Expected assets directory %s, got %s", dirAssets, site.DirAssets)
	}
	if !site.CreatedAt.Equal(createdAt) {
		t.Errorf("Expected created at %v, got %v", createdAt, site.CreatedAt)
	}
	if site.T != tmpl {
		t.Errorf("Expected template %v, got %v", tmpl, site.T)
	}
	if len(site.Categories) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(site.Categories))
	}
	if len(site.CategoriesMap) != 0 {
		t.Errorf("Expected 0 categories in map, got %d", len(site.CategoriesMap))
	}
	if site.NerdStats != nerdStats {
		t.Errorf("Expected nerd stats %v, got %v", nerdStats, site.NerdStats)
	}
}

func TestSite_Init(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dirBuild := tmpDir + "/build"
	dirAssets := tmpDir + "/assets"

	err = os.Mkdir(dirAssets, 0755)
	if err != nil {
		t.Fatal(err)
	}

	assetFile := dirAssets + "/style.css"
	err = os.WriteFile(assetFile, []byte("body {}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	s := &MySite{
		DirBuild:  dirBuild,
		DirAssets: dirAssets,
	}

	err = s.Init()
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if _, err := os.Stat(dirBuild); os.IsNotExist(err) {
		t.Errorf("Build directory was not created")
	}

	if _, err := os.Stat(dirBuild + "/style.css"); os.IsNotExist(err) {
		t.Errorf("Asset file was not copied")
	}
}

func TestSite_Build(t *testing.T) {
	s := &MySite{
		Categories: []*model.Category{
			{
				Title: "Category 1",
				Entries: []*model.Entry{
					{Title: "Entry 1", CreatedAt: time.Now()},
					{Title: "Entry 2", CreatedAt: time.Now().Add(time.Hour)},
				},
			},
		},
		NerdStats: &model.NerdStats{},
	}

	pages := s.Build()

	if len(pages) != 5 {
		t.Errorf("Expected 5 pages, got %d", len(pages))
	}

	if s.NerdStats.PageCount != 5 {
		t.Errorf("Expected page count to be 5, got %d", s.NerdStats.PageCount)
	}
}

func TestSite_Render(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dirBuild := tmpDir + "/build"

	tmpl, err := template.New("page").Parse("{{.Title}}")
	if err != nil {
		t.Fatal(err)
	}

	s := &MySite{
		DirBuild: dirBuild,
		T:        tmpl,
	}

	pages := []site.Page{
		&testPage{title: "Page 1", filename: "page1.html"},
		&testPage{title: "Page 2", filename: "page2.html"},
	}

	err = s.Render(pages)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	for _, p := range pages {
		content, err := os.ReadFile(dirBuild + "/" + p.FileName())
		if err != nil {
			t.Errorf("Could not read rendered file: %v", err)
		}

		if string(content) != p.Title() {
			t.Errorf("Expected content %q, got %q", p.Title(), string(content))
		}
	}
}

type testPage struct {
	title    string
	filename string
}

func (p *testPage) Title() string {
	return p.title
}

func (p *testPage) FileName() string {
	return p.filename
}

func (p *testPage) TemplateName() string {
	return "page"
}
