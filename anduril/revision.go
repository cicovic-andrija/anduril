package anduril

// Revision is a version (identified by Hash) of a set of data files (articles),
// stored on file system in ContainerPath.
type Revision struct {
	Articles      map[string]*Article
	Tags          map[string][]*Article
	ContainerPath string
	Hash          string
}
