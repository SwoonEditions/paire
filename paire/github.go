package paire

import (
	"strings"
	"regexp"
)

type Repository struct {
	Owner, Repository string
}

type GithubInterface interface {
	CurrentRepository(remotes string) Repository
}

type Github struct {
	Api GithubApiInterface
}

func (github Github) CurrentRepository(allRemotes string) Repository {
	remote := strings.Split(allRemotes, "\n")
	expr := regexp.MustCompilePOSIX("git@github.com:(.+)/(.+).git")
	matches := expr.FindStringSubmatch(remote[0])
	return Repository{
		Owner: (matches[1]),
		Repository: (matches[2]),
	}
}
