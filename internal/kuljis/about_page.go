package kuljis

import "path/filepath"

type AboutPage struct {
	Site *MySite
}

func NewAboutPage(site *MySite) AboutPage {
	p := AboutPage{
		Site: site,
	}
	return p
}

func (p AboutPage) Title() string {
	return "About"
}

func (p AboutPage) FileName() string {
	return filepath.Join("/", "about", "index.html")
}

func (p AboutPage) TemplateName() string {
	return "_about.html"
}
