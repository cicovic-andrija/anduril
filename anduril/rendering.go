package anduril

import (
	"fmt"
	"html/template"
	"io"
)

const (
	// Top-level template for rendering HTML pages.
	PageTemplate = "page-v2.html"
	// Template for the Articles page.
	ArticlesTemplate = "articles.html"
	// Template for the Not Found page.
	NotFoundTemplate = "404.html"
	// Dynamic article content template name.
	ContentPlaceholderTemplate = "content"
	// String format for dynamic templates that render article content.
	ContentPlaceholderTemplateFmt = "{{ template \"%s\" . }}"
)

type Page struct {
	Key             string
	Title           string
	Sidebar         Sidebar
	Tags            []string
	HighlightedTags []string
	Articles        []*Article
	FooterText      string
	contentTemplate string
	alreadyCompiled bool
}

type Sidebar struct {
	ArticlesHighlighted       bool
	GroupedByTitleHighlighted bool
	GroupedByDateHighlighted  bool
	TagsHighlighted           bool
}

func (p *Page) IsHighlighted(tag string) bool {
	for _, highlighted := range p.HighlightedTags {
		if highlighted == tag {
			return true
		}
	}
	return false
}

func (s *WebServer) renderArticle(w io.Writer, article *Article, revision *Revision) error {
	footerText := fmt.Sprintf("Last updated on %s", article.ModifiedTime.Format("Jan 2 2006."))
	if article.Comment != "" {
		footerText = article.Comment
	}

	return s.renderPage(w, &Page{
		Key:             article.Key,
		Title:           article.Title,
		Tags:            revision.SortedTags,
		HighlightedTags: append([]string{}, article.Tags...),
		FooterText:      footerText,
		contentTemplate: VersionedHTMLTemplate(article.Key, revision.Hash),
		alreadyCompiled: true,
	})
}

func (s *WebServer) renderArticleList(w io.Writer, revision *Revision) error {
	return s.renderPage(w, &Page{
		Key:   "articles",
		Title: "Articles",
		Sidebar: Sidebar{
			ArticlesHighlighted: true,
		},
		Tags:            revision.SortedTags,
		FooterText:      fmt.Sprintf("There are %d articles listed.", 0), // TODO
		contentTemplate: ArticlesTemplate,
	})
}

func (s *WebServer) renderArticleListForTag(w io.Writer, tag string, articles []*Article, revision *Revision) error {
	return s.renderPage(w, &Page{
		Key:   tag,
		Title: tag,
		Sidebar: Sidebar{
			TagsHighlighted: true,
		},
		Tags:            revision.SortedTags,
		HighlightedTags: []string{tag},
		Articles:        articles,
		FooterText:      fmt.Sprintf("There are %d articles listed.", len(articles)),
		contentTemplate: ArticlesTemplate,
	})
}

func (s *WebServer) renderPage(w io.Writer, page *Page) error {
	t, err := template.ParseFiles(s.env.TemplatePath(PageTemplate))
	if err == nil {
		if page.contentTemplate != "" {
			contentPlaceholder := fmt.Sprintf(ContentPlaceholderTemplateFmt, page.contentTemplate)
			contentTemplatePath := s.env.CompiledTemplatePath(page.contentTemplate)
			if !page.alreadyCompiled {
				contentTemplatePath = s.env.TemplatePath(page.contentTemplate)
			}
			t.New(ContentPlaceholderTemplate).Parse(contentPlaceholder)
			_, err = t.ParseFiles(contentTemplatePath)
		} else {
			t.New(ContentPlaceholderTemplate).Parse("")
		}
	}
	if err != nil {
		return fmt.Errorf("failed to parse one or more template files: %v", err)
	}
	return t.ExecuteTemplate(w, PageTemplate, page)
}

func VersionedHTMLTemplate(baseName string, versionHash string) string {
	return fmt.Sprintf("%s_%s.html", baseName, versionHash)
}
