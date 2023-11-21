package eventfilter

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/teambition/rrule-go"
)

func ResolveIntervals(ef *Rule, rrule *rrule.RRule, now time.Time, location *time.Location) {
	if rrule == nil {
		return
	}

	if ef.ResolvedStart == nil || ef.ResolvedStart.IsZero() {
		ef.ResolvedStart = &datetime.CpsTime{Time: ef.Start.Time}
	}

	if ef.ResolvedStop == nil || ef.ResolvedStop.IsZero() {
		ef.ResolvedStop = &datetime.CpsTime{Time: ef.Stop.Time}
	}

	interval := ef.Stop.Sub(ef.Start.Time)

	if ef.NextResolvedStart == nil || ef.NextResolvedStop == nil || ef.NextResolvedStart.IsZero() || ef.NextResolvedStop.IsZero() {
		v := rrule.After(ef.ResolvedStop.Add(-interval).In(location), false)
		if !v.IsZero() {
			ef.NextResolvedStart = &datetime.CpsTime{Time: v}
			ef.NextResolvedStop = &datetime.CpsTime{Time: v.Add(interval)}
		}
	}

	for now.After(ef.ResolvedStop.Time.In(location)) && ef.NextResolvedStart != nil && ef.NextResolvedStop != nil {
		ef.ResolvedStart = ef.NextResolvedStart
		ef.ResolvedStop = ef.NextResolvedStop

		v := rrule.After(ef.ResolvedStop.Add(-interval).In(location), false)
		if !v.IsZero() {
			ef.NextResolvedStart = &datetime.CpsTime{Time: v}
			ef.NextResolvedStop = &datetime.CpsTime{Time: v.Add(interval)}
		} else {
			ef.NextResolvedStart = nil
			ef.NextResolvedStop = nil
		}
	}
}
