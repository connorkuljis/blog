package site

type HomePage struct {
	Site *MySite
}

func NewHomePage(site *MySite) HomePage {
	return HomePage{Site: site}
}

func (p HomePage) Filepath() string {
	return "public/index.html"
}

func (p HomePage) TemplateName() string {
	return "_index.html"
}
