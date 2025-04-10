package model

type Page interface {
	Filepath() string
	TemplateName() string
}
