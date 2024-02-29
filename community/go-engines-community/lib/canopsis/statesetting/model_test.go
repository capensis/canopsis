package statesetting

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func TestGetJunitState(t *testing.T) {
	var dataSets = []struct {
		testName      string
		thresholds    JUnitThresholds
		skipped       int64
		errors        int64
		failures      int64
		total         int64
		expectedState int
	}{
		{
			testName: "percentage simple critical",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       16,
			errors:        6,
			failures:      8,
			total:         38,
			expectedState: types.AlarmStateCritical,
		},
		{
			testName: "percentage simple major",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        6,
			failures:      8,
			total:         38,
			expectedState: types.AlarmStateMajor,
		},
		{
			testName: "percentage simple minor",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        6,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateMinor,
		},
		{
			testName: "percentage simple ok",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        1,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "percentage always ok",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       38,
			errors:        38,
			failures:      38,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "percentage always minor",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    0,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateMinor,
		},
		{
			testName: "percentage shouldn't be minor if 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    0,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       0,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "percentage always major",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    0,
					Critical: 100,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateMajor,
		},
		{
			testName: "percentage shouldn't be major if 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       0,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "percentage always critical",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 0,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateCritical,
		},
		{
			testName: "percentage shouldn't be major if 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 0,
					Type:     TypePercentage,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypePercentage,
				},
			},
			skipped:       0,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "number simple critical",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
			},
			skipped:       16,
			errors:        6,
			failures:      31,
			total:         38,
			expectedState: types.AlarmStateCritical,
		},
		{
			testName: "number simple major",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
			},
			skipped:       1,
			errors:        6,
			failures:      21,
			total:         38,
			expectedState: types.AlarmStateMajor,
		},
		{
			testName: "number simple minor",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
			},
			skipped:       1,
			errors:        11,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateMinor,
		},
		{
			testName: "number simple ok",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
			},
			skipped:       1,
			errors:        1,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "number always ok",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       38,
			errors:        38,
			failures:      38,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "number always minor",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    0,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       1,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateMinor,
		},
		{
			testName: "number shouldn't be minor if 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    0,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       0,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "number always major",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    0,
					Critical: 100,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       1,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateMajor,
		},
		{
			testName: "number shouldn't be major if 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       0,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "number always critical",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 0,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       1,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateCritical,
		},
		{
			testName: "number shouldn't be major if 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 0,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
				Failures: JUnitThreshold{
					Minor:    100,
					Major:    100,
					Critical: 100,
					Type:     TypeNumber,
				},
			},
			skipped:       0,
			errors:        0,
			failures:      0,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "mixed critical",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       11,
			errors:        11,
			failures:      15,
			total:         38,
			expectedState: types.AlarmStateCritical,
		},
		{
			testName: "mixed major",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       11,
			errors:        11,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateMajor,
		},
		{
			testName: "mixed minor",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       11,
			errors:        1,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateMinor,
		},
		{
			testName: "mixed ok",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        1,
			failures:      1,
			total:         38,
			expectedState: types.AlarmStateOK,
		},
		{
			testName: "ok if total = 0",
			thresholds: JUnitThresholds{
				Skipped: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypeNumber,
				},
				Errors: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
				Failures: JUnitThreshold{
					Minor:    10,
					Major:    20,
					Critical: 30,
					Type:     TypePercentage,
				},
			},
			skipped:       1,
			errors:        1,
			failures:      1,
			total:         0,
			expectedState: types.AlarmStateOK,
		},
	}

	for _, dataSet := range dataSets {
		t.Run(dataSet.testName, func(t *testing.T) {
			resultState := dataSet.thresholds.GetState(dataSet.skipped, dataSet.errors, dataSet.failures, dataSet.total)
			if resultState != dataSet.expectedState {
				t.Errorf("Expected state = %d, got = %d", dataSet.expectedState, resultState)
			}
		})
	}
}
