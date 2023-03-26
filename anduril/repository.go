package anduril

import (
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type RepositoryProcessor struct {
	RepositoryConfig
	repo       *git.Repository
	trace      TraceCallback
	commitHash string
}

func (r *RepositoryProcessor) OpenOrClone(path string) error {
	// Try to open an existing local repository if there is one, otherwise clone it from remote location.
	local, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			r.trace("local repository not found on: %s", path)
		} else {
			return fmt.Errorf("open local repository %s: operation failed with error: %v", path, err)
		}
	} else {
		r.repo = local
		return r.ValidateState()
	}

	clone, err := git.PlainClone(
		path,
		false, // bare
		&git.CloneOptions{
			URL:           r.URL,
			RemoteName:    r.Remote,
			ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
			Progress:      io.Discard,
		},
	)

	if err != nil {
		return fmt.Errorf("clone repository: remote %q: operation failed with error: %v", r.URL, err)
	}

	r.trace("repository cloned from remote location %q to local path %s", r.URL, path)
	r.repo = clone
	return r.ValidateState()
}

func (r *RepositoryProcessor) ValidateState() error {
	current, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to obtain a reference to the current commit: %v", err)
	}
	r.trace("current commit hash is %s", current.Hash())

	branchRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(r.Branch), false /* resolved */)
	if err != nil {
		return fmt.Errorf("failed to obtain a reference to the %q branch: %v", r.Branch, err)
	}
	r.trace("latest %q branch commit hash is %s", r.Branch, branchRef.Hash())

	if current.Hash().String() != branchRef.Hash().String() {
		return fmt.Errorf("current commit is not at the tip of the expected branch: %q", r.Branch)
	}

	r.commitHash = current.Hash().String()
	return nil
}
