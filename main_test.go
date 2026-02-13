package main

import (
	"bytes"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantURL  bool
		wantArgs []string
	}{
		{
			name:     "no flags",
			args:     []string{},
			wantURL:  false,
			wantArgs: []string{},
		},
		{
			name:     "--url flag",
			args:     []string{"--url"},
			wantURL:  true,
			wantArgs: []string{},
		},
		{
			name:     "-u flag",
			args:     []string{"-u"},
			wantURL:  true,
			wantArgs: []string{},
		},
		{
			name:     "flag with positional arg",
			args:     []string{"--url", "main..feature"},
			wantURL:  true,
			wantArgs: []string{"main..feature"},
		},
		{
			name:     "positional arg only",
			args:     []string{"main..feature"},
			wantURL:  false,
			wantArgs: []string{"main..feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := parseFlags(tt.args)

			if opts.printURL != tt.wantURL {
				t.Errorf("printURL = %v, want %v", opts.printURL, tt.wantURL)
			}

			if len(opts.args) != len(tt.wantArgs) {
				t.Fatalf("args length = %d, want %d", len(opts.args), len(tt.wantArgs))
			}

			for i, arg := range opts.args {
				if arg != tt.wantArgs[i] {
					t.Errorf("args[%d] = %q, want %q", i, arg, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestHandleURL(t *testing.T) {
	t.Run("--url prints URL to stdout", func(t *testing.T) {
		var buf bytes.Buffer
		opts := options{printURL: true}

		err := handleURL("https://github.com/owner/repo/compare/feature", opts, &buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := "https://github.com/owner/repo/compare/feature\n"
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}
	})
}
