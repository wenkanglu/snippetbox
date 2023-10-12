package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/wenkanglu/snippetbox/internal/models"
	"github.com/wenkanglu/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

const TemplateSuffix = ".go.html"

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, withTemplateSuffix("html/pages/*"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			withTemplateSuffix("html/base"),
			withTemplateSuffix("html/partials/*"),
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func withTemplateSuffix(name string) string {
	return fmt.Sprintf("%s%s", name, TemplateSuffix)
}
