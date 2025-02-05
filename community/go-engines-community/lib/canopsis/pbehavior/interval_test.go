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
	f := func(event pbehavior.Event, viewSpan timespan.Span, expected []timespan.Span) {
		t.Helper()
		output, err := pbehavior.GetTimeSpans(event, viewSpan)
		if len(output) != len(expected) || len(output) > 0 && !reflect.DeepEqual(output, expected) {
			t.Errorf(
				"expected result: (len %d)\n%s\nbut got: (len %d)\n%s\n",
				len(expected),
				sprintSpans(expected),
				len(output),
				sprintSpans(output),
			)
		}

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}

	// single event
	event := pbehavior.NewEvent(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00"))
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")), // current day
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("02-06-2020 00:00"), genTime("02-06-2020 23:59")), // tomorrow
		[]timespan.Span{},
	)
	f(
		event,
		timespan.New(genTime("31-05-2020 00:00"), genTime("31-05-2020 23:59")), // yesterday
		[]timespan.Span{},
	)

	// workday event
	event = pbehavior.NewRecEvent(
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
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 00:00")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
			timespan.New(genTime("02-06-2020 09:00"), genTime("02-06-2020 18:00")),
			timespan.New(genTime("03-06-2020 09:00"), genTime("03-06-2020 18:00")),
			timespan.New(genTime("04-06-2020 09:00"), genTime("04-06-2020 18:00")),
			timespan.New(genTime("05-06-2020 09:00"), genTime("05-06-2020 18:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("08-06-2020 00:00"), genTime("14-06-2020 00:00")), // next week
		[]timespan.Span{
			timespan.New(genTime("08-06-2020 09:00"), genTime("08-06-2020 18:00")),
			timespan.New(genTime("09-06-2020 09:00"), genTime("09-06-2020 18:00")),
			timespan.New(genTime("10-06-2020 09:00"), genTime("10-06-2020 18:00")),
			timespan.New(genTime("11-06-2020 09:00"), genTime("11-06-2020 18:00")),
			timespan.New(genTime("12-06-2020 09:00"), genTime("12-06-2020 18:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("25-05-2020 00:00"), genTime("31-05-2020 00:00")), // prev week
		[]timespan.Span{},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")), // current day
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 18:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("06-06-2020 00:00"), genTime("06-06-2020 23:59")), // weekend
		[]timespan.Span{},
	)

	// workday event for 1 and half week
	event = pbehavior.NewRecEvent(
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
	)
	f(
		event,
		timespan.New(genTime("08-06-2020 00:00"), genTime("14-06-2020 23:59")), // next week
		[]timespan.Span{
			timespan.New(genTime("08-06-2020 09:00"), genTime("08-06-2020 18:00")),
			timespan.New(genTime("09-06-2020 09:00"), genTime("09-06-2020 18:00")),
		},
	)

	// midnight/midday event
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 00:00"),
		genTime("01-06-2020 02:00"),
		&rrule.ROption{
			Freq:     rrule.HOURLY,
			Interval: 12,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
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
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")), // current day
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
			timespan.New(genTime("01-06-2020 12:00"), genTime("01-06-2020 14:00")),
		},
	)

	// recurrent event for 2 days
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 00:00"),
		genTime("01-06-2020 02:00"),
		&rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 2,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
			timespan.New(genTime("03-06-2020 00:00"), genTime("03-06-2020 02:00")),
			timespan.New(genTime("05-06-2020 00:00"), genTime("05-06-2020 02:00")),
			timespan.New(genTime("07-06-2020 00:00"), genTime("07-06-2020 02:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")), // current day
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 02:00")),
		},
	)
	f(
		event,
		timespan.New(genTime("02-06-2020 00:00"), genTime("02-06-2020 23:59")), // next day
		[]timespan.Span{},
	)

	// midday to next midday event for 3 days
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 12:00"),
		genTime("02-06-2020 12:00"),
		&rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 3,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 12:00"), genTime("02-06-2020 12:00")),
			timespan.New(genTime("04-06-2020 12:00"), genTime("05-06-2020 12:00")),
			timespan.New(genTime("07-06-2020 12:00"), genTime("07-06-2020 23:59")),
		},
	)

	// 2 days event for 3 days
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 12:00"),
		genTime("03-06-2020 12:00"),
		&rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 3,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 12:00"), genTime("03-06-2020 12:00")),
			timespan.New(genTime("04-06-2020 12:00"), genTime("06-06-2020 12:00")),
			timespan.New(genTime("07-06-2020 12:00"), genTime("07-06-2020 23:59")),
		},
	)

	// 3 days event for 4 days
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 12:00"),
		genTime("04-06-2020 12:00"),
		&rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 4,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 12:00"), genTime("04-06-2020 12:00")),
			timespan.New(genTime("05-06-2020 12:00"), genTime("07-06-2020 23:59")),
		},
	)

	// full day event for 3 days
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 00:00"),
		genTime("01-06-2020 23:59"),
		&rrule.ROption{
			Freq:     rrule.DAILY,
			Interval: 3,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("07-06-2020 23:59")), // current week
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")),
			timespan.New(genTime("04-06-2020 00:00"), genTime("04-06-2020 23:59")),
			timespan.New(genTime("07-06-2020 00:00"), genTime("07-06-2020 23:59")),
		},
	)

	// event each 30 minutes for 4 times
	event = pbehavior.NewRecEvent(
		genTime("01-06-2020 00:00"),
		genTime("01-06-2020 00:20"),
		&rrule.ROption{
			Freq:     rrule.MINUTELY,
			Interval: 30,
			Count:    4,
		},
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 23:59")), // current day
		[]timespan.Span{
			timespan.New(genTime("01-06-2020 00:00"), genTime("01-06-2020 00:20")),
			timespan.New(genTime("01-06-2020 00:30"), genTime("01-06-2020 00:50")),
			timespan.New(genTime("01-06-2020 01:00"), genTime("01-06-2020 01:20")),
			timespan.New(genTime("01-06-2020 01:30"), genTime("01-06-2020 01:50")),
		},
	)

	// monthly event from 2020-06-27 17:00:00 UTC to 2020-06-30 16:59:59 UTC
	event = pbehavior.NewRecEvent(
		genTime("27-06-2020 17:00"),
		genTime("30-06-2020 16:59"),
		resolveROption("FREQ=MONTHLY;BYMONTHDAY=27"),
	)
	f(
		event,
		timespan.New(genTime("01-06-2020 00:00"), genTime("01-08-2020 16:59")), // 2 months
		[]timespan.Span{
			timespan.New(genTime("27-06-2020 17:00"), genTime("30-06-2020 16:59")),
			timespan.New(genTime("27-07-2020 17:00"), genTime("30-07-2020 16:59")),
		},
	)

	// event during summer/winter time change
	parisTZ, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatalf("failed to load timezone %v", err)
	}

	event = pbehavior.NewRecEvent(
		genTime("01-01-2020 09:00", parisTZ),
		genTime("01-01-2020 18:00", parisTZ),
		&rrule.ROption{
			Freq: rrule.DAILY,
		},
	)
	f(
		event,
		timespan.New(genTime("26-03-2020 00:00"), genTime("02-04-2020 00:00")), // summer change
		[]timespan.Span{
			timespan.New(genTime("26-03-2020 09:00", parisTZ), genTime("26-03-2020 18:00", parisTZ)),
			timespan.New(genTime("27-03-2020 09:00", parisTZ), genTime("27-03-2020 18:00", parisTZ)),
			timespan.New(genTime("28-03-2020 09:00", parisTZ), genTime("28-03-2020 18:00", parisTZ)),
			timespan.New(genTime("29-03-2020 09:00", parisTZ), genTime("29-03-2020 18:00", parisTZ)),
			timespan.New(genTime("30-03-2020 09:00", parisTZ), genTime("30-03-2020 18:00", parisTZ)),
			timespan.New(genTime("31-03-2020 09:00", parisTZ), genTime("31-03-2020 18:00", parisTZ)),
			timespan.New(genTime("01-04-2020 09:00", parisTZ), genTime("01-04-2020 18:00", parisTZ)),
		},
	)
	f(
		event,
		timespan.New(genTime("22-10-2020 00:00"), genTime("29-10-2020 00:00")), // winter change
		[]timespan.Span{
			timespan.New(genTime("22-10-2020 09:00", parisTZ), genTime("22-10-2020 18:00", parisTZ)),
			timespan.New(genTime("23-10-2020 09:00", parisTZ), genTime("23-10-2020 18:00", parisTZ)),
			timespan.New(genTime("24-10-2020 09:00", parisTZ), genTime("24-10-2020 18:00", parisTZ)),
			timespan.New(genTime("25-10-2020 09:00", parisTZ), genTime("25-10-2020 18:00", parisTZ)),
			timespan.New(genTime("26-10-2020 09:00", parisTZ), genTime("26-10-2020 18:00", parisTZ)),
			timespan.New(genTime("27-10-2020 09:00", parisTZ), genTime("27-10-2020 18:00", parisTZ)),
			timespan.New(genTime("28-10-2020 09:00", parisTZ), genTime("28-10-2020 18:00", parisTZ)),
		},
	)
}

func genTime(value string, l ...*time.Location) time.Time {
	format := "02-01-2006 15:04"
	location := time.UTC
	if len(l) > 0 {
		location = l[0]
	}

	date, err := time.ParseInLocation(format, value, location)
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
