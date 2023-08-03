package pbehavior

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
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
func GetTimeSpans(event Event, view timespan.Span) ([]timespan.Span, error) {
	if event.rOption == nil {
		return getSpansForSingleEvent(event.span, view), nil
	}

	return getTimeSpansForRecEvent(event, view)
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
		from := utils.MaxTime(start, view.From())
		to := utils.MinTime(start.Add(duration), view.To())

		if from.Before(to) {
			res = append(res, timespan.New(from, to))
		}
	}

	return res, nil
}
