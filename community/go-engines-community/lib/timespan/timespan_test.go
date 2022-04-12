package timespan

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

const format = "02-01-2006 15:04:05"

func TestNew(t *testing.T) {
	today := time.Now()
	tomorrow := time.Now().Add(24 * time.Hour)
	span := New(today, tomorrow)

	if span.From() != today || span.To() != tomorrow {
		t.Errorf("expected span [%v, %v] but got [%v %v]",
			tomorrow, tomorrow, span.From(), span.To())
	}

	span = New(tomorrow, today)

	if span.From() != today || span.To() != tomorrow {
		t.Errorf("expected span [%v, %v] but got [%v %v]",
			tomorrow, tomorrow, span.From(), span.To())
	}
}

func TestIntersect(t *testing.T) {
	weekago := time.Now().Add(-7 * 24 * time.Hour)
	yesterday := time.Now().Add(-24 * time.Hour)
	today := time.Now()
	tomorrow := time.Now().Add(24 * time.Hour)

	dataSets := map[string]struct {
		left     Span
		right    Span
		expected *Span
	}{
		"Given left is before right Should return right start to left end span": {
			left:  New(weekago, today),
			right: New(yesterday, tomorrow),
			expected: &Span{
				from: yesterday,
				to:   today,
			},
		},
		"Given right is before left Should return left start to right end span": {
			left:  New(yesterday, tomorrow),
			right: New(weekago, today),
			expected: &Span{
				from: yesterday,
				to:   today,
			},
		},
		"Given left adjoins to left Should return output start is equal to its end": {
			left:  New(yesterday, today),
			right: New(today, tomorrow),
			expected: &Span{
				from: today,
				to:   today,
			},
		},
		"Given left is long before right Should return no span": {
			left:     New(weekago, yesterday),
			right:    New(today, tomorrow),
			expected: nil,
		},
		"Given right is before left Should return no span": {
			left:     New(today, tomorrow),
			right:    New(weekago, yesterday),
			expected: nil,
		},
	}

	for name, data := range dataSets {
		output := Intersect(data.left, data.right)

		if ((output == nil || data.expected == nil) && output != data.expected) ||
			(output != nil && data.expected != nil && (output.From() != data.expected.From() || output.To() != data.expected.To())) {
			t.Errorf("%s: expected %s but got %s", name, sprintSpan(data.expected), sprintSpan(output))
		}
	}
}

func TestSpan_Duration(t *testing.T) {
	from := genTime("01-06-2020 10:05")
	to := genTime("02-06-2020 10:06")
	expected := 24*time.Hour + time.Minute
	span := New(from, to)

	if span.Duration() != expected {
		t.Errorf("expected %v but got %v", expected, span.Duration())
	}

	// daylight saving time started 2020-03-29 02:00 CET
	from = genTime("28-03-2020 10:05")
	to = genTime("29-03-2020 10:06")
	expected = 23*time.Hour + time.Minute
	span = New(from, to)

	if span.Duration() != expected {
		t.Errorf("expected %v but got %v", expected, span.Duration())
	}
}

func TestSpan_Diff(t *testing.T) {
	weekago := time.Now().Add(-7 * 24 * time.Hour)
	yesterday := time.Now().Add(-24 * time.Hour)
	today := time.Now()
	tomorrow := time.Now().Add(24 * time.Hour)

	dataSets := map[string]struct {
		left     Span
		right    Span
		expected []Span
	}{
		"Given left is before right and they do not overlap Should return left span": {
			left:     New(weekago, yesterday),
			right:    New(today, tomorrow),
			expected: []Span{New(weekago, yesterday)},
		},
		"Given right is before left and they do not overlap Should return left span": {
			left:     New(today, tomorrow),
			right:    New(weekago, yesterday),
			expected: []Span{New(today, tomorrow)},
		},
		"Given left is before right and they overlap Should return left start to right start span": {
			left:     New(weekago, today),
			right:    New(yesterday, tomorrow),
			expected: []Span{New(weekago, yesterday)},
		},
		"Given right is before left and they overlap Should return right end to left end span": {
			left:     New(yesterday, tomorrow),
			right:    New(weekago, today),
			expected: []Span{New(today, tomorrow)},
		},
		"Given left is in right Should return no span": {
			left:     New(yesterday, today),
			right:    New(weekago, tomorrow),
			expected: []Span{},
		},
		"Given right is in left Should return left start to right start span and right end to left end span": {
			left:  New(weekago, tomorrow),
			right: New(yesterday, today),
			expected: []Span{
				New(weekago, yesterday),
				New(today, tomorrow),
			},
		},
	}

	for name, data := range dataSets {
		output := data.left.Diff(data.right)

		if len(output) != len(data.expected) || len(output) > 0 && !reflect.DeepEqual(output, data.expected) {
			t.Errorf("%s: expected %s but got %s", name, sprintSpans(data.expected), sprintSpans(output))
		}
	}
}

func genTime(value string) time.Time {
	const timeZoneNaiveLayout = "02-01-2006 15:04"

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}

	date, err := time.ParseInLocation(timeZoneNaiveLayout, value, loc)

	if err != nil {
		panic(err)
	}

	return date
}

func sprintSpan(v *Span) string {
	if v == nil {
		return "nil"
	}

	return fmt.Sprintf("[%v %v]", v.From().Format(format), v.To().Format(format))
}

func sprintSpans(list []Span) string {
	res := make([]string, len(list))

	for i, s := range list {
		res[i] = sprintSpan(&s)
	}

	return strings.Join(res, "\n")
}
