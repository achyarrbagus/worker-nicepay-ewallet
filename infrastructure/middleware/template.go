package middleware

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	Templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.Templates.ExecuteTemplate(w, name, data)
}
