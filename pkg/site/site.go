package site

type Site interface {
	Init() error
	Build() []Page
	Render([]Page) error
}

type Page interface {
	Title() string
	FileName() string
	TemplateName() string
}
