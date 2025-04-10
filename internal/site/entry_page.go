package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type EntryPage struct {
	Site     *MySite
	Category model.Category
	Entry    model.Entry
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", p.Entry.Slug(), "index.html")
}

func (p EntryPage) TemplateName() string {
	return "_entry.html"
}
