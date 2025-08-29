package site

type HomePage struct {
	Site *MySite
}

func NewHomePage(s *MySite) HomePage {
	return HomePage{Site: s}
}

func (p HomePage) Title() string {
	return "Home"
}

func (p HomePage) FileName() string {
	return "index.html"
}

func (p HomePage) TemplateName() string {
	return "index.html"
}
