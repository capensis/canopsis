package correlation_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
)

func randomIndexes(n int) []int {
	sequence := make([]int, n)
	for i := 0; i < n; i++ {
		sequence[i] = i
	}

	for i := range sequence {
		j := rand.Intn(i + 1)
		sequence[i], sequence[j] = sequence[j], sequence[i]
	}

	return sequence
}

func TestMetaAlarmState_PushChild_BasicShift(t *testing.T) {
	now := time.Now()
	timeInterval := int64(60)

	times := make(map[int]int64)
	times[0] = now.Unix()
	times[1] = now.Add(time.Second * 10).Unix()
	times[2] = now.Add(time.Second * 20).Unix()
	times[3] = now.Add(time.Second * 30).Unix()
	times[4] = now.Add(time.Second * 40).Unix()
	times[5] = now.Add(time.Second * 50).Unix()

	//fill alarm group randomly, should be sorted by time
	state := correlation.MetaAlarmState{}
	for _, idx := range randomIndexes(len(times)) {
		state.PushChild("test_alarm_"+strconv.Itoa(idx), times[idx], timeInterval)
	}

	if len(state.ChildrenEntityIDs) != 6 {
		t.Fatalf("children length should be %d, but got %d", 6, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Unix(), state.GetChildrenOpenTime())
	}

	for idx := 0; idx < 6; idx++ {
		if state.ChildrenTimestamps[idx] != times[idx] {
			t.Fatal("children timestamps have wrong order")
		}

		if state.ChildrenEntityIDs[idx] != "test_alarm_"+strconv.Itoa(idx) {
			t.Fatal("children ids have wrong order")
		}
	}

	//This call should shift time interval => so the storage should delete the first alarm in the Group
	// and update 'create' time, since the first alarm won't be in the minute time window anymore.
	state.PushChild("test_alarm_6", now.Add(time.Second*65).Unix(), timeInterval)
	if len(state.ChildrenEntityIDs) != 6 {
		t.Fatalf("children length should be %d, but got %d", 6, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*10).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*10).Unix(), state.GetChildrenOpenTime())
	}

	state.PushChild("test_alarm_7", now.Add(time.Second*300).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 1 {
		t.Fatalf("children length should be %d, but got %d", 1, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*300).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*300).Unix(), state.GetChildrenOpenTime())
	}
}

func TestMetaAlarmState_PushChild_BackInTimeShift(t *testing.T) {
	now := time.Now()
	timeInterval := int64(60)

	state := correlation.MetaAlarmState{}
	state.PushChild("test_alarm", now.Unix(), timeInterval)
	state.PushChild("test_alarm_2", now.Add(time.Second*-10).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 2 {
		t.Fatalf("children length should be %d, but got %d", 2, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*-10).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*-10).Unix(), state.GetChildrenOpenTime())
	}
}

func TestMetaAlarmState_PushChild_BackInTimeShiftWithoutAlarmLoss(t *testing.T) {
	now := time.Now()
	timeInterval := int64(60)

	state := correlation.MetaAlarmState{}
	state.PushChild("test_alarm", now.Unix(), timeInterval)
	state.PushChild("test_alarm_2", now.Add(time.Second*10).Unix(), timeInterval)
	state.PushChild("test_alarm_3", now.Add(time.Second*20).Unix(), timeInterval)
	state.PushChild("test_alarm_4", now.Add(time.Second*40).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 4 {
		t.Fatalf("children length should be %d, but got %d", 4, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Unix(), state.GetChildrenOpenTime())
	}

	//test_alarm_4 will be missed, so there shouldn't be any interval shifting
	state.PushChild("test_alarm_5", now.Add(time.Second*-30).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 4 {
		t.Fatalf("children length should be %d, but got %d", 4, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Unix(), state.GetChildrenOpenTime())
	}

	//Interval can be shifted, since none alarm will be lost
	state.PushChild("test_alarm_6", now.Add(time.Second*-5).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 5 {
		t.Fatalf("children length should be %d, but got %d", 5, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*-5).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*-5).Unix(), state.GetChildrenOpenTime())
	}
}

func TestMetaAlarmState_PushChild_ShiftWithUpdates(t *testing.T) {
	now := time.Now()
	timeInterval := int64(60)

	state := correlation.MetaAlarmState{}
	state.PushChild("test_alarm", now.Unix(), timeInterval)
	state.PushChild("test_alarm_2", now.Add(time.Second*5).Unix(), timeInterval)
	state.PushChild("test_alarm_3", now.Add(time.Second*10).Unix(), timeInterval)
	state.PushChild("test_alarm_4", now.Add(time.Second*15).Unix(), timeInterval)
	state.PushChild("test_alarm_5", now.Add(time.Second*20).Unix(), timeInterval)
	state.PushChild("test_alarm_6", now.Add(time.Second*25).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 6 {
		t.Fatalf("children length should be %d, but got %d", 4, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Unix(), state.GetChildrenOpenTime())
	}

	state.PushChild("test_alarm", now.Add(time.Second*40).Unix(), timeInterval)
	state.PushChild("test_alarm_2", now.Add(time.Second*55).Unix(), timeInterval)
	state.PushChild("test_alarm_3", now.Add(time.Second*45).Unix(), timeInterval)
	state.PushChild("test_alarm_4", now.Add(time.Second*30).Unix(), timeInterval)
	state.PushChild("test_alarm_5", now.Add(time.Second*35).Unix(), timeInterval)
	state.PushChild("test_alarm_6", now.Add(time.Second*40).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 6 {
		t.Fatalf("children length should be %d, but got %d", 4, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*30).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*30).Unix(), state.GetChildrenOpenTime())
	}

	//This call should shift time interval, but no alarms should be deleted, since they belong to the new time interval.
	state.PushChild("test_alarm_7", now.Add(time.Second*65).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 7 {
		t.Fatalf("children length should be %d, but got %d", 7, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*30).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*30).Unix(), state.GetChildrenOpenTime())
	}
}

func TestMetaAlarmState_RemoveChildrenBefore(t *testing.T) {
	now := time.Now()
	timeInterval := int64(60)

	state := correlation.MetaAlarmState{}
	state.PushChild("test_alarm", now.Unix(), timeInterval)
	state.PushChild("test_alarm_2", now.Add(time.Second*5).Unix(), timeInterval)
	state.PushChild("test_alarm_3", now.Add(time.Second*10).Unix(), timeInterval)
	state.PushChild("test_alarm_4", now.Add(time.Second*15).Unix(), timeInterval)
	state.PushChild("test_alarm_5", now.Add(time.Second*20).Unix(), timeInterval)
	state.PushChild("test_alarm_6", now.Add(time.Second*25).Unix(), timeInterval)

	if len(state.ChildrenEntityIDs) != 6 {
		t.Fatalf("children length should be %d, but got %d", 4, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Unix(), state.GetChildrenOpenTime())
	}

	state.RemoveChildrenBefore(now.Add(time.Second * 20).Unix())

	if len(state.ChildrenEntityIDs) != 2 {
		t.Fatalf("children length should be %d, but got %d", 2, len(state.ChildrenEntityIDs))
	}

	if state.GetChildrenOpenTime() != now.Add(time.Second*20).Unix() {
		t.Fatalf("children open time should be %d, but got %d", now.Add(time.Second*20).Unix(), state.GetChildrenOpenTime())
	}
}
