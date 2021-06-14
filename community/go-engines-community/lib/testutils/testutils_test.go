package testutils_test

import (
	"os"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnvBackup(t *testing.T) {
	Convey("Given a manually set env var", t, func() {
		envname := "CPS_TEST_TESTUTILS"
		envval := "CoolWorld"
		err := os.Setenv(envname, envval)
		So(err, ShouldBeNil)

		Convey("I can create an EnvBackup, save and restore en env var", func() {
			e, err := testutils.NewEnvBackup(envname, "MadWorld")

			So(err, ShouldBeNil)
			So(os.Getenv(envname), ShouldEqual, "MadWorld")

			err = e.Restore()

			So(err, ShouldBeNil)
			So(os.Getenv(envname), ShouldEqual, envval)
		})
	})
}
