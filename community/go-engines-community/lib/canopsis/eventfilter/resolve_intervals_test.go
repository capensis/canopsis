package eventfilter

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/teambition/rrule-go"
)

func TestUpdateIntervals(t *testing.T) {
	dataSets := map[string]struct {
		ef                        Rule
		rrule                     string
		now                       time.Time
		expectedResolvedStart     types.CpsTime
		expectedResolvedStop      types.CpsTime
		expectedNextResolvedStart *types.CpsTime
		expectedNextResolvedStop  *types.CpsTime
	}{
		"simple case, before the interval, shouldn't be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"simple case, in the beginning of the interval, shouldn't be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"simple case, in the middle of the interval, shouldn't be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"simple case, after the interval, should be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC)},
		},
		"before the interval, but next interval is not calculated, should be updated": {
			ef: Rule{
				Start:         &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"in the interval, but next interval is not calculated, should be updated": {
			ef: Rule{
				Start:         &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"after the interval, but next interval is not calculated, should be updated": {
			ef: Rule{
				Start:         &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC)},
		},
		"before the interval, but next interval is not calculated, cross intervals, should be updated": {
			ef: Rule{
				Start:         &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
		},
		"in the interval, cross intervals, shouldn't be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
		},
		"after the interval, cross intervals, should be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 3, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
		},
		"after the interval, with counts, next interval should be nil, because of exceeded count, should be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;COUNT=5",
			now:                       time.Date(2022, 1, 6, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: nil,
			expectedNextResolvedStop:  nil,
		},
		"after the interval, with counts, next interval should be nil, because of exceeded count, shouldn't be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: nil,
				NextResolvedStop:  nil,
			},
			rrule:                     "FREQ=DAILY;COUNT=5",
			now:                       time.Date(2022, 1, 7, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: nil,
			expectedNextResolvedStop:  nil,
		},
		"after the interval, long cross intervals with counts, should be updated": {
			ef: Rule{
				Start:             &types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &types.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &types.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: nil,
				NextResolvedStop:  nil,
			},
			rrule:                     "FREQ=DAILY;COUNT=5",
			now:                       time.Date(2022, 1, 7, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"test with old intervals, they should be be updated": {
			ef: Rule{
				Start:         &types.CpsTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &types.CpsTime{Time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &types.CpsTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &types.CpsTime{Time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
			expectedResolvedStart:     types.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &types.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &types.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
		},
	}

	for testName, dataSet := range dataSets {
		t.Run(testName, func(t *testing.T) {
			opt, err := rrule.StrToROption(dataSet.rrule)
			if err != nil {
				t.Fatalf("error is not expected = %s", err.Error())
			}

			r, err := rrule.NewRRule(*opt)
			if err != nil {
				t.Fatalf("error is not expected = %s", err.Error())
			}

			if opt.Count != 0 {
				r.DTStart(dataSet.ef.Start.Time)
			} else {
				r.DTStart(dataSet.ef.ResolvedStart.Time)
			}

			ResolveIntervals(&dataSet.ef, r, dataSet.now, time.UTC)

			if !dataSet.expectedResolvedStart.Equal(dataSet.ef.ResolvedStart.Time) {
				t.Errorf("expected start = %v, got = %v", dataSet.expectedResolvedStart, dataSet.ef.ResolvedStart)
			}

			if !dataSet.expectedResolvedStop.Equal(dataSet.ef.ResolvedStop.Time) {
				t.Errorf("expected stop = %v, got = %v", dataSet.expectedResolvedStop, dataSet.ef.ResolvedStop)
			}

			if dataSet.expectedNextResolvedStart == nil && dataSet.ef.NextResolvedStart != nil {
				t.Fatalf("expected next start = nil, but got = %v", dataSet.ef.NextResolvedStart)
			}

			if dataSet.expectedNextResolvedStart != nil && dataSet.ef.NextResolvedStart == nil {
				t.Fatalf("expected next start = %v, but got nil", dataSet.expectedNextResolvedStart)
			}

			if dataSet.expectedNextResolvedStop == nil && dataSet.ef.NextResolvedStop != nil {
				t.Fatalf("expected next stop = nil, but got = %v", dataSet.ef.NextResolvedStop)
			}

			if dataSet.expectedNextResolvedStop != nil && dataSet.ef.NextResolvedStop == nil {
				t.Fatalf("expected next stop = %v, but got nil", dataSet.expectedNextResolvedStop)
			}

			if dataSet.expectedNextResolvedStart != nil && !dataSet.expectedNextResolvedStart.Equal(dataSet.ef.NextResolvedStart.Time) {
				t.Errorf("expected next start = %v, got = %v", dataSet.expectedNextResolvedStart, dataSet.ef.NextResolvedStart)
			}

			if dataSet.expectedNextResolvedStop != nil && !dataSet.expectedNextResolvedStop.Equal(dataSet.ef.NextResolvedStop.Time) {
				t.Errorf("expected next stop = %v, got = %v", dataSet.expectedNextResolvedStop, dataSet.ef.NextResolvedStop)
			}
		})
	}
}
