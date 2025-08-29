package site

import (
	"path/filepath"
)

type ArchivePage struct {
	Site *MySite
}

func NewArchivePage(s *MySite) *ArchivePage {
	return &ArchivePage{Site: s}
}

func (p *ArchivePage) Title() string {
	return "Archive"
}

func (p *ArchivePage) FileName() string {
	return "/" + filepath.Join("archive", "index.html")
}

func (p *ArchivePage) TemplateName() string {
	return "archive.html"
}
