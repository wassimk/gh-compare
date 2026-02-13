package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/wassimk/gh-compare/internal/git"
	"github.com/wassimk/gh-compare/internal/github"
)

type options struct {
	printURL bool
	args     []string
}

func parseFlags(args []string) options {
	fs := flag.NewFlagSet("gh-compare", flag.ExitOnError)

	var opts options
	fs.BoolVar(&opts.printURL, "url", false, "Print the compare URL to standard output")
	fs.BoolVar(&opts.printURL, "u", false, "Print the compare URL to standard output")

	fs.Parse(args)
	opts.args = fs.Args()

	return opts
}

func handleURL(url string, opts options, stdout io.Writer) error {
	if opts.printURL {
		fmt.Fprintln(stdout, url)
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

	return handleURL(url, opts, stdout)
}

func main() {
	opts := parseFlags(os.Args[1:])

	if err := run(opts, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
