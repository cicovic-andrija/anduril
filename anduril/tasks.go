package anduril

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cicovic-andrija/go-util"
)

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

	if !found && s.repository.repo != nil {
		new, err := s.repository.Pull()
		if err != nil {
			return err
		}
		found = new
	}

	if found {
		revision := &Revision{
			Articles:       make(map[string]*Article),
			SortedArticles: make([]*Article, 0),
			Tags:           make(map[string][]*Article),
			SortedTags:     make([]string, 0),
			ContainerPath:  s.repository.ContentRoot(),
			Hash:           s.repository.LatestCommitShortHash(),
		}
		trace("new revision found with hash %s", revision.Hash)

		err := s.processRevision(revision)
		if err != nil {
			return fmt.Errorf("failed to process new revision %s: %v", revision.Hash, err)
		}

		s.revisionLock.Lock()
		s.latestRevision = revision
		s.revisionLock.Unlock()
		trace("latest revision updated to %s", revision.Hash)
	}

	return nil
}

func (s *WebServer) cleanUpCompiledFiles(trace TraceCallback, v ...interface{}) error {
	trace("checking for stale files ready for cleanup...")

	if s.latestRevision == nil {
		trace("aborting search because latest revision is unknown")
		return nil
	}

	latestVersionSuffix := VersionedArticleTemplateSuffix(s.latestRevision.Hash)
	failed := []string{}
	cleanedUp := 0

	if err := util.EnumerateDirectory(
		filepath.Join(s.env.WorkDirectoryPath(), compiledSubdir),
		func(fileName string) {
			if !strings.HasSuffix(fileName, latestVersionSuffix) {
				if err := os.Remove(filepath.Join(s.env.WorkDirectoryPath(), compiledSubdir, fileName)); err == nil {
					cleanedUp += 1
					trace("%s was cleaned up", fileName)
				} else {
					failed = append(failed, fileName)
				}
			}
		},
	); err != nil {
		return fmt.Errorf("failed to enumerate directory for stale file cleanup: %v", err)
	}

	if cleanedUp > 0 || len(failed) > 0 {
		trace("successfully cleaned up %d stale files; failed to clean up %d stale files", cleanedUp, len(failed))
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed to clean up %d stale files: %v", len(failed), failed)
	}
	return nil
}
