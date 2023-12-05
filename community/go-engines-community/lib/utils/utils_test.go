package utils_test

import (
	"regexp"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
)

type Station struct {
	Name       string
	Faction    string
	Population int
}

func TestFindStringSubmatchMapWithRegexExpression(t *testing.T) {
	Convey("Given regex that's compatible with builtin regex", t, func() {
		re, err := utils.NewRegexExpression("abc-(?P<sub1>.*)-def-(?P<sub2>\\d+)")
		So(err, ShouldBeNil)
		So(re, ShouldHaveSameTypeAs, utils.WrapperBuiltInRegex{Regexp: nil})

		Convey("A map containing the values of the subexpressions is returned for strings that match the regex", func() {
			match := utils.FindStringSubmatchMapWithRegexExpression(re, "abc-test-def-123")
			So(match, ShouldNotBeNil)
			So(match["sub1"], ShouldEqual, "test")
			So(match["sub2"], ShouldEqual, "123")
		})

		Convey("nil is returned for strings that do not match the regex", func() {
			match := utils.FindStringSubmatchMapWithRegexExpression(re, "abc-test-def-d23")
			So(match, ShouldBeNil)
		})
	})

	Convey("Given regex that's compatible with regex2 library", t, func() {
		re, err := utils.NewRegexExpression("^(?!resource_CPU).*$")
		So(err, ShouldBeNil)
		So(re, ShouldHaveSameTypeAs, utils.WrapperRegex2{Regexp: nil})

		Convey("nil is returned for strings that do not match the regex2", func() {
			match := utils.FindStringSubmatchMapWithRegexExpression(re, "resource_RAM")
			So(match, ShouldBeEmpty)
		})

		Convey("empty map is returned for strings that matches the regex2", func() {
			match := utils.FindStringSubmatchMapWithRegexExpression(re, "resource_CPU")
			So(match, ShouldBeNil)
		})
	})

	Convey("Given regex that has unnamed group", t, func() {
		re, err := utils.NewRegexExpression(`(?P<name1>name)(a)(?P<name2>name)(b)`)
		So(err, ShouldBeNil)
		So(re, ShouldHaveSameTypeAs, utils.WrapperBuiltInRegex{Regexp: nil})

		Convey("the unnamed group should replace with index", func() {
			match := utils.FindStringSubmatchMapWithRegexExpression(re, "nameanameb")
			So(match, ShouldContainKey, "1")
			So(match, ShouldContainKey, "2")
			So(match["1"], ShouldEqual, "a")
			So(match["2"], ShouldEqual, "b")
		})
	})
}

func TestFindStringSubmatchMap(t *testing.T) {
	Convey("Given a regular expression", t, func() {
		re := regexp.MustCompile(`abc-(?P<sub1>.*)-def-(?P<sub2>\d+)`)

		Convey("A map containing the values of the subexpressions is returned for strings that match the regex", func() {
			match := utils.FindStringSubmatchMap(re, "abc-test-def-123")
			So(match, ShouldNotBeNil)
			So(match["sub1"], ShouldEqual, "test")
			So(match["sub2"], ShouldEqual, "123")
		})

		Convey("nil is returned for strings that do not match the regex", func() {
			match := utils.FindStringSubmatchMap(re, "abc-test-def-d23")
			So(match, ShouldBeNil)
		})
	})
}

func TestGetStringField(t *testing.T) {
	Convey("Given an alarm", t, func() {
		s := Station{"Loubyanka", "Red", 1234}

		Convey("We can read a field", func() {
			So(utils.GetStringField(s, "Faction"), ShouldEqual, "Red")
		})
		Convey("We cannot read an int field", func() {
			So(utils.GetStringField(s, "Population"), ShouldEqual, "<int Value>")
		})
		Convey("We cannot read a missing field", func() {
			So(utils.GetStringField(s, "Evolution"), ShouldEqual, "")
		})
	})
}

func TestAsString(t *testing.T) {
	Convey("Given a string", t, func() {
		s := "a string"

		Convey("Calling AsString on this string should succeed and return its value", func() {
			value, success := utils.AsString(s)
			So(success, ShouldBeTrue)
			So(value, ShouldEqual, "a string")
		})

		Convey("Calling AsString on a reference to this string should succeed and return its value", func() {
			value, success := utils.AsString(&s)
			So(success, ShouldBeTrue)
			So(value, ShouldEqual, "a string")
		})
	})

	Convey("Given an integer", t, func() {
		i := 3

		Convey("Calling AsString on this integer should fail", func() {
			_, success := utils.AsString(i)
			So(success, ShouldBeFalse)
		})

		Convey("Calling AsString on a reference to this integer should fail", func() {
			_, success := utils.AsString(&i)
			So(success, ShouldBeFalse)
		})
	})
}

func TestTruncateString(t *testing.T) {
	Convey("Given a string", t, func() {
		s := "string to truncate"

		Convey("Calling TruncateString where chars < 1 shouldn't truncate string ", func() {
			So(utils.TruncateString(s, 0), ShouldEqual, s)
			So(utils.TruncateString(s, -1), ShouldEqual, s)
		})

		Convey("Calling TruncateString where chars > 0 should truncate string", func() {
			So(utils.TruncateString(s, 6), ShouldEqual, "string")
		})
	})
}
