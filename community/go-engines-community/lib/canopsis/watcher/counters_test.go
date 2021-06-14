package watcher_test

import (
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"testing"
)

// This example corresponds to the resolution of a critical (and acknowledged)
// alarm.
// In this case, the Alarms, Critical and Acknowledged counters are decremented
// (since there is one less acknowledged critical alarm that impacts the
// watcher), and Info is incremented (since there is one more entity in this
// state).
func ExampleGetCountersIncrementsFromStates_alarmResolution() {
	previousState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers"},
		HasAlarm:          true,
		AlarmState:        3,
		AlarmAcknowledged: true,
		IsEntityActive:    true,
	}
	currentState := watcher.DependencyState{
		EntityID:         "cpu/server1",
		ImpactedWatchers: []string{"servers"},
		HasAlarm:         false,
		IsEntityActive:   true,
	}

	increments := watcher.GetCountersIncrementsFromStates(
		previousState,
		currentState)

	json, err := json.MarshalIndent(increments, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
	// Output:
	// {
	//   "servers": {
	//     "All": -1,
	//     "Alarms": -1,
	//     "State": {
	//       "Critical": -1,
	//       "Major": 0,
	//       "Minor": 0,
	//       "Info": 1
	//     },
	//     "Acknowledged": -1,
	//     "NotAcknowledged": 0,
	//     "PbehaviorCounters": {}
	//   }
	// }
}

// This example correspond to the creation of a new watcher that is impacted by
// an existing entity (with a major alarm).
// In this case, the counters of the "servers" watcher are not updated (since
// the entity's state did not change). The Alarms, Major and NotAcknowledged
// counters of the "critical_servers" watcher are incremented, since there is a
// new non-acknowledged major alarm impacting this watcher.
func ExampleGetCountersIncrementsFromStates_newWatcher() {
	previousState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers"},
		HasAlarm:          true,
		AlarmState:        2,
		AlarmAcknowledged: false,
		IsEntityActive:    true,
	}
	currentState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers", "critical_servers"},
		HasAlarm:          true,
		AlarmState:        2,
		AlarmAcknowledged: false,
		IsEntityActive:    true,
	}

	increments := watcher.GetCountersIncrementsFromStates(
		previousState,
		currentState)

	json, err := json.MarshalIndent(increments, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
	// Output:
	// {
	//   "critical_servers": {
	//     "All": 1,
	//     "Alarms": 1,
	//     "State": {
	//       "Critical": 0,
	//       "Major": 1,
	//       "Minor": 0,
	//       "Info": 0
	//     },
	//     "Acknowledged": 0,
	//     "NotAcknowledged": 1,
	//     "PbehaviorCounters": null
	//   }
	// }
}

// This example corresponds to a pbehavior becoming inactive on an entity.
func ExampleGetCountersIncrementsFromStates_pbehaviorEnd() {
	previousState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers"},
		HasAlarm:          true,
		AlarmState:        1,
		AlarmAcknowledged: true,
		IsEntityActive:    false,
	}
	currentState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers"},
		HasAlarm:          true,
		AlarmState:        1,
		AlarmAcknowledged: true,
		IsEntityActive:    true,
	}

	increments := watcher.GetCountersIncrementsFromStates(
		previousState,
		currentState)

	json, err := json.MarshalIndent(increments, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
	// Output:
	// {
	//   "servers": {
	//     "All": 0,
	//     "Alarms": 1,
	//     "State": {
	//       "Critical": 0,
	//       "Major": 0,
	//       "Minor": 1,
	//       "Info": -1
	//     },
	//     "Acknowledged": 1,
	//     "NotAcknowledged": 0,
	//     "PbehaviorCounters": {}
	//   }
	// }
}

func TestGetCountersIncrementsFromStatesPbhCountersIssue(t *testing.T) {
	previousState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers"},
		HasAlarm:          true,
		AlarmState:        1,
		AlarmAcknowledged: true,
		IsEntityActive:    false,
		PbehaviorType:     "pbh1",
	}
	currentState := watcher.DependencyState{
		EntityID:          "cpu/server1",
		ImpactedWatchers:  []string{"servers"},
		HasAlarm:          true,
		AlarmState:        1,
		AlarmAcknowledged: true,
		IsEntityActive:    false,
		PbehaviorType:     "pbh2",
	}

	increment := watcher.GetCountersIncrementsFromStates(previousState, currentState)

	value, ok := increment["servers"]
	if !ok {
		t.Error("increment map should contain servers key")
	}

	pbh1incr, ok := value.PbehaviorCounters["pbh1"]
	if !ok {
		t.Fatal("pbehavior counters should contain key for pbh1")
	}

	if pbh1incr != -1 {
		t.Errorf("pbh1 increment should be -1, got %d", pbh1incr)
	}

	pbh2incr, ok := value.PbehaviorCounters["pbh2"]
	if !ok {
		t.Fatal("pbehavior counters should contain key for pbh2")
	}

	if pbh2incr != 1 {
		t.Errorf("pbh2 increment should be 1, got %d", pbh2incr)
	}
}