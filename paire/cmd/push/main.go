package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/asarturas/paire/paire"
	"strings"
)

func main() {
	var pkges pkge
	flag.Var(&pkges, "package", "specify at least one package file (example: -package one.zip -package another.zip)")
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	flag.Parse()
	if len(pkges) == 0 {
		fmt.Println("Package is required (see -package flag)")
		os.Exit(1)
	}
	if !exists {
		fmt.Println("Github token should be set with $GITHUB_TOKEN environment variable")
		os.Exit(1)
	}
	releasePackages := []paire.ReleasePackage{}
	for _, pkg := range pkges {
		if _, err := os.Stat(pkg); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		releasePackages = append(releasePackages, paire.ReleasePackage{Name: pkg})
	}
	git := new(paire.Git)
	release := paire.Release{
		Commit: git.CurrentCommit(),
		Tag: git.CurrentTag(),
		Name: git.CurrentTag(),
		Packages: releasePackages,
	}
	github := new(paire.Github)
	githubApi := paire.NewGithubApi(github.CurrentRepository(git.CurrentRemotes()), token)
	githubRelease, err := githubApi.Release(release)
	if err != nil {
		fmt.Printf("there was a problem with release: %s", err)
		os.Exit(1)
	} else {
		fmt.Println("successfully released: " + githubRelease.Name)
		os.Exit(0)
	}
}

type pkge []string

func (i *pkge) String() string {
	return fmt.Sprint(*i)
}

func (i *pkge) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		*i = append(*i, dt)
	}
	return nil
}
