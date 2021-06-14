package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func getEvent() types.Event {
	var event = types.Event{
		Connector:     "red",
		ConnectorName: "is",
		Component:     "dead",
		Resource:      "adieu_yourï",
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		LongOutput:    "",
		State:         types.AlarmStateMajor,
		Output:        "",
		PerfData:      nil,
		PerfDataArray: nil,
		ID:            nil,
		Status:        nil,
		StateType:     nil,
		Alarm:         nil,
		Entity:        nil,
	}
	event.Format()

	return event
}

// Generate a bad event
func GetBadEvent() types.Event {

	status := types.CpsNumber(-3)

	return types.Event{
		Connector:     "red",
		ConnectorName: "is",
		Component:     "nil !",
		Resource:      "",
		EventType:     "wrong_event_type",
		PerfData:      nil,
		PerfDataArray: nil,
		ID:            nil,
		Status:        &status,
		StateType:     nil,
		SourceType:    "wrong_source_type",
		LongOutput:    "",
		State:         -1,
		Output:        "",
		Alarm:         nil,
		Entity:        nil,
	}
}

func GetBadEvent2() types.Event {
	evt := GetBadEvent()
	evt.EventType = ""
	return evt
}

func TestPartialID(t *testing.T) {
	Convey("HeartBeatItem and Event", t, func() {
		li := heartbeat.NewItem(time.Second * 2)
		li.AddMapping("connector", "test")
		li.AddMapping("component", "testc")

		bevent := []byte(`{"connector": "test", "component": "testc"}`)
		var ievent types.GenericEvent

		Convey("GenericEvent can JSONUnmarshal", func() {
			So(ievent.JSONUnmarshal(bevent), ShouldBeNil)

			Convey("I must have a valid PartialID without error", func() {
				id, err := ievent.PartialID(li)
				So(id, ShouldEqual, "component:testc.connector:test")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestEmptyPartialID(t *testing.T) {
	Convey("Empty PartialID", t, func() {
		li := heartbeat.NewItem(0)
		bevent := []byte(`{"connector":"test"}`)
		var ievent types.GenericEvent

		Convey("GenericEvent can JSONUnmarshal", func() {
			So(ievent.JSONUnmarshal(bevent), ShouldBeNil)

			Convey("GenericEvent can PartialID", func() {
				id, err := ievent.PartialID(li)
				So(id, ShouldEqual, "")
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestBadJSONPartialID(t *testing.T) {
	Convey("Bad Partial ID", t, func() {
		var ievent types.GenericEvent
		sevent := `"bla"`
		li := heartbeat.NewItem(0)
		li.AddMapping("connector", "bla")

		Convey("GenericEvent can JSONUnmarshal", func() {
			err := ievent.JSONUnmarshal([]byte(sevent))
			So(err, ShouldBeNil)

			Convey("PartialID raises an error", func() {
				_, err = ievent.PartialID(li)

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "cannot fetch data")
			})
		})

	})
}

func TestBadJSON2PartialID(t *testing.T) {
	Convey("Given a bad json document", t, func() {
		var ievent types.GenericEvent
		ev := []byte(`{"bla": ["bla"]}`)

		Convey("I must unmarshal it", func() {
			err := ievent.JSONUnmarshal(ev)
			So(err, ShouldBeNil)

			Convey("But i cannot get a partial id from a bad key", func() {
				li := heartbeat.NewItem(0)
				li.AddMapping("bla", "vert")
				id, err := ievent.PartialID(li)

				So(err, ShouldBeNil)
				So(id, ShouldEqual, "bla:bla")
				So(li.ID(), ShouldNotEqual, id)
				So(li.ID(), ShouldEqual, "bla:vert")
			})
		})
	})
}

func TestNoFieldPartialID(t *testing.T) {
	var ievent types.GenericEvent
	ev := []byte(`{"connector": "bla"}`)

	err := ievent.JSONUnmarshal(ev)
	if err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}

	li := heartbeat.NewItem(0)
	li.AddMapping("ZBLEURRGANRSTUIEBL", "^_^")
	id, err := ievent.PartialID(li)
	if err == nil {
		t.Fatal("no error returned")
	}
	if id != "" {
		t.Fatalf("got non empty id: %s", id)
	}
}

func TestEventTypeIsReplaced(t *testing.T) {
	Convey("Given an event with EventType empty", t, func() {
		evt := GetBadEvent2()

		Convey("EventType goes to EventTypeCheck when Format() is called", func() {
			evt.Format()
			So(evt.EventType, ShouldEqual, types.EventTypeCheck)
		})
	})
}

func TestEventIsGood(t *testing.T) {

	Convey("Given an event", t, func() {
		e := getEvent()

		Convey("Then the event is valid", func() {
			So(e.IsValid(), ShouldBeNil)
		})

		Convey("Then the event is correctly formatted", func() {
			e.Format()
			So(e.Timestamp, ShouldNotBeNil)
			So(e.EventType, ShouldNotBeBlank)
		})

		Convey("Then the event is catched correctly", func() {
			fields := []string{"Fear", "Resource", "Component"}
			So(e.IsMatched(".*dead", fields), ShouldBeTrue)
		})
	})
}

func TestEventIsBad(t *testing.T) {

	Convey("Given a bad event", t, func() {
		e := GetBadEvent()

		Convey("The event is not valid", func() {
			So(e.IsValid(), ShouldNotBeNil)
		})

		Convey("The event is correctly formatted anyway", func() {
			e.Format()
			So(e.Timestamp, ShouldNotBeNil)
			So(e.EventType, ShouldNotBeNil)
			So(e.EventType, ShouldNotBeBlank)
			So(e.Output, ShouldNotBeNil)
		})
	})

	Convey("Given an event without resource", t, func() {
		e := getEvent()
		e.Resource = ""

		Convey("The event is not valid", func() {
			So(e.IsValid(), ShouldNotBeNil)
		})
	})

	Convey("Given a malformed perf event", t, func() {
		e := getEvent()
		e.EventType = types.EventTypePerf

		Convey("The event is not valid", func() {
			So(e.IsValid(), ShouldNotBeNil)
		})
	})
}

func TestListTypeFields(t *testing.T) {
	Convey("Given an event struct", t, func() {
		event := getEvent()
		Convey("Then basic fields are listed by listTypeFields()", func() {
			tags := event.GetRequiredKeys()
			So(tags, ShouldContain, "connector")
			So(tags, ShouldContain, "connector_name")
			So(tags, ShouldContain, "component")
		})
	})
}

func TestGenerateContextIds(t *testing.T) {

	Convey("Given an event", t, func() {
		e := getEvent()

		Convey("The generated context informations are corrects", func() {
			infos := e.GenerateContextInformations()
			So(infos[0].ID, ShouldEqual, "red/is")
			So(infos[0].Impacts, ShouldResemble, []string{"adieu_yourï/dead"})
			So(infos[0].Depends, ShouldResemble, []string{"dead"})
			So(infos[1].ID, ShouldEqual, "dead")
			So(infos[1].Impacts, ShouldResemble, []string{"red/is"})
			So(infos[1].Depends, ShouldResemble, []string{"adieu_yourï/dead"})
			So(infos[2].ID, ShouldEqual, "adieu_yourï/dead")
			So(infos[2].Type, ShouldEqual, types.SourceTypeResource)
			So(infos[2].Impacts, ShouldResemble, []string{"dead"})
			So(infos[2].Depends, ShouldResemble, []string{"red/is"})
		})

		Convey("The generated context informations are corrects, event without resource", func() {
			e.SourceType = types.SourceTypeComponent
			e.Resource = ""

			infos := e.GenerateContextInformations()
			So(infos[0].ID, ShouldEqual, "red/is")
			So(infos[0].Impacts, ShouldResemble, []string{})
			So(infos[0].Depends, ShouldResemble, []string{"dead"})
			So(infos[1].ID, ShouldEqual, "dead")
			So(infos[1].Type, ShouldEqual, types.SourceTypeComponent)
			So(infos[1].Impacts, ShouldResemble, []string{"red/is"})
			So(infos[1].Depends, ShouldResemble, []string{})
		})
	})
}

var ctxInfos []types.ContextInformation

func BenchmarkGenerateContextInformations(b *testing.B) {
	e := getEvent()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctxInfos = e.GenerateContextInformations()
		if len(ctxInfos) == 0 {
			b.Fatal("no infos")
		}
	}
}

func TestInjectExtraInfos(t *testing.T) {
	Convey("Setup", t, func() {
		evt := `{
		"event_type":"check",
		"component":"bla",
		"resource":"blurk",
		"state":3,
		"connector":"bla",
		"connector_name":"bla",
		"personnemeconnait":"ulyss31",
		"personnemeconnait2":"ulyss62"
	}`
		var e types.Event
		So(json.Unmarshal([]byte(evt), &e), ShouldBeNil)

		Convey("Extra informations are into ExtraInfos", func() {
			So(e.InjectExtraInfos([]byte(evt)), ShouldBeNil)
			So(e.ExtraInfos, ShouldContainKey, "personnemeconnait")
			So(e.ExtraInfos, ShouldContainKey, "personnemeconnait2")
			So(e.ExtraInfos["personnemeconnait"], ShouldEqual, "ulyss31")
			So(e.ExtraInfos["personnemeconnait2"], ShouldEqual, "ulyss62")
		})
	})
}

func TestSetField(t *testing.T) {
	Convey("Given an event", t, func() {
		evt := `{
			"event_type":"check",
			"component":"bla",
			"resource":"blurk",
			"state":3,
			"connector":"bla",
			"connector_name":"bla",
			"extra_info":"ulyss31"
		}`
		var e types.Event
		So(json.Unmarshal([]byte(evt), &e), ShouldBeNil)
		So(e.InjectExtraInfos([]byte(evt)), ShouldBeNil)

		Convey("Setting a field that does not exist sets it in ExtraInfos", func() {
			So(e.SetField("extra_info", 12), ShouldBeNil)
			So(e.ExtraInfos["extra_info"], ShouldEqual, 12)

			So(e.SetField("new_info", "twelve"), ShouldBeNil)
			So(e.ExtraInfos, ShouldContainKey, "new_info")
			So(e.ExtraInfos["new_info"], ShouldEqual, "twelve")
		})

		Convey("Setting a field containing a CpsNumber works", func() {
			So(e.SetField("State", 2), ShouldBeNil)
			So(e.State, ShouldEqual, 2)
		})

		Convey("Setting a field containing a *CpsNumber works", func() {
			So(e.SetField("Status", 2), ShouldBeNil)
			So(*e.Status, ShouldEqual, 2)
		})

		Convey("Setting a field containing a CpsTime works", func() {
			So(e.SetField("Timestamp", 12), ShouldBeNil)
			So(e.Timestamp.Time.Unix(), ShouldEqual, 12)
		})

		Convey("Setting a field containing a string works", func() {
			So(e.SetField("EventType", "test"), ShouldBeNil)
			So(e.EventType, ShouldEqual, "test")
		})

		Convey("Setting a field containing a *string works", func() {
			So(e.SetField("ID", "test"), ShouldBeNil)
			So(*e.ID, ShouldEqual, "test")
		})

		Convey("Setting a field containing a bool works", func() {
			So(e.SetField("Debug", true), ShouldBeNil)
			So(e.Debug, ShouldEqual, true)
		})

		Convey("Setting a field containing an entity works", func() {
			entity := types.Entity{
				ID: "eid",
			}
			So(e.SetField("Entity", entity), ShouldBeNil)
			So(e.Entity, ShouldNotBeNil)
			So(e.Entity.ID, ShouldEqual, "eid")
		})
	})
}
