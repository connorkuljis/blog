package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/dto"
)

type EntryPage struct {
	Site     *MySite
	Category *dto.Category
	Entry    *dto.Entry
	Next     *dto.Entry
	Prev     *dto.Entry
}

func NewEntryPage(
	site *MySite,
	category *dto.Category,
	entry *dto.Entry,
	next *dto.Entry,
	prev *dto.Entry,
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
