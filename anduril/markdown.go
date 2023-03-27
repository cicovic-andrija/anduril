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

// Structures and routines for processing data files (articles) in markdown format.

type Article struct {
	Title        string    `yaml:"title"`
	Tags         []string  `yaml:"tags"`
	Created      string    `yaml:"created"`
	CreatedTime  time.Time `yaml:"-"`
	Modified     string    `yaml:"modified"`
	ModifiedTime time.Time `yaml:"-"`
	File         string    `yaml:"-"`
	Key          string    `yaml:"-"`
}

func (s *WebServer) processRevision(revision *Revision) error {
	if err := util.EnumerateDirectory(
		revision.ContainerPath,
		func(fileName string) {
			err := s.processDataFile(revision, fileName)
			if err != nil {
				s.warn("failed to process data file %s: %v", fileName, err)
			}
		},
	); err != nil {
		return fmt.Errorf("failed to process data batch: %v", err)
	}

	for _, article := range revision.Articles {
		err := s.executor.ConvertMarkdownToHTML(
			filepath.Join(revision.ContainerPath, article.File),
			filepath.Join(s.env.WorkDirectoryPath(), compiledSubdir, fmt.Sprintf("%s_%s.html", article.Key, revision.Hash)),
		)
		if err != nil {
			s.warn("failed to convert %s to HTML: %v", article.File, err)
		}
	}

	return nil
}

func (s *WebServer) processDataFile(revision *Revision, fileName string) error {
	file, err := util.OpenFile(filepath.Join(revision.ContainerPath, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	article := &Article{
		File: fileName,
	}

	if err := yfm.Parse(file, article); err != nil {
		return fmt.Errorf("failed to parse metadata: %v", err)
	}

	if err := article.Normalize(); err != nil {
		return fmt.Errorf("invalid metadata: %v", err)
	}

	// Cache article by key.
	revision.Articles[article.Key] = article

	// Cache tags.
	for _, tag := range article.Tags {
		revision.Tags[tag] = append(revision.Tags[tag], article)
	}

	s.trace(
		MarkdownProcessorTag,
		"%s => [%s]: %q, tags:%v, created:%s, modified:%s",
		article.File,
		article.Key,
		article.Title,
		article.Tags,
		article.Created,
		article.Modified,
	)

	return nil
}

func (md *Article) Normalize() (err error) {
	if md.Title == "" {
		err = errors.New("empty title")
		return
	}

	if md.Created != "" {
		md.CreatedTime, err = time.Parse(time.RFC3339, md.Created)
		if err != nil {
			err = fmt.Errorf("failed to parse 'created' timestamp: %v", err)
			return
		}
	}

	if md.Modified != "" {
		md.ModifiedTime, err = time.Parse(time.RFC3339, md.Modified)
		if err != nil {
			err = fmt.Errorf("failed to parse 'modified' timestamp: %v", err)
			return
		}
	}

	md.Key = strings.ReplaceAll(strings.ToLower(strings.TrimSuffix(md.File, ".md")), " ", "-")

	return nil
}
