package anduril

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type RepositoryProcessor struct {
	RepositoryConfig
	repo     *git.Repository
	repoRoot string
	tipHash  string
	trace    TraceCallback
}

func (r *RepositoryProcessor) OpenOrClone(path string) error {
	// Try to open an existing local repository if there is one, otherwise clone it from remote location.
	local, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			r.trace("local repository not found on path %s", path)
		} else {
			return fmt.Errorf("open local repository %s: operation failed with error: %v", path, err)
		}
	} else {
		r.trace("opened local repository on path %s", path)
		r.repo = local
		r.repoRoot = path
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
	r.repoRoot = path
	return r.ValidateState()
}

func (r *RepositoryProcessor) Pull() (new bool, err error) {
	w, err := r.repo.Worktree()
	if err != nil {
		err = fmt.Errorf("pull: failed to obtain worktree: %v", err)
		return
	}

	err = w.Pull(&git.PullOptions{RemoteName: r.Remote})
	switch err {
	case nil:
		r.trace("pull: successfully fetched and merged latest changes")
		new = true
		err = r.ValidateState()
		return
	case git.NoErrAlreadyUpToDate:
		r.trace("pull: already up-to-date")
		new = false
		err = nil
		return
	default:
		new = false
		err = fmt.Errorf("pull operation failed: %v", err)
		return
	}
}

func (r *RepositoryProcessor) ValidateState() error {
	current, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to obtain a reference to the current commit: %v", err)
	}

	branchRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(r.Branch), false /* resolved */)
	if err != nil {
		return fmt.Errorf("failed to obtain a reference to the %q branch: %v", r.Branch, err)
	}

	r.trace(
		"current checked out commit hash is %s, and latest %q branch commit hash is %s",
		current.Hash(),
		r.Branch,
		branchRef.Hash(),
	)

	if current.Hash().String() != branchRef.Hash().String() {
		return fmt.Errorf("current commit is not at the tip of the expected branch: %q", r.Branch)
	}

	r.tipHash = current.Hash().String()
	return nil
}

func (r *RepositoryProcessor) Root() string {
	return r.repoRoot
}

func (r *RepositoryProcessor) ContentRoot() string {
	return filepath.Join(r.repoRoot, r.RelativeContentPath)
}

func (r *RepositoryProcessor) LatestCommitShortHash() string {
	return r.tipHash[:10] // this should be safe
}
