package model

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type Entry struct {
	ID         int64     `db:"id"`
	CategoryID int64     `db:"category_id" yaml:"category_id"`
	Title      string    `db:"title" yaml:"title"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at" yaml:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" yaml:"updated_at"`
	Publish    int       `db:"publish" yaml:"publish"`

	CategoryTitle       string `db:"category_title"`
	CategoryDescription string `db:"category_description"`
}

func NewEntry(categoryID int64, title string) *Entry {
	return &Entry{CategoryID: categoryID, Title: title}
}

func (e Entry) String() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("category_id: %d\n", e.CategoryID))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", e.CreatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", e.UpdatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("publish: %d\n", e.Publish))
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("%s", e.Content))

	return sb.String()
}

func (e *Entry) LoadFromContentString(r io.Reader) error {
	b, err := frontmatter.Parse(r, e)
	if err != nil {
		return err
	}

	e.Content = string(b)

	return nil
}
