package site

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/model"
)

type TagsPage struct {
	Site *MySite
	Tags []*model.Tag
}

func NewTagsPage(s *MySite, t []*model.Tag) *TagsPage {
	return &TagsPage{
		Site: s,
		Tags: t,
	}
}

func (p *TagsPage) Title() string {
	return "Tags"
}

func (p *TagsPage) FileName() string {
	return "/" + filepath.Join("tags", "index.html")
}

func (p *TagsPage) TemplateName() string {
	return "tags.html"
}
