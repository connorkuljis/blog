package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
	"github.com/urfave/cli/v3"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:  "content",
		Usage: "My content management system to store markdown entries in sqlite. Portable, Simple, Isolated",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			fmt.Println("Reminder: The substance of the writing is more important that the code.")
			return context.Background(), nil
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "entries",
				Usage: "Operations for creating, editing, deleting and listing entries.",
				Action: func(ctx context.Context, c *cli.Command) error {
					entries := model.NewEntryRepository(db)
					all, err := entries.ReadAllEntries()
					if err != nil {
						return err
					}

					for _, entry := range all {
						fmt.Printf("id: %d, title: %s\n", entry.ID, entry.Title)
					}

					return nil
				},
				Commands: []*cli.Command{
					&cli.Command{
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
						Action: func(ctx context.Context, c *cli.Command) error {
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
						},
					},
					&cli.Command{
						Name:  "edit",
						Usage: "Edit an entry by id.",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							id := c.Int("id")

							repo := model.NewEntryRepository(db)

							currentEntry, err := repo.ReadEntryByID(id)
							if err != nil {
								return err
							}

							f, err := os.CreateTemp("", "*.md")
							if err != nil {
								log.Fatal(err)
							}

							_, err = f.WriteString(currentEntry.String())
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

							err = currentEntry.LoadFromContentString(bytes.NewReader(b))
							if err != nil {
								return err
							}

							err = repo.UpdateEntry(currentEntry)
							if err != nil {
								fmt.Println("Something went wrong! Your entry was not saved.")
								fmt.Println("Backup at:", f.Name())
								log.Fatal(err)
							}

							fmt.Println("Saved entry:")
							fmt.Printf("id: %d, title: %s\n", currentEntry.ID, currentEntry.Title)

							return nil
						},
					},
					&cli.Command{
						Name:  "delete",
						Usage: "Delete an entry by id.",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
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
						},
					},
				},
			},
			&cli.Command{
				Name:  "categories",
				Usage: "operations on categories",
				Action: func(ctx context.Context, c *cli.Command) error {
					repo := model.NewCategoryRepository(db)
					categories, err := repo.ReadAllCategoriesWithEntries()
					if err != nil {
						return err
					}
					for _, category := range categories {
						fmt.Printf("'%s': %s. (%d)\n", category.Title, category.Description, len(category.Entries))
					}
					return nil
				},
				Commands: []*cli.Command{
					&cli.Command{
						Name:  "new",
						Usage: "Create a new category.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "title",
								Usage: "Category title. It must be unique to other categories.",
							},
							&cli.StringFlag{
								Name:    "description",
								Aliases: []string{"d", "desc"},
								Usage:   "Category description.",
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
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

							err = model.NewCategoryRepository(db).CreateCategory(category)
							if err != nil {
								return err
							}

							fmt.Println("Created category:", category.Title)

							return nil
						},
					},
					&cli.Command{
						Name:  "delete",
						Usage: "Delete a category.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "title",
								Required: true,
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							title := c.String("title")

							repo := model.NewCategoryRepository(db)
							if err != nil {
								return err
							}

							category, err := repo.ReadCategoryByTitle(title)
							if err != nil {
								return err
							}

							err = repo.DeleteCategoryByTitle(title)
							if err != nil {
								return err
							}

							fmt.Println("Deleted category.")
							fmt.Printf("'%s': %s\n", category.Title, category.Description)

							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
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
