package anduril

import (
	"fmt"
	"html/template"
	"io"
)

const (
	// Top-level template for rendering HTML pages.
	PageTemplate = "page.html"

	// String format for dynamic templates that render unique HTML page content.
	ContentTemplateText = "{{ template \"%s\" . }}"
)

func (s *WebServer) renderArticle(w io.Writer, articleTemplate string) error {
	return s.renderPage(w, articleTemplate)
}

func (s *WebServer) renderPage(w io.Writer, contentTemplate string) error {
	t, err := template.ParseFiles(
		s.env.TemplatePath(PageTemplate),
		s.env.CompiledTemplatePath(contentTemplate),
	)
	if err != nil {
		return fmt.Errorf("failed to parse one or more template files: %v", err)
	}

	t.New("content").Parse(fmt.Sprintf(ContentTemplateText, contentTemplate))
	return t.ExecuteTemplate(w, PageTemplate, nil)
}
