package paire_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/asarturas/paire/paire"
	"fmt"
)

func TestPullRelease(t *testing.T) {
	Convey("Given there is a successful build for current commit", t, func() {
		releasePackages := []paire.ReleasePackage{
			paire.ReleasePackage{
				Name: "package1.zip",
				Downloaded: true,
			},
			paire.ReleasePackage{
				Name: "package2.zip",
				Downloaded: true,
			},
		}
		git := paire.GitMock{}
		git.On("CurrentCommit").Return("sha123456")
		git.On("CurrentTag").Return("sha1234")
		release := paire.Release{
			Commit: git.CurrentCommit(),
			Tag: git.CurrentTag(),
			Name: git.CurrentTag(),
			Packages: releasePackages,
			PreRelease: true,
			Id: 13,
		}
		githubApi := new(paire.GithubApiMock)
		githubApi.On("Download", "sha1234").Return(release, nil)
		Convey("When I pull this package from github", func() {
			githubRelease, err := githubApi.Download(git.CurrentTag(), ".")
			assert.Nil(t, err, fmt.Sprintf("download was not successful %s", err))
			Convey("Then I should have pre-release downloaded for current commit", func() {
				So(githubRelease.Packages[0].Downloaded, ShouldBeTrue)
				So(githubRelease.Packages[1].Downloaded, ShouldBeTrue)
			})
		})
	})
}
