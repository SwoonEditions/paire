package paire_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/SwoonEditions/paire/paire"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestPushSimpleCommit(t *testing.T) {
	Convey("Given I have a successfully built release package for current commit", t, func() {
		releasePackages := []paire.ReleasePackage {
			paire.ReleasePackage{
				Name: "package1.zip",
			},
			paire.ReleasePackage{
				Name: "package2.zip",
			},
		}
		git := paire.GitMock{}
		git.On("CurrentCommit").Return("sha123456")
		git.On("CurrentTag").Return("sha1234")
		Convey("When I push this package to github", func() {
			release := paire.Release{
				Commit: git.CurrentCommit(),
				Tag: git.CurrentTag(),
				Name:git.CurrentTag(),
				Packages: releasePackages,
			}
			githubApi := new(paire.GithubApiMock)
			githubApi.
				On("Release", release).
				Return(
					paire.Release{
						Commit: release.Commit,
						Tag: release.Tag,
						Name: release.Name,
						Packages: []paire.ReleasePackage{
							paire.ReleasePackage{
								Name: "package1.zip",
								Pushed: true,
							},
							paire.ReleasePackage{
								Name: "package2.zip",
								Pushed: true,
							},
						},
						Pushed: true,
						PreRelease: true,
						Id: 12,
					},
					nil,
				)

			githubRelease, err := githubApi.Release(release)
			assert.Nil(t, err, fmt.Sprintf("release was not successful %s", err))
			Convey("Then I should have pre-release with my package for current commit", func() {
				So(githubRelease.Pushed, ShouldBeTrue)
				So(githubRelease.PreRelease, ShouldBeTrue)
				So(githubRelease.Commit, ShouldEqual, release.Commit)
				So(githubRelease.Packages[0].Pushed, ShouldBeTrue)
				So(githubRelease.Packages[1].Pushed, ShouldBeTrue)
			})
		})
	})
}
