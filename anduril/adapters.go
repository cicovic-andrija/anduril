package anduril

import (
	"net/http"

	"github.com/cicovic-andrija/libgo/https"
)

// ReadLockRevision is an https.Adapter used to read-lock the current revision
// until the underlying handler completes processing the request.
func (s *WebServer) ReadLockRevision(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.latestRevision == nil {
			s.PageNotFoundHandler(w, r)
			return
		}
		s.revisionLock.RLock()
		defer s.revisionLock.RUnlock()

		// Call the next handler in the chain.
		h.ServeHTTP(w, r)
	})
}

// FindAndReadLockRevision is an https.Adapter generator used to make adapters for requests
// that require access to a specific object. If the object identified by URL path
// exists, the current revision is read-locked until the underlying handler completes
// processing the request.
func (s *WebServer) FindAndReadLockRevision(objectType ObjectType) https.Adapter {
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
