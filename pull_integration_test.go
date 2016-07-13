package paire_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/SwoonEditions/paire/paire"
	"fmt"
	"os"
)

func TestPullReleaseIntegration(t *testing.T) {
	Convey("Given there already is a successful build for current commit", t, func() {
		git := paire.Git{}
		github := new(paire.Github)
		token, _ := os.LookupEnv("GITHUB_TOKEN")
		githubApi := paire.NewGithubApi(
			github.CurrentRepository(git.CurrentRemotes()),
			token,
		)
		Convey("When I pull this package from github", func() {
			githubRelease, err := githubApi.Download("0.1.0", ".")
			assert.Nil(t, err, fmt.Sprintf("download was not successful %s", err))
			Convey("Then I should have pre-release downloaded for current commit", func() {
				So(githubRelease.PreRelease, ShouldBeFalse)
				So(githubRelease.Packages[0].Downloaded, ShouldBeTrue)
			})
		})
	})
}
