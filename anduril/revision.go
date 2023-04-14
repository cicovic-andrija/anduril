package anduril

import (
	"errors"
)

// ObjectType represents a type of object within a revision.
type ObjectType int

// Object types enum.
const (
	ArticleObject ObjectType = iota
	TagObject
)

var ErrNotFound = errors.New("content not found")

var DummyRevision = &Revision{}

// Revision is a version (identified by Hash) of a set of objects
// that represent a set of data files (articles) stored on
// the file system location written in ContainerPath.
type Revision struct {
	Articles       map[string]*Article
	SortedArticles []*Article
	Tags           map[string][]*Article
	SortedTags     []string
	ContainerPath  string
	Hash           string
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
