package anduril

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cicovic-andrija/anduril/yfm"
	"github.com/cicovic-andrija/libgo/fs"
	"github.com/cicovic-andrija/libgo/slice"
)

// Structures and routines for processing data files (articles) in markdown format.

const MarkdownExtension = ".md"

// Tags which trigger special behavior or different way of rendering.
const (
	MetaPageTag        = "meta"
	DraftTag           = "draft"
	OutdatedTag        = "outdated"
	PrivateArticleTag  = "private"
	PersonalArticleTag = "my"
	HiddenArticleTag   = "hidden"
)

type Article struct {
	Title        string    `yaml:"title"`
	Comment      string    `yaml:"comment"`
	Tags         []string  `yaml:"tags"`
	Created      string    `yaml:"created"`
	CreatedTime  time.Time `yaml:"-"`
	Modified     string    `yaml:"modified"`
	ModifiedTime time.Time `yaml:"-"`
	File         string    `yaml:"-"`
	Key          string    `yaml:"-"`
}

func (s *WebServer) processRevision(revision *Revision) error {
	if err := fs.EnumerateDirectory(
		revision.ContainerPath,
		func(fileName string) {
			if !strings.HasSuffix(fileName, MarkdownExtension) {
				s.warn("file %s does not have the expected extension %q and will not been processed", fileName, MarkdownExtension)
				return
			}
			if err := s.scanDataFile(revision, fileName); err != nil {
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
			filepath.Join(revision.ContainerPath, article.File),                                             // input file
			filepath.Join(s.env.CompiledWorkDirectory(), VersionedHTMLTemplate(article.Key, revision.Hash)), // output file
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

func (s *WebServer) scanDataFile(revision *Revision, fileName string) error {
	file, err := fs.OpenFile(filepath.Join(revision.ContainerPath, fileName))
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

	if !s.settings.PublishPrivateArticles && slice.ContainsString(article.Tags, PrivateArticleTag) {
		return nil
	}

	if !s.settings.PublishPersonalArticles && slice.ContainsString(article.Tags, PersonalArticleTag) {
		return nil
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
		article.Tags = []string{DraftTag}
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

func (s *WebServer) generateSupportingPages() {
	// TODO
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

	return nil
}
