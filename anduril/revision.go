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

// Revision is a version (identified by Hash) of a set of data files (articles),
// stored on the file system in ContainerPath.
type Revision struct {
	Articles      map[string]*Article
	Tags          map[string][]*Article
	ContainerPath string
	Hash          string
}

func (r *Revision) FindArticle(key string) *Article {
	if article, exists := r.Articles[key]; exists {
		return article
	}
	return nil
}

func (r *Revision) SearchByTag(tag string) []*Article {
	if articles, exists := r.Tags[tag]; exists {
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
		// pull
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
		s.latestRevision = revision
	}

	return nil
}
