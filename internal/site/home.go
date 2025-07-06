package site

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/connorkuljis/blog/internal/model"
)

type HomePage struct {
	Site        *MySite
	LatestEntry *model.Entry
	Banner      [][]rune
}

func NewHomePage(s *MySite) HomePage {
	f := figure.NewFigure("connorkuljis", "larry3d", true)

	lines := f.Slicify()
	var banner [][]rune
	for _, line := range lines {
		banner = append(banner, []rune(line))
	}

	var latestEntry *model.Entry
	for _, c := range s.Categories {
		for _, e := range c.Entries {
			if latestEntry == nil || e.CreatedAt.After(latestEntry.CreatedAt) {
				latestEntry = e
			}
		}
	}

	p := HomePage{Site: s, LatestEntry: latestEntry, Banner: banner}
	return p
}

func (p HomePage) Title() string {
	return "Home"
}

func (p HomePage) FileName() string {
	return "index.html"
}

func (p HomePage) TemplateName() string {
	return "_index.html"
}
