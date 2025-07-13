package templates

import (
	"html/template"

	"github.com/connorkuljis/blog/internal/util"
)

var (
	funcMap = template.FuncMap{
		"runeToString": func(r rune) string {
			return string(r)
		},
		"slugify": util.Slugify,
	}
)

func NewTemplate() *template.Template {
	return template.Must(template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html"))
}
