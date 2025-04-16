package kuljis

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type EntryPage struct {
	Site     *MySite
	Category *model.Category
	Entry    *model.Entry
}

func NewEntryPage(site *MySite, category *model.Category, entry *model.Entry) EntryPage {
	return EntryPage{
		Site:     site,
		Category: category,
		Entry:    entry,
	}
}

func (p EntryPage) Filepath() string {
	return filepath.Join("public", p.Entry.Permalink(), "index.html")
}

func (p EntryPage) TemplateName() string {
	return "_entry.html"
}
