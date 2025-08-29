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

// ToStoreTag maps this model.Tag to a store.Tag.
func (t *Tag) ToStoreTag() *store.Tag {
	if t == nil || t.Tag == nil {
		return nil
	}
	return &store.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}
