package errt_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrT(t *testing.T) {
	Convey("Setup", t, func() {
		Convey("DocumentNotFound", func() {
			err := errt.NewNotFound(nil)
			So(err, ShouldBeNil)

			referr := errors.New("not found")
			err = errt.NewNotFound(referr)

			ednf, ok := err.(errt.NotFound)
			So(ok, ShouldBeTrue)
			_, ok = err.(errt.IOError)
			So(ok, ShouldBeFalse)

			So(ednf.Error(), ShouldEqual, "not found")
			So(ednf.Err(), ShouldEqual, referr)
		})

		Convey("IOError", func() {
			err := errt.NewIOError(nil)
			So(err, ShouldBeNil)

			referr := errors.New("io error")
			err = errt.NewIOError(referr)

			eio, ok := err.(errt.IOError)
			So(ok, ShouldBeTrue)
			_, ok = err.(errt.NotFound)
			So(ok, ShouldBeFalse)

			So(eio.Error(), ShouldEqual, "io error")
			So(eio.Err(), ShouldEqual, referr)
		})
	})
}
