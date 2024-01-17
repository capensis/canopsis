package entitycounters_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func TestEntityCounters_GetWorstState_GivenNoRule(t *testing.T) {
	datasets := []struct {
		name               string
		counters           entitycounters.EntityCounters
		expectedWorstState int
	}{
		{
			name:               "no counters, worst state should be ok",
			counters:           entitycounters.EntityCounters{},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "only ok counters, worst state should be ok",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 4,
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "ok and minor counters, worst state should minor",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    4,
					Minor: 3,
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "ok, minor and major counters, worst state should major",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    4,
					Minor: 3,
					Major: 2,
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "ok, minor, major and critical counters, worst state should critical",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       4,
					Minor:    3,
					Major:    2,
					Critical: 1,
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
	}

	for _, dSet := range datasets {
		t.Run(dSet.name, func(t *testing.T) {
			resultWorstState := dSet.counters.GetWorstState()

			if resultWorstState != dSet.expectedWorstState {
				t.Errorf("expected worst state to be %d, but got %d", dSet.expectedWorstState, resultWorstState)
			}
		})
	}
}

func TestEntityCounters_GetWorstState_GivenDependencyRule(t *testing.T) {
	datasets := []struct {
		name               string
		counters           entitycounters.EntityCounters
		expectedWorstState int
	}{
		// only ok
		{
			name: "counters with rule for ok gt method by ok number, state should be ok, because ok number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for ok gt method by ok number, state should be minor, because ok number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for ok lt method by ok number, state should be minor, because ok number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for ok lt method by ok number, state should be ok, because ok number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for ok gt method by ok share, state should be ok, because ok share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 5,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for ok gt method by ok share, state should be critical, because ok share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       1,
					Minor:    3,
					Major:    3,
					Critical: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for ok lt method by ok share, state should be ok, because ok share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       1,
					Minor:    3,
					Major:    3,
					Critical: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for ok lt method by ok share, state should be major, because ok share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 5,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for minor gt method by ok number, state should be minor, because ok number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for major gt method by ok number, state should be major, because ok number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for critical gt method by ok number, state should be critical, because ok number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    3,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		// only minor
		{
			name: "counters with rule for minor gt method by minor number, state should be minor, because minor number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 3,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for minor gt method by minor number, state should be major, because minor number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for minor lt method by minor number, state should be major, because minor number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 3,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for minor lt method by minor number, state should be minor, because minor number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for minor gt method by minor share, state should be minor, because minor share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    5,
					Minor: 3,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for minor gt method by minor share, state should be critical, because minor share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       3,
					Minor:    1,
					Major:    3,
					Critical: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for minor lt method by minor share, state should be minor, because minor share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       3,
					Minor:    1,
					Major:    3,
					Critical: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for minor lt method by minor share, state should be major, because minor share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    5,
					Minor: 3,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for ok gt method by minor number, state should be ok, because minor number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 3,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for major gt method by minor number, state should be major, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 3,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for critical gt method by minor number, state should be critical, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 3,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		// only major
		{
			name: "counters with rule for major gt method by major number, state should be major, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for major gt method by major number, state should be critical, because major number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for major lt method by major number, state should be critical, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for major lt method by major number, state should be major, because major number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for major gt method by major share, state should be major, because major share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor:    5,
					Major:    3,
					Critical: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for major gt method by major share, state should be critical, because major share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       3,
					Minor:    3,
					Major:    1,
					Critical: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for major lt method by major share, state should be major, because major share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       3,
					Minor:    3,
					Major:    1,
					Critical: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for major lt method by major share, state should be critical, because major share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor:    5,
					Major:    3,
					Critical: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for ok gt method by major number, state should be ok, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for minor gt method by major number, state should be minor, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for critical gt method by major number, state should be critical, because major number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		// only critical
		{
			name: "counters with rule for critical gt method by critical number, state should be critical, because critical number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Critical: 3,
					Major:    1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for critical gt method by critical number, state should be major, because critical number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for critical lt method by critical number, state should be major, because critical number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Critical: 3,
					Major:    1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for critical lt method by critical number, state should be critical, because critical number is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondLT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for critical gt method by critical share, state should be critical, because critical share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor:    5,
					Critical: 3,
					Major:    2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for critical gt method by critical share, state should be major, because critical share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       3,
					Minor:    3,
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for critical lt method by critical share, state should be critical, because critical share is less than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:       3,
					Minor:    3,
					Major:    3,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for critical lt method by critical share, state should be major, because critical share is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor:    5,
					Critical: 3,
					Major:    2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondLT,
							Value:  20,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for ok gt method by critical number, state should be ok, because critical number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Critical: 3,
					Major:    1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for minor gt method by critical number, state should be ok, because critical number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Critical: 3,
					Major:    1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for major gt method by critical number, state should be ok, because critical number is greater than a threshold",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Critical: 3,
					Major:    1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleCritical,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		// some overlapping cases
		{
			name: "counters with rule for minor state, state should be ok, because rule condition is not matched",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for ok and minor state, state should be minor despite the rule, because all rule condition are not matched, then worst is taken",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for major state, state should be minor, because rule condition is not matched",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for minor and major state, state should be ok, because both rule condition is not matched",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateOK,
		},
		{
			name: "counters with rule for ok, minor and major state, state should be major despite the rule, because all rule condition are not matched, then worst is taken",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		{
			name: "counters with rule for critical state, state should be minor, because rule condition is not matched",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMinor,
		},
		{
			name: "counters with rule for ok, minor, major and critical state, state should be critical despite the rule, because all rule condition are not matched, then worst is taken",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Critical: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateCritical,
		},
		{
			name: "counters with rule for minor and major state, state should be major, because all rule condition are matched, then the worst from them is taken",
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 3,
					Major: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Minor: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMinor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
						Major: &statesetting.StateThreshold{
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Cond:   statesetting.CalculationCondGT,
							Value:  2,
						},
					},
				},
			},
			expectedWorstState: types.AlarmStateMajor,
		},
		// border cases
	}

	for _, dSet := range datasets {
		t.Run(dSet.name, func(t *testing.T) {
			resultWorstState := dSet.counters.GetWorstState()

			if resultWorstState != dSet.expectedWorstState {
				t.Errorf("expected worst state to be %d, but got %d", dSet.expectedWorstState, resultWorstState)
			}
		})
	}
}
