package paire_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"fmt"
	"os/exec"
	"github.com/asarturas/paire/paire"
	"os"
)

func TestUsingPushBinary(t *testing.T) {
	Convey("Given I have a binary built for current commit", t, func() {
		pkg1 := "testdata/package1.zip"
		pkg2 := "testdata/package2.zip"
		Convey("When I push this package to github", func() {
			result, err := exec.Command("./paire/cmd/push/paire-push_linux_amd64", "-package", pkg1, "-package", pkg2).Output()
			assert.Nil(t, err, fmt.Sprintf("There was a problem running the binary: %s", err))
			Convey("Then I should have pre-release with my package for current commit", func() {
				So(string(result[:]), ShouldContainSubstring, "successfully released")
			})
			git := paire.Git{}
			release := paire.Release{
				Commit: git.CurrentCommit(),
				Tag: git.CurrentTag(),
				Name:git.CurrentTag(),
			}
			github := new (paire.Github)
			token, _ := os.LookupEnv("GITHUB_TOKEN")
			githubApi := paire.NewGithubApi(
				github.CurrentRepository(git.CurrentRemotes()),
				token,
			)
			githubApi.Cleanup(release)
		})
	})
}
