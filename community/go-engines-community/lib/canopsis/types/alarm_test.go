package types_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func GetStep() types.AlarmStep {
	return types.AlarmStep{
		Type:      "",
		Timestamp: types.CpsTime{Time: time.Now()},
		Author:    "Hodor",
		Message:   "HODOR",
		Value:     types.AlarmStateMajor,
	}
}

func GetAlarmValue() types.AlarmValue {
	stateStep := GetStep()

	return types.AlarmValue{
		Component:      "test",
		Connector:      "test_c",
		ConnectorName:  "test_c_name",
		Resource:       "",
		State:          &stateStep,
		Status:         nil,
		ACK:            nil,
		Canceled:       nil,
		CreationDate:   types.CpsTime{Time: time.Now()},
		Extra:          make(map[string]interface{}),
		HardLimit:      nil,
		InitialOutput:  "",
		LastUpdateDate: types.CpsTime{Time: time.Now()},
		Resolved:       nil,
		Snooze:         nil,
		Steps:          make(types.AlarmSteps, 0),
		Tags:           make([]string, 0),
		Ticket:         nil,
	}
}

func getAlarm() types.Alarm {
	v := GetAlarmValue()

	a := types.Alarm{
		ID:       "abcde",
		Time:     types.CpsTime{Time: time.Now()},
		EntityID: "ab",
		Value:    v,
	}

	return a
}

func TestAlarmFromEvents(t *testing.T) {
	Convey("Given an alarm from an event", t, func() {
		e := getEvent()
		c := testutils.GetTestConf()

		a, err := types.NewAlarm(e, c)

		if err != nil {
			t.Fatal(err)
		}

		Convey("The state is Major", func() {
			So(a.CurrentState(), ShouldEqual, types.AlarmStateMajor)
		})
		Convey("The status is Ongoing", func() {
			So(a.CurrentStatus(c), ShouldEqual, types.AlarmStatusOngoing)
		})

		Convey("Given another event on the same alarm", func() {
			e.State = types.AlarmStateMinor
			a.Update(e, c)

			Convey("The state is Minor", func() {
				So(a.CurrentState(), ShouldEqual, types.AlarmStateMinor)
			})
			Convey("The state type is decreased", func() {
				So(a.Value.State.Type, ShouldEqual, types.AlarmStepStateDecrease)
			})
		})

		Convey("Given a locked alarm and a critical event", func() {
			a.Value.State.Type = types.AlarmStepChangeState

			e.State = types.AlarmStateCritical
			a.Update(e, c)

			Convey("The state is still Major", func() {
				So(a.CurrentState(), ShouldEqual, types.AlarmStateMajor)
			})
		})

		Convey("Given a locked alarm and a ok event", func() {
			e.State = types.AlarmStateOK
			a.Update(e, c)

			Convey("Then the state is now OK", func() {
				So(a.CurrentState(), ShouldEqual, types.AlarmStateOK)
			})
			Convey("Then the alarm is not locked anymore", func() {
				So(a.IsStateLocked(), ShouldBeFalse)
			})
		})

		Convey("Given a locked alarm & ok event", func() {
			e.State = types.AlarmStateCritical
			a, err := types.NewAlarm(e, c)

			Convey("No error on new alarm", func() {
				So(err, ShouldBeNil)
			})

			e.State = types.AlarmStateOK
			a.Value.State.Type = types.AlarmStepChangeState
			a.Update(e, c)
		})

		Convey("Then it has a valide display name", func() {
			So(a.Value.DisplayName, ShouldNotBeEmpty)
			So(string(a.Value.DisplayName[2]), ShouldEqual, "-")
		})
	})
}

func TestAlarmFromEventsNotCheck(t *testing.T) {
	Convey("Setup", t, func() {
		c := testutils.GetTestConf()
		eCheck := getEvent()
		a, err := types.NewAlarm(eCheck, c)
		if err != nil {
			t.Fatal(err)
		}

		eNotCheck := getEvent()
		eNotCheck.State = types.AlarmStateCritical
		eNotCheck.EventType = types.EventTypeAck
		stepsnum := len(a.Value.Steps)

		Convey("An event not of type check should not alter an alarm", func() {
			a.Update(eNotCheck, c)
			So(len(a.Value.Steps), ShouldEqual, stepsnum)
			So(a.Value.State.Value, ShouldEqual, types.AlarmStateMajor)
			So(a.Value.Status, ShouldNotBeNil)
			So(a.Value.Status.Value, ShouldEqual, types.AlarmStatusOngoing)
		})
	})
}

func TestAlarmComputeStatus(t *testing.T) {
	Convey("Given an alarm and a config", t, func() {
		a := getAlarm()
		c := testutils.GetTestConf()

		Convey("Then the status is Ongoing", func() {
			So(a.ComputeStatus(c), ShouldEqual, types.AlarmStatusOngoing)
		})

		Convey("and given a checked alarm", func() {
			a.Value.State.Value = types.AlarmStateOK
			Convey("Then the status is Off", func() {
				So(a.ComputeStatus(c), ShouldEqual, types.AlarmStatusOff)
			})
		})

		Convey("and given a cancelled alarm", func() {
			a.Value.Canceled = a.Value.State
			Convey("Then the status is Cancelled", func() {
				So(a.ComputeStatus(c), ShouldEqual, types.AlarmStatusCancelled)
				So(a.IsCanceled(), ShouldBeTrue)
			})
		})
	})
}

func TestAlarmComputeStatusWhenFlapping(t *testing.T) {
	Convey("Given an alarm and a config and a flapping alarm", t, func() {
		a := getAlarm()
		c := testutils.GetTestConf()

		Convey("When the alarm is Flapping", func() {
			c.Alarm.FlappingInterval = 3600
			c.Alarm.FlappingFreqLimit = 1
			So(c.Alarm.FlappingFreqLimit, ShouldEqual, 1)

			event := getEvent()
			event.State = types.AlarmStateCritical
			a.Update(event, c)
			event = getEvent()
			a.Update(event, c)

			Convey("Then the status is Flapping", func() {
				So(a.ComputeStatus(c), ShouldEqual, types.AlarmStatusFlapping)
			})
			c.Alarm.FlappingInterval = 0
			c.Alarm.FlappingFreqLimit = 0
		})
	})
}

func TestAlarmComputeStatusWhenStealthy(t *testing.T) {
	Convey("Given an alarm and a config", t, func() {
		a := getAlarm()
		c := testutils.GetTestConf()

		Convey("When the alarm is stealthy", func() {
			c.Alarm.StealthyInterval = 3600
			So(c.Alarm.StealthyInterval, ShouldEqual, 3600)

			event := getEvent()
			a.Update(event, c)
			event.State = types.AlarmStateOK
			event.Timestamp = types.CpsTime{Time: time.Now()}
			a.Update(event, c)

			Convey("Then the status is Stealthy", func() {
				So(a.ComputeStatus(c), ShouldEqual, types.AlarmStatusStealthy)
			})
			c.Alarm.StealthyInterval = 0
		})
	})
}

func TestAlarmCurrentStateWithState0(t *testing.T) {

	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("When state is Major", func() {
			Convey("Then GetCurrentState should be Major", func() {
				So(a.CurrentState(), ShouldEqual, types.AlarmStateMajor)
			})
		})

		Convey("When CurrentState Is < 0", func() {
			a.Value.State.Value = -1
			Convey("Then current state should be Ok", func() {
				So(a.CurrentState(), ShouldEqual, types.AlarmStateOK)
			})
		})
	})
}

func TestAlarmCurrentStatus(t *testing.T) {

	Convey("Given an alarm", t, func() {
		a := getAlarm()
		c := testutils.GetTestConf()

		So(a.CurrentStatus(c), ShouldEqual, types.AlarmStatusOngoing)

		Convey("When CurrentState Is < 0", func() {
			var statusStep = GetStep()
			a.Value.Status = &statusStep
			a.Value.Status.Value = -2
			Convey("Then A.getCurrentStatus should be Off", func() {
				So(a.CurrentStatus(c), ShouldEqual, types.AlarmStatusOff)
			})
		})

	})

}

func TestAlarmIsStateLocked(t *testing.T) {

	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("When the alarm is locked", func() {
			a.Value.State.Type = types.AlarmStepChangeState
			Convey("Then IsStateLocked should be true", func() {
				So(a.IsStateLocked(), ShouldBeTrue)
			})
		})

		Convey("The the alarm has another state", func() {
			a.Value.State.Type = "ack"

			Convey("Then IsStateLocked should be false", func() {
				So(a.IsStateLocked(), ShouldBeFalse)
			})
		})

		Convey("When the alarm has no state", func() {
			a.Value.State = nil

			Convey("Then IsStateLocked should be false", func() {
				So(a.IsStateLocked(), ShouldBeFalse)
			})
		})

	})

}

func TestAlarmFromEventMissingTS(t *testing.T) {
	Convey("Given a poorly created event, without timestamp", t, func() {
		evt := types.Event{}

		Convey("I ask to create an alarm from that event", func() {
			_, err := types.NewAlarm(evt, testutils.GetTestConf())
			So(err, ShouldNotBeNil)
		})
	})
}

func TestSnoozedAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()
		now := time.Now()

		Convey("Then an never snoozed alarm should not be snoozed", func() {
			So(a.Value.Snooze, ShouldBeNil)
			So(a.IsSnoozed(), ShouldBeFalse)
		})

		Convey("Then an just snoozed alarm should be snoozed", func() {
			inThreeYears := getEvent()
			duration := types.CpsNumber(now.Unix()) + types.CpsNumber(12*365*24*60*60)
			inThreeYears.Duration = &duration
			a.SnoozeFromEvent(inThreeYears)

			So(a.IsSnoozed(), ShouldBeTrue)
		})

		Convey("Then an obsolete snoozed alarm should not be snoozed anymore", func() {
			eventDate := types.CpsTime{
				Time: now.AddDate(0, 0, -24*60*60),
			}
			yesterday := getEvent()
			yesterday.Timestamp = eventDate
			duration := types.CpsNumber(eventDate.Unix()) + types.CpsNumber(1)
			yesterday.Duration = &duration
			a.SnoozeFromEvent(yesterday)

			So(a.IsSnoozed(), ShouldBeFalse)
		})
	})
}

func TestAcknowledgedAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("Then an never acked alarm should not be acked", func() {
			So(a.Value.ACK, ShouldBeNil)
			So(a.IsAck(), ShouldBeFalse)
		})

		Convey("Then an just acknowledged alarm should be acknowledged", func() {
			ackEvent := getEvent()
			ackEvent.EventType = types.EventTypeAck
			a.Ack(ackEvent)

			So(a.IsAck(), ShouldBeTrue)

			Convey("Then we should be able to unacknowledged the alarm", func() {
				ackEvent := getEvent()
				ackEvent.EventType = types.EventTypeAckremove
				So(a.IsAck(), ShouldBeTrue)

				a.Unack(ackEvent)
				So(a.IsAck(), ShouldBeFalse)
			})
		})
	})
}

func TestCanceledAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()
		now := types.CpsTime{Time: time.Now()}

		Convey("Then an never canceled alarm should not be canceled", func() {
			So(a.Value.Canceled, ShouldBeNil)
			So(a.IsCanceled(), ShouldBeFalse)
		})

		Convey("Then an just canceled alarm should be canceled", func() {
			cancelEvent := getEvent()
			cancelEvent.EventType = types.EventTypeCancel
			a.Cancel(cancelEvent)

			So(a.IsCanceled(), ShouldBeTrue)

			Convey("Then we should be able to uncanceled the alarm", func() {
				cancelEvent := getEvent()
				cancelEvent.EventType = types.EventTypeUncancel
				So(a.IsCanceled(), ShouldBeTrue)

				a.Uncancel(cancelEvent)
				So(a.IsCanceled(), ShouldBeFalse)
			})

			Convey("and a resolved alarm should not be canceled anymore", func() {
				a.ResolveCancel(&now)

				So(a.IsCanceled(), ShouldBeFalse)
			})
		})
	})
}

func TestChangeStateAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("Then we should be able to change his state", func() {
			changeEvent := getEvent()
			changeEvent.EventType = types.EventTypeChangestate
			step := a.ChangeState(changeEvent)

			So(step.Type, ShouldEqual, types.EventTypeChangestate)
			So(step.Value, ShouldEqual, changeEvent.State)
		})
	})
}

func TestCommentAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("Then we should be able to generate a comment on it", func() {
			commentEvent := getEvent()
			commentEvent.EventType = types.EventTypeComment
			step := a.Comment(commentEvent)

			So(step.Type, ShouldEqual, types.EventTypeComment)
			So(step.Message, ShouldEqual, commentEvent.Output)
		})
	})
}

func TestDoneAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("Then we should be able to mark it as done", func() {
			commentEvent := getEvent()
			commentEvent.EventType = types.EventTypeDone
			step := a.Done(commentEvent)

			So(step.Type, ShouldEqual, types.EventTypeDone)
			So(a.Done, ShouldNotBeNil)
		})
	})
}

func TestAssocTicketAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("Then a never ticketed alarm should not be ticketed", func() {
			So(a.Value.Ticket, ShouldBeNil)
		})

		Convey("Then a just associated ticketed alarm should be ticketed", func() {
			assocTicketEvent := getEvent()
			assocTicketEvent.EventType = types.EventTypeAssocTicket
			assocTicketEvent.Ticket = "Liet Kynes"
			a.AssocTicket(assocTicketEvent)

			So(a.Value.Ticket, ShouldNotBeNil)
			So(a.Value.Ticket.Message, ShouldEqual, "Liet Kynes")
			So(a.Value.Ticket.Value, ShouldEqual, "Liet Kynes") // TODO: remove that
		})
	})
}

func TestIsMatchedAlarm(t *testing.T) {
	Convey("Given an alarm", t, func() {
		a := getAlarm()

		Convey("Then it should be matched with a correct regex", func() {
			fields := []string{"Fear", "Resource", "Component"}
			So(a.IsMatched(".*test", fields), ShouldBeTrue)
		})
	})
}

func TestBagot(t *testing.T) {

	Convey("Given an alarm from an event", t, func() {
		e := getEvent()
		c := config.CanopsisConf{}
		a, err := types.NewAlarm(e, c)

		So(a.Value.StateChangesSinceStatusUpdate, ShouldEqual, types.CpsNumber(0))
		So(a.Value.TotalStateChanges, ShouldEqual, types.CpsNumber(1))

		if err != nil {
			t.Fatal(err)
		}

		Convey("Given another event on the same alarm", func() {
			e.State = types.AlarmStateMinor
			a.Update(e, c)

			So(a.Value.StateChangesSinceStatusUpdate, ShouldEqual, types.CpsNumber(1))
			So(a.Value.TotalStateChanges, ShouldEqual, types.CpsNumber(2))

			Convey("Given another event on the same alarm", func() {
				e.State = types.AlarmStateCritical
				a.Update(e, c)

				So(a.Value.StateChangesSinceStatusUpdate, ShouldEqual, types.CpsNumber(2))
				So(a.Value.TotalStateChanges, ShouldEqual, types.CpsNumber(3))

				Convey("Given an event OK to fnish the alarm", func() {
					e.State = types.AlarmStateOK
					a.Update(e, c)

					So(a.Value.StateChangesSinceStatusUpdate, ShouldEqual, types.CpsNumber(0))
					So(a.Value.TotalStateChanges, ShouldEqual, types.CpsNumber(4))

				})

			})
		})

	})

}

func TestHasSingleAck(t *testing.T) {
	Convey("Given an alarm that has not been acknowledged", t, func() {
		alarm := getAlarm()

		Convey("HasSingleAck should return false", func() {
			So(alarm.HasSingleAck(), ShouldBeFalse)
		})
	})

	Convey("Given an alarm that has been acknowledged once", t, func() {
		alarm := getAlarm()

		ackEvent := getEvent()
		ackEvent.EventType = types.EventTypeAck
		alarm.Value.Steps.Add(alarm.Ack(ackEvent))

		Convey("HasSingleAck should return true", func() {
			So(alarm.HasSingleAck(), ShouldBeTrue)
		})
	})

	Convey("Given an alarm that has been acknowledged once, with an ackremove", t, func() {
		alarm := getAlarm()

		ackEvent := getEvent()
		ackEvent.EventType = types.EventTypeAck
		alarm.Value.Steps.Add(alarm.Ack(ackEvent))

		ackRemoveEvent := getEvent()
		ackRemoveEvent.EventType = types.EventTypeAckremove
		alarm.Value.Steps.Add(alarm.Unack(ackRemoveEvent))

		Convey("HasSingleAck should return true", func() {
			So(alarm.HasSingleAck(), ShouldBeTrue)
		})
	})
	Convey("Given an alarm that has been acknowledged twice", t, func() {
		alarm := getAlarm()

		ackEvent := getEvent()
		ackEvent.EventType = types.EventTypeAck
		alarm.Value.Steps.Add(alarm.Ack(ackEvent))
		alarm.Value.Steps.Add(alarm.Ack(ackEvent))

		Convey("HasSingleAck should return false", func() {
			So(alarm.HasSingleAck(), ShouldBeFalse)
		})
	})

	Convey("Given an alarm that has been acknowledged twice, with an ackremove", t, func() {
		alarm := getAlarm()

		ackEvent := getEvent()
		ackEvent.EventType = types.EventTypeAck
		alarm.Value.Steps.Add(alarm.Ack(ackEvent))

		ackRemoveEvent := getEvent()
		ackRemoveEvent.EventType = types.EventTypeAckremove
		alarm.Value.Steps.Add(alarm.Unack(ackRemoveEvent))

		alarm.Value.Steps.Add(alarm.Ack(ackEvent))

		Convey("HasSingleAck should return false", func() {
			So(alarm.HasSingleAck(), ShouldBeFalse)
		})
	})
}
