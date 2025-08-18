package templates

import (
	"html/template"
	"os"

	"github.com/connorkuljis/blog/internal/util"
)

var (
	funcMap = template.FuncMap{
		"runeToString": func(r rune) string {
			return string(r)
		},
		"slugify":  util.Slugify,
		"truncate": util.Truncate,
		"add":      func(a, b int) int { return a + b },
	}
)

func NewTemplate() *template.Template {
	tmpl := template.New("").Funcs(funcMap).Option("missingkey=error")

	// Parse base layout first
	template.Must(tmpl.ParseGlob("templates/_layout.html"))

	// Parse all component files (they use internal {{ define }} names).
	template.Must(tmpl.ParseGlob("templates/components/*.html"))

	// Parse other base templates (standalone filenames) if any
	// We'll also register specific base fragments under short names expected by the layout.
	template.Must(tmpl.ParseGlob("templates/*.html"))

	// Ensure 'head' and 'header' are available as named templates by registering their file
	// contents under those names if they aren't already defined.
	if tmpl.Lookup("head") == nil {
		if b, err := os.ReadFile("templates/head.html"); err == nil {
			template.Must(tmpl.New("head").Parse(string(b)))
		}
	}
	if tmpl.Lookup("header") == nil {
		if b, err := os.ReadFile("templates/header.html"); err == nil {
			template.Must(tmpl.New("header").Parse(string(b)))
		}
	}

	// Note: Do NOT parse pages here; pages will be parsed per-page into clones of TRoot.

	return tmpl
}

// RenderTemplate is unused in current flow. Keep for future use or remove.
func RenderTemplate(tmpl *template.Template, pageName string, layoutName ...string) (*template.Template, error) {
	// Intentionally return the cloned template so callers may execute the layout name.
	pageTmpl := template.Must(tmpl.Clone())
	return pageTmpl, nil
}
