package render

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"path"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	Templates    *template.Template
	BuildVersion string
	PublicFS     fs.FS
}

func NewRenderer(fs fs.FS) (*TemplateRenderer, error) {
	// Parse Templates
	modelsGlob := path.Base("*.html")
	partialsGlob := path.Join("partials", "*.html")
	tmpl, err := template.ParseFS(fs, modelsGlob, partialsGlob)
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{
		Templates: tmpl,
	}, nil
}

func (r *TemplateRenderer) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return r.Templates.ExecuteTemplate(w, name, data)
}

func (r *TemplateRenderer) RenderToString(name string, data interface{}) (string, error) {
	var buf []byte
	w := bytes.NewBuffer(buf)

	err := r.Render(w, name, data, nil)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
