package kuljis

import "github.com/connorkuljis/blog/internal/model"

type HomePage struct {
	Title string

	Site        *MySite
	LatestEntry *model.Entry
}

func NewHomePage(site *MySite, title string) HomePage {
	p := HomePage{Site: site, Title: title}

	var latestEntry *model.Entry
	for _, c := range p.Site.Categories {
		for _, e := range c.Entries {
			if latestEntry == nil || e.CreatedAt.After(latestEntry.CreatedAt) {
				latestEntry = e
			}
		}
	}

	p.LatestEntry = latestEntry

	return p
}

func (p HomePage) FileName() string {
	return "index.html"
}

func (p HomePage) TemplateName() string {
	return "_index.html"
}
