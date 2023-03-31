package types_test

import (
	"encoding/json"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
