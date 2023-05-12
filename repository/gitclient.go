package repository

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/cicovic-andrija/anduril/service"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepository struct {
	Config
	repo    *git.Repository
	root    string
	tipHash string
	trace   service.TraceCallback
}

func (r *GitRepository) Initialize(root string, trace service.TraceCallback) error {
	r.root = root
	r.trace = trace
	return r.openOrClone()
}

func (r *GitRepository) Sync() (bool, error) {
	if r.Empty() {
		return false, ErrNotInitialized
	}
	return r.pull()
}

func (r *GitRepository) Root() string {
	if r.Empty() {
		return ""
	}
	return r.root
}

func (r *GitRepository) ContentRoot() string {
	if r.Empty() {
		return ""
	}
	return filepath.Join(r.root, r.RelativeContentPath)
}

func (r *GitRepository) LatestRevisionID() string {
	if r.Empty() {
		return ""
	}
	return r.tipHash[:10] // this should be safe
}

func (r *GitRepository) Empty() bool {
	return r.repo == nil
}

// Try to open an existing local repository if there is one, otherwise clone it from remote location.
func (r *GitRepository) openOrClone() error {
	local, err := git.PlainOpen(r.root)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			r.trace("local repository not found on path %s", r.root)
		} else {
			return fmt.Errorf("open local repository %s: operation failed with error: %v", r.root, err)
		}
	} else {
		r.trace("opened local repository on path %s", r.root)
		return r.validateRefs(local)
	}

	clone, err := git.PlainClone(
		r.root,
		false, // bare
		&git.CloneOptions{
			URL:           r.URL(),
			RemoteName:    r.Remote,
			ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
			Progress:      io.Discard,
		},
	)

	if err != nil {
		return fmt.Errorf("clone repository: remote %q: operation failed with error: %v", r.URL(), err)
	}

	r.trace("repository cloned from remote location %q to local path %s", r.URL(), r.root)
	return r.validateRefs(clone)
}

func (r *GitRepository) pull() (new bool, err error) {
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
		err = r.validateRefs(r.repo)
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

func (r *GitRepository) validateRefs(repo *git.Repository) error {
	current, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to obtain a reference to the current commit: %v", err)
	}

	branchRef, err := repo.Reference(plumbing.NewBranchReferenceName(r.Branch), false /* resolved */)
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

	r.repo = repo
	r.tipHash = current.Hash().String()
	return nil
}
