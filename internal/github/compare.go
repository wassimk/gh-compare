// Package github provides GitHub-specific functionality for generating compare URLs.
package github

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/wassimk/gh-compare/internal/git"
)

type CompareService struct {
	repo *git.Repository
}

func NewCompareService(repo *git.Repository) *CompareService {
	return &CompareService{
		repo: repo,
	}
}

func (s *CompareService) GenerateCompareURL(args []string) (string, error) {
	compareRequest, err := s.buildCompareRequest(args)
	if err != nil {
		return "", err
	}

	ghRepo, err := repository.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get GitHub repository info: %w", err)
	}

	compareArg := compareRequest.BuildArgument()
	return fmt.Sprintf("https://%s/%s/%s/compare/%s",
		ghRepo.Host, ghRepo.Owner, ghRepo.Name, compareArg), nil
}

func (s *CompareService) buildCompareRequest(args []string) (*git.CompareRequest, error) {
	request := &git.CompareRequest{
		Repository: s.repo,
		HeadBranch: s.repo.CurrentBranch,
	}

	if len(args) > 0 {
		arg := args[0]
		if strings.Contains(arg, "..") {
			request.CustomFormat = arg
		} else {
			request.BaseBranch = arg
		}
	}

	return request, nil
}
