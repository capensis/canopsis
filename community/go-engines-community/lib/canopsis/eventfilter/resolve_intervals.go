package eventfilter

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/teambition/rrule-go"
)

func ResolveIntervals(ef *Rule, rrule *rrule.RRule, now time.Time, location *time.Location) {
	if rrule == nil {
		return
	}

	if ef.ResolvedStart == nil || ef.ResolvedStart.IsZero() {
		ef.ResolvedStart = &types.CpsTime{Time: ef.Start.Time}
	}

	if ef.ResolvedStop == nil || ef.ResolvedStop.IsZero() {
		ef.ResolvedStop = &types.CpsTime{Time: ef.Stop.Time}
	}

	interval := ef.Stop.Sub(ef.Start.Time)

	if ef.NextResolvedStart == nil || ef.NextResolvedStop == nil || ef.NextResolvedStart.IsZero() || ef.NextResolvedStop.IsZero() {
		v := rrule.After(ef.ResolvedStop.Add(-interval).In(location), false)
		if !v.IsZero() {
			ef.NextResolvedStart = &types.CpsTime{Time: v}
			ef.NextResolvedStop = &types.CpsTime{Time: v.Add(interval)}
		}
	}

	for now.After(ef.ResolvedStop.Time.In(location)) && ef.NextResolvedStart != nil && ef.NextResolvedStop != nil {
		ef.ResolvedStart = ef.NextResolvedStart
		ef.ResolvedStop = ef.NextResolvedStop

		v := rrule.After(ef.ResolvedStop.Add(-interval).In(location), false)
		if !v.IsZero() {
			ef.NextResolvedStart = &types.CpsTime{Time: v}
			ef.NextResolvedStop = &types.CpsTime{Time: v.Add(interval)}
		} else {
			ef.NextResolvedStart = nil
			ef.NextResolvedStop = nil
		}
	}
}
