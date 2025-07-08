package git

import (
	"testing"
)

func TestRepository_HasRemote(t *testing.T) {
	repo := &Repository{
		Remotes: []Remote{
			{Name: "origin", URL: "git@github.com:wassimk/gh-compare.git"},
			{Name: "upstream", URL: "git@github.com:cli/cli.git"},
		},
	}

	tests := []struct {
		name       string
		remoteName string
		expected   bool
	}{
		{
			name:       "Existing remote",
			remoteName: "origin",
			expected:   true,
		},
		{
			name:       "Another existing remote",
			remoteName: "upstream",
			expected:   true,
		},
		{
			name:       "Non-existing remote",
			remoteName: "nonexistent",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.HasRemote(tt.remoteName)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRepository_GetRemote(t *testing.T) {
	repo := &Repository{
		Remotes: []Remote{
			{Name: "origin", URL: "git@github.com:wassimk/gh-compare.git"},
			{Name: "upstream", URL: "git@github.com:cli/cli.git"},
		},
	}

	tests := []struct {
		name        string
		remoteName  string
		expectedURL string
		expectError bool
	}{
		{
			name:        "Get existing remote",
			remoteName:  "origin",
			expectedURL: "git@github.com:wassimk/gh-compare.git",
			expectError: false,
		},
		{
			name:        "Get non-existing remote",
			remoteName:  "nonexistent",
			expectedURL: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote, err := repo.GetRemote(tt.remoteName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if remote.URL != tt.expectedURL {
				t.Errorf("Expected URL %s, got %s", tt.expectedURL, remote.URL)
			}
		})
	}
}

func TestCompareRequest_BuildArgument(t *testing.T) {
	tests := []struct {
		name     string
		request  *CompareRequest
		expected string
	}{
		{
			name: "Custom format takes precedence",
			request: &CompareRequest{
				CustomFormat: "feature...main",
				BaseBranch:   "other",
				HeadBranch:   "head",
				Repository: &Repository{
					CurrentBranch: "current",
					IsForked:      false,
				},
			},
			expected: "feature...main",
		},
		{
			name: "Forked repository",
			request: &CompareRequest{
				Repository: &Repository{
					CurrentBranch: "feature-branch",
					IsForked:      true,
					Remotes: []Remote{
						{Name: "origin", URL: "git@github.com:wassimk/gh-compare.git"},
					},
				},
			},
			expected: "main...wassimk:feature-branch",
		},
		{
			name: "Base branch specified",
			request: &CompareRequest{
				BaseBranch: "develop",
				HeadBranch: "feature",
				Repository: &Repository{
					CurrentBranch: "current",
					IsForked:      false,
				},
			},
			expected: "develop...feature",
		},
		{
			name: "Default behavior",
			request: &CompareRequest{
				Repository: &Repository{
					CurrentBranch: "current-branch",
					IsForked:      false,
				},
			},
			expected: "current-branch",
		},
		{
			name: "Forked repository with invalid origin",
			request: &CompareRequest{
				Repository: &Repository{
					CurrentBranch: "feature-branch",
					IsForked:      true,
					Remotes: []Remote{
						{Name: "origin", URL: "invalid-url"},
					},
				},
			},
			expected: "feature-branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.BuildArgument()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
