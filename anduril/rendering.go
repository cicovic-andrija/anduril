package anduril

import (
	"fmt"
	"html/template"
	"io"
)

const (
	// Top-level template for rendering HTML pages.
	PageTemplate = "page-v2.html"
	// Dynamic article content template name.
	ContentPlaceholderTemplate = "content"
	// String format for dynamic templates that render article content.
	ContentPlaceholderTemplateFmt = "{{ template \"%s\" . }}"
)

type Page struct {
	Title           string
	Articles        []*Article
	Tags            []string
	HighlightedTags []string
	Is404           bool
	FooterText      string
	contentTemplate string
}

func (p *Page) SlicedArticlesFirstPage() []*Article {
	return p.Articles[0 : len(p.Articles)/2+len(p.Articles)%2]
}

func (p *Page) SlicedArticlesSecondPage() []*Article {
	return p.Articles[len(p.Articles)/2+len(p.Articles)%2 : len(p.Articles)]
}

func (p *Page) IsHighlighted(tag string) bool {
	for _, highlighted := range p.HighlightedTags {
		if highlighted == tag {
			return true
		}
	}
	return false
}

func (p *Page) IsHighlightedRed(tag string) bool {
	return tag == DraftTag || tag == OutdatedTag
}

func (p *Page) IsArticleListVisible() bool {
	return (len(p.Articles) > 1) || (len(p.Articles) == 1 && len(p.HighlightedTags) == 1 && p.contentTemplate == "")
}

func (p *Page) IsTagListVisible() bool {
	return len(p.Tags) > 0 && !(len(p.Articles) == 1 && len(p.HighlightedTags) == 1 && p.HighlightedTags[0] == MetaPageTag)
}

func (s *WebServer) renderArticle(w io.Writer, article *Article, revision *Revision) error {
	footerText := fmt.Sprintf("Last updated on %s", article.ModifiedTime.Format("Jan 2 2006."))
	if article.Comment != "" {
		footerText = article.Comment
	}

	return s.renderPage(w, &Page{
		Title:           article.Title,
		Articles:        []*Article{article},
		HighlightedTags: append([]string{}, article.Tags...),
		Tags:            revision.SortedTags,
		FooterText:      footerText,
		contentTemplate: VersionedHTMLTemplate(article.Key, revision.Hash),
	})
}

func (s *WebServer) renderListOfAllArticles(w io.Writer, revision *Revision) error {
	return s.renderPage(w, &Page{
		Title:      "Articles",
		Articles:   revision.SortedArticles,
		FooterText: fmt.Sprintf("There are %d articles listed.", len(revision.SortedArticles)),
	})
}

func (s *WebServer) renderListOfAllArticlesForTag(w io.Writer, tag string, articles []*Article, revision *Revision) error {
	return s.renderPage(w, &Page{
		Title:           tag,
		Articles:        articles,
		HighlightedTags: []string{tag},
		Tags:            revision.SortedTags,
		FooterText:      fmt.Sprintf("There are %d articles listed.", len(articles)),
	})
}

func (s *WebServer) renderListOfTaggedArticles(w io.Writer, revision *Revision) error {
	return s.renderPage(w, &Page{
		Title:      "Tags",
		Articles:   revision.SortedArticles,
		Tags:       revision.SortedTags,
		FooterText: fmt.Sprintf("There are %d tags listed.", len(revision.SortedTags)),
	})
}

// TODO: Handle situations where Footer text is ""
func (s *WebServer) renderPage(w io.Writer, page *Page) error {
	t, err := template.ParseFiles(s.env.TemplatePath(PageTemplate))
	if err == nil {
		if page.contentTemplate != "" {
			contentPlaceholder := fmt.Sprintf(ContentPlaceholderTemplateFmt, page.contentTemplate)
			t.New(ContentPlaceholderTemplate).Parse(contentPlaceholder)
			_, err = t.ParseFiles(s.env.CompiledTemplatePath(page.contentTemplate))
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
