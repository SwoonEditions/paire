package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/SwoonEditions/paire/paire"
)

func main() {
	destination := flag.String("destination", ".", "specify where to store assets (defaults to current directory)")
	targetTag := flag.String("version", "", "specify version to download (default to sha of current commit)")
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	flag.Parse()
	if len(*destination) == 0 {
		fmt.Println("Destination is required (see -destination flag)")
		os.Exit(1)
	}
	if !exists {
		fmt.Println("Github token should be set with $GITHUB_TOKEN environment variable")
		os.Exit(1)
	}
	git := new(paire.Git)
	if len(*targetTag) == 0 {
		*targetTag = git.CurrentTag()
	}
	github := new(paire.Github)
	githubApi := paire.NewGithubApi(github.CurrentRepository(git.CurrentRemotes()), token)
	githubRelease, err := githubApi.Download(*targetTag, *destination)
	if err != nil {
		fmt.Printf("there was a problem with download: %s", err)
		os.Exit(1)
	} else {
		fmt.Printf("successfully downloaded release: %s", githubRelease.Name)
		os.Exit(0)
	}
}