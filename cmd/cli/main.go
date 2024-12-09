package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/repo"
	"github.com/jmoiron/sqlx"
)

const (
	storeDir = "store"
	schema   = "create_tables.sql"
	db       = "content.sqlite3"
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

	err = ListEntriesAction(db)
	if err != nil {
		log.Fatal(err)
	}
}

// eg: NewCategoryAction(db, "blog", "all blog posts")
func NewCategoryAction(db *sqlx.DB, title, description string) error {
	category := model.NewCategory(title, description)

	categories := repo.NewCategoryRepository(db)

	err := categories.CreateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

// eg: NewEntryAction(db, "foo", "Great Green Forest")
func NewEntryAction(db *sqlx.DB, category, title string) error {
	entry := model.NewEntry(category, title)

	entries := repo.NewEntryRepository(db)
	err := entries.CreateEntry(entry)
	if err != nil {
		return err
	}

	return nil
}

func ListEntriesAction(db *sqlx.DB) error {
	entries := repo.NewEntryRepository(db)
	all, err := entries.ReadAllEntries()
	if err != nil {
		return err
	}

	for _, e := range all {
		fmt.Printf("%d: [%s] '%s'\n", e.Id, e.Category, e.Title)
	}

	return nil
}

// eg: OpenEntryAction(db, 23)
func OpenEntryAction(db *sqlx.DB, id int64) error {
	entries := repo.NewEntryRepository(db)
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
}
