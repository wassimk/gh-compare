package github

import (
	"strings"
	"testing"

	"github.com/wassimk/gh-compare/internal/git"
)

func TestCompareService_buildCompareRequest(t *testing.T) {
	repo := &git.Repository{
		CurrentBranch: "feature-branch",
		IsForked:      false,
	}

	service := NewCompareService(repo)

	tests := []struct {
		name           string
		args           []string
		expectedFormat string
		expectedBase   string
		expectedHead   string
	}{
		{
			name:           "No arguments",
			args:           []string{},
			expectedFormat: "",
			expectedBase:   "",
			expectedHead:   "feature-branch",
		},
		{
			name:           "Single branch argument",
			args:           []string{"main"},
			expectedFormat: "",
			expectedBase:   "main",
			expectedHead:   "feature-branch",
		},
		{
			name:           "Compare format with two dots",
			args:           []string{"main..feature"},
			expectedFormat: "main..feature",
			expectedBase:   "",
			expectedHead:   "feature-branch",
		},
		{
			name:           "Compare format with three dots",
			args:           []string{"main...feature"},
			expectedFormat: "main...feature",
			expectedBase:   "",
			expectedHead:   "feature-branch",
		},
		{
			name:           "Commit hash comparison",
			args:           []string{"abc123..def456"},
			expectedFormat: "abc123..def456",
			expectedBase:   "",
			expectedHead:   "feature-branch",
		},
	}

	t.Run("Too many arguments", func(t *testing.T) {
		_, err := service.buildCompareRequest([]string{"branch1", "branch2"})
		if err == nil {
			t.Error("Expected error for too many arguments, got nil")
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := service.buildCompareRequest(tt.args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if request.CustomFormat != tt.expectedFormat {
				t.Errorf("Expected CustomFormat %s, got %s", tt.expectedFormat, request.CustomFormat)
			}

			if request.BaseBranch != tt.expectedBase {
				t.Errorf("Expected BaseBranch %s, got %s", tt.expectedBase, request.BaseBranch)
			}

			if request.HeadBranch != tt.expectedHead {
				t.Errorf("Expected HeadBranch %s, got %s", tt.expectedHead, request.HeadBranch)
			}
		})
	}
}

func TestNewCompareService(t *testing.T) {
	repo := &git.Repository{
		CurrentBranch: "test-branch",
	}

	service := NewCompareService(repo)

	if service == nil {
		t.Error("Expected service to be created")
		return
	}

	if service.repo != repo {
		t.Error("Expected service to hold reference to repository")
	}
}

func TestCompareService_GenerateCompareURL_Integration(t *testing.T) {
	repo := &git.Repository{
		CurrentBranch: "feature-branch",
		IsForked:      false,
	}

	service := NewCompareService(repo)

	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "Default behavior",
			args:     []string{},
			contains: "feature-branch",
		},
		{
			name:     "With base branch",
			args:     []string{"main"},
			contains: "main...feature-branch",
		},
		{
			name:     "With custom format",
			args:     []string{"main..feature"},
			contains: "main..feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.GenerateCompareURL(tt.args)

			if err != nil {
				t.Logf("Note: This test may fail if not run in a GitHub repository context")
				return
			}

			if !strings.Contains(url, tt.contains) {
				t.Errorf("Expected URL to contain %s, got %s", tt.contains, url)
			}

			if !strings.HasPrefix(url, "https://") {
				t.Errorf("Expected URL to start with https://, got %s", url)
			}

			if !strings.Contains(url, "/compare/") {
				t.Errorf("Expected URL to contain /compare/, got %s", url)
			}
		})
	}
}
