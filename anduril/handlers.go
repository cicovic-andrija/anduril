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

func (s *WebServer) ArticleRootHandler(w http.ResponseWriter, r *http.Request) {
	revision := DummyRevision
	if s.latestRevision != nil {
		s.revisionLock.RLock()
		defer s.revisionLock.RUnlock()
		revision = s.latestRevision
	}

	err := s.renderListOfAllArticles(w, revision)
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

func (s *WebServer) TagRootHandler(w http.ResponseWriter, r *http.Request) {
	revision := DummyRevision
	if s.latestRevision != nil {
		s.revisionLock.RLock()
		defer s.revisionLock.RUnlock()
		revision = s.latestRevision
	}

	err := s.renderListOfTaggedArticles(w, revision)
	if err != nil {
		s.warn("failed to render list of all articles and tags: %v", err)
	}
}

func (s *WebServer) TagHandlerLocked(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Path
	articles := s.latestRevision.SearchByTag(tag)
	if articles == nil {
		panic(s.error("impossible server state: WebServer.TagHandlerLocked: articles must exist for tag but not found: tag: %s", tag))
	}
	err := s.renderListOfAllArticlesForTag(w, tag, articles, s.latestRevision)
	if err != nil {
		s.warn("failed to render list of all articles for tag %q: %v", tag, err)
	}
}

func (s *WebServer) PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, &Page{
		Title:           "Not Found",
		FooterText:      "Page Not Found",
		contentTemplate: NotFoundTemplate,
		isStatic:        true,
	})
}

// FindAndLock is an https.Adapter generator used to make adapters for requests
// that require access to a specific object. If the object identified by URL path
// exists, the current revision is read-locked until the underlying handler completes
// processing the request.
func (s *WebServer) FindAndLock(objectType ObjectType) https.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.revisionLock.RLock()
			if s.latestRevision == nil || !s.latestRevision.FindObject(r.URL.Path, objectType) {
				s.revisionLock.RUnlock()
				s.PageNotFoundHandler(w, r)
				return
			}

			// Call the next handler in the chain.
			h.ServeHTTP(w, r)
			s.revisionLock.RUnlock()
		})
	}
}

func (s *WebServer) registerHandlers() {
	s.httpsServer.Handle(
		"/",
		http.HandlerFunc(s.RootHandler),
	)

	s.httpsServer.Handle(
		"/home",
		https.Adapt(
			http.HandlerFunc(s.ArticleHandlerLocked),
			s.FindAndLock(ArticleObject),
			https.StripPrefix("/"),
		),
	)

	s.httpsServer.Handle(
		"/tags",
		http.HandlerFunc(s.TagRootHandler),
	)

	s.httpsServer.Handle(
		"/tags/",
		https.Adapt(
			http.HandlerFunc(s.TagHandlerLocked),
			s.FindAndLock(TagObject),
			https.StripPrefix("/tags/"),
			https.RedirectRootToParentTree,
		),
	)

	s.httpsServer.Handle(
		"/articles",
		http.HandlerFunc(s.ArticleRootHandler),
	)

	s.httpsServer.Handle(
		"/articles/",
		https.Adapt(
			http.HandlerFunc(s.ArticleHandlerLocked),
			s.FindAndLock(ArticleObject),
			https.StripPrefix("/articles/"),
			https.RedirectRootToParentTree,
		),
	)

	s.httpsServer.Handle(
		"/about",
		https.Adapt(
			http.HandlerFunc(s.ArticleHandlerLocked),
			s.FindAndLock(ArticleObject),
			https.StripPrefix("/"),
		),
	)
}
