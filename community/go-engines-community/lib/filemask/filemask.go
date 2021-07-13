package filemask

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	FileMaskTestCaseName = "%test_case%"
	FileMaskYear         = "YYYY"
	FileMaskMonth        = "MM"
	FileMaskDay          = "DD"
	FileMaskHour         = "hh"
	FileMaskMinute       = "mm"
	FileMaskSecond       = "ss"
	FileMaskTimezone     = "Z"
)

type FileMask struct {
	mask string
	re   *regexp.Regexp
}

func NewFileMask(v string) (FileMask, error) {
	maskPatternReplace := map[string]string{
		FileMaskTestCaseName: "(?P<name>.+)",
		FileMaskYear:         "(?P<year>\\d{4})",
		FileMaskMonth:        "(?P<month>\\d{2})",
		FileMaskDay:          "(?P<day>\\d{2})",
		FileMaskHour:         "(?P<hour>\\d{2})",
		FileMaskMinute:       "(?P<minute>\\d{2})",
		FileMaskSecond:       "(?P<second>\\d{2})",
		FileMaskTimezone:     "(?P<tz>[\\w\\+\\-]+)",
	}

	expr := regexp.QuoteMeta(v)

	for maskVar, pattern := range maskPatternReplace {
		newExpr := strings.Replace(expr, maskVar, pattern, 1)
		if maskVar != FileMaskTimezone && newExpr == expr {
			return FileMask{}, fmt.Errorf("invalid mask %q : %q is missing", v, maskVar)
		}

		expr = newExpr
	}

	expr = fmt.Sprintf("^%s$", expr)
	re, err := regexp.Compile(expr)
	if err != nil {
		return FileMask{}, err
	}

	return FileMask{re: re, mask: v}, nil
}

func (m *FileMask) Parse(v string) (string, time.Time, error) {
	var t time.Time

	if m.re == nil {
		return "", t, errors.New("file mask is not defined")
	}

	matches := m.re.FindStringSubmatch(v)
	if len(matches) == 0 {
		return "", t, errors.New("value doesn't match file mask")
	}

	names := m.re.SubexpNames()
	res := make(map[string]string, len(names))

	for _, name := range names {
		if name == "" {
			continue
		}

		i := m.re.SubexpIndex(name)
		if i < 0 {
			return "", t, fmt.Errorf("value doesn't match file mask : %q is missing", name)
		}

		res[name] = matches[i]
	}

	var err error
	format := "2006-01-02 15:04:05"
	tval := fmt.Sprintf("%s-%s-%s %s:%s:%s", res["year"], res["month"], res["day"],
		res["hour"], res["minute"], res["second"])
	if res["tz"] != "" {
		format += " MST"
		tval += " " + res["tz"]
	}

	t, err = time.Parse(format, tval)
	if err != nil {
		return "", t, err
	}

	return res["name"], t, nil
}

func (m FileMask) String() string {
	return m.mask
}
