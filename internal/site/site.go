package site

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/content/internal/database"
	"github.com/connorkuljis/content/internal/model"
)

var root = `
{{define "root"}}
<!DOCTYPE html>
<html lang="en">
	{{ template "head" . }}

	<body>
		{{ template "layout" . }}
	</body>
</html>
{{ end }}`

var head = `
{{ define "head" }}
<head>
<title>Content Blog</title>
</head>
{{ end }}
`

var categoryLayout = `
{{ define "layout" }}
<h1>{{ .Category.Title }}</h1>

	{{ range .Category.Entries }}
	<article>
	<h2>{{ .Title }}</h2>
	<p>{{ .CreatedAt.Format "2006-01-02"}}</p>
	</article>
	{{ end }}

{{ end }}
	`
var entryLayout = `
{{ define "layout" }}
<h1>{{ .Entry.Title }}</h1>
<p>{{ .CreatedAt.Format "2006-01-02"}}</p>

<article>
	<p style="white-space: pre-line">{{ .Entry.Content}}</p>
</article>
{{ end }}
	`

func Render() error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	categoryRepo := model.NewCategoryRepository(db)
	categories, err := categoryRepo.ReadAllCategoriesWithEntries()
	if err != nil {
		return err
	}

	categoryTemplate := template.Must(template.New("").Parse(root + head + categoryLayout))
	entryTemplate := template.Must(template.New("").Parse(root + head + entryLayout))
	for _, category := range categories {
		f, err := os.Create(filepath.Join("public", category.Title+".html"))
		if err != nil {
			panic(err)
		}

		err = categoryTemplate.ExecuteTemplate(f, "root", map[string]any{"Category": category})
		if err != nil {
			return err
		}

		for _, entry := range category.Entries {
			f, err := os.Create(filepath.Join("public", entry.Title+".html"))
			if err != nil {
				panic(err)
			}

			err = entryTemplate.ExecuteTemplate(f, "root", map[string]any{"Entry": entry})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
