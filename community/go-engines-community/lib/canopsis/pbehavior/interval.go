package pbehavior

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"github.com/teambition/rrule-go"
)

// Event represents a recurrent calendar event.
type Event struct {
	span    timespan.Span
	rOption *rrule.ROption
}

// NewEvent creates a new event with the given start and end times.
func NewEvent(startAt, endAt time.Time) Event {
	return Event{
		span:    timespan.New(startAt, endAt),
		rOption: nil,
	}
}

// NewRecEvent creates a new recurrent event with the given start and end times
// and recurrent rule.
func NewRecEvent(startAt, endAt time.Time, rOption *rrule.ROption) Event {
	return Event{
		span:    timespan.New(startAt, endAt),
		rOption: rOption,
	}
}

// GetTimeSpans returns all time spans of event within view.
func GetTimeSpans(event Event, view timespan.Span, exdates ...[]timespan.Span) ([]timespan.Span, error) {
	var ex, res []timespan.Span
	if len(exdates) == 1 {
		ex = exdates[0]
	} else if len(exdates) > 1 {
		panic("too much arguments")
	}

	if event.rOption == nil {
		res = getSpansForSingleEvent(event.span, view)
	} else {
		var err error
		res, err = getTimeSpansForRecEvent(event, view)
		if err != nil {
			return nil, err
		}
	}

	return removeExdates(res, ex), nil
}

// GetDateSpans returns all date spans of event within view.
func GetDateSpans(event Event, view timespan.Span, exdates ...[]timespan.Span) ([]timespan.Span, error) {
	v := timespan.New(
		dateOf(view.From()),
		endDateOf(view.To()),
	)
	spans, err := GetTimeSpans(event, v, exdates...)
	if err != nil {
		return nil, err
	}

	res := make([]timespan.Span, len(spans))
	for i := range spans {
		res[i] = dateSpanOf(spans[i])
	}

	return timespan.Compact(res, 24*time.Hour), nil
}

func getSpansForSingleEvent(span, view timespan.Span) []timespan.Span {
	intersect := timespan.Intersect(span, view)
	if intersect == nil {
		return []timespan.Span{}
	}

	return []timespan.Span{*intersect}
}

func getTimeSpansForRecEvent(event Event, view timespan.Span) ([]timespan.Span, error) {
	rOption := *event.rOption
	rOption.Dtstart = event.span.From()
	r, err := rrule.NewRRule(rOption)
	if err != nil {
		return nil, err
	}

	duration := event.span.Duration()
	after := view.From().Add(-duration)
	before := view.To()
	startList := r.Between(after, before, true)
	res := make([]timespan.Span, 0, len(startList))

	for _, start := range startList {
		from := maxTime(start, view.From())
		to := minTime(start.Add(duration), view.To())

		if from.Before(to) {
			res = append(res, timespan.New(from, to))
		}
	}

	return res, nil
}

func removeExdates(spans, exdates []timespan.Span) []timespan.Span {
	res := spans

	for _, ex := range exdates {
		diff := make([]timespan.Span, 0)
		for _, v := range res {
			diff = append(diff, v.Diff(ex)...)
		}

		res = diff
		if len(res) == 0 {
			break
		}
	}

	return res
}

// dateSpanOf returns time span without hour, minute, and second.
func dateSpanOf(v timespan.Span) timespan.Span {
	return timespan.New(dateOf(v.From()), dateOf(v.To()))
}

// dateOf returns time without hour, minute, and second.
func dateOf(v time.Time) time.Time {
	y, m, d := v.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, v.Location())
}

// endDateOf returns date with 23:59:59 time.
func endDateOf(v time.Time) time.Time {
	y, m, d := v.Date()
	return time.Date(y, m, d, 23, 59, 59, 0, v.Location())
}

// minTime returns minimal time between arguments.
func minTime(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}

	return right
}

// maxTime returns max time between arguments.
func maxTime(left, right time.Time) time.Time {
	if left.After(right) {
		return left
	}

	return right
}
