package anduril

import (
	"errors"
	"fmt"
	"path/filepath"
)

// Named errors.
var (
	ErrNotFound = errors.New("content not found")
)

// ObjectType represents a type of object within a revision.
type ObjectType int

// Object types enum.
const (
	ArticleObject ObjectType = iota
	TagObject
)

// Revision is a version (identified by Hash) of a set of objects
// that represent a set of data files (articles) stored on
// the file system location written in ContainerPath.
type Revision struct {
	Articles      map[string]*Article
	Tags          map[string][]*Article
	ContainerPath string
	Hash          string
}

func (r *Revision) FindObject(key string, objectType ObjectType) (found bool) {
	switch objectType {
	case ArticleObject:
		_, found = r.Articles[key]
		return
	case TagObject:
		_, found = r.Tags[key]
		return
	default:
		return
	}
}

func (r *Revision) GetArticle(key string) *Article {
	if article, exists := r.Articles[key]; exists {
		return article
	}
	return nil
}

func (r *Revision) SearchByTag(key string) []*Article {
	if articles, exists := r.Tags[key]; exists {
		return articles
	}
	return nil
}

func (s *WebServer) checkForNewRevision(trace TraceCallback, v ...interface{}) error {
	var (
		found bool
	)

	trace("checking for new content revision...")

	// First iteration will initialize the repository.
	if s.repository.repo == nil {
		s.repository.trace = trace
		err := s.repository.OpenOrClone(filepath.Join(s.env.WorkDirectoryPath(), repositorySubdir))
		if err != nil {
			return err
		}
		found = true
	}

	if !found {
		// TODO: pull
	}

	if found {
		revision := &Revision{
			Articles:      make(map[string]*Article),
			Tags:          make(map[string][]*Article),
			ContainerPath: s.repository.ContentRoot(),
			Hash:          s.repository.LatestCommitShortHash(),
		}
		trace("new revision %s found", revision.Hash)

		err := s.processRevision(revision)
		if err != nil {
			return fmt.Errorf("failed to process new revision %s: %v", revision.Hash, err)
		}

		s.revisionLock.Lock()
		s.latestRevision = revision
		s.revisionLock.Unlock()
	}

	return nil
}
