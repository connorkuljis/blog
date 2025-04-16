package kuljis

import "github.com/connorkuljis/blog/internal/model"

type HomePage struct {
	Site        *MySite
	LatestEntry *model.Entry
}

func NewHomePage(site *MySite) HomePage {
	p := HomePage{Site: site}
	p.SetLatestEntry()
	return p
}

func (p HomePage) Filepath() string {
	return "public/index.html"
}

func (p HomePage) TemplateName() string {
	return "_index.html"
}

func (p HomePage) SetLatestEntry() {
	if len(p.Site.Categories) < 0 || len(p.Site.Categories[0].Entries) < 0 {
		return
	}

	p.LatestEntry = p.Site.Categories[0].Entries[0]
	for _, c := range p.Site.Categories {
		for _, e := range c.Entries {
			if e.CreatedAt.After(p.LatestEntry.CreatedAt) {
				p.LatestEntry = e
			}
		}
	}
}
