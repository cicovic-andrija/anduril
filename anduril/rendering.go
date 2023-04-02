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

type Page struct {
	Title           string
	Articles        []*Article
	Tags            []string
	HighlightedTags []string
	contentTemplate string
}

func NewArticlePage(article *Article, revision *Revision) *Page {
	page := &Page{
		Title:           article.Title,
		Articles:        []*Article{article},
		HighlightedTags: append([]string{}, article.Tags...),
		Tags:            revision.SortedTags,
		contentTemplate: article.VersionedHTMLTemplate(revision.Hash),
	}

	return page
}

func (p *Page) IsHighlighted(tag string) (highlighted bool) {
	for _, highlighted := range p.HighlightedTags {
		if highlighted == tag {
			return true
		}
	}
	return false
}

func (p *Page) Layout() string {
	if len(p.Articles) == 0 {
		if len(p.HighlightedTags) == 1 {
			return "ArticleListWithTags"
		} else {
			return "NoArticleWithTags"
		}
	} else if len(p.Articles) == 1 {
		return "ArticleWithTags"
	} else {
		return "ArticleListNoTags"
	}
}

func (s *WebServer) renderArticle(w io.Writer, article *Article, revision *Revision) error {
	return s.renderPage(w, NewArticlePage(article, revision))
}

func (s *WebServer) renderPage(w io.Writer, page *Page) error {
	t, err := template.ParseFiles(
		s.env.TemplatePath(PageTemplate),
		s.env.CompiledTemplatePath(page.contentTemplate),
	)
	if err != nil {
		return fmt.Errorf("failed to parse one or more template files: %v", err)
	}
	t.New("content").Parse(fmt.Sprintf(ContentTemplateText, page.contentTemplate))
	return t.ExecuteTemplate(w, PageTemplate, page)
}
