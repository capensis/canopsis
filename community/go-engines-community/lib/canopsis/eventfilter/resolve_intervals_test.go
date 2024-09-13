package eventfilter

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/teambition/rrule-go"
)

func TestUpdateIntervals(t *testing.T) {
	dataSets := map[string]struct {
		ef                        Rule
		rrule                     string
		now                       time.Time
		expectedResolvedStart     datetime.CpsTime
		expectedResolvedStop      datetime.CpsTime
		expectedNextResolvedStart *datetime.CpsTime
		expectedNextResolvedStop  *datetime.CpsTime
	}{
		"simple case, before the interval, shouldn't be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"the case, where resolved start and resolved stop is nil, should be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     nil,
				ResolvedStop:      nil,
				NextResolvedStart: nil,
				NextResolvedStop:  nil,
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"simple case, in the beginning of the interval, shouldn't be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"simple case, in the middle of the interval, shouldn't be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"simple case, after the interval, should be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC)},
		},
		"before the interval, but next interval is not calculated, should be updated": {
			ef: Rule{
				Start:         &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"in the interval, but next interval is not calculated, should be updated": {
			ef: Rule{
				Start:         &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"after the interval, but next interval is not calculated, should be updated": {
			ef: Rule{
				Start:         &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;INTERVAL=7",
			now:                       time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 16, 0, 0, 0, 0, time.UTC)},
		},
		"before the interval, but next interval is not calculated, cross intervals, should be updated": {
			ef: Rule{
				Start:         &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
		},
		"in the interval, cross intervals, shouldn't be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
		},
		"after the interval, cross intervals, should be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 3, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
		},
		"after the interval, with counts, next interval should be nil, because of exceeded count, should be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY;COUNT=5",
			now:                       time.Date(2022, 1, 6, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: nil,
			expectedNextResolvedStop:  nil,
		},
		"after the interval, with counts, next interval should be nil, because of exceeded count, shouldn't be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: nil,
				NextResolvedStop:  nil,
			},
			rrule:                     "FREQ=DAILY;COUNT=5",
			now:                       time.Date(2022, 1, 7, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: nil,
			expectedNextResolvedStop:  nil,
		},
		"after the interval, long cross intervals with counts, should be updated": {
			ef: Rule{
				Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
				ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				NextResolvedStart: nil,
				NextResolvedStop:  nil,
			},
			rrule:                     "FREQ=DAILY;COUNT=5",
			now:                       time.Date(2022, 1, 7, 0, 0, 0, 1, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 4, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 9, 0, 0, 0, 0, time.UTC)},
		},
		"test with old intervals, they should be be updated": {
			ef: Rule{
				Start:         &datetime.CpsTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				Stop:          &datetime.CpsTime{Time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
				ResolvedStart: &datetime.CpsTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				ResolvedStop:  &datetime.CpsTime{Time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=DAILY",
			now:                       time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
		},
		"weird test": {
			ef: Rule{
				Start:         &datetime.CpsTime{Time: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC)},
				Stop:          &datetime.CpsTime{Time: time.Date(2022, 1, 1, 9, 20, 0, 0, time.UTC)},
				ResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC)},
				ResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 1, 9, 20, 0, 0, time.UTC)},
			},
			rrule:                     "FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16",
			now:                       time.Date(2022, 1, 4, 10, 30, 0, 0, time.UTC),
			expectedResolvedStart:     datetime.CpsTime{Time: time.Date(2022, 1, 4, 10, 20, 0, 0, time.UTC)},
			expectedResolvedStop:      datetime.CpsTime{Time: time.Date(2022, 1, 4, 10, 40, 0, 0, time.UTC)},
			expectedNextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 4, 10, 40, 0, 0, time.UTC)},
			expectedNextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 4, 11, 0, 0, 0, time.UTC)},
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

			if opt.Count != 0 || dataSet.ef.ResolvedStart == nil || dataSet.ef.ResolvedStart.IsZero() {
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

func BenchmarkCalculateRRuleAllCalculated(b *testing.B) {
	ef := &Rule{
		Start:             &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
		Stop:              &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
		ResolvedStart:     &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
		ResolvedStop:      &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
		NextResolvedStart: &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
		NextResolvedStop:  &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
	}

	now := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	opt, err := rrule.StrToROption("FREQ=DAILY")
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	r, err := rrule.NewRRule(*opt)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	r.DTStart(ef.Start.Time)

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = &datetime.CpsTime{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)}
		ef.ResolvedStop = &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)}
		ef.NextResolvedStart = &datetime.CpsTime{Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)}
		ef.NextResolvedStop = &datetime.CpsTime{Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)}

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleDailyInTheBeginningOfTheInterval(b *testing.B) {
	now := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		"FREQ=DAILY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleDailyOneDayGap(b *testing.B) {
	now := time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		"FREQ=DAILY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleDailyOneWeekGap(b *testing.B) {
	now := time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		"FREQ=DAILY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleDailyOneMonthGap(b *testing.B) {
	now := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		"FREQ=DAILY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleDailyOneYearGap(b *testing.B) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		"FREQ=DAILY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleHourlyOneDayGap(b *testing.B) {
	now := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
		"FREQ=HOURLY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleHourlyOneWeekGap(b *testing.B) {
	now := time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
		"FREQ=HOURLY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleHourlyOneMonthGap(b *testing.B) {
	now := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
		"FREQ=HOURLY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleHourlyOneYearGap(b *testing.B) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
		"FREQ=HOURLY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleMinutelyOneDayGap(b *testing.B) {
	now := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 0, 1, 0, 0, time.UTC),
		"FREQ=MINUTELY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleMinutelyOneWeekGap(b *testing.B) {
	now := time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 0, 1, 0, 0, time.UTC),
		"FREQ=MINUTELY",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleComplexOneDayGap(b *testing.B) {
	now := time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 9, 20, 0, 0, time.UTC),
		"FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleComplexOneWeekGap(b *testing.B) {
	now := time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 9, 20, 0, 0, time.UTC),
		"FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleComplexOneMonthGap(b *testing.B) {
	now := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 9, 20, 0, 0, time.UTC),
		"FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func BenchmarkCalculateRRuleComplexOneYearGap(b *testing.B) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	ef, r, err := getBenchmarkData(
		time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 1, 9, 20, 0, 0, time.UTC),
		"FREQ=MINUTELY;INTERVAL=20;BYHOUR=9,10,11,12,13,14,15,16",
	)
	if err != nil {
		b.Fatalf("error is not expected = %s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ef.ResolvedStart = nil
		ef.ResolvedStop = nil
		ef.NextResolvedStart = nil
		ef.NextResolvedStop = nil

		ResolveIntervals(ef, r, now, time.UTC)
	}
}

func getBenchmarkData(start, stop time.Time, rruleString string) (*Rule, *rrule.RRule, error) {
	ef := Rule{
		Start: &datetime.CpsTime{Time: start},
		Stop:  &datetime.CpsTime{Time: stop},
	}

	opt, err := rrule.StrToROption(rruleString)
	if err != nil {
		return nil, nil, err
	}

	r, err := rrule.NewRRule(*opt)
	if err != nil {
		return nil, nil, err
	}

	r.DTStart(ef.Start.Time)

	return &ef, r, nil
}
