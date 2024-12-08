package main

import (
	"bytes"
	_ "embed"
	"log"
	"os"
	"os/exec"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/repo"
	"github.com/jmoiron/sqlx"
)

const (
	dbName = "content.sqlite3"
)

//go:embed create_tables.sql
var schema string

func main() {
	db, err := database.Connect(dbName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(db, schema)
	if err != nil {
		log.Fatal(err)
	}

	// err = NewCategoryAction(db, "blog", "all blog posts")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = NewEntryAction(db, "blog", "Great Green Forest")
	if err != nil {
		log.Fatal(err)
	}
}

func NewCategoryAction(db *sqlx.DB, title, description string) error {
	category := model.NewCategory(title, description)

	categories := repo.NewCategoryRepository(db)

	err := categories.CreateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

func NewEntryAction(db *sqlx.DB, category, title string) error {
	entry := model.NewEntry(category, title)

	entries := repo.NewEntryRepository(db)
	err := entries.CreateEntry(entry)
	if err != nil {
		return err
	}

	return nil
}

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

	_, err = f.WriteString(entry.String())
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
		log.Fatal(err)
	}

	return nil
}
