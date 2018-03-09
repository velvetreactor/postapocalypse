package web

import (
	"html/template"
	"io"

	"gopkg.in/labstack/echo.v3"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(tmplPath string) *Renderer {
	return &Renderer{
		templates: template.Must(template.ParseGlob(tmplPath)),
	}
}

func (renderer *Renderer) Render(w io.Writer, templateName string, data interface{}, ctx echo.Context) error {
	return renderer.templates.ExecuteTemplate(w, templateName, data)
}
