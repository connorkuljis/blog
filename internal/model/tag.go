package model

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/store"
)

type Tag struct {
	*store.Tag

	Permalink string
	Entries   []*Entry
}

func NewTag(t *store.Tag) *Tag {
	return &Tag{
		Tag:       t,
		Permalink: "/" + filepath.Join("tags", t.Name),
	}
}
