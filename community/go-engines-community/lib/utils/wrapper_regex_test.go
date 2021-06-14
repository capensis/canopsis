package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRegexExpression(t *testing.T) {
	Convey("Given regex that's compatible with builtin regex", t, func() {
		re, err := NewRegexExpression("abc-(?P<sub1>.*)-def-(?P<sub2>\\d+)")
		So(err, ShouldBeNil)
		So(re, ShouldHaveSameTypeAs, WrapperBuiltInRegex{Regexp: nil})
	})

	Convey("Given regex that's compatible with regex2 library", t, func() {
		re, err := NewRegexExpression("^(?!resource_CPU).*$")
		So(err, ShouldBeNil)
		So(re, ShouldHaveSameTypeAs, WrapperRegex2{Regexp: nil})
	})
}

func TestWrapperRegex2_Match(t *testing.T) {
	Convey("Test regex2 match []byte", t, func() {
		re, err := NewRegexExpression("^(?!resource_CPU).*$")
		So(err, ShouldBeNil)
		So(re, ShouldHaveSameTypeAs, WrapperRegex2{Regexp: nil})
		So(re.Match([]byte("resource_CPU")), ShouldBeFalse)
		So(re.Match([]byte("resource_CPU_not_available")), ShouldBeFalse)
		So(re.Match([]byte("not start with resource_CPU")), ShouldBeTrue)
	})
}
