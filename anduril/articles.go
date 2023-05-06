package anduril

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cicovic-andrija/anduril/yfm"
	"github.com/cicovic-andrija/go-util"
)

// Structures and routines for processing data files (articles) in markdown format.

// Tags which trigger special behavior or different way of rendering.
const (
	MetaPageTag = "meta"
	WIPTag      = "wip"
	OutdatedTag = "outdated"
)

type Article struct {
	Title        string    `yaml:"title"`
	Subtitle     string    `yaml:"subtitle"`
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

	for tag := range revision.Tags {
		revision.SortedTags = append(revision.SortedTags, tag)
	}
	sort.Strings(revision.SortedTags)

	for _, article := range revision.Articles {
		err := s.executor.ConvertMarkdownToHTML(
			filepath.Join(revision.ContainerPath, article.File),
			filepath.Join(s.env.WorkDirectoryPath(), compiledSubdir, article.VersionedHTMLTemplate(revision.Hash)),
		)
		if err != nil {
			s.warn("failed to convert %s to HTML: %v", article.File, err)
		}
	}

	sort.Slice(revision.SortedArticles, func(i, j int) bool {
		return revision.SortedArticles[i].Title < revision.SortedArticles[j].Title
	})

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

	// Note: will be sorted later.
	revision.SortedArticles = append(revision.SortedArticles, article)

	// Ensure every article is tagged; articles without tags are viewed as incomplete.
	if len(article.Tags) == 0 {
		article.Tags = []string{WIPTag}
	}

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

func (a *Article) Normalize() (err error) {
	if a.Title == "" {
		err = errors.New("empty title")
		return
	}

	if a.Created != "" {
		a.CreatedTime, err = time.Parse(time.RFC3339, a.Created)
		if err != nil {
			err = fmt.Errorf("failed to parse 'created' timestamp: %v", err)
			return
		}
	}

	if a.Modified != "" {
		a.ModifiedTime, err = time.Parse(time.RFC3339, a.Modified)
		if err != nil {
			err = fmt.Errorf("failed to parse 'modified' timestamp: %v", err)
			return
		}
	}

	a.Key = strings.ReplaceAll(strings.ToLower(strings.TrimSuffix(a.File, ".md")), " ", "-")

	if a.Subtitle == "" {
		a.Subtitle = fmt.Sprintf("Last updated on %s", a.ModifiedTime.Format("Jan 2 2006."))
	}

	return nil
}

func (a *Article) VersionedHTMLTemplate(versionHash string) string {
	return fmt.Sprintf("%s%s", a.Key, VersionedArticleTemplateSuffix(versionHash))
}

func VersionedArticleTemplateSuffix(versionHash string) string {
	return fmt.Sprintf("_%s.html", versionHash)
}
