package anduril

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/cicovic-andrija/anduril/yfm"
	"github.com/cicovic-andrija/libgo/fs"
	"github.com/cicovic-andrija/libgo/slice"
)

// Structures and routines for processing data files (articles) in markdown format.

const MarkdownExtension = ".md"

// Tags which trigger special behavior or different way of rendering.
const (
	InProgressTag     = "in-progress"
	PrivateArticleTag = "private"
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

	// Axiom: There is at least one article.

	for _, article := range revision.Articles {
		err := s.executor.ConvertMarkdownToHTML(
			filepath.Join(revision.ContainerPath, article.File),                                             // input file
			filepath.Join(s.env.CompiledWorkDirectory(), VersionedHTMLTemplate(article.Key, revision.Hash)), // output file
		)
		if err != nil {
			s.warn("failed to convert %s to HTML: %v", article.File, err)
		}
	}

	// Sort out tags and associated articles.
	for tag, articles := range revision.Tags {
		sort.Slice(articles, func(i, j int) bool {
			return articles[i].ModifiedTime.After(articles[j].ModifiedTime)
		})

		revision.SortedTags = append(revision.SortedTags, tag)
	}
	sort.Strings(revision.SortedTags)
	revision.DefaultTag = revision.SortedTags[0]

	// Sort out articles into groups.
	groupByDate(revision)
	groupByTitle(revision)

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

	if err := article.Normalize(); err != nil {
		return fmt.Errorf("invalid metadata: %v", err)
	}

	// Cache article by key.
	revision.Articles[article.Key] = article

	// Ensure every article is tagged; articles without tags are viewed as incomplete.
	if len(article.Tags) == 0 {
		article.Tags = []string{InProgressTag}
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

	return nil
}

func (a *Article) ModifiedDate() string {
	return a.ModifiedTime.Format("Jan 2 2006.")
}

func groupByDate(revision *Revision) {
	groupMap := make(map[string]ArticleGroup)
	for _, article := range revision.Articles {
		groupTitle := article.ModifiedTime.Format("January 2006.")
		group, found := groupMap[groupTitle]
		if !found {
			group = ArticleGroup{
				GroupTitle: groupTitle,
			}
		}
		group.Articles = append(group.Articles, article)
		groupMap[groupTitle] = group
	}

	revision.GroupsByDate = make([]ArticleGroup, 0, len(groupMap))
	for _, group := range groupMap {
		revision.GroupsByDate = append(revision.GroupsByDate, group)
	}
	sort.Slice(revision.GroupsByDate, func(i, j int) bool {
		return revision.GroupsByDate[i].Articles[0].ModifiedTime.After(revision.GroupsByDate[j].Articles[0].ModifiedTime)
	})
}

func groupByTitle(revision *Revision) {
	groupMap := make(map[string]ArticleGroup)
	for _, article := range revision.Articles {
		groupTitle := string(unicode.ToUpper([]rune(article.Title)[0]))
		group, found := groupMap[groupTitle]
		if !found {
			group = ArticleGroup{
				GroupTitle: groupTitle,
			}
		}
		group.Articles = append(group.Articles, article)
		groupMap[groupTitle] = group
	}

	revision.GroupsByTitle = make([]ArticleGroup, 0, len(groupMap))
	for _, group := range groupMap {
		revision.GroupsByTitle = append(revision.GroupsByTitle, group)
	}
	sort.Slice(revision.GroupsByTitle, func(i, j int) bool {
		return revision.GroupsByTitle[i].GroupTitle < revision.GroupsByTitle[j].GroupTitle
	})
}
