package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
	"github.com/urfave/cli/v3"
)

const (
	storeDir = "store"
	schema   = "create_tables.sql"
	db       = "content.db"
)

func main() {
	db, err := database.Connect(filepath.Join(storeDir, db))
	if err != nil {
		log.Fatal(err)
	}

	b, _ := os.ReadFile(filepath.Join(storeDir, schema))
	_, err = database.Exec(db, string(b))
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:  "cms",
		Usage: "a content management system",
		Commands: []*cli.Command{
			&cli.Command{
				Name: "new",
				Commands: []*cli.Command{
					&cli.Command{
						Name: "entry",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "category",
								Aliases:  []string{"c"},
								Required: true,
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							fmt.Println("### Creating New Entry")

							category := c.String("category")

							if c.NArg() < 1 {
								return fmt.Errorf("invalid args...")
							}
							title := c.Args().First()

							e := model.NewEntry(category, title)
							r := model.NewEntryRepository(db)
							err := r.CreateEntry(e)
							if err != nil {
								return fmt.Errorf("error creating entry: %w", err)
							}

							fmt.Println("Successfully created (1) entry in entries.")
							fmt.Println("id:", e.Id)
							fmt.Println("title:", e.Title)
							fmt.Println("category:", e.Category)

							return nil
						},
					},
					&cli.Command{
						Name: "category",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "description",
								Aliases:  []string{"d", "desc"},
								Required: true,
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							fmt.Println("### Creating New Category")

							desc := c.String("description")

							if c.NArg() < 1 {
								return fmt.Errorf("invalid args...")
							}
							title := c.Args().First()

							newCategory := model.NewCategory(title, desc)
							r := model.NewCategoryRepository(db)
							err := r.CreateCategory(newCategory)
							if err != nil {
								return fmt.Errorf("error creating category: %w", err)
							}

							fmt.Println("Successfully created (1) entry in categories.")
							fmt.Println("title:", newCategory.Title)
							fmt.Println("description:", newCategory.Description)

							return nil
						},
					},
				},
			},
			&cli.Command{
				Name: "list",
				Action: func(ctx context.Context, c *cli.Command) error {
					entries := model.NewEntryRepository(db)
					all, err := entries.ReadAllEntries()
					if err != nil {
						return err
					}

					for _, e := range all {
						fmt.Printf("%d: [%s] '%s'\n", e.Id, e.Category, e.Title)
					}

					return nil
				},
			},
			&cli.Command{
				Name: "edit",
				Action: func(ctx context.Context, c *cli.Command) error {

					if c.NArg() < 1 {
						return fmt.Errorf("invalid args...")
					}
					id, err := strconv.ParseInt(c.Args().First(), 10, 64)
					if err != nil {
						return fmt.Errorf("bad argument '%s': %w", c.Args().First(), err)
					}

					entries := model.NewEntryRepository(db)
					entry, err := entries.ReadEntryByID(id)
					if err != nil {
						return err
					}

					f, err := os.CreateTemp("", "*")
					if err != nil {
						log.Fatal(err)
					}
					tempFile := f.Name()

					body := entry.String()
					_, err = f.WriteString(body)
					if err != nil {
						log.Fatal(err)
					}
					f.Close()

					cmd := exec.Command(os.Getenv("EDITOR"), tempFile)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout

					err = cmd.Run()
					if err != nil {
						log.Fatal(err)
					}

					b, err := os.ReadFile(tempFile)
					if err != nil {
						log.Fatal(err)
					}

					err = model.ParseEntryFromFrontMatter(bytes.NewReader(b), entry)
					if err != nil {
						log.Fatal(err)
					}

					err = entries.UpdateEntry(entry)
					if err != nil {
						fmt.Println("Something went wrong! Your entry was not saved.")
						fmt.Println("Backup:")
						fmt.Println(entry.String())
						log.Fatal(err)
					}
					fmt.Printf("%d: [%s] '%s' saved.\n", entry.Id, entry.Category, entry.Title)

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
