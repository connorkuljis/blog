package model

import "html/template"

// Site data is available to every page.
type Site interface {
	Init() error
	Build() []Page
	Render([]Page) error
	RootDir() string
	Template() *template.Template
}
