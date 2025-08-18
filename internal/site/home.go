package site

import (
	"sort"

	"github.com/common-nighthawk/go-figure"
	"github.com/connorkuljis/blog/internal/model"
)

type HomePage struct {
	Site        *MySite
	Posts       []*model.Entry
	RecentPosts []*model.Entry
	Banner      [][]rune
}

func NewHomePage(s *MySite) HomePage {
	f := figure.NewFigure("connorkuljis", "larry3d", true)

	lines := f.Slicify()
	var banner [][]rune
	for _, line := range lines {
		banner = append(banner, []rune(line))
	}

	var allEntries []*model.Entry
	for _, c := range s.Categories {
		for _, e := range c.Entries {
			allEntries = append(allEntries, e)
		}
	}

	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].CreatedAt.After(allEntries[j].CreatedAt)
	})

	posts := allEntries

	recentPosts := posts
	if len(posts) > 5 {
		recentPosts = posts[:5]
	}

	p := HomePage{Site: s, Posts: posts, RecentPosts: recentPosts, Banner: banner}
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
