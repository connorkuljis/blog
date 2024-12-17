package site

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

type Page struct {
	Title        string
	Filepath     string
	ResourcePath string

	BaseTemplate       string
	LayoutTemplate     string
	HeadTemplate       string
	ViewTemplate       string
	ComponentTemplates []string

	Data map[string]any

	Template *template.Template
	Writer   io.Writer
}

func (p *Page) ParseTemplate() (*template.Template, error) {
	var templateStrings []string
	templateStrings = append(templateStrings, p.BaseTemplate)
	templateStrings = append(templateStrings, p.HeadTemplate)
	templateStrings = append(templateStrings, p.LayoutTemplate)
	templateStrings = append(templateStrings, p.ViewTemplate)
	templateStrings = append(templateStrings, p.ComponentTemplates...)

	templateString := ""
	for _, str := range templateStrings {
		templateString += str
	}

	tpl, err := template.New(p.Title).Option("missingkey=error").Parse(templateString)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func (p *Page) Render() error {
	dir, _ := filepath.Split(p.Filepath)
	os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(p.Filepath)
	if err != nil {
		return err
	}

	tpl, err := p.ParseTemplate()
	if err != nil {
		fmt.Errorf("template error")
	}

	err = tpl.ExecuteTemplate(f, "base", p.Data)
	if err != nil {
		return err
	}

	return nil
}
