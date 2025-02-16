package site

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/connorkuljis/content/internal/model"
	"github.com/connorkuljis/content/internal/util"
)

var funcMap = template.FuncMap{
	"slugify":  util.Slugify,
	"truncate": util.Truncate,
}

type Site struct {
	Title      string
	Categories []model.Category

	RecentEntries []model.Entry
	Socials       []Social
}

type Social struct {
	URL         string
	Description string
	IconPath    string
}

func (s *Site) Init() error {
	err := removePublicDir()
	if err != nil {
		return err
	}
	err = makePublicDir()
	if err != nil {
		return err
	}
	err = copyStaticFilesIntoPublic()
	if err != nil {
		return err
	}
	return nil
}

func removePublicDir() error {
	return os.RemoveAll("public")
}

func makePublicDir() error {
	return os.MkdirAll("public", os.ModePerm)
}

func copyStaticFilesIntoPublic() error {
	return os.CopyFS("public", os.DirFS("static"))
}

func (s *Site) BuildPages() []Page {
	pages := []Page{
		HomePage{Site: s},
	}

	for _, category := range s.Categories {
		pages = append(pages, CategoryPage{
			Site:     s,
			Category: category,
		})
		for _, entry := range category.Entries {
			pages = append(pages, EntryPage{
				Site:         s,
				Category:     category,
				CurrentEntry: entry,
			})
		}
	}

	return pages
}

func (s *Site) RenderPages(pages []Page) error {
	tpls, err := template.New("").Funcs(funcMap).Option("missingkey=error").ParseGlob("templates/*.html")
	if err != nil {
		return err
	}

	for _, page := range pages {
		dir := filepath.Dir(page.Filepath())
		os.MkdirAll(dir, os.ModePerm)

		f, err := os.Create(page.Filepath())
		if err != nil {
			return err
		}
		defer f.Close()

		err = tpls.ExecuteTemplate(f, page.TemplateName(), page)
		if err != nil {
			return err
		}
	}

	return nil
}
