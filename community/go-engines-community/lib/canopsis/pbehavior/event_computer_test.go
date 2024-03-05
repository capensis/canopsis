package pbehavior_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
)

func BenchmarkEventComputer_Compute(b *testing.B) {
	typesByID := map[string]pbehavior.Type{
		"test-inactive": {
			ID:       "test-active",
			Name:     "active",
			Type:     pbehavior.TypeInactive,
			Priority: 1,
		},
		"test-active": {
			ID:       "test-active",
			Name:     "active",
			Type:     pbehavior.TypeActive,
			Priority: 1,
		},
		"test-pause": {
			ID:       "test-pause",
			Name:     "active",
			Type:     pbehavior.TypePause,
			Priority: 2,
		},
		"test-maintenance": {
			ID:       "test-maintenance",
			Name:     "active",
			Type:     pbehavior.TypeMaintenance,
			Priority: 3,
		},
	}
	defaultTypes := map[string]string{
		pbehavior.TypeActive:   "test-active",
		pbehavior.TypeInactive: "test-inactive",
	}
	now := time.Now()
	params := pbehavior.PbhEventParams{
		Start: datetime.CpsTime{Time: now.Add(1 * time.Hour)},
		End:   datetime.CpsTime{Time: now.Add(3 * time.Hour)},
		RRule: "FREQ=DAILY;BYDAY=WE,TH,FR;BYHOUR=1,14,17",
		Type:  "test-active",
		Exdates: []pbehavior.Exdate{
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(-3 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(-time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(24 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(26 * time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(48 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(50 * time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(49 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(53 * time.Hour)},
				},
				Type: "test-pause",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(72 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(120 * time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(-3 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(-time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(24 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(26 * time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(48 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(50 * time.Hour)},
				},
				Type: "test-maintenance",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(49 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(53 * time.Hour)},
				},
				Type: "test-pause",
			},
			{
				Exdate: types.Exdate{
					Begin: datetime.CpsTime{Time: now.Add(72 * time.Hour)},
					End:   datetime.CpsTime{Time: now.Add(120 * time.Hour)},
				},
				Type: "test-maintenance",
			},
		},
	}
	computer := pbehavior.NewEventComputer(typesByID, defaultTypes)
	span := timespan.New(now, now.Add(7*24*time.Hour))

	for i := 0; i < b.N; i++ {
		_, err := computer.Compute(params, span)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
