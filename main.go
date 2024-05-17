package main

import (
	"fmt"
	"os"
	"github.com/cli/go-gh/v2/pkg/browser"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/go-git/go-git/v5"
	"log"
)

func main() {
	var arg string

	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else {
    gitRepo, err := git.PlainOpen(".")
    if err != nil {
      log.Fatal(err)
    }

		ref, err := gitRepo.Head()
		if err != nil {
			log.Fatal(err)
		}

		arg = ref.Name().Short()
	}

	ghRepo, err := repository.Current()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://%s/%s/%s/compare/%s", ghRepo.Host, ghRepo.Owner, ghRepo.Name, arg)

	browser := browser.New("", os.Stdout, os.Stderr)

	err = browser.Browse(url)
	if err != nil {
		log.Fatal(err)
	}
}

