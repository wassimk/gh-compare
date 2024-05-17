package main

import (
	"fmt"
	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"strings"
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

	if len(os.Args) == 2 {
		compareArgument = os.Args[1]

		if !(strings.Contains(compareArgument, "..") || strings.Contains(compareArgument, "...")) {
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
