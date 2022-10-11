package template_test

import (
	"bytes"
	"io"
	"testing"
	"text/template"
	"time"

	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func executeTemplate(tmpl *template.Template, payload interface{}) (string, error) {
	var b bytes.Buffer
	err := tmpl.Execute(io.Writer(&b), payload)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func TestTemplate_Split(t *testing.T) {
	Convey("template use function should be good", t, func() {
		Convey("split function", func() {
			templateText := `{{ split .Sep .Index .Input }}`
			tmpl, err := template.New("func-split-test").Funcs(libtemplate.GetFunctions(nil)).Parse(templateText)
			if err != nil {
				t.Fatalf("parsing: %s", err)
			}
			testCase := struct {
				Sep   string
				Index int
				Input string
			}{
				Sep: ",", Index: 0, Input: "NgocHa,MinhNghia,Minh",
			}
			output, err := executeTemplate(tmpl, testCase)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "NgocHa")

			testCase.Index = -1
			output, err = executeTemplate(tmpl, testCase)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "")

			testCase.Index = 1
			testCase.Sep = " "
			testCase.Input = "This is space"
			output, err = executeTemplate(tmpl, testCase)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "is")
		})

		Convey("trim function", func() {
			templateText := `{{ trim . }}`
			tmpl, err := template.New("func-trim-test").Funcs(libtemplate.GetFunctions(nil)).Parse(templateText)
			if err != nil {
				t.Fatalf("parsing: %s", err)
			}

			output, err := executeTemplate(tmpl, "  ")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "")

			output, err = executeTemplate(tmpl, " kratos ")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "kratos")

			output, err = executeTemplate(tmpl, "\tkratos\n ")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "kratos")
		})

		Convey("upper function", func() {
			templateText := `{{ uppercase . }}`
			tmpl, err := template.New("func-upper-test").Funcs(libtemplate.GetFunctions(nil)).Parse(templateText)
			if err != nil {
				t.Fatalf("parsing: %s", err)
			}

			output, err := executeTemplate(tmpl, "  ")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "  ")

			output, err = executeTemplate(tmpl, "kratos")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "KRATOS")

			output, err = executeTemplate(tmpl, "KraTos")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "KRATOS")
		})

		Convey("lowercase function", func() {
			templateText := `{{ lowercase . }}`
			tmpl, err := template.New("func-lower-test").Funcs(libtemplate.GetFunctions(nil)).Parse(templateText)
			if err != nil {
				t.Fatalf("parsing: %s", err)
			}

			output, err := executeTemplate(tmpl, "  ")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "  ")

			output, err = executeTemplate(tmpl, "kratos")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "kratos")

			output, err = executeTemplate(tmpl, "KraTos")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "kratos")

			output, err = executeTemplate(tmpl, "KRATOS")
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "kratos")
		})
	})
}

func TestTemplate_localtime(t *testing.T) {
	Convey("template use function should be good", t, func() {
		Convey("localtime function", func() {
			templateText := `{{ .TestDate | localtime "Mon, 02 Jan 2006 15:04:05 MST" "Australia/Queensland" }}`
			tmpl, err := template.New("func-localtime-test").Funcs(libtemplate.GetFunctions(nil)).Parse(templateText)
			if err != nil {
				t.Fatalf("parsing: %s", err)
			}
			testCase := struct {
				TestDate types.CpsTime
			}{
				TestDate: types.CpsTime{
					Time: time.Date(2021, time.October, 28, 7, 5, 0, 0, time.UTC),
				},
			}
			output, err := executeTemplate(tmpl, testCase)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "Thu, 28 Oct 2021 17:05:00 AEST")
		})

		Convey("localtime function with no timezone provided", func() {
			loc, _ := time.LoadLocation("Australia/Queensland")
			templateText := `{{ .TestDate | localtime "Mon, 02 Jan 2006 15:04:05 MST" }}`
			tmpl, err := template.New("func-localtime-test").Funcs(libtemplate.GetFunctions(loc)).Parse(templateText)
			if err != nil {
				t.Fatalf("parsing: %s", err)
			}
			testCase := struct {
				TestDate types.CpsTime
			}{
				TestDate: types.CpsTime{
					Time: time.Date(2021, time.October, 28, 7, 5, 0, 0, time.UTC),
				},
			}
			output, err := executeTemplate(tmpl, testCase)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "Thu, 28 Oct 2021 17:05:00 AEST")
		})
	})
}
