package site

import (
	"sort"

	"github.com/connorkuljis/blog/internal/model"
)

type HomePage struct {
	Site        *MySite
	Posts       []*model.Entry
	RecentPosts []*model.Entry
}

func NewHomePage(s *MySite) HomePage {
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

	p := HomePage{Site: s, Posts: posts, RecentPosts: recentPosts}
	return p
}

func (p HomePage) Title() string {
	return "Home"
}

func (p HomePage) FileName() string {
	return "index.html"
}

func (p HomePage) TemplateName() string {
	return "index.html"
}
