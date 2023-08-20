package anduril

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cicovic-andrija/anduril/service"
	"github.com/cicovic-andrija/libgo/fs"
)

func (s *WebServer) syncRepository(trace service.TraceCallback, v ...interface{}) error {
	var (
		found bool
		err   error
	)

	trace("checking for new content revision...")

	// First iteration will initialize the repository.
	if s.repository.Empty() {
		repoRoot := s.env.RepositoryWorkingDirectory()
		if err = s.repository.Initialize(repoRoot, trace); err != nil {
			return err
		}
		found = true
	}

	if !found && !s.repository.Empty() {
		found, err = s.repository.Sync()
		if err != nil {
			return err
		}
	}

	if found {
		revision := &Revision{
			Articles:      make(map[string]*Article),
			Tags:          make(map[string][]*Article),
			SortedTags:    make([]string, 0),
			ContainerPath: s.repository.ContentRoot(),
			Hash:          s.repository.LatestRevisionID(),
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

func (s *WebServer) cleanUpStaleFiles(trace service.TraceCallback, v ...interface{}) error {
	trace("checking for stale files ready for cleanup...")

	if s.latestRevision == nil {
		trace("aborting search because latest revision is unknown")
		return nil
	}

	latestVersionSuffix := fmt.Sprintf("_%s.html", s.latestRevision.Hash)
	failed := []string{}
	cleanedUp := 0

	if err := fs.EnumerateDirectory(
		filepath.Join(s.env.CompiledWorkDirectory()),
		func(fileName string) {
			if !strings.HasSuffix(fileName, latestVersionSuffix) {
				if err := os.Remove(s.env.CompiledTemplatePath(fileName)); err == nil {
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
