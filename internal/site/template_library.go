package site

import (
	"os"
	"path/filepath"
)

type HTMLTemplateLibrary struct {
	Base       map[string]string // maps filename -> html string
	Components map[string]string // eg: header.html -> <h1>...</h1>
	Layouts    map[string]string
	Views      map[string]string
}

func NewTemplateLibrary() (*HTMLTemplateLibrary, error) {
	loadTemplates := func(dir string) (map[string]string, error) {
		t := make(map[string]string)

		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".html" {
				filename := filepath.Join(dir, file.Name())
				b, err := os.ReadFile(filename)
				if err != nil {
					return nil, err
				}

				t[file.Name()] = string(b)
			}
		}

		return t, nil
	}

	// load template strings into maps
	base, err := loadTemplates("templates")
	if err != nil {
		return nil, err
	}
	layouts, err := loadTemplates("templates/layouts")
	if err != nil {
		return nil, err
	}
	components, err := loadTemplates("templates/components")
	if err != nil {
		return nil, err
	}
	views, err := loadTemplates("templates/views")
	if err != nil {
		return nil, err
	}

	library := &HTMLTemplateLibrary{
		Base:       base,
		Layouts:    layouts,
		Components: components,
		Views:      views,
	}

	return library, nil
}

func (library *HTMLTemplateLibrary) AllComponents() []string {
	var components []string
	for _, cmpt := range library.Components {
		components = append(components, cmpt)
	}
	return components
}
