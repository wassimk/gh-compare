package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantURL  bool
		wantCopy bool
		wantArgs []string
	}{
		{
			name:     "no flags",
			args:     []string{},
			wantURL:  false,
			wantCopy: false,
			wantArgs: []string{},
		},
		{
			name:     "--url flag",
			args:     []string{"--url"},
			wantURL:  true,
			wantCopy: false,
			wantArgs: []string{},
		},
		{
			name:     "-u flag",
			args:     []string{"-u"},
			wantURL:  true,
			wantCopy: false,
			wantArgs: []string{},
		},
		{
			name:     "--copy flag",
			args:     []string{"--copy"},
			wantURL:  false,
			wantCopy: true,
			wantArgs: []string{},
		},
		{
			name:     "-c flag",
			args:     []string{"-c"},
			wantURL:  false,
			wantCopy: true,
			wantArgs: []string{},
		},
		{
			name:     "both -u and -c",
			args:     []string{"-u", "-c"},
			wantURL:  true,
			wantCopy: true,
			wantArgs: []string{},
		},
		{
			name:     "flag with positional arg",
			args:     []string{"--url", "main..feature"},
			wantURL:  true,
			wantCopy: false,
			wantArgs: []string{"main..feature"},
		},
		{
			name:     "copy with positional arg",
			args:     []string{"-c", "main..feature"},
			wantURL:  false,
			wantCopy: true,
			wantArgs: []string{"main..feature"},
		},
		{
			name:     "positional arg only",
			args:     []string{"main..feature"},
			wantURL:  false,
			wantCopy: false,
			wantArgs: []string{"main..feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := parseFlags(tt.args)

			if opts.printURL != tt.wantURL {
				t.Errorf("printURL = %v, want %v", opts.printURL, tt.wantURL)
			}

			if opts.copyToClip != tt.wantCopy {
				t.Errorf("copyToClip = %v, want %v", opts.copyToClip, tt.wantCopy)
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
	testURL := "https://github.com/owner/repo/compare/feature"
	noopClip := func(string) error { return nil }

	t.Run("--url prints URL to stdout", func(t *testing.T) {
		var buf bytes.Buffer
		opts := options{printURL: true}

		err := handleURL(testURL, opts, &buf, noopClip)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := testURL + "\n"
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}
	})

	t.Run("--copy prints URL and calls clipboard func", func(t *testing.T) {
		var buf bytes.Buffer
		var clipped string
		mockClip := func(text string) error {
			clipped = text
			return nil
		}
		opts := options{copyToClip: true}

		err := handleURL(testURL, opts, &buf, mockClip)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := testURL + "\n"
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}

		if clipped != testURL {
			t.Errorf("clipboard text = %q, want %q", clipped, testURL)
		}
	})

	t.Run("both flags prints URL and copies", func(t *testing.T) {
		var buf bytes.Buffer
		var clipped string
		mockClip := func(text string) error {
			clipped = text
			return nil
		}
		opts := options{printURL: true, copyToClip: true}

		err := handleURL(testURL, opts, &buf, mockClip)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := testURL + "\n"
		if buf.String() != want {
			t.Errorf("output = %q, want %q", buf.String(), want)
		}

		if clipped != testURL {
			t.Errorf("clipboard text = %q, want %q", clipped, testURL)
		}
	})

	t.Run("copy error propagates", func(t *testing.T) {
		var buf bytes.Buffer
		failClip := func(string) error {
			return errors.New("clipboard unavailable")
		}
		opts := options{copyToClip: true}

		err := handleURL(testURL, opts, &buf, failClip)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, errors.Unwrap(err)) && err.Error() != "failed to copy to clipboard: clipboard unavailable" {
			t.Errorf("error = %q, want wrapped clipboard error", err)
		}
	})
}
