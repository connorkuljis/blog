package model

import (
	"github.com/connorkuljis/blog/internal/store"
)

type Author struct {
	*store.Author
}

func NewAuthor(a *store.Author) *Author {
	return &Author{
		Author: a,
	}
}