package anduril

import (
	"net/http"

	"github.com/cicovic-andrija/https"
)

func (s *WebServer) RootHandler(w http.ResponseWriter, r *http.Request) {
	s.log("RootHandler")
}

func (s *WebServer) TopicRootHandler(w http.ResponseWriter, r *http.Request) {
	s.log("TopicRootHandler")
}

func (s *WebServer) TopicHandler(w http.ResponseWriter, r *http.Request) {
	s.log("TopicHandler")
}

func (s *WebServer) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	s.log("ArticleHandler")
}

func (s *WebServer) PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.log("PageNotFoundHandler")
}

func (s *WebServer) registerHandlers() {
	s.httpsServer.Handle(
		"/",
		http.HandlerFunc(s.RootHandler),
	)

	s.httpsServer.Handle(
		"/topics",
		http.HandlerFunc(s.TopicRootHandler),
	)

	s.httpsServer.Handle(
		"/topics/",
		https.Adapt(
			http.HandlerFunc(s.TopicHandler),
			https.RedirectRootToParentTree,
		),
	)

	s.httpsServer.Handle(
		"/articles",
		http.HandlerFunc(s.PageNotFoundHandler),
	)

	s.httpsServer.Handle(
		"/articles/",
		https.Adapt(
			http.HandlerFunc(s.ArticleHandler),
			https.RedirectRootToParentTree,
		),
	)
}
