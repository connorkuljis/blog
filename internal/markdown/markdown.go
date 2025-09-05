package markdown

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Renderer struct {
	Markdown goldmark.Markdown
}

func NewRenderer() *Renderer {
	highlighter := highlighting.NewHighlighting(
		highlighting.WithStyle("solarized-dark"),
	)

	return &Renderer{
		Markdown: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				highlighter,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
		),
	}
}

func (r *Renderer) RenderHTML(content string) template.HTML {
	var buf bytes.Buffer
	err := r.Markdown.Convert([]byte(content), &buf)
	if err != nil {
		panic(err)
	}
	return template.HTML(buf.String())
}
