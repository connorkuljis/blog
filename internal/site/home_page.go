package site

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/connorkuljis/blog/internal/dto"
)

type HomePage struct {
	Site        *MySite
	LatestEntry *dto.Entry
	Banner      [][]rune
}

func NewHomePage(site *MySite, latestEntry *dto.Entry) HomePage {
	// f := figure.NewFigure("connorkuljis", "isometric1", true)
	// f := figure.NewFigure("connorkuljis", "slant", true)
	// f := figure.NewFigure("connorkuljis", "univers", true)
	f := figure.NewFigure("connorkuljis", "larry3d", true)

	lines := f.Slicify()
	var banner [][]rune
	for _, line := range lines {
		banner = append(banner, []rune(line))
	}

	p := HomePage{Site: site, LatestEntry: latestEntry, Banner: banner}
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
