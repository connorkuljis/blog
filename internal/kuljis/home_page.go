package kuljis

import "github.com/connorkuljis/blog/internal/model"

type HomePage struct {
	Site        *MySite
	LatestEntry *model.Entry
}

func NewHomePage(site *MySite, latestEntry *model.Entry) HomePage {
	p := HomePage{Site: site, LatestEntry: latestEntry}
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
