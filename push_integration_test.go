package paire_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/asarturas/paire/paire"
	"os"
	"fmt"
)

func TestPushCurrentCommit(t *testing.T) {
	Convey("Given I have a successfully built release package for current commit", t, func() {
		releasePackages := []paire.ReleasePackage {
			paire.ReleasePackage{
				Name: "testdata/package1.zip",
			},
			paire.ReleasePackage{
				Name: "testdata/package2.zip",
			},
		}
		git := paire.Git{}
		Convey("When I push this package to github", func() {
			release := paire.Release{
				Commit: git.CurrentCommit(),
				Tag: git.CurrentTag(),
				Name:git.CurrentTag(),
				Packages: releasePackages,
			}
			github := new (paire.Github)
			token, _ := os.LookupEnv("GITHUB_TOKEN")
			githubApi := paire.NewGithubApi(
				github.CurrentRepository(git.CurrentRemotes()),
				token,
			)

			githubRelease, err := githubApi.Release(release)
			assert.Nil(t, err, fmt.Sprintf("got an error while releasing %s", err))
			Convey("Then I should have release with my package for current commit", func() {
				So(githubRelease.Pushed, ShouldBeTrue)
				if (git.Shorten(git.CurrentCommit()) == git.CurrentTag()) {
					So(githubRelease.PreRelease, ShouldBeTrue)
				} else {
					So(githubRelease.PreRelease, ShouldBeFalse)
				}
				So(githubRelease.Commit, ShouldEqual, release.Commit)
				assert.True(t, githubRelease.Packages[0].Pushed, "first package doesn't seem to be pushed")
				assert.True(t, githubRelease.Packages[1].Pushed, "second package doesn't seem to be pushed")
			})
			githubApi.Cleanup(release)
		})
	})
}
