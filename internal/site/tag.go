package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type TagPage struct {
	Site *MySite
	Tag  *model.Tag
}

func NewTagPage(s *MySite, t *model.Tag) *TagPage {
	return &TagPage{
		Site: s,
		Tag:  t,
	}
}

func (p *TagPage) Title() string {
	return p.Tag.Name
}

func (p *TagPage) FileName() string {
	return filepath.Join(p.Tag.Permalink, "index.html")
}

func (p *TagPage) TemplateName() string {
	return "_tag.html"
}
