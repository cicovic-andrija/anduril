package anduril

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cicovic-andrija/anduril/yfm"
	"github.com/cicovic-andrija/go-util"
)

// Structure and routines for processing data files in markdown format.

type MarkdownMetadata struct {
	Title        string    `yaml:"title"`
	Tags         []string  `yaml:"tags"`
	Created      string    `yaml:"created"`
	CreatedTime  time.Time `yaml:"-"`
	Modified     string    `yaml:"modified"`
	ModifiedTime time.Time `yaml:"-"`
}

func (s *WebServer) processBatch(workDirPath string) error {
	if err := util.EnumerateDirectory(
		workDirPath,
		func(dataFileName string) {
			err := s.processDataFile(workDirPath, dataFileName)
			if err != nil {
				s.warn("failed to process data file %s: %v", dataFileName, err)
			}
		},
	); err != nil {
		return fmt.Errorf("failed to process data batch: %v", err)
	}
	return nil
}

func (s *WebServer) processDataFile(workDirPath string, dataFileName string) error {
	file, err := util.OpenFile(filepath.Join(workDirPath, dataFileName))
	if err != nil {
		return err
	}

	articleMetadata := &MarkdownMetadata{}
	if err := yfm.Parse(file, articleMetadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %v", err)
	}

	s.checkpoint(
		MarkdownProcessorTag,
		"%s: %q: %v, created:%s modified:%s",
		dataFileName,
		articleMetadata.Title,
		articleMetadata.Tags,
		articleMetadata.Created,
		articleMetadata.Modified,
	)

	return nil
}
