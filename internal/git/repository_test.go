package git

import (
	"testing"
)

func TestParseRepoOwnerFromURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expected    string
		expectError bool
	}{
		{
			name:        "SSH URL",
			url:         "git@github.com:wassimk/gh-compare.git",
			expected:    "wassimk",
			expectError: false,
		},
		{
			name:        "HTTPS URL",
			url:         "https://github.com/wassimk/gh-compare.git",
			expected:    "wassimk",
			expectError: false,
		},
		{
			name:        "HTTPS URL without .git",
			url:         "https://github.com/wassimk/gh-compare",
			expected:    "wassimk",
			expectError: false,
		},
		{
			name:        "Empty URL",
			url:         "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid URL",
			url:         "not-a-valid-url",
			expected:    "",
			expectError: true,
		},
		{
			name:        "SSH URL without colon",
			url:         "git@github.com",
			expected:    "",
			expectError: true,
		},
		{
			name:        "HTTPS URL with only domain",
			url:         "https://github.com/",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRepoOwnerFromURL(tt.url)

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

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestHasUpstreamRemote(t *testing.T) {
	tests := []struct {
		name     string
		remotes  []Remote
		expected bool
	}{
		{
			name:     "No remotes",
			remotes:  []Remote{},
			expected: false,
		},
		{
			name: "Only origin remote",
			remotes: []Remote{
				{Name: "origin", URL: "git@github.com:wassimk/gh-compare.git"},
			},
			expected: false,
		},
		{
			name: "Has upstream remote",
			remotes: []Remote{
				{Name: "origin", URL: "git@github.com:wassimk/gh-compare.git"},
				{Name: "upstream", URL: "git@github.com:cli/cli.git"},
			},
			expected: true,
		},
		{
			name: "Multiple remotes without upstream",
			remotes: []Remote{
				{Name: "origin", URL: "git@github.com:wassimk/gh-compare.git"},
				{Name: "fork", URL: "git@github.com:other/gh-compare.git"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasUpstreamRemote(tt.remotes)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
