package msgpack_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/msgpack"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

var globRes []byte
var alarmdec types.Alarm

func BenchmarkMsgPackEncoder(b *testing.B) {
	now := types.CpsTime{Time: time.Now()}
	event := types.Event{
		EventType:  types.EventTypeCheck,
		SourceType: types.SourceTypeResource,
		Component:  "benchmark",
		Resource:   "msgpack.enc",
		State:      types.AlarmStateMajor,
		Timestamp:  now,
	}
	event.Format()
	alarm, err := types.NewAlarm(event, testutils.GetTestConf())

	if err != nil {
		b.Fatal(err)
	}

	encoder := msgpack.NewEncoder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globRes, err = encoder.Encode(alarm)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMsgPackDecoder(b *testing.B) {
	now := types.CpsTime{Time: time.Now()}
	event := types.Event{
		EventType:  types.EventTypeCheck,
		SourceType: types.SourceTypeResource,
		Component:  "benchmark",
		Resource:   "msgpack.dec",
		State:      types.AlarmStateMajor,
		Timestamp:  now,
	}
	event.Format()
	alarm, err := types.NewAlarm(event, testutils.GetTestConf())

	if err != nil {
		b.Fatal(err)
	}

	encoder := msgpack.NewEncoder()
	decoder := msgpack.NewDecoder()
	bres, err := encoder.Encode(alarm)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := decoder.Decode(bres, &alarmdec); err != nil {
			b.Fatal(err)
		}
	}
}
