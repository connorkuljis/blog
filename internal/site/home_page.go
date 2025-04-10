package site

type HomePage struct {
	Site *MySite
}

func (p HomePage) Filepath() string {
	return "public/index.html"
}

func (p HomePage) TemplateName() string {
	return "_index.html"
}
