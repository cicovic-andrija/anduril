package repository

import (
	"errors"

	"github.com/cicovic-andrija/anduril/service"
)

// Package errors.
var (
	ErrNotInitialized  = errors.New("repository not initialized")
	ErrInvalidProtocol = errors.New("invalid protocol, allowed values are 'https' and 'ssh'")
)

// Repository is a local set of files managed by a version control system.
type Repository interface {
	// Root returns the absolute path of repository's root directory on the local file system.
	Root() string

	// ContentRoot returns the absolute path of the directory which contains content files,
	// which naturally must be a subdirectory of the directory returned by Root, or the Root itself.
	ContentRoot() string

	// Empty returns a value indicating whether the repository has been initialized (downloaded and/or opened).
	Empty() bool

	// Initialize either downloads the repository from a remote location to Root and/or opens an already
	// downloaded repository found in Root, or returns an error value in case of failure.
	Initialize(root string, trace service.TraceCallback) error

	// Sync attemps to download a new revision of the repository from the remote location, and returns
	// a value indicating whether a new version was found and downloaded, or an error value in case of failure.
	Sync() (bool, error)

	// LatestRevisionID returns an ID value of the repository's latest revision.
	LatestRevisionID() string
}
