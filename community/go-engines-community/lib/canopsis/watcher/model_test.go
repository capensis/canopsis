package watcher_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOutput(t *testing.T) {
	Convey("Given a watcher with an invalid template", t, func() {
		w := watcher.Watcher{
			OutputTemplate: "{{",
		}

		Convey("GetOutput should return an error", func() {
			_, err := w.GetOutput(watcher.AlarmCounters{})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a watcher with a template using undefined fields", t, func() {
		w := watcher.Watcher{
			OutputTemplate: "{{.NoSuchField}}",
		}

		Convey("GetOutput should return an error", func() {
			_, err := w.GetOutput(watcher.AlarmCounters{})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a watcher with a valid template", t, func() {
		w := watcher.Watcher{
			OutputTemplate: "{{.Alarms}},{{.State.Critical}},{{.State.Major}},{{.State.Minor}},{{.State.Info}},{{.Acknowledged}},{{.NotAcknowledged}}",
		}

		Convey("GetOutput should execute this template correctly", func() {
			output, err := w.GetOutput(watcher.AlarmCounters{
				Alarms: 1,
				State: watcher.StateCounters{
					Critical: 2,
					Major:    3,
					Minor:    4,
					Info:     5,
				},
				Acknowledged:    6,
				NotAcknowledged: 7,
			})
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "1,2,3,4,5,6,7")
		})
	})
}

func TestGetState(t *testing.T) {
	Convey("Given a watcher without method", t, func() {
		w := watcher.Watcher{
			State: map[string]interface{}{},
		}

		Convey("GetState should return an error", func() {
			_, err := w.GetState(watcher.AlarmCounters{})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a watcher with an unknown method", t, func() {
		w := watcher.Watcher{
			State: map[string]interface{}{
				"method": "unknown",
			},
		}

		Convey("GetState should return an error", func() {
			_, err := w.GetState(watcher.AlarmCounters{})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a watcher with the worst method", t, func() {
		w := watcher.Watcher{
			State: map[string]interface{}{
				"method": "worst",
			},
		}

		Convey("GetState should return the state of the worst alarm", func() {
			state, err := w.GetState(watcher.AlarmCounters{
				Alarms: 4,
				State: watcher.StateCounters{
					Critical: 2,
					Major:    1,
					Minor:    1,
					Info:     4,
				},
				Acknowledged:    2,
				NotAcknowledged: 2,
			})
			So(err, ShouldBeNil)
			So(state, ShouldEqual, types.AlarmStateCritical)

			state, err = w.GetState(watcher.AlarmCounters{
				Alarms: 2,
				State: watcher.StateCounters{
					Critical: 0,
					Major:    1,
					Minor:    1,
					Info:     4,
				},
				Acknowledged:    1,
				NotAcknowledged: 1,
			})
			So(err, ShouldBeNil)
			So(state, ShouldEqual, types.AlarmStateMajor)

			state, err = w.GetState(watcher.AlarmCounters{
				Alarms: 1,
				State: watcher.StateCounters{
					Critical: 0,
					Major:    0,
					Minor:    1,
					Info:     4,
				},
				Acknowledged:    1,
				NotAcknowledged: 0,
			})
			So(err, ShouldBeNil)
			So(state, ShouldEqual, types.AlarmStateMinor)

			state, err = w.GetState(watcher.AlarmCounters{
				Alarms: 1,
				State: watcher.StateCounters{
					Critical: 0,
					Major:    0,
					Minor:    0,
					Info:     4,
				},
				Acknowledged:    1,
				NotAcknowledged: 0,
			})
			So(err, ShouldBeNil)
			So(state, ShouldEqual, types.AlarmStateOK)
		})
	})
}
