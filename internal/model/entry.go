package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"html/template"

	"github.com/connorkuljis/blog/internal/util"
	"github.com/yuin/goldmark"
)

type Entry struct {
	ID               int64          `db:"id"`
	CategoryID       int64          `db:"category_id"`
	Title            string         `db:"title"`
	Content          string         `db:"content"`
	Description      sql.NullString `db:"description"`
	FeaturedImageURL sql.NullString `db:"featured_image_url"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	IsDraft          int            `db:"is_draft"`

	// below fields are computed values, that may or may not be populated.
	Category  Category
	Markdown  template.HTML
	WordCount int
}

func NewEntry(categoryID int64, title string) *Entry {
	now := time.Now()
	return &Entry{
		CategoryID: categoryID,
		Title:      title,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsDraft:    util.BoolToInt(true),
	}
}

func (e Entry) String() string {
	var sb strings.Builder

	// frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", e.Title))
	sb.WriteString(fmt.Sprintf("category_id: %d\n", e.CategoryID))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", e.CreatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", e.UpdatedAt.UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("is_draft: %d\n", e.IsDraft))
	sb.WriteString("---\n")
	// ./frontmatter

	sb.WriteString(fmt.Sprintf("%s", e.Content))

	return sb.String()
}

func (e *Entry) AddCategory(category Category) {
	e.Category = category
}

func (e *Entry) ToHTML(parser goldmark.Markdown) error {
	var buf bytes.Buffer
	err := parser.Convert([]byte(e.Content), &buf)
	if err != nil {
		return err
	}

	e.Markdown = template.HTML(buf.String())

	e.CalculateWordCount()

	return nil
}

func (e Entry) Slug() string {
	return e.Category.Slug() + "/" + util.Slugify(e.Title)
}

func (e *Entry) CalculateWordCount() {
	e.WordCount = len(strings.Split(e.Content, " "))
}

func GroupByYear(entries []Entry) [][]Entry {
	// Handle empty or nil input slice gracefully
	if len(entries) == 0 {
		return [][]Entry{} // Return an empty slice, not nil
	}

	// 1. Use a map to group entries by year.
	//    Key: Year (int)
	//    Value: Slice of entries for that year ([]Entry)
	groups := make(map[int][]Entry)

	for _, entry := range entries {
		year := entry.CreatedAt.Year()
		// Append the entry to the slice associated with its year.
		// If the key (year) doesn't exist yet, it will be created with a new slice.
		groups[year] = append(groups[year], entry)
	}

	// 2. Get the years (keys) from the map to sort them.
	years := make([]int, 0, len(groups))
	for year := range groups {
		years = append(years, year)
	}

	// 3. Sort the years chronologically.
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	// 4. Build the final result slice, ordered by the sorted years.
	//    Pre-allocate capacity for efficiency.
	result := make([][]Entry, 0, len(years))
	for _, year := range years {
		// Append the slice of entries for the current year (from the map)
		// to the result slice.
		result = append(result, groups[year])
	}

	return result
}
