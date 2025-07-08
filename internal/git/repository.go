package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
)

func NewRepository(path string) (*Repository, error) {
	gitRepo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotAGitRepository, err)
	}

	currentBranch, err := getCurrentBranch(gitRepo)
	if err != nil {
		return nil, err
	}

	remotes, err := getRemotes(gitRepo)
	if err != nil {
		return nil, err
	}

	isForked := hasUpstreamRemote(remotes)

	return &Repository{
		Path:          path,
		CurrentBranch: currentBranch,
		Remotes:       remotes,
		IsForked:      isForked,
	}, nil
}

func getCurrentBranch(repo *git.Repository) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoCurrentBranch, err)
	}
	return ref.Name().Short(), nil
}

func getRemotes(repo *git.Repository) ([]Remote, error) {
	gitRemotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoRemoteFound, err)
	}

	var remotes []Remote
	for _, remote := range gitRemotes {
		config := remote.Config()
		if len(config.URLs) > 0 {
			remotes = append(remotes, Remote{
				Name: config.Name,
				URL:  config.URLs[0],
			})
		}
	}

	return remotes, nil
}

func hasUpstreamRemote(remotes []Remote) bool {
	for _, remote := range remotes {
		if remote.Name == "upstream" {
			return true
		}
	}
	return false
}

func parseRepoOwnerFromURL(url string) (string, error) {
	if url == "" {
		return "", ErrInvalidRemoteURL
	}

	if strings.HasPrefix(url, "git@github.com:") {
		parts := strings.Split(url, ":")
		if len(parts) != 2 {
			return "", ErrInvalidRemoteURL
		}
		ownerRepo := strings.TrimSuffix(parts[1], ".git")
		parts = strings.Split(ownerRepo, "/")
		if len(parts) != 2 {
			return "", ErrInvalidRemoteURL
		}
		return parts[0], nil
	}

	if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimPrefix(url, "https://github.com/")
		url = strings.TrimSuffix(url, ".git")
		parts := strings.Split(url, "/")
		if len(parts) < 2 {
			return "", ErrInvalidRemoteURL
		}
		return parts[0], nil
	}

	return "", ErrInvalidRemoteURL
}

