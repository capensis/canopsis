package action

import (
	"reflect"
	"testing"
)

func TestPriorityIntervalsSet(t *testing.T) {
	var dataSets = []struct {
		testName          string
		priorities        []int
		expectedIntervals map[int]int
		expectedKeys      []int
	}{
		{
			testName:          "always in lower bound",
			priorities:        []int{1, 2, 3, 4, 5},
			expectedIntervals: map[int]int{6: 0},
			expectedKeys:      []int{6},
		},
		{
			testName:          "double set 1",
			priorities:        []int{1, 1},
			expectedIntervals: map[int]int{2: 0},
			expectedKeys:      []int{2},
		},
		{
			testName:          "double set 2",
			priorities:        []int{5, 5},
			expectedIntervals: map[int]int{1: 4, 6: 0},
			expectedKeys:      []int{1, 6},
		},
		{
			testName:   "in the middle",
			priorities: []int{10},
			expectedIntervals: map[int]int{
				1:  9,
				11: 0,
			},
			expectedKeys: []int{1, 11},
		},
		{
			testName:   "mixed",
			priorities: []int{10, 100, 5, 9},
			expectedIntervals: map[int]int{
				1:   4,
				6:   8,
				11:  99,
				101: 0,
			},
			expectedKeys: []int{1, 6, 11, 101},
		},
		{
			testName:   "backwards",
			priorities: []int{3, 2, 1},
			expectedIntervals: map[int]int{
				4: 0,
			},
			expectedKeys: []int{4},
		},
		{
			testName:   "with point creation",
			priorities: []int{3, 1},
			expectedIntervals: map[int]int{
				2: 2,
				4: 0,
			},
			expectedKeys: []int{2, 4},
		},
		{
			testName:   "with point creation and point destruction",
			priorities: []int{3, 1, 2},
			expectedIntervals: map[int]int{
				4: 0,
			},
			expectedKeys: []int{4},
		},
		{
			testName:   "do nothing with bad values",
			priorities: []int{0, -1},
			expectedIntervals: map[int]int{
				1: 0,
			},
			expectedKeys: []int{1},
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			pi := NewPriorityIntervals()

			for _, v := range dataset.priorities {
				pi.Take(v)
			}

			if !reflect.DeepEqual(dataset.expectedIntervals, pi.GetIntervals()) {
				t.Errorf("expected intervals = %v, got = %v", dataset.expectedIntervals, pi.GetIntervals())
			}

			if !reflect.DeepEqual(dataset.expectedKeys, pi.GetSortedKeys()) {
				t.Errorf("expected keys = %v, got = %v", dataset.expectedKeys, pi.GetSortedKeys())
			}
		})
	}
}

func TestPriorityIntervalsRestore(t *testing.T) {
	var dataSets = []struct {
		testName          string
		setPriorities     []int
		restorePriorities []int
		expectedIntervals map[int]int
		expectedKeys      []int
	}{
		{
			testName:          "set and restore",
			setPriorities:     []int{1},
			restorePriorities: []int{1},
			expectedIntervals: map[int]int{1: 0},
			expectedKeys:      []int{1},
		},
		{
			testName:          "double restore 1",
			setPriorities:     []int{1},
			restorePriorities: []int{1, 1},
			expectedIntervals: map[int]int{1: 0},
			expectedKeys:      []int{1},
		},
		{
			testName:          "double restore 2",
			setPriorities:     []int{5},
			restorePriorities: []int{5, 5},
			expectedIntervals: map[int]int{1: 0},
			expectedKeys:      []int{1},
		},
		{
			testName:          "restore, when left interval doesn't exist",
			setPriorities:     []int{1, 2, 3, 4, 10},
			restorePriorities: []int{4},
			expectedIntervals: map[int]int{4: 9, 11: 0},
			expectedKeys:      []int{4, 11},
		},
		{
			testName:          "point, which shouldn't be restored",
			setPriorities:     []int{},
			restorePriorities: []int{10},
			expectedIntervals: map[int]int{1: 0},
			expectedKeys:      []int{1},
		},
		{
			testName:          "with point creation",
			setPriorities:     []int{1, 2},
			restorePriorities: []int{1},
			expectedIntervals: map[int]int{1: 1, 3: 0},
			expectedKeys:      []int{1, 3},
		},
		{
			testName:          "between two intervals, merge both intervals",
			setPriorities:     []int{3},
			restorePriorities: []int{3},
			expectedIntervals: map[int]int{1: 0},
			expectedKeys:      []int{1},
		},
		{
			testName:          "between two intervals, merge with left interval",
			setPriorities:     []int{3, 4, 5},
			restorePriorities: []int{3},
			expectedIntervals: map[int]int{1: 3, 6: 0},
			expectedKeys:      []int{1, 6},
		},
		{
			testName:          "between two intervals, merge with right interval",
			setPriorities:     []int{3, 4, 5},
			restorePriorities: []int{5},
			expectedIntervals: map[int]int{1: 2, 5: 0},
			expectedKeys:      []int{1, 5},
		},
		{
			testName:          "between two intervals, without merge",
			setPriorities:     []int{3, 4, 5, 6, 7},
			restorePriorities: []int{4, 6},
			expectedIntervals: map[int]int{1: 2, 4: 4, 6: 6, 8: 0},
			expectedKeys:      []int{1, 4, 6, 8},
		},
		{
			testName:          "mixed",
			setPriorities:     []int{10, 100, 5, 9},
			restorePriorities: []int{10, 100, 5, 9},
			expectedIntervals: map[int]int{1: 0},
			expectedKeys:      []int{1},
		},
		{
			testName:          "do nothing with bad values",
			setPriorities:     []int{1},
			restorePriorities: []int{0, -1},
			expectedIntervals: map[int]int{2: 0},
			expectedKeys:      []int{2},
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			pi := NewPriorityIntervals()

			for _, v := range dataset.setPriorities {
				pi.Take(v)
			}

			for _, v := range dataset.restorePriorities {
				pi.Restore(v)
			}

			if !reflect.DeepEqual(dataset.expectedIntervals, pi.GetIntervals()) {
				t.Errorf("expected intervals = %v, got = %v", dataset.expectedIntervals, pi.GetIntervals())
			}

			if !reflect.DeepEqual(dataset.expectedKeys, pi.GetSortedKeys()) {
				t.Errorf("expected keys = %v, got = %v", dataset.expectedKeys, pi.GetSortedKeys())
			}
		})
	}
}
