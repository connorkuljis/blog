package model

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type Entry struct {
	Id        int64     `db:"id"`
	Category  string    `db:"category"`
	Title     string    `db:"title" yaml:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `db:"updated_at" yaml:"updated_at"`
	Publish   int       `db:"publish" yaml:"publish"`
}

func NewEntry(category string, title string) *Entry {
	return &Entry{Title: title, Category: category}
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", e.CreatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", e.UpdatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("publish: %d\n", e.Publish))
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("%s", e.Content))

	return sb.String()
}

func ParseEntryFromFrontMatter(r io.Reader, entry *Entry) error {
	b, err := frontmatter.Parse(r, entry)
	if err != nil {
		return err
	}

	entry.Content = string(b)

	return nil
}
