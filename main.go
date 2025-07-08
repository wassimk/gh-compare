package main

import (
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/wassimk/gh-compare/internal/git"
	"github.com/wassimk/gh-compare/internal/github"
)

func main() {
	repo, err := git.NewRepository(".")
	if err != nil {
		log.Fatal(err)
	}

	compareService := github.NewCompareService(repo)

	args := os.Args[1:]
	url, err := compareService.GenerateCompareURL(args)
	if err != nil {
		log.Fatal(err)
	}

	browser := browser.New("", os.Stdout, os.Stderr)
	if err = browser.Browse(url); err != nil {
		log.Fatal(err)
	}
}
