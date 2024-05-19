package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/go-git/go-git/v5"
)

func main() {
	gitRepo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	ref, err := gitRepo.Head()
	if err != nil {
		log.Fatal(err)
	}
	currentBranch := ref.Name().Short()

	var compareArgument string
	var forkedRepo bool
	var originRepo string

	remotes, err := gitRepo.Remotes()
	if err != nil {
		log.Fatal(err)
	}

	for _, remote := range remotes {
		if remote.Config().Name == "upstream" {
			forkedRepo = true
		}

		if remote.Config().Name == "origin" {
			originRepo = parseRepoOriginOwnerFromURL(remote.Config().URLs[0])
		}
	}

	if forkedRepo {
		compareArgument = fmt.Sprintf("main...%s:%s", originRepo, currentBranch)
	} else if len(os.Args) == 2 {
		compareArgument = os.Args[1]

		if !strings.Contains(compareArgument, "..") && !strings.Contains(compareArgument, "...") {
			compareArgument = compareArgument + "..." + currentBranch
		}
	} else {
		compareArgument = currentBranch
	}

	ghRepo, err := repository.Current()
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("https://%s/%s/%s/compare/%s", ghRepo.Host, ghRepo.Owner, ghRepo.Name, compareArgument)

	browser := browser.New("", os.Stdout, os.Stderr)

	if err = browser.Browse(url); err != nil {
		log.Fatal(err)
	}
}

func parseRepoOriginOwnerFromURL(url string) string {
	parts := strings.Split(url, ":")
	return strings.Split(parts[1], "/")[0]
}
