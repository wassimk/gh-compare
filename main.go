package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/wassimk/gh-compare/internal/clipboard"
	"github.com/wassimk/gh-compare/internal/git"
	"github.com/wassimk/gh-compare/internal/github"
)

type options struct {
	printURL   bool
	copyToClip bool
	args       []string
}

func parseFlags(args []string) options {
	fs := flag.NewFlagSet("gh-compare", flag.ExitOnError)

	var opts options
	fs.BoolVar(&opts.printURL, "url", false, "Print the compare URL to standard output")
	fs.BoolVar(&opts.printURL, "u", false, "Print the compare URL to standard output")
	fs.BoolVar(&opts.copyToClip, "copy", false, "Copy the compare URL to the clipboard")
	fs.BoolVar(&opts.copyToClip, "c", false, "Copy the compare URL to the clipboard")

	fs.Parse(args)
	opts.args = fs.Args()

	return opts
}

func handleURL(url string, opts options, stdout io.Writer, clipFn func(string) error) error {
	if opts.printURL || opts.copyToClip {
		fmt.Fprintln(stdout, url)
	}

	if opts.copyToClip {
		if err := clipFn(url); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
	}

	if opts.printURL || opts.copyToClip {
		return nil
	}

	b := browser.New("", os.Stdout, os.Stderr)
	return b.Browse(url)
}

func run(opts options, stdout io.Writer) error {
	repo, err := git.NewRepository(".")
	if err != nil {
		return err
	}

	compareService := github.NewCompareService(repo)

	url, err := compareService.GenerateCompareURL(opts.args)
	if err != nil {
		return err
	}

	return handleURL(url, opts, stdout, clipboard.Write)
}

func main() {
	opts := parseFlags(os.Args[1:])

	if err := run(opts, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
