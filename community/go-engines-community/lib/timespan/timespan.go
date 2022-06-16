package timespan

import (
	"encoding/json"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// Span represents interval between two times.
type Span struct {
	from, to time.Time
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

// Intersect returns intersect of spans or nil if they don't overlap.
func Intersect(left, right Span) *Span {
	if left.to.Before(right.from) || right.to.Before(left.from) {
		return nil
	}

	return &Span{
		from: utils.MaxTime(left.from, right.from),
		to:   utils.MinTime(left.to, right.to),
	}
}

// From returns time at the start of span.
func (s *Span) From() time.Time {
	return s.from
}

// To returns time at the end of span.
func (s *Span) To() time.Time {
	return s.to
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

func (s Span) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int64{s.from.Unix(), s.to.Unix()})
}

func (s *Span) UnmarshalJSON(b []byte) error {
	var v []int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if len(v) != 2 {
		return errors.New("invalid data")
	}

	s.from = time.Unix(v[0], 0)
	s.to = time.Unix(v[1], 0)

	return nil
}
