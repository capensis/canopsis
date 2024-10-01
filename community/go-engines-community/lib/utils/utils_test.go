package utils_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/kylelemons/godebug/pretty"
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

		s = "こんにちは"
		Convey("Calling TruncateString where chars > 0 should truncate unicode string", func() {
			So(utils.TruncateString(s, 3), ShouldEqual, "こんに")
		})
	})
}

func BenchmarkTruncateStringUnicode(b *testing.B) {
	s := "こんにちは"
	for i := 0; i < b.N; i++ {
		utils.TruncateString(s, 3)
	}
}

func BenchmarkTruncateStringShort(b *testing.B) {
	s := "string to truncate"
	for i := 0; i < b.N; i++ {
		utils.TruncateString(s, 30)
	}
}

func BenchmarkTruncateStringASCII(b *testing.B) {
	s := "string to truncate"
	for i := 0; i < b.N; i++ {
		utils.TruncateString(s, 6)
	}
}

var findAllStringSubmatchMapCases = []struct {
	name     string
	regex    string
	input    string
	expected []map[string]string
}{
	{
		name:  "empty input",
		regex: "abc-(?P<sub1>.*)-def-(?P<sub2>\\d+)",
		input: "",
	},
	{
		name:  "regex with no subexpressions",
		regex: "abc-test-def-123",
		input: "abc-test-def-123",
	},
	{
		name:  "regex with one subexpression",
		regex: "abc-(?P<sub1>.*)-def-123",
		input: "abc-test-def-123",
		expected: []map[string]string{
			{
				"sub1": "test",
			},
		},
	},
	{
		name:  "regex with two subexpressions",
		regex: "abc-(?P<sub1>.*)-def-(?P<sub2>\\d+)",
		input: "abc-test-def-123",
		expected: []map[string]string{
			{
				"sub1": "test",
				"sub2": "123",
			},
		},
	},
	{
		name:  "regex with two subexpressions and multiple matches",
		regex: "abc-(?P<sub1>.*?)-def-(?P<sub2>\\d+)",
		input: "abc-test-def-123 abc-test-def-456",
		expected: []map[string]string{
			{
				"sub1": "test",
				"sub2": "123",
			},
			{
				"sub1": "test",
				"sub2": "456",
			},
		},
	},
	{
		name:  "regex with unnamed group",
		regex: "(name)(a)(name)(b)",
		input: "nameanameb",
		expected: []map[string]string{
			{
				"1": "name",
				"2": "a",
				"3": "name",
				"4": "b",
			},
		},
	},
	{
		name:  "regexp2 empty input",
		regex: "(?!ce)?abc-(?P<sub1>.*)-def-(?P<sub2>\\d+)",
		input: "",
	},
	{
		name:  "regexp2 with no subexpressions",
		regex: "(?!ce)?abc-test-def-123",
		input: "abc-test-def-123",
	},
	{
		name:  "regexp2 with one subexpression",
		regex: "(?!ce)?abc-(?P<sub1>.*)-def-123",
		input: "abc-test-def-123",
		expected: []map[string]string{
			{
				"sub1": "test",
			},
		},
	},
	{
		name:  "regexp2 with two subexpressions",
		regex: "(?!ce)?abc-(?P<sub1>.*)-def-(?P<sub2>\\d+)",
		input: "abc-test-def-123",
		expected: []map[string]string{
			{
				"sub1": "test",
				"sub2": "123",
			},
		},
	},
	{
		name:  "regexp2 with two subexpressions and multiple matches",
		regex: "(?!ce)?abc-(?P<sub1>.*?)-def-(?P<sub2>\\d+)",
		input: "abc-test-def-123 abc-test-def-456",
		expected: []map[string]string{
			{
				"sub1": "test",
				"sub2": "123",
			},
			{
				"sub1": "test",
				"sub2": "456",
			},
		},
	},
	{
		name:  "regexp2 with unnamed group",
		regex: "(?!ce)(name)(a)(name)(b)",
		input: "nameanameb",
		expected: []map[string]string{
			{
				"1": "name",
				"2": "a",
				"3": "name",
				"4": "b",
			},
		},
	},
}

func TestFindAllStringSubmatchMapWithRegexExpression(t *testing.T) {
	for _, tc := range findAllStringSubmatchMapCases {
		t.Run(tc.name, func(t *testing.T) {
			re, err := utils.NewRegexExpression(tc.regex)
			if err != nil {
				t.Errorf("error creating regex expression: %s", err)
			}
			actual := utils.FindAllStringSubmatchMapWithRegexExpression(re, tc.input)
			if diff := pretty.Compare(tc.expected, actual); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}

// Benchmark test for FindStringSubmatchMapWithRegexExpression
func BenchmarkFindStringSubmatchMapWithRegexExpression(b *testing.B) {
	const input = "abc-test-def-123 abc-test-def-456"
	re, err := utils.NewRegexExpression("abc-(?P<sub1>.*?)-def-(?P<sub2>\\d+)")
	if err != nil {
		b.Errorf("error creating regex expression: %s", err)
	}
	re2, err := utils.NewRegexExpression("(?!ba)abc-(?P<sub1>.*?)-def-(?P<sub2>\\d+)")
	if err != nil {
		b.Errorf("error creating regex expression: %s", err)
	}
	findAllStringSubmatchMap := func(b *testing.B, re utils.RegexExpression, input string) {
		b.Helper()
		utils.FindAllStringSubmatchMapWithRegexExpression(re, input)
	}
	findStringSubmatchMap := func(b *testing.B, re utils.RegexExpression, input string) {
		b.Helper()
		utils.FindStringSubmatchMapWithRegexExpression(re, input)
	}
	benchmarkCases := []struct {
		name string
		re   utils.RegexExpression
		all  bool
		fn   func(*testing.B, utils.RegexExpression, string)
	}{
		{
			name: "builtin single",
			re:   re,
			all:  false,
			fn:   findStringSubmatchMap,
		},
		{
			name: "builtin all",
			re:   re,
			all:  true,
			fn:   findAllStringSubmatchMap,
		},
		{
			name: "regexp2 single",
			re:   re2,
			all:  false,
			fn:   findStringSubmatchMap,
		},
		{
			name: "regexp2 all",
			re:   re2,
			all:  true,
			fn:   findAllStringSubmatchMap,
		},
	}
	for _, tc := range benchmarkCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tc.fn(b, tc.re, input)
			}
		})
	}
}
