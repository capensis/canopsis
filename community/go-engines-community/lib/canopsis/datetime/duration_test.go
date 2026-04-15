package datetime_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

func TestShortDurationWithUnit_AsDuration(t *testing.T) {
	testCases := []struct {
		duration         datetime.ShortDurationWithUnit
		expectedDuration time.Duration
		expectedIsZero   bool
	}{
		{
			duration: datetime.ShortDurationWithUnit{
				Value: 100,
				Unit:  "ms",
			},
			expectedDuration: 100 * time.Millisecond,
			expectedIsZero:   false,
		},
		{
			duration: datetime.ShortDurationWithUnit{
				Value: 5,
				Unit:  "s",
			},
			expectedDuration: 5 * time.Second,
			expectedIsZero:   false,
		},
		{
			duration: datetime.ShortDurationWithUnit{
				Value: 100,
			},
			expectedDuration: 100 * time.Second,
			expectedIsZero:   false,
		},
		{
			duration:         datetime.ShortDurationWithUnit{},
			expectedDuration: 0,
			expectedIsZero:   true,
		},
	}

	for _, tc := range testCases {
		actual := tc.duration.AsDuration()
		if actual != tc.expectedDuration {
			t.Errorf("expected %v but got %v", tc.expectedDuration, actual)
		}
		isZero := tc.duration.IsZero()
		if isZero != tc.expectedIsZero {
			t.Errorf("expected %t but got %t", tc.expectedIsZero, isZero)
		}
	}
}
