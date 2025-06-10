package kuljis

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type EntryPage struct {
	Site     *MySite
	Category *model.Category
	Entry    *model.Entry
	Next     *model.Entry
	Prev     *model.Entry
}

func NewEntryPage(
	site *MySite,
	category *model.Category,
	entry *model.Entry,
	next *model.Entry,
	prev *model.Entry,
) EntryPage {
	return EntryPage{
		Site:     site,
		Category: category,
		Entry:    entry,
		Next:     next,
		Prev:     prev,
	}
}

func (p EntryPage) Title() string {
	return p.Entry.Title
}

func (p EntryPage) FileName() string {
	return filepath.Join(p.Entry.Permalink(), "index.html")
}

func (p EntryPage) TemplateName() string {
	return "_entry.html"
}
