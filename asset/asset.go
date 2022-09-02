package asset

import (
	"embed"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

//go:embed *.css
var Static embed.FS

//go:embed *.html
var html embed.FS
var templates = template.Must(template.ParseFS(html, "*"))

var Renderer echo.Renderer = &renderer{templates: templates}

type renderer struct {
	templates *template.Template
}

func (t *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
