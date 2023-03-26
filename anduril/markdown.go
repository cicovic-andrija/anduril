package anduril

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/cicovic-andrija/anduril/yfm"
	"github.com/cicovic-andrija/go-util"
)

// Structure and routines for processing data files in markdown format.

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

func (s *WebServer) processDataFile(workDirPath string, fileName string) error {
	file, err := util.OpenFile(filepath.Join(workDirPath, fileName))
	if err != nil {
		return err
	}

	articleMetadata := &ArticleMetadata{}
	if err := yfm.Parse(file, articleMetadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %v", err)
	}

	articleMetadata.File = fileName
	if err := articleMetadata.Normalize(); err != nil {
		return fmt.Errorf("invalid metadata: %v", err)
	}

	s.trace(
		MarkdownProcessorTag,
		"%s (%s): %q: %v, created:%s modified:%s",
		articleMetadata.File,
		articleMetadata.Key,
		articleMetadata.Title,
		articleMetadata.Tags,
		articleMetadata.Created,
		articleMetadata.Modified,
	)

	return nil
}

type ArticleMetadata struct {
	Title        string    `yaml:"title"`
	Tags         []string  `yaml:"tags"`
	Created      string    `yaml:"created"`
	CreatedTime  time.Time `yaml:"-"`
	Modified     string    `yaml:"modified"`
	ModifiedTime time.Time `yaml:"-"`
	File         string    `yaml:"-"`
	Key          string    `yaml:"-"`
}

func (md *ArticleMetadata) Normalize() (err error) {
	if md.Title == "" {
		err = errors.New("empty title")
		return
	}

	if md.Modified == "" {
		err = errors.New("empty 'modified' timestamp")
		return
	} else {
		md.ModifiedTime, err = time.Parse(time.RFC3339, md.Modified)
		if err != nil {
			err = fmt.Errorf("failed to parse 'modified' timestamp: %v", err)
			return
		}
	}

	if md.Created != "" {
		md.CreatedTime, err = time.Parse(time.RFC3339, md.Created)
		if err != nil {
			err = fmt.Errorf("failed to parse 'created' timestamp: %v", err)
			return
		}
	}

	md.Key = strings.ReplaceAll(strings.ToLower(strings.TrimSuffix(md.File, ".md")), " ", "-")

	return nil
}
