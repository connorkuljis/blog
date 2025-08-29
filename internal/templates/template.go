package templates

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/blog/internal/util"
)

var (
	funcMap = template.FuncMap{
		"runeToString": func(r rune) string { return string(r) },
		"slugify":      util.Slugify,
		"truncate":     util.Truncate,
		"add":          func(a, b int) int { return a + b },
	}
)

// Renderer encapsulates all template rendering concerns for the site package.
// It owns the root template and a prebuilt map of page template sets.
// The site package should only call BuildSet once and then Render() pages through this type.
type Renderer struct {
	TRoot *template.Template
	TSet  map[string]*template.Template
}

// NewRenderer parses the shared layout and components and returns a Renderer.
func NewRenderer() *Renderer {
	root := template.New("").Funcs(funcMap).Option("missingkey=error")
	// Parse base layout first
	template.Must(root.ParseGlob("templates/_layout.html"))
	// Parse all component files (they use internal {{ define }} names).
	template.Must(root.ParseGlob("templates/components/*.html"))
	// Parse other base templates (standalone filenames) if any. Do not parse pages here.
	template.Must(root.ParseGlob("templates/*.html"))

	pages, err := filepath.Glob("templates/pages/*.html")
	if err != nil {
		panic(err)
	}

	tSet := make(map[string]*template.Template)
	for _, page := range pages {
		name := filepath.Base(page)
		root, err := root.Clone()
		if err != nil {
			panic(err)
		}
		content, err := os.ReadFile(page)
		if err != nil {
			panic(err)
		}
		if _, err := root.New(name).Parse(string(content)); err != nil {
			panic(err)
		}
		tSet[name] = root
	}

	return &Renderer{TRoot: root, TSet: tSet}
}
