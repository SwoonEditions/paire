package paire_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"fmt"
	"os/exec"
	"os"
)

func TestUsingPullBinary(t *testing.T) {
	Convey("Given there is a successful build for current commit", t, func() {
		Convey("When I pull this package from github", func() {
			result, err := exec.Command("./paire/cmd/pull/paire-pull_linux_amd64", "-destination", ".", "-version", "0.3.1").Output()
			assert.Nil(t, err, fmt.Sprintf("There was a problem running the binary: %s", err))
			Convey("Then I should have pre-release downloaded for current commit", func() {
				So(string(result[:]), ShouldContainSubstring, "successfully downloaded release")
				_, fileError := os.Stat("./testdata.package.zip");
				So(fileError, ShouldBeNil)
			})
		})
	})
}
