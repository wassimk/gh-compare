package main

import (
	"fmt"
	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
)

func main() {
	var arg string

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else {
		// Get the current branch
		ref, err := repo.Head()
		if err != nil {
			log.Fatal(err)
		}

		arg = ref.Name().Short()
	}

	cfg, err := repo.Config()
	if err != nil {
		log.Fatal(err)
	}
	remoteURL := cfg.Remotes["origin"].URLs[0]

	url := fmt.Sprintf("%s/compare/%s", remoteURL, arg)

	browser := browser.New("", os.Stdout, os.Stderr)

	err = browser.Browse(url)
	if err != nil {
		log.Fatal(err)
	}
}
