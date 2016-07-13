package paire_test

import (
	"testing"
	"github.com/SwoonEditions/paire/paire"
)

func Test_it_gets_current_repository_from_remotes_string(t *testing.T) {
	gitHub := new(paire.Github)
	repo := gitHub.CurrentRepository(
		"origin	git@github.com:SwoonEditions/paire.git (fetch)\n" +
		"origin	git@github.com:SwoonEditions/paire.git (push)\n",
	)
	if repo.Owner != "SwoonEditions" {
		t.Errorf("Unexpected owner, got '%s'", repo.Owner)
	}
	if repo.Repository != "paire" {
		t.Errorf("Unexpected repository, got '%s'", repo.Repository)
	}
}
