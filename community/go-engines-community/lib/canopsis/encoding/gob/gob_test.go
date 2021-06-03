package gob_test

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"testing"
	"time"

	gogob "encoding/gob"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/gob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var gobRes []byte

func BenchmarkGOBEncoder(b *testing.B) {
	now := types.CpsTime{Time: time.Now()}
	event := types.Event{
		EventType:  types.EventTypeCheck,
		SourceType: types.SourceTypeResource,
		Component:  "benchmark",
		Resource:   "gob",
		State:      types.AlarmStateMajor,
		Timestamp:  now,
	}
	event.Format()
	alarm, err := types.NewAlarm(event, config.AlarmConfig{})

	if err != nil {
		b.Fatal(err)
	}

	gogob.Register(types.Alarm{})

	encoder := gob.NewEncoder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gobRes, _ = encoder.Encode(alarm)
	}
}
func BenchmarkGOBDecoder(b *testing.B) {
	now := types.CpsTime{Time: time.Now()}
	event := types.Event{
		EventType:  types.EventTypeCheck,
		SourceType: types.SourceTypeResource,
		Component:  "benchmark",
		Resource:   "gob",
		State:      types.AlarmStateMajor,
		Timestamp:  now,
	}
	event.Format()
	alarm, err := types.NewAlarm(event, config.AlarmConfig{})

	if err != nil {
		b.Fatal(err)
	}

	gogob.Register(types.Alarm{})

	encoder := gob.NewEncoder()
	decoder := gob.NewDecoder()
	bres, err := encoder.Encode(alarm)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := decoder.Decode(bres, &alarm); err != nil {
			b.Fatal(err)
		}
	}
}