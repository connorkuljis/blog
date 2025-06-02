package model

type Page interface {
	FileName() string
	TemplateName() string
}
