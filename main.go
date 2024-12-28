package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/site"
	"github.com/connorkuljis/content/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v3"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const sqlxKey = "db"

func main() {
	cmd := &cli.Command{
		Name:  "content",
		Usage: "My content management system to store markdown entries in sqlite. Portable, Simple, Isolated",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			db, err := store.Connect()
			if err != nil {
				log.Fatal(err)
			}

			ctx = context.WithValue(ctx, sqlxKey, db)

			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name: "render",
				Action: func(ctx context.Context, c *cli.Command) error {
					db := ctx.Value(sqlxKey).(*sqlx.DB)

					md := goldmark.New(
						goldmark.WithExtensions(
							extension.GFM,
						),
						goldmark.WithParserOptions(
							parser.WithAutoHeadingID(),
						),
						goldmark.WithRendererOptions(
							html.WithHardWraps(),
							html.WithXHTML(),
						),
					)

					s := site.Site{
						Title:          "Connor's Blog",
						DB:             db,
						MarkdownParser: md,
					}

					err := s.Init()
					if err != nil {
						return err
					}

					err = s.Render()
					if err != nil {
						return err
					}

					return nil
				},
			},

			{
				Name:  "entries",
				Usage: "Operations for creating, editing, deleting and listing entries.",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all entries.",
						Action: listEntries,
					},
					{
						Name:   "print",
						Usage:  "Print entry content",
						Action: printEntry,
					},
					{
						Name:   "new",
						Usage:  "Create a new entry.",
						Action: createEntry,
					},
					{
						Name:   "edit",
						Usage:  "Edit an entry by id.",
						Action: editEntry,
					},
					{
						Name:   "delete",
						Usage:  "Delete an entry by id.",
						Action: deleteEntry,
					},
				},
			},
			{
				Name:  "categories",
				Usage: "operations on categories",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List all categories.",
						Action: listCategories,
					},
					{
						Name:   "new",
						Usage:  "Create a new category.",
						Action: createCategory,
					},
					{
						Name:  "delete",
						Usage: "Delete a category.",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name: "id",
							},
						},
						// Action: deleteCategory,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// listEntries lists all entries.
func listEntries(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)

	entryRepo := store.NewEntryRepository(db)

	entries, err := entryRepo.ReadAllEntriesWithCategoryData()
	if err != nil {
		return err
	}

	printEntries(entries)

	return nil
}

// createEntry creates a new entry.
func createEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	categoryRepo := store.NewCategoryRepository(db)
	entryRepo := store.NewEntryRepository(db)

	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}
	printCategories(categories)

	var index int
	fmt.Printf("category index: ")
	fmt.Scanf("%d", &index)

	if index < 0 || index > len(categories)-1 {
		fmt.Errorf("invalid input")
	}

	category := categories[index]

	reader := bufio.NewReader(os.Stdin)

	var title string
	fmt.Println("Please enter entry title, or leave blank for current timestamp")
	fmt.Printf("title: ")
	title, _ = reader.ReadString('\n')
	title = strings.TrimSpace(title)

	if title == "" {
		title = time.Now().Format(time.RFC3339)
	}

	entry := model.NewEntry(category.ID, title)
	err = entryRepo.CreateEntry(entry)
	if err != nil {
		return fmt.Errorf("error creating entry: %w", err)
	}

	var choice string
	fmt.Printf("Open '%s' in editor? [y/N]", entry.Title)
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)

	switch choice {
	case "y":
		err := editEntryContent(entry)
		if err != nil {
			return err
		}
		err = entryRepo.UpdateEntry(entry)
		if err != nil {
			return err
		}
	case "n", "":
		fmt.Println("Done.")
	default:
		return fmt.Errorf("bad input")
	}

	return nil
}

func printEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	entryRepo := store.NewEntryRepository(db)

	entries, err := entryRepo.ReadAllEntries()
	if err != nil {
		return err
	}
	printEntries(entries)

	var index int
	fmt.Printf("entry index: ")
	fmt.Scanf("%d", &index)

	if index < 0 || index > len(entries)-1 {
		fmt.Errorf("invalid input")
	}

	entry := entries[index]

	fmt.Println(entry.Content)

	return nil
}

func editEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)

	entryRepo := store.NewEntryRepository(db)
	entries, err := entryRepo.ReadAllEntriesWithCategoryData()
	if err != nil {
		return err
	}
	printEntries(entries)

	var index int
	fmt.Printf("index: ")
	fmt.Scanf("%d", &index)

	if index < 0 || index > len(entries)-1 {
		return fmt.Errorf("Invalid index")
	}

	entry := &entries[index]

	err = editEntryContent(entry)
	if err != nil {
		return err
	}

	err = entryRepo.UpdateEntry(entry)
	if err != nil {
		return err
	}

	fmt.Println("Updated entry:", entry.Title)

	return nil
}

// deleteEntry deletes an entry by id.
func deleteEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	entryRepo := store.NewEntryRepository(db)

	entries, err := entryRepo.ReadAllEntriesWithCategoryData()
	if err != nil {
		return err
	}
	printEntries(entries)

	// handle user input
	var index int
	fmt.Printf("index: ")
	fmt.Scanf("%d", &index)

	if index < 0 || index > len(entries)-1 {
		return fmt.Errorf("Invalid index")
	}

	entry := &entries[index]

	reader := bufio.NewReader(os.Stdin)

	var choice string
	fmt.Printf("Are you sure you want to delete '%s' [y/N]", entry.Title)
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)

	switch choice {
	case "y":
		err := entryRepo.DeleteEntryByID(entry.ID)
		if err != nil {
			return err
		}
		fmt.Printf("deleted: '%s'\n", entry.Title)
	case "n", "":
		fmt.Println("exiting...")
		return nil
	default:
		return fmt.Errorf("bad input")
	}

	return nil
}

// listCategories lists all categories.
func listCategories(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	categoryRepo := store.NewCategoryRepository(db)

	categories, err := categoryRepo.ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}
	printCategories(categories)

	return nil
}

// createCategory creates a new category.
func createCategory(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	categoryRepo := store.NewCategoryRepository(db)

	reader := bufio.NewReader(os.Stdin)

	var title string
	fmt.Println("Please enter a category title: (Must be unique)")
	fmt.Printf("title: ")
	title, _ = reader.ReadString('\n')
	title = strings.TrimSpace(title)

	var description string
	fmt.Printf("Please enter a short description for '%s'\n", title)
	fmt.Printf("description: ")
	description, _ = reader.ReadString('\n')
	description = strings.TrimSpace(description)

	category := model.NewCategory(title, description)
	err := categoryRepo.CreateCategory(category)
	if err != nil {
		return err
	}

	categories, err := categoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}
	printCategories(categories)

	return nil
}

// deleteCategory deletes a category.
// TODO: check flags or first argument.
// func deleteCategory(ctx context.Context, c *cli.Command) error {
// 	db := ctx.Value(sqlxKey).(*sqlx.DB)

// 	repo := store.NewCategoryRepository(db)

// 	category, err := repo.ReadCategoryByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	err = repo.DeleteCategoryByTitle(id)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Deleted category.")
// 	fmt.Printf("'%s': %s\n", category.Title, category.Description)

// 	return nil
// }

// GetInputWithPrompt prints a prompt to the user and returns the input string from the keyboard
func GetInputWithPrompt(prompt string) (string, error) {
	// Print the prompt
	fmt.Print(prompt)

	// Create a new reader from standard input (keyboard)
	reader := bufio.NewReader(os.Stdin)

	// Read the input line (until newline)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err) // Wrap for context
	}

	// Trim the newline character and any leading/trailing spaces
	input = strings.TrimSpace(input)

	return input, nil // Return the string without the newline character
}

func printEntries(entries []model.Entry) {
	t := tabby.New()
	t.AddHeader("INDEX", "CATEGORY", "TITLE", "CREATED", "CHAR")
	for i, entry := range entries {
		t.AddLine(i, entry.CategoryTitle, entry.Title, entry.CreatedAt.Format("2006-01-02"), len(entry.Content))
	}
	t.Print()
}

func printCategories(categories []model.Category) {
	t := tabby.New()
	t.AddHeader("INDEX", "TITLE", "DESCRIPTION")
	for i, category := range categories {
		t.AddLine(i, category.Title, category.Description)
	}
	t.Print()
}

func editEntryContent(entry *model.Entry) error {
	f, err := os.CreateTemp("/tmp", entry.Title+"*.md")
	if err != nil {
		return err
	}
	fmt.Println("Created temporary file:", f.Name())

	_, err = f.WriteString(entry.Content)
	if err != nil {
		return err
	}
	f.Close() // close the file

	cmd := exec.Command(os.Getenv("EDITOR"), f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error: unable to execute command '%s': %w", cmd.String(), err)
	}

	// returned from editing, open the file again.
	b, err := os.ReadFile(f.Name())
	if err != nil {
		return fmt.Errorf("Error: unable to read from '%s': %w", f.Name(), err)
	}

	entry.Content = string(b)

	return nil
}

// - common functions:
// - select a valid entry from list
