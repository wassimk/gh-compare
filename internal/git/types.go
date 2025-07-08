package git

import (
	"errors"
	"fmt"
)

var (
	ErrNotAGitRepository = errors.New("not a git repository")
	ErrNoCurrentBranch   = errors.New("no current branch found")
	ErrNoRemoteFound     = errors.New("no remote found")
	ErrInvalidRemoteURL  = errors.New("invalid remote URL")
)

type Repository struct {
	Path          string
	CurrentBranch string
	Remotes       []Remote
	IsForked      bool
}

type Remote struct {
	Name string
	URL  string
}

type CompareRequest struct {
	BaseBranch   string
	HeadBranch   string
	Repository   *Repository
	CustomFormat string
}

func (r *Repository) HasRemote(name string) bool {
	for _, remote := range r.Remotes {
		if remote.Name == name {
			return true
		}
	}
	return false
}

func (r *Repository) GetRemote(name string) (*Remote, error) {
	for _, remote := range r.Remotes {
		if remote.Name == name {
			return &remote, nil
		}
	}
	return nil, fmt.Errorf("remote '%s' not found", name)
}

func (r *Repository) GetOriginOwner() (string, error) {
	originRemote, err := r.GetRemote("origin")
	if err != nil {
		return "", err
	}
	return parseRepoOwnerFromURL(originRemote.URL)
}

func (c *CompareRequest) BuildArgument() string {
	if c.CustomFormat != "" {
		return c.CustomFormat
	}

	if c.Repository.IsForked {
		owner, err := c.Repository.GetOriginOwner()
		if err != nil {
			return c.Repository.CurrentBranch
		}
		return fmt.Sprintf("main...%s:%s", owner, c.Repository.CurrentBranch)
	}

	if c.BaseBranch != "" {
		return fmt.Sprintf("%s...%s", c.BaseBranch, c.HeadBranch)
	}

	return c.Repository.CurrentBranch
}

