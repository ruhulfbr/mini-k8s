package web

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewRenderer() *TemplateRenderer {
	tmpl := template.Must(template.New("").ParseGlob("views/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("views/**/*.html"))

	for _, t := range tmpl.Templates() {
		fmt.Println("Loaded template:", t.Name())
	}

	return &TemplateRenderer{
		templates: tmpl,
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
