package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/connorkuljis/blog/internal/markdown"
	"github.com/connorkuljis/blog/internal/model"
	"github.com/connorkuljis/blog/internal/site"
	"github.com/connorkuljis/blog/internal/store"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v3"
)

type App struct {
	DB         *sqlx.DB
	Entries    []*model.Entry
	Categories []*model.Category
	Tags       []*model.Tag
}

const KeyApp = "app"

var ErrSelectionCancelled = errors.New("selection cancelled by user")

func main() {
	cmd := &cli.Command{
		Name:  "cms",
		Usage: "A portable and simple content management system.",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {

			cfg, err := site.LoadConfig("config.json")
			if err != nil {
				log.Fatal(err)
			}
			db, err := store.Connect(cfg.SqliteURI)
			if err != nil {
				log.Fatal(err)
			}
			markdownRenderer := markdown.NewRenderer()

			app := &App{
				DB: db,
			}

			tags, err := store.NewTagRepo(db).GetTags()
			if err != nil {
				log.Fatal(err)
			}
			for _, t := range tags {
				mTag := model.NewTag(t)
				app.Tags = append(app.Tags, mTag)
			}

			categories, err := store.NewCategoryRepo(db).ReadAllCategories()
			if err != nil {
				log.Fatal(err)
			}

			for _, c := range categories {
				mCategory := model.NewCategory(c)
				app.Categories = append(app.Categories, mCategory)

				entries, err := store.NewEntryRepo(db).ReadAllByCategoryID(c.ID, true)
				if err != nil {
					log.Fatal(err)
				}
				for _, e := range entries {
					tags, err := store.NewTagRepo(db).GetTagsForEntry(e.ID)
					if err != nil {
						log.Fatal(err)
					}
					var mTags []*model.Tag
					for _, t := range tags {
						mTags = append(mTags, model.NewTag(t))
					}

					mEntry := model.NewEntry(e, mCategory, mTags)
					mEntry.Markdown = markdownRenderer.RenderHTML(mEntry.Content)
					mCategory.AddEntry(mEntry)
					app.Entries = append(app.Entries, mEntry)
				}
			}

			// Sort entries by creation date (newest first)
			sort.Slice(app.Entries, func(i, j int) bool {
				return app.Entries[i].CreatedAt.After(app.Entries[j].CreatedAt)
			})

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
						Action: listEntries,
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
						Name:   "update",
						Usage:  "Update an existing category.",
						Action: updateCategory,
					},
					{
						Name:   "delete",
						Usage:  "Delete a category.",
						Action: deleteCategory,
					},
				},
			},
			{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "operations on tags",
				Action:  listTags,
				Commands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new tag.",
						Action: createTag,
					},
					{
						Name:   "update",
						Usage:  "Update an existing tag.",
						Action: updateTag,
					},
					{
						Name:   "delete",
						Usage:  "Delete a tag.",
						Action: deleteTag,
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

	entries := app.Entries

	fmt.Println("ID | Created | Title ")
	for _, e := range entries {
		fmt.Printf("%d | %s | %s\n", e.ID, e.CreatedAt.Format("2006-01-02"), e.Title)
	}

	return nil
}

func createEntry(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	categories := app.Categories

	reader := bufio.NewReader(os.Stdin)
	category, err := selectCategory(reader, categories)
	if err != nil {
		return err
	}
	fmt.Println(category)

	fmt.Println("Please enter entry title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	title = strings.TrimSpace(title)

	entry := &store.Entry{
		CategoryID: category.ID,
		Title:      title,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = store.NewEntryRepo(app.DB).CreateEntry(entry)
	if err != nil {
		return fmt.Errorf("error creating entry: %w", err)
	}

	fmt.Printf("Created entry: %s\n", entry.Title)

	return nil
}

func updateEntry(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	entries := app.Entries

	reader := bufio.NewReader(os.Stdin)
	entry, err := selectEntry(reader, entries)
	if err != nil {
		return err
	}

	for {
		fmt.Println("Select a field to update.")
		fmt.Println("1. Title")
		fmt.Println("2. Description")
		fmt.Println("3. Content")
		fmt.Println("4. Featured Image Url")
		fmt.Println("5. Manage Tags")

		fmt.Printf("Enter the number (%d-%d) or 'q' to quit: ", 1, 5)
		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		choice = strings.TrimSpace(choice)

		if choice == "q" {
			return ErrSelectionCancelled
		}

		choiceNum, _ := strconv.Atoi(choice) // let error be handled by default case below

		switch choiceNum {
		case 1:
			title, err := readContentFromEditor(entry.Title)
			if err != nil {
				return err
			}
			entry.Title = title
		case 2:
			description, err := readContentFromEditor(entry.Description)
			if err != nil {
				return err
			}
			entry.Description = description
		case 3:
			content, err := readContentFromEditor(entry.Content)
			if err != nil {
				return err
			}
			entry.Content = content
		case 4:
			featuredImageUrl, err := readContentFromEditor(entry.FeaturedImageURL)
			if err != nil {
				return err
			}
			entry.FeaturedImageURL = featuredImageUrl
		case 5:
			err := manageTagsForEntry(reader, app, entry)
			if err != nil {
				return err
			}
		default:
			fmt.Println("Bad input, must be between 1 and 5: got:", choice)
			continue
		}

		err = store.NewEntryRepo(app.DB).UpdateEntry(entry.ToStoreEntry())
		if err != nil {
			return err
		}

		fmt.Println("Updated entry:", entry.Title)
	}
}

// deleteEntry deletes an entry by id.
func deleteEntry(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)
	reader := bufio.NewReader(os.Stdin)

	entries := app.Entries

	entry, err := selectEntry(reader, entries)
	if err != nil {
		return err
	}

	fmt.Printf("Confirm delete '%s' [y/N]: ", entry.Title)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err) // Wrap for context
	}

	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "y":
		if err := store.NewEntryRepo(app.DB).DeleteEntryByID(entry.ID); err != nil {
			return fmt.Errorf("failed to delete entry: %w", err)
		}
		fmt.Printf("Deleted: '%s'\n", entry.Title)
	case "n", "":
		fmt.Println("Exiting...")
		return nil
	default:
		return fmt.Errorf("invalid input: '%s', expected 'y' or 'n'", choice)
	}

	return nil
}

// listCategories lists all categories.
func listCategories(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	categories := app.Categories

	for i, c := range categories {
		fmt.Printf("%d. %s\n", i, c.Title)
	}

	return nil
}

// createCategory creates a new category.
func createCategory(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter a category title: (Must be unique)")
	fmt.Printf("title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	title = strings.TrimSpace(title)

	fmt.Printf("Please enter a short description for '%s'\n", title)
	fmt.Printf("description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	description = strings.TrimSpace(description)

	category := &store.Category{
		Title:       title,
		Description: description,
	}
	err = store.NewCategoryRepo(app.DB).CreateCategory(category)
	if err != nil {
		return err
	}

	fmt.Printf("Created category: '%s'\n", category.Title)

	return nil
}

func updateCategory(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)

	categories := app.Categories

	reader := bufio.NewReader(os.Stdin)
	selectedCategory, err := selectCategory(reader, categories)
	if err != nil {
		return err
	}

	for {
		fmt.Println("Select a field to update.")
		fmt.Println("1. Title")
		fmt.Println("2. Description")

		fmt.Printf("Enter the number (%d-%d) or 'q' to quit: ", 1, 2)
		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		choice = strings.TrimSpace(choice)

		if choice == "q" {
			return ErrSelectionCancelled
		}

		choiceNum, err := strconv.Atoi(choice)
		if err != nil {
			return err
		}

		switch choiceNum {
		case 1:
			title, err := readContentFromEditor(selectedCategory.Title)
			if err != nil {
				return err
			}
			selectedCategory.Title = title
		case 2:
			description, err := readContentFromEditor(selectedCategory.Description)
			if err != nil {
				return err
			}
			selectedCategory.Description = description
		default:
			fmt.Println("Bad input, must be between 1 and 4: got:", choiceNum)
			continue
		}

		err = store.NewCategoryRepo(app.DB).UpdateCategory(selectedCategory.ToStoreCategory())
		if err != nil {
			return err
		}

		fmt.Println("Updated category:", selectedCategory.Title)
	}
}

// deleteCategory deletes a category.
func deleteCategory(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)
	reader := bufio.NewReader(os.Stdin)

	categories := app.Categories

	category, err := selectCategory(reader, categories)
	if err != nil {
		return err
	}

	fmt.Printf("Confirm delete '%s' [y/N]: ", category.Title)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "y":
		if err = store.NewCategoryRepo(app.DB).DeleteCategoryByID(category.ID); err != nil {
			return fmt.Errorf("failed to delete category: %w", err)
		}
		fmt.Printf("Deleted: '%s'\n", category.Title)
	case "n", "":
		fmt.Println("Exiting...")
		return nil
	default:
		return fmt.Errorf("invalid input: '%s', expected 'y' or 'n'", choice)
	}

	return nil
}

func readContentFromEditor(content string) (string, error) {
	f, err := os.CreateTemp("/tmp", "*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}
	f.Close() // close the file because we are going to read from it again.

	cmd := exec.Command(os.Getenv("EDITOR"), f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error: unable to execute command '%s': %w", cmd.String(), err)
	}

	b, err := os.ReadFile(f.Name()) // re-read the file again
	if err != nil {
		return "", fmt.Errorf("Error: unable to read from '%s': %w", f.Name(), err)
	}

	result := strings.TrimRight(string(b), "\n\r")

	return result, nil
}

func selectEntry(reader *bufio.Reader, entries []*model.Entry) (*model.Entry, error) {
	for i, e := range entries {
		fmt.Printf("%d. %s\n", i+1, e.Title)
	}

	idx, err := getInputBetween(reader, 1, len(entries))
	if err != nil {
		return nil, err
	}

	return entries[idx-1], nil
}

func selectCategory(reader *bufio.Reader, categories []*model.Category) (*model.Category, error) {
	for i, c := range categories {
		fmt.Printf("%d. %s\n", i+1, c.Title)
	}

	idx, err := getInputBetween(reader, 1, len(categories))
	if err != nil {
		return nil, err
	}

	return categories[idx-1], nil
}

func getInputBetween(reader *bufio.Reader, min int, max int) (int, error) {
	var choiceNum int

	for {
		fmt.Printf("Enter the number (%d-%d) or 'q' to quit: ", min, max)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read input: %w", err)
			continue
		}

		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "q" {
			fmt.Println("\nSelection cancelled.")
			return 0, ErrSelectionCancelled // Return specific error for cancellation
		}

		choiceNum, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number or 'q'.")
			continue // Ask the user again
		}

		if choiceNum < min || choiceNum > max {
			fmt.Printf("Invalid choice. Please enter a number between %d and %d or 'q'.\n", min, max)
			continue
		}

		break
	}

	return choiceNum, nil
}

func listTags(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)
	tags := app.Tags
	for i, t := range tags {
		fmt.Printf("%d. %s\n", i, t.Name)
	}
	return nil
}

func createTag(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter a tag name: (Must be unique)")
	fmt.Printf("name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)

	tag := &store.Tag{
		Name: name,
	}
	err = store.NewTagRepo(app.DB).CreateTag(tag)
	if err != nil {
		return err
	}

	fmt.Printf("Created tag: '%s'\n", tag.Name)
	return nil
}

func updateTag(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)
	reader := bufio.NewReader(os.Stdin)

	tags := app.Tags
	tag, err := selectTag(reader, tags)
	if err != nil {
		return err
	}

	fmt.Printf("Enter the new name for '%s': ", tag.Name)
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	tag.Name = name

	err = store.NewTagRepo(app.DB).UpdateTag(tag.ToStoreTag())
	if err != nil {
		return err
	}

	fmt.Printf("Updated tag: '%s'\n", tag.Name)
	return nil
}

func deleteTag(ctx context.Context, c *cli.Command) error {
	app := ctx.Value(KeyApp).(*App)
	reader := bufio.NewReader(os.Stdin)

	tags := app.Tags
	tag, err := selectTag(reader, tags)
	if err != nil {
		return err
	}

	fmt.Printf("Confirm delete '%s' [y/N]: ", tag.Name)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "y":
		if err = store.NewTagRepo(app.DB).DeleteTag(tag.ID); err != nil {
			return fmt.Errorf("failed to delete tag: %w", err)
		}
		fmt.Printf("Deleted: '%s'\n", tag.Name)
	case "n", "":
		fmt.Println("Exiting...")
		return nil
	default:
		return fmt.Errorf("invalid input: '%s', expected 'y' or 'n'", choice)
	}

	return nil
}

func selectTag(reader *bufio.Reader, tags []*model.Tag) (*model.Tag, error) {
	for i, t := range tags {
		fmt.Printf("%d. %s\n", i+1, t.Name)
	}

	idx, err := getInputBetween(reader, 1, len(tags))
	if err != nil {
		return nil, err
	}

	return tags[idx-1], nil
}

func manageTagsForEntry(reader *bufio.Reader, app *App, entry *model.Entry) error {
	fmt.Println("1. Add Tag")
	fmt.Println("2. Remove Tag")
	fmt.Printf("Enter the number (1-2) or 'q' to quit: ")

	choice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	choice = strings.TrimSpace(choice)

	if choice == "q" {
		return ErrSelectionCancelled
	}

	choiceNum, _ := strconv.Atoi(choice)

	switch choiceNum {
	case 1:
		// Filter out tags already on the entry
		var availableTags []*model.Tag
		for _, t := range app.Tags {
			found := false
			for _, et := range entry.Tags {
				if et.ID == t.ID {
					found = true
					break
				}
			}
			if !found {
				availableTags = append(availableTags, t)
			}
		}

		if len(availableTags) == 0 {
			fmt.Println("No new tags to add.")
			return nil
		}

		tag, err := selectTag(reader, availableTags)
		if err != nil {
			return err
		}
		err = store.NewEntryRepo(app.DB).AddTagToEntry(entry.ID, tag.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Added tag '%s' to entry '%s'\n", tag.Name, entry.Title)

	case 2:
		if len(entry.Tags) == 0 {
			fmt.Println("No tags to remove from this entry.")
			return nil
		}

		tag, err := selectTag(reader, entry.Tags)
		if err != nil {
			return err
		}
		err = store.NewEntryRepo(app.DB).RemoveTagFromEntry(entry.ID, tag.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Removed tag '%s' from entry '%s'\n", tag.Name, entry.Title)
	default:
		fmt.Println("Bad input, must be 1 or 2: got:", choice)
	}

	return nil
}
