package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"

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
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "entries",
				Usage: "Operations for creating, editing, deleting and listing entries.",
				Commands: []*cli.Command{
					&cli.Command{
						Name:  "new",
						Usage: "Create a new entry.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "category",
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
							category := c.String("category")
							title := c.String("title")

							entry := model.NewEntry(category, title)

							err := model.NewEntryRepository(db).CreateEntry(entry)
							if err != nil {
								return fmt.Errorf("error creating entry: %w", err)
							}

							fmt.Println("Created entry:")
							fmt.Printf("%d: [%s] '%s'\n", entry.Id, entry.Category, entry.Title)

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
								fmt.Println("Backup:")
								fmt.Println(currentEntry.String())
								log.Fatal(err)
							}

							fmt.Println("Saved entry:")
							fmt.Printf("%d: [%s] '%s'\n", currentEntry.Id, currentEntry.Category, currentEntry.Title)

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
							fmt.Printf("%d: [%s] '%s'\n", entry.Id, entry.Category, entry.Title)

							return nil
						},
					},
					&cli.Command{
						Name:  "list",
						Usage: "List all entries.",
						Action: func(ctx context.Context, c *cli.Command) error {
							entries := model.NewEntryRepository(db)
							all, err := entries.ReadAllEntries()
							if err != nil {
								return err
							}

							for _, entry := range all {
								fmt.Printf("%d: [%s] '%s'\n", entry.Id, entry.Category, entry.Title)
							}

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
