package timespan

import (
	"encoding/json"
	"errors"
	"time"
)

// Span represents interval between two times.
type Span struct {
	from, to time.Time
	stype    string
}

// New creates new span with given start and end times.
func New(from, to time.Time) Span {
	if to.Before(from) {
		from, to = to, from
	}

	return Span{
		from: from,
		to:   to,
	}
}

// TypedNew creates new span of given type with start and end times.
func TypedNew(from, to time.Time, stype string) Span {
	s := New(from, to)
	s.stype = stype
	return s
}

// Intersect returns intersect of spans or nil if they don't overlap.
func Intersect(left, right Span) *Span {
	if left.to.Before(right.from) || right.to.Before(left.from) {
		return nil
	}

	return &Span{
		from: maxTime(left.from, right.from),
		to:   minTime(left.to, right.to),
	}
}

// Compact merges adjacent spans if end of prev span is less
// than start of next span or gap less than diff.
func Compact(list []Span, diff time.Duration) []Span {
	if len(list) == 0 {
		return list
	}

	res := make([]Span, 0, len(list))
	newFrom := list[0].from

	if len(list) == 1 {
		res = append(res, list[0])
	}

	for i := range list {
		if i == 0 {
			continue
		}

		if list[i].from.Sub(list[i-1].to) > diff {
			res = append(res, Span{
				from: newFrom,
				to:   list[i-1].to,
			})
			newFrom = list[i].from
		}

		if i == len(list)-1 {
			res = append(res, Span{
				from: newFrom,
				to:   list[i].to,
			})
		}
	}

	return res
}

// From returns time at the start of span.
func (s *Span) From() time.Time {
	return s.from
}

// To returns time at the end of span.
func (s *Span) To() time.Time {
	return s.to
}

// Type returns type of span.
func (s *Span) Type() string {
	return s.stype
}

// Duration returns duration of span.
func (s *Span) Duration() time.Duration {
	return s.to.Sub(s.from)
}

// In checks if the time moment is in the time span.
func (s *Span) In(v time.Time) bool {
	return !v.Before(s.from) && !v.After(s.to)
}

// Diff returns difference which is in s but not in d.
func (s *Span) Diff(d Span) []Span {
	// sFrom sTo dFrom dTo => sFrom sTo
	// dFrom dTo sFrom sTo => sFrom sTo
	if s.to.Before(d.from) || s.from.After(d.to) {
		return []Span{*s}
	}

	res := make([]Span, 0)
	// sFrom dFrom sTo dTo => sFrom dFrom
	// sFrom dFrom dTo sTo => sFrom dFrom
	if s.from.Before(d.from) {
		res = append(res, Span{
			from: s.from,
			to:   d.from,
		})
	}

	// sFrom dFrom dTo sTo => dTo sTo
	// dFrom sFrom dTo sTo => dTo sTo
	if s.to.After(d.to) {
		res = append(res, Span{
			from: d.to,
			to:   s.to,
		})
	}

	// dFrom sFrom sTo dTo => nil
	return res
}

type BySpans []Span

func (a BySpans) Len() int      { return len(a) }
func (a BySpans) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySpans) Less(i, j int) bool {
	return a[i].from.UTC().Before(a[j].from.UTC()) ||
		(a[i].from.UTC().Equal(a[j].from) && a[i].to.UTC().Before(a[j].to.UTC()))
}

// minTime returns maximal time between arguments.
func maxTime(left, right time.Time) time.Time {
	if left.After(right) {
		return left
	}

	return right
}

// minTime returns minimal time between arguments.
func minTime(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}

	return right
}

func (s Span) MarshalJSON() ([]byte, error) {
	return json.Marshal([]time.Time{s.from, s.to})
}

func (s *Span) UnmarshalJSON(b []byte) error {
	var v []time.Time
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if len(v) != 2 {
		return errors.New("invalid data")
	}

	s.from = v[0]
	s.to = v[1]

	return nil
}
