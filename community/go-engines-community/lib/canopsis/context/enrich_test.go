package context_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildEnrichOnFields(t *testing.T) {
	Convey("Setup", t, func() {
		Convey("Includes only", func() {
			include := "sauron,nazgul,saruman"
			exclude := "frodo,galadriel,boromir"

			r := context.NewEnrichFields(include, exclude)
			So(r.Allow("sauron"), ShouldBeTrue)
			So(r.Allow("nazgul"), ShouldBeTrue)
			So(r.Allow("saruman"), ShouldBeTrue)
			So(r.Allow("frodo"), ShouldBeFalse)
			So(r.Allow("galadriel"), ShouldBeFalse)
			So(r.Allow("boromir"), ShouldBeFalse)

			// check unlisted
			So(r.Allow("gollum"), ShouldBeFalse)

			// add previously unlisted
			r.AddInclude("gollum")
			So(r.Allow("gollum"), ShouldBeTrue)
		})

		Convey("Excludes only", func() {
			include := ""
			exclude := "frodo,galadriel,boromir"

			r := context.NewEnrichFields(include, exclude)
			So(r.Allow("sauron"), ShouldBeTrue)
			So(r.Allow("nazgul"), ShouldBeTrue)
			So(r.Allow("saruman"), ShouldBeTrue)
			So(r.Allow("frodo"), ShouldBeFalse)
			So(r.Allow("galadriel"), ShouldBeFalse)
			So(r.Allow("boromir"), ShouldBeFalse)

			// check unlisted
			So(r.Allow("gollum"), ShouldBeTrue)

			r.AddExclude("gollum")
			So(r.Allow("gollum"), ShouldBeFalse)
		})
	})
}
