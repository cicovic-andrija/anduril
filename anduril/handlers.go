package anduril

import "net/http"

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
