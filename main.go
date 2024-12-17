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

	"github.com/cheynewallace/tabby"
	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/site"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v3"
)

const sqlxKey = "db"

func main() {
	cmd := &cli.Command{
		Name:  "content",
		Usage: "My content management system to store markdown entries in sqlite. Portable, Simple, Isolated",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			fmt.Println()
			fmt.Println("Reminder: The substance of the writing is more important that the code.")
			fmt.Println()

			db, err := database.Connect()
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
					s := site.Site{
						Title: "Connor's Blog",
					}
					err := s.Render()
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
						Name:  "new",
						Usage: "Create a new entry.",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "category-id",
								Aliases:  []string{"c"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Aliases:  []string{"t"},
								Required: true,
							},
						},
						Action: createEntry,
					},
					{
						Name:  "edit",
						Usage: "Edit an entry by id.",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: editEntry,
					},
					{
						Name:  "delete",
						Usage: "Delete an entry by id.",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Required: true,
							},
						},
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
								Name:     "id",
								Required: true,
							},
						},
						Action: deleteCategory,
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

	entries, err := model.NewEntryRepository(db).ReadAllJoinCategories()
	if err != nil {
		return err
	}

	t := tabby.New()
	t.AddHeader("INDEX", "CATEGORY", "TITLE", "CREATED", "CHAR")
	for i, entry := range entries {
		t.AddLine(i, entry.CategoryTitle, entry.Title, entry.CreatedAt.Format("2006-01-02"), len(entry.Content))
	}

	t.Print()
	return nil
}

// createEntry creates a new entry.
func createEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	category := c.Int("category-id")
	title := c.String("title")

	entry := model.NewEntry(category, title)

	err := model.NewEntryRepository(db).CreateEntry(entry)
	if err != nil {
		return fmt.Errorf("error creating entry: %w", err)
	}

	fmt.Println("Created entry:")
	fmt.Printf("id: %d, title: %s\n", entry.ID, entry.Title)

	return nil
}

// editEntry edits an entry by id.
func editEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	id := c.Int("id")

	repo := model.NewEntryRepository(db)

	entry, err := repo.ReadEntryByID(id)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("/tmp", entry.Title+"*.md")
	if err != nil {
		return err
	}

	_, err = f.WriteString(entry.Content)
	if err != nil {
		return err
	}

	// close the file
	f.Close()

	cmd := exec.Command(os.Getenv("EDITOR"), f.Name())

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	// open the file again.
	b, err := os.ReadFile(f.Name())
	if err != nil {
		return err
	}

	entry.Content = string(b)

	err = repo.UpdateEntry(entry)
	if err != nil {
		fmt.Println("Something went wrong! Your entry was not saved.")
		fmt.Println("Backup at:", f.Name())
		return err
	}

	fmt.Println("Saved entry:")
	fmt.Printf("id: %d, title: %s\n", entry.ID, entry.Title)

	return nil
}

// deleteEntry deletes an entry by id.
func deleteEntry(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	id := c.Int("id")

	repo := model.NewEntryRepository(db)

	entry, err := repo.ReadEntryByID(id)
	if err != nil {
		return err
	}

	err = repo.DeleteEntryByID(id)
	if err != nil {
		return err
	}

	fmt.Println("Deleted entry:")
	fmt.Printf("id: %d, title: %s\n", entry.ID, entry.Title)

	return nil
}

// listCategories lists all categories.
func listCategories(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	repo := model.NewCategoryRepository(db)
	categories, err := repo.ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}
	for _, category := range categories {
		fmt.Printf("[%d] %s(%d)\n", category.ID, category.Title, len(category.Entries))
	}
	return nil
}

// createCategory creates a new category.
func createCategory(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	var title string
	var description string

	if !c.IsSet("title") && !c.IsSet("description") {
		var err error
		title, err = GetInputWithPrompt("Title: ")
		if err != nil {
			return err
		}
		description, err = GetInputWithPrompt("Description: ")
		if err != nil {
			return err
		}
	} else {
		title = c.String("title")
		description = c.String("description")
	}

	category := model.NewCategory(title, description)

	err := model.NewCategoryRepository(db).CreateCategory(category)
	if err != nil {
		return err
	}

	fmt.Println("Created category:", category.Title)

	return nil
}

// deleteCategory deletes a category.
func deleteCategory(ctx context.Context, c *cli.Command) error {
	db := ctx.Value(sqlxKey).(*sqlx.DB)
	id := c.Int("id")

	repo := model.NewCategoryRepository(db)

	category, err := repo.ReadCategoryByID(id)
	if err != nil {
		return err
	}

	err = repo.DeleteCategoryByTitle(id)
	if err != nil {
		return err
	}

	fmt.Println("Deleted category.")
	fmt.Printf("'%s': %s\n", category.Title, category.Description)

	return nil
}

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
