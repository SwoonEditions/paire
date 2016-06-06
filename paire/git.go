package paire

import (
	"os/exec"
	"strings"
)

type GitInterface interface {
	CurrentCommit() string
	CurrentTag() string
	CurrentRemotes() string
	Shorten(string) string
}

type Git struct {}

func (git Git) CurrentCommit() string {
	commit, _ := exec.Command("git", "rev-list", "-n", "1", "HEAD").Output()
	return strings.Trim(string(commit[:]), " \n\r\t")
}

func (git Git) CurrentTag() string {
	commit := git.CurrentCommit()
	allTags := git.allTags()
	for _, tagString := range allTags {
		if strings.HasPrefix(tagString, commit) {
			return tagString[strings.LastIndex(tagString, "refs/tags/") + len("refs/tags/"):]
		}
	}
	return git.Shorten(commit)
}

func (git Git) allTags() []string {
	allTags, _ := exec.Command("git", "show-ref", "--tags").Output()
	return strings.Split(string(allTags[:]), "\n")
}

func (git Git) CurrentRemotes() string {
	remotes, _ := exec.Command("git", "remote", "-v").Output()
	return string(remotes[:])
}

func (git Git) Shorten(full string) string {
	return full[0:7]
}
