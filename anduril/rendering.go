package anduril

import (
	"fmt"
	"html/template"
	"io"
)

const (
	// Top-level template for rendering HTML pages.
	PageTemplate = "page-v2.html"
	// Content placeholder template name.
	ContentPlaceholderTemplate = "content"
	// Content placeholder template format.
	ContentPlaceholderTemplateFmt = "{{ template \"%s\" . }}"
)

type Page struct {
	Key               string
	Title             string
	Sidebar           Sidebar
	Tags              []string
	HighlightedTags   []string
	Articles          []*Article
	ArticleGroups     []ArticleGroup
	HeaderText        string
	FooterText        string
	contentTemplate   string
	isCompiledContent bool
}

type Sidebar struct {
	ArticlesHighlighted       bool
	GroupedByTitleHighlighted bool
	GroupedByDateHighlighted  bool
	GroupedByTypeHighlighted  bool
	TagsHighlighted           bool
}

var (
	StaticPages = map[string]*Page{
		"home": {
			Key:        "home",
			Title:      "The L-Archive",
			FooterText: "The universe is transformation: life is opinion.",
		},
		"about": {
			Key:        "about",
			Title:      "About",
			FooterText: "Made by Andrija Cicović, 2023.",
		},
		"search": {
			Key:        "search",
			Title:      "Search Results",
			FooterText: "You can use the sidebar to explore the website.",
		},
		"404": {
			Key:        "404",
			Title:      "Not Found",
			FooterText: "Page not found.",
		},
	}
)

func (p *Page) IsHighlighted(tag string) bool {
	for _, highlighted := range p.HighlightedTags {
		if highlighted == tag {
			return true
		}
	}
	return false
}

func (s *WebServer) renderArticle(w io.Writer, article *Article, revision *Revision) error {
	footerText := fmt.Sprintf("Last updated on %s", article.ModifiedTime.Format("January 2 2006."))
	if article.Comment != "" {
		footerText = article.Comment
	}

	return s.renderPage(w, &Page{
		Key:               article.Key,
		Title:             article.Title,
		Tags:              revision.SortedTags,
		HighlightedTags:   append([]string{}, article.Tags...),
		HeaderText:        fmt.Sprintf("%s | %s", article.Type, article.CreatedTime.Format("January 2 2006.")),
		FooterText:        footerText,
		contentTemplate:   compiledHTMLTemplate(article.Key, revision.Hash),
		isCompiledContent: true,
	})
}

func (s *WebServer) renderArticleList(w io.Writer, revision *Revision, groupBy string) error {
	var (
		articleGroups []ArticleGroup
		sidebar       = Sidebar{ArticlesHighlighted: true}
	)

	if groupBy == "date" {
		sidebar.GroupedByDateHighlighted = true
		articleGroups = revision.GroupsByDate
	} else if groupBy == "title" {
		sidebar.GroupedByTitleHighlighted = true
		articleGroups = revision.GroupsByTitle
	} else { // "type"
		sidebar.GroupedByTypeHighlighted = true
		articleGroups = revision.GroupsByType
	}

	return s.renderPage(w, &Page{
		Key:           "articles",
		Title:         "Articles",
		Sidebar:       sidebar,
		Tags:          revision.SortedTags,
		ArticleGroups: articleGroups,
		FooterText:    fmt.Sprintf("There are %d articles listed.", len(revision.Articles)),
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
		contentTemplate: htmlTemplate("articles"),
	})
}

func (s *WebServer) renderPage(w io.Writer, page *Page) error {
	t, err := template.ParseFiles(s.env.TemplatePath(PageTemplate))
	if err == nil {
		if page.contentTemplate == "" {
			page.contentTemplate = htmlTemplate(page.Key)
			if page.isCompiledContent {
				panic(s.error("impossible server state: WebServer.renderPage: static page marked as isCompiledContent"))
			}
		}
		contentTemplatePath := s.env.TemplatePath(page.contentTemplate)
		if page.isCompiledContent {
			contentTemplatePath = s.env.CompiledTemplatePath(page.contentTemplate)
		}
		t.New(ContentPlaceholderTemplate).Parse(fmt.Sprintf(ContentPlaceholderTemplateFmt, page.contentTemplate))
		_, err = t.ParseFiles(contentTemplatePath)
	}
	if err != nil {
		return fmt.Errorf("failed to parse one or more template files: %v", err)
	}
	return t.ExecuteTemplate(w, PageTemplate, page)
}

func htmlTemplate(key string) string {
	return fmt.Sprintf("%s.html", key)
}

func compiledHTMLTemplate(key string, versionHash string) string {
	return fmt.Sprintf("%s_%s.html", key, versionHash)
}
