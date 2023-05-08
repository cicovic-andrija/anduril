package anduril

import (
	"fmt"
	"html/template"
	"io"
)

const (
	// Top-level template for rendering HTML pages.
	PageTemplate = "page.html"
	// Dynamic article content template name.
	ArticleContentTextTemplate = "articleContent"
	// String format for dynamic templates that render article content.
	ArticleContentTextTemplateFmt = "{{ template \"%s\" . }}"
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

func (p *Page) IsHighlighted(tag string) bool {
	for _, highlighted := range p.HighlightedTags {
		if highlighted == tag {
			return true
		}
	}
	return false
}

func (p *Page) RedHighlightedTags() []string {
	var reds []string = nil
	for _, highlighted := range p.HighlightedTags {
		if p.ShouldHighlightRed(highlighted) {
			reds = append(reds, highlighted)
		}
	}
	return reds
}

func (p *Page) ShouldHighlightRed(tag string) bool {
	return tag == WIPTag || tag == OutdatedTag
}

func (p *Page) ShowArticleListInsteadOfContent() bool {
	return (len(p.Articles) > 1) || (len(p.Articles) == 1 && len(p.HighlightedTags) == 1 && p.contentTemplate == "")
}

func (p *Page) IsTagListVisible() bool {
	return len(p.Tags) > 0 && !(len(p.Articles) == 1 && len(p.HighlightedTags) == 1 && p.HighlightedTags[0] == MetaPageTag)
}

func (p *Page) ArticleListHeader() string {
	if len(p.HighlightedTags) == 1 {
		return fmt.Sprintf("Articles tagged %q", p.HighlightedTags[0])
	} else if len(p.Tags) > 0 {
		return "Tagged articles"
	} else {
		return "All articles"
	}
}

func (s *WebServer) renderArticle(w io.Writer, article *Article, revision *Revision) error {
	footerText := article.LastModificationDateMessage()
	if article.Comment != "" {
		footerText = article.Comment
	}

	return s.renderPage(w, &Page{
		Title:           article.Title,
		Articles:        []*Article{article},
		HighlightedTags: append([]string{}, article.Tags...),
		Tags:            revision.SortedTags,
		FooterText:      footerText,
		contentTemplate: article.VersionedHTMLTemplate(revision.Hash),
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

func (s *WebServer) renderPage(w io.Writer, page *Page) error {
	t, err := template.ParseFiles(s.env.TemplatePath(PageTemplate))
	if err == nil {
		if page.contentTemplate != "" {
			t.New(ArticleContentTextTemplate).Parse(fmt.Sprintf(ArticleContentTextTemplateFmt, page.contentTemplate))
			_, err = t.ParseFiles(s.env.CompiledTemplatePath(page.contentTemplate))
		} else {
			t.New(ArticleContentTextTemplate).Parse("")
		}
	}
	if err != nil {
		return fmt.Errorf("failed to parse one or more template files: %v", err)
	}
	return t.ExecuteTemplate(w, PageTemplate, page)
}
