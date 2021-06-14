package utils_test

import (
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFailOnError(t *testing.T) {
	Convey("Setup", t, func() {
		So(func() { utils.FailOnError(nil, "nope") }, ShouldNotPanic)
		So(func() { utils.FailOnError(fmt.Errorf("hargh"), "panic") }, ShouldPanicWith, "panic: hargh")
	})
}
