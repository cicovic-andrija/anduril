package anduril

import (
	"net/http"
	"net/url"

	"github.com/cicovic-andrija/https"
)

func (s *WebServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		u := &url.URL{
			Path:     "/articles",
			RawQuery: r.URL.RawQuery,
		}
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	} else {
		s.PageNotFoundHandler(w, r)
	}
}

func (s *WebServer) ArticleRootHandler(w http.ResponseWriter, r *http.Request) {
	s.log("ArticleRootHandler")
}

func (s *WebServer) ArticleHandlerLocked(w http.ResponseWriter, r *http.Request) {
	s.log("ArticleHandler")
}

func (s *WebServer) TagRootHandler(w http.ResponseWriter, r *http.Request) {
	s.log("TagRootHandler")
}

func (s *WebServer) TagHandlerLocked(w http.ResponseWriter, r *http.Request) {
	s.log("TagHandler")
}

func (s *WebServer) PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.log("PageNotFoundHandler")
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
}
