package model

type Site interface {
	Init() error
	Build() []Page
	Render([]Page) error
}
