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

			var notFoundErr errt.NotFound
			ok := errors.As(err, &notFoundErr)
			So(ok, ShouldBeTrue)
			var ioError errt.IOError
			ok = errors.As(err, &ioError)
			So(ok, ShouldBeFalse)

			So(notFoundErr.Error(), ShouldEqual, "not found")
			So(notFoundErr.Err(), ShouldEqual, referr)
		})

		Convey("IOError", func() {
			err := errt.NewIOError(nil)
			So(err, ShouldBeNil)

			referr := errors.New("io error")
			err = errt.NewIOError(referr)

			var eio errt.IOError
			ok := errors.As(err, &eio)
			So(ok, ShouldBeTrue)
			var notFoundErr errt.NotFound
			ok = errors.As(err, &notFoundErr)
			So(ok, ShouldBeFalse)

			So(eio.Error(), ShouldEqual, "io error")
			So(eio.Err(), ShouldEqual, referr)
		})
	})
}
