package anduril

import (
	"net/http"
	"net/url"

	"github.com/cicovic-andrija/libgo/https"
)

func (s *WebServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		u := &url.URL{
			Path:     "/home",
			RawQuery: r.URL.RawQuery,
		}
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	} else {
		s.PageNotFoundHandler(w, r)
	}
}

func (s *WebServer) ArticleRootHandlerLocked(w http.ResponseWriter, r *http.Request) {
	// By default, group articles by date.
	groupArticlesBy := "date"
	if r.URL.Query().Get("group-by") == "title" {
		groupArticlesBy = "title"
	}
	if r.URL.Query().Get("group-by") == "type" {
		groupArticlesBy = "type"
	}

	err := s.renderArticleList(w, s.latestRevision, groupArticlesBy)
	if err != nil {
		s.warn("failed to render list of all articles: %v", err)
	}
}

func (s *WebServer) ArticleHandlerLocked(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	article := s.latestRevision.GetArticle(key)
	if article == nil {
		panic(s.error("impossible server state: WebServer.ArticleHandlerLocked: article must exist but not found: key: %s", key))
	}
	err := s.renderArticle(w, article, s.latestRevision)
	if err != nil {
		s.warn("failed to render article: %v", err)
	}
}

func (s *WebServer) TagRootHandlerLocked(w http.ResponseWriter, r *http.Request) {
	tag := s.latestRevision.DefaultTag
	articles := s.latestRevision.SearchByTag(tag)
	if articles == nil {
		panic(s.error("impossible server state: WebServer.TagRootHandlerLocked: articles must exist for tag but not found: tag: %s", tag))
	}
	err := s.renderArticleListForTag(w, tag, articles, s.latestRevision)
	if err != nil {
		s.warn("failed to render list of all articles for tag %q: %v", tag, err)
	}
}

func (s *WebServer) TagHandlerLocked(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Path
	articles := s.latestRevision.SearchByTag(tag)
	if articles == nil {
		panic(s.error("impossible server state: WebServer.TagHandlerLocked: articles must exist for tag but not found: tag: %s", tag))
	}
	err := s.renderArticleListForTag(w, tag, articles, s.latestRevision)
	if err != nil {
		s.warn("failed to render list of all articles for tag %q: %v", tag, err)
	}
}

func (s *WebServer) StaticPageHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	page, found := StaticPages[key]
	if !found {
		panic(s.error("impossible server state: WebServer.StaticPageHandler: key not found: %s", key))
	}
	s.renderPage(w, page)
}

func (s *WebServer) StaticPageRequestHandler() http.Handler {
	return https.Adapt(
		http.HandlerFunc(s.StaticPageHandler),
		https.StripPrefix("/"),
	)
}

func (s *WebServer) PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "404"
	s.StaticPageHandler(w, r)
}

func (s *WebServer) registerHandlers() {
	s.httpsServer.Handle(
		"/",
		http.HandlerFunc(s.RootHandler),
	)

	s.httpsServer.Handle(
		"/home",
		s.StaticPageRequestHandler(),
	)

	s.httpsServer.Handle(
		"/tags",
		https.Adapt(
			http.HandlerFunc(s.TagRootHandlerLocked),
			s.ReadLockRevision,
		),
	)

	s.httpsServer.Handle(
		"/tags/",
		https.Adapt(
			http.HandlerFunc(s.TagHandlerLocked),
			s.FindAndReadLockRevision(TagObject),
			https.StripPrefix("/tags/"),
			https.RedirectRootToParentTree,
		),
	)

	s.httpsServer.Handle(
		"/articles",
		https.Adapt(
			http.HandlerFunc(s.ArticleRootHandlerLocked),
			s.ReadLockRevision,
		),
	)

	s.httpsServer.Handle(
		"/articles/",
		https.Adapt(
			http.HandlerFunc(s.ArticleHandlerLocked),
			s.FindAndReadLockRevision(ArticleObject),
			https.StripPrefix("/articles/"),
			https.RedirectRootToParentTree,
		),
	)

	s.httpsServer.Handle(
		"/about",
		s.StaticPageRequestHandler(),
	)

	s.httpsServer.Handle(
		"/search",
		s.StaticPageRequestHandler(),
	)

	s.httpsServer.Handle(
		"/look-and-feel",
		s.StaticPageRequestHandler(),
	)
}
