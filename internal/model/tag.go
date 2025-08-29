package model

import (
	"path/filepath"

	"github.com/connorkuljis/blog/internal/store"
)

type Tag struct {
	ID      int64
	Name    string
	Entries []*Entry
}

func NewTag(t *store.Tag) *Tag {
	return &Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (t *Tag) Permalink() string {
	return "/" + filepath.Join("tags", t.Name)
}

func (t *Tag) ToStoreTag() *store.Tag {
	return &store.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}
