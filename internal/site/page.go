package site

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var funcMap = template.FuncMap{
	"slugify": slugify,
}

type Page struct {
	Title        string
	Filepath     string
	ResourcePath string

	T    *template.Template
	Data map[string]any
}

type PageView struct {
	Name string

	BaseTemplate       string
	LayoutTemplate     string
	HeadTemplate       string
	ViewTemplate       string
	ComponentTemplates []string
}

func (p *PageView) Template() (*template.Template, error) {
	var sb strings.Builder

	sb.WriteString(p.BaseTemplate)
	sb.WriteString(p.HeadTemplate)
	sb.WriteString(p.LayoutTemplate)
	sb.WriteString(p.ViewTemplate)

	for _, component := range p.ComponentTemplates {
		sb.WriteString(component)
	}

	tpl, err := template.New(p.Name).Funcs(funcMap).Option("missingkey=error").Parse(sb.String())
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func (p *Page) Render(data any) error {
	dir, _ := filepath.Split(p.Filepath)
	os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(p.Filepath)
	if err != nil {
		return err
	}

	err = p.T.ExecuteTemplate(f, "base", data)
	if err != nil {
		return err
	}

	return nil
}

func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace non-alphanumeric characters with a hyphen
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
