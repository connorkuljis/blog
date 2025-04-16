package main

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v3"
)

type App struct {
	DB           *sqlx.DB
	EntryRepo    *store.EntryRepo
	CategoryRepo *store.CategoryRepo
}

const KeyApp = "app"

// ErrSelectionCancelled is a specific error returned when the user quits.
var ErrSelectionCancelled = errors.New("selection cancelled by user")

func main() {
	cmd := &cli.Command{
		Name:  "cms",
		Usage: "A portable and simple content management system.",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			db, err := store.Connect()
			if err != nil {
				log.Fatal(err)
			}

			app := &App{
				DB:           db,
				EntryRepo:    store.NewEntryRepo(db),
				CategoryRepo: store.NewCategoryRepo(db),
			}

			ctx = context.WithValue(ctx, KeyApp, app)

			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:    "entries",
				Aliases: []string{"e"},
				Usage:   "Operations for creating, editing, deleting and listing entries.",
				Action:  listEntries,
				Commands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create entry.",
						Action: createEntry,
					},
					{
						Name:   "update",
						Usage:  "Update existing entries.",
						Action: updateEntry,
					},
					{
						Name:   "list",
						Usage:  "List entries.",
						Action: updateEntry,
					},
					{
						Name:   "delete",
						Usage:  "Delete entry",
						Action: deleteEntry,
					},
				},
			},
			{
				Name:    "categories",
				Aliases: []string{"c"},
				Usage:   "operations on categories",
				Action:  listCategories,
				Commands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new category.",
						Action: createCategory,
					},
					{
						Name:   "delete",
						Usage:  "Delete a category.",
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
	app := ctx.Value(KeyApp).(*App)

	categories, err := app.CategoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}

	for _, c := range categories {
		fmt.Printf("[%s]\n", c.Title)

		entries, err := app.EntryRepo.ReadAllByCategoryID(c.ID, true)
		if err != nil {
			return err
		}

		for _, e := range entries {
			fmt.Printf("\t%s\n", e.Title)
		}
	}

	return nil
}

func createEntry(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	categories, err := app.CategoryRepo.ReadAllCategories()
	if err != nil {
		return err
	}

	var index int
	fmt.Printf("category index: ")
	fmt.Scanf("%d", &index)

	if index < 0 || index > len(categories)-1 {
		return fmt.Errorf("invalid input")
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
	err = app.EntryRepo.CreateEntry(entry)
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
		err = app.EntryRepo.UpdateEntry(entry)
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

func updateEntry(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	reader := bufio.NewReader(os.Stdin)
	var selectedItem model.Entry

	for {
		entries, err := app.EntryRepo.ReadAllEntries(true)
		if err != nil {
			return err
		}

		fmt.Println("\nPlease select an item from the list")
		for i, entry := range entries {
			fmt.Printf("%d. %s\n", i+1, entry.Title)
		}

		fmt.Printf("Enter the number (1-%d) or 'q' to quit: ", len(entries))

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "q" {
			fmt.Println("\nSelection cancelled.")
			return ErrSelectionCancelled // Return specific error for cancellation
		}

		choiceNum, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number or 'q'.")
			continue // Ask the user again
		}

		if choiceNum < 1 && choiceNum > len(entries) {
			fmt.Printf("Invalid choice. Please enter a number between 1 and %d or 'q'.\n", len(entries))
			continue
		}

		selectedIndex := choiceNum - 1
		selectedItem = entries[selectedIndex]

		err = editEntryContent(&selectedItem)
		if err != nil {
			return err
		}

		err = app.EntryRepo.UpdateEntry(&selectedItem)
		if err != nil {
			return err
		}

		fmt.Println("Updated entry:", selectedItem.Title)
	}
}

// deleteEntry deletes an entry by id.
func deleteEntry(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	entries, err := app.EntryRepo.ReadAllEntries(true)
	if err != nil {
		return err
	}

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
		err := app.EntryRepo.DeleteEntryByID(entry.ID)
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
	// app := ctx.Value(appKey).(*App)

	return nil
}

// createCategory creates a new category.
func createCategory(ctx context.Context, c *cli.Command) error {
	// app := ctx.Value(appKey).(*App)
	//
	// reader := bufio.NewReader(os.Stdin)
	//
	// var title string
	// fmt.Println("Please enter a category title: (Must be unique)")
	// fmt.Printf("title: ")
	// title, _ = reader.ReadString('\n')
	// title = strings.TrimSpace(title)
	//
	// var description string
	// fmt.Printf("Please enter a short description for '%s'\n", title)
	// fmt.Printf("description: ")
	// description, _ = reader.ReadString('\n')
	// description = strings.TrimSpace(description)
	//
	// category := model.NewCategory(title, description)
	// err := app.CategoryRepo.CreateCategory(category)
	// if err != nil {
	// 	return err
	// }
	//
	// categories, err := app.CategoryRepo.ReadAllCategories()
	// if err != nil {
	// 	return err
	// }
	//
	return nil
}

// deleteCategory deletes a category.
func deleteCategory(ctx context.Context, c *cli.Command) error {
	// db := ctx.Value(sqlxKey).(*sqlx.DB)
	//
	// first := c.Args().First()
	// if first == "" {
	// 	// TODO: define errors such as missing argument, invalid argument ect...
	// 	return fmt.Errorf("error: missing 1 positional argument: id")
	// }
	//
	// repo := store.NewCategoryRepository(db)
	//
	// category, err := repo.ReadCategoryByID(id)
	// if err != nil {
	// 	return err
	// }
	//
	// fmt.Println("Are you sure you want to delete category '" + category.Title + "'")
	//
	// err = repo.DeleteCategoryByTitle(id)
	// if err != nil {
	// 	return err
	// }
	//
	// fmt.Println("Deleted category.")
	// fmt.Printf("'%s': %s\n", category.Title, category.Description)

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
