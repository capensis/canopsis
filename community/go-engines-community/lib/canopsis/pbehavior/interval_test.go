package pbehavior_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/teambition/rrule-go"
)

func TestGetTimeSpans(t *testing.T) {
	for suiteName, data := range dataSetsForGetTimeSpans() {
		for caseName, caseData := range data.cases {
			output, err := pbehavior.GetTimeSpans(data.event, caseData.viewSpan)

			if len(output) != len(caseData.expected) || len(output) > 0 && !reflect.DeepEqual(output, caseData.expected) {
				t.Errorf(
					"%s %s: expected result: (len %d)\n%s\nbut got: (len %d)\n%s\n",
					suiteName,
					caseName,
					len(caseData.expected),
					sprintSpans(caseData.expected),
					len(output),
					sprintSpans(output),
				)
			}

			if err != nil {
				t.Errorf("%s %s: expected no error but got %v", suiteName, caseName, err)
			}
		}
	}
}

type intervalSuitDataSet struct {
	event pbehavior.Event
	cases map[string]intervalCaseDataSet
}

type intervalCaseDataSet struct {
	viewSpan timespan.Span
	expected []timespan.Span
}

func dataSetsForGetTimeSpans() map[string]intervalSuitDataSet {
	return map[string]intervalSuitDataSet{
		"Given single event": {
			event: pbehavior.NewEvent(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return one day span": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
					},
				},
				"and current day span Should return one day span": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
					},
				},
				"and tomorrow span Should return 0 spans": {
					viewSpan: timespan.New(genTime("02-06-2020 00:00"), genTime("02-06-2020 23:59")),
					expected: []timespan.Span{},
				},
				"and yesterday span Should return 0 spans": {
					viewSpan: timespan.New(genTime("31-05-2020 00:00"), genTime("31-05-2020 23:59")),
					expected: []timespan.Span{},
				},
			},
		},
		"Given workday event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 09:00"),
				genTime("01-06-2020 18:00"),
				&rrule.ROption{
					Freq: rrule.DAILY,
					Byweekday: []rrule.Weekday{
						rrule.MO,
						rrule.TU,
						rrule.WE,
						rrule.TH,
						rrule.FR,
					},
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return 5 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 00:00")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
						timespan.New(genTime("02-06-2020 09:00"), genTime("02-06-2020 18:00")),
						timespan.New(genTime("03-06-2020 09:00"), genTime("03-06-2020 18:00")),
						timespan.New(genTime("04-06-2020 09:00"), genTime("04-06-2020 18:00")),
						timespan.New(genTime("05-06-2020 09:00"), genTime("05-06-2020 18:00")),
					},
				},
				"and next week span Should return 5 spans": {
					viewSpan: timespan.New(genTime("08-06-2020 00:00"), genTime("14-06-2020 00:00")),
					expected: []timespan.Span{
						timespan.New(genTime("08-06-2020 09:00"), genTime("08-06-2020 18:00")),
						timespan.New(genTime("09-06-2020 09:00"), genTime("09-06-2020 18:00")),
						timespan.New(genTime("10-06-2020 09:00"), genTime("10-06-2020 18:00")),
						timespan.New(genTime("11-06-2020 09:00"), genTime("11-06-2020 18:00")),
						timespan.New(genTime("12-06-2020 09:00"), genTime("12-06-2020 18:00")),
					},
				},
				"and prev week span Should return 0 spans": {
					viewSpan: timespan.New(genTime("25-05-2020 00:00"), genTime("31-05-2020 00:00")),
					expected: []timespan.Span{},
				},
				"and current day span Should return 1 span": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
					},
				},
				"and weekend span Should return 0 spans": {
					viewSpan: timespan.New(genTime("06-06-2020 00:00"), genTime("06-06-2020 23:59")),
					expected: []timespan.Span{},
				},
			},
		},
		"Given workday event for 1 and half week": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 09:00"),
				genTime("01-06-2020 18:00"),
				&rrule.ROption{
					Freq:  rrule.DAILY,
					Count: 7,
					Byweekday: []rrule.Weekday{
						rrule.MO,
						rrule.TU,
						rrule.WE,
						rrule.TH,
						rrule.FR,
					},
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and next week span Should return 2 spans": {
					viewSpan: timespan.New(genTime("08-06-2020 00:00"), genTime("14-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("08-06-2020 09:00"), genTime("08-06-2020 18:00")),
						timespan.New(genTime("09-06-2020 09:00"), genTime("09-06-2020 18:00")),
					},
				},
			},
		},
		"Given midnight/midday event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 00:00"),
				genTime("01-06-2020 02:00"),
				&rrule.ROption{
					Freq:     rrule.HOURLY,
					Interval: 12,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return current 14 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
						timespan.New(genTime("01-06-2020 12:00"), genTime("01-06-2020 14:00")),
						timespan.New(genTime("02-06-2020 00:00"), genTime("02-06-2020 02:00")),
						timespan.New(genTime("02-06-2020 12:00"), genTime("02-06-2020 14:00")),
						timespan.New(genTime("03-06-2020 00:00"), genTime("03-06-2020 02:00")),
						timespan.New(genTime("03-06-2020 12:00"), genTime("03-06-2020 14:00")),
						timespan.New(genTime("04-06-2020 00:00"), genTime("04-06-2020 02:00")),
						timespan.New(genTime("04-06-2020 12:00"), genTime("04-06-2020 14:00")),
						timespan.New(genTime("05-06-2020 00:00"), genTime("05-06-2020 02:00")),
						timespan.New(genTime("05-06-2020 12:00"), genTime("05-06-2020 14:00")),
						timespan.New(genTime("06-06-2020 00:00"), genTime("06-06-2020 02:00")),
						timespan.New(genTime("06-06-2020 12:00"), genTime("06-06-2020 14:00")),
						timespan.New(genTime("07-06-2020 00:00"), genTime("07-06-2020 02:00")),
						timespan.New(genTime("07-06-2020 12:00"), genTime("07-06-2020 14:00")),
					},
				},
				"and current day span Should return current 2 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
						timespan.New(genTime("01-06-2020 12:00"), genTime("01-06-2020 14:00")),
					},
				},
			},
		},
		"Given getTimeSpansForSingleEvent in 2 days event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 00:00"),
				genTime("01-06-2020 02:00"),
				&rrule.ROption{
					Freq:     rrule.DAILY,
					Interval: 2,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return 4 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
						timespan.New(genTime("03-06-2020 00:00"), genTime("03-06-2020 02:00")),
						timespan.New(genTime("05-06-2020 00:00"), genTime("05-06-2020 02:00")),
						timespan.New(genTime("07-06-2020 00:00"), genTime("07-06-2020 02:00")),
					},
				},
				"and current day span Should return 1 span": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
					},
				},
				"and next day span Should return 0 spans": {
					viewSpan: timespan.New(genTime("02-06-2020 00:00"), genTime("02-06-2020 23:59")),
					expected: []timespan.Span{},
				},
			},
		},
		"Given midday to next midday getTimeSpansForSingleEvent in 3 days event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 12:00"),
				genTime("02-06-2020 12:00"),
				&rrule.ROption{
					Freq:     rrule.DAILY,
					Interval: 3,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return 3 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 12:00"), genTime("02-06-2020 12:00")),
						timespan.New(genTime("04-06-2020 12:00"), genTime("05-06-2020 12:00")),
						timespan.New(genTime("07-06-2020 12:00"), genTime("07-06-2020 23:59")),
					},
				},
			},
		},
		"Given 2 days getTimeSpansForSingleEvent in 3 days event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 12:00"),
				genTime("03-06-2020 12:00"),
				&rrule.ROption{
					Freq:     rrule.DAILY,
					Interval: 3,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return 3 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 12:00"), genTime("03-06-2020 12:00")),
						timespan.New(genTime("04-06-2020 12:00"), genTime("06-06-2020 12:00")),
						timespan.New(genTime("07-06-2020 12:00"), genTime("07-06-2020 23:59")),
					},
				},
			},
		},
		"Given 3 days getTimeSpansForSingleEvent in 4 days event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 12:00"),
				genTime("04-06-2020 12:00"),
				&rrule.ROption{
					Freq:     rrule.DAILY,
					Interval: 4,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return 2 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 12:00"), genTime("04-06-2020 12:00")),
						timespan.New(genTime("05-06-2020 12:00"), genTime("07-06-2020 23:59")),
					},
				},
			},
		},
		"Given full day getTimeSpansForSingleEvent in 3 days event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 00:00"),
				genTime("01-06-2020 23:59"),
				&rrule.ROption{
					Freq:     rrule.DAILY,
					Interval: 3,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current week span Should return 2 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
						timespan.New(genTime("04-06-2020 00:00"), genTime("04-06-2020 23:59")),
						timespan.New(genTime("07-06-2020 00:00"), genTime("07-06-2020 23:59")),
					},
				},
			},
		},
		"Given getTimeSpansForSingleEvent in 30 minutes 4 times event": {
			event: pbehavior.NewRecEvent(
				genTime("01-06-2020 00:00"),
				genTime("01-06-2020 00:20"),
				&rrule.ROption{
					Freq:     rrule.MINUTELY,
					Interval: 30,
					Count:    4,
				},
			),
			cases: map[string]intervalCaseDataSet{
				"and current day span Should return 4 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
					expected: []timespan.Span{
						timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 00:20")),
						timespan.New(genTime("01-06-2020 00:30"), genTime("01-06-2020 00:50")),
						timespan.New(genTime("01-06-2020 01:00"), genTime("01-06-2020 01:20")),
						timespan.New(genTime("01-06-2020 01:30"), genTime("01-06-2020 01:50")),
					},
				},
			},
		},
		"Given MONTHLY event from 2020-06-27 17:00:00 UTC to 2020-06-30 16:59:59 UTC": {
			event: pbehavior.NewRecEvent(
				genTime("27-06-2020 17:00"),
				genTime("30-06-2020 16:59"),
				resolveROption("FREQ=MONTHLY;BYMONTHDAY=27"),
			),
			cases: map[string]intervalCaseDataSet{
				"and current day span Should return 2 spans": {
					viewSpan: timespan.New(genTime("01-06-2020 00:00"), genTime("01-08-2020 16:59")),
					expected: []timespan.Span{
						timespan.New(genTime("27-06-2020 17:00"), genTime("30-06-2020 16:59")),
						timespan.New(genTime("27-07-2020 17:00"), genTime("30-07-2020 16:59")),
					},
				},
			},
		},
	}
}

func genTime(value string) time.Time {
	format := "02-01-2006 15:04"
	date, err := time.Parse(format, value)

	if err != nil {
		panic(err)
	}

	return date
}

func resolveROption(s string) *rrule.ROption {
	r, err := rrule.StrToROption(s)
	if err != nil {
		panic(err)
	}

	return r
}

func sprintSpans(list []timespan.Span) string {
	res := make([]string, len(list))

	for i, s := range list {
		res[i] = fmt.Sprintf("[%v %v]", s.From(), s.To())
	}

	return strings.Join(res, "\n")
}
