package types_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAlarmStepsCrop(t *testing.T) {
	Convey("Setup", t, func() {
		steps := make(types.AlarmSteps, 0)
		st := types.AlarmStep{
			Type:    types.AlarmStepStatusIncrease,
			Value:   1,
			Author:  "test",
			Message: "test",
		}
		s1 := types.AlarmStep{
			Type:  types.AlarmStepStateIncrease,
			Value: 2,
		}
		s2 := types.AlarmStep{
			Type:  types.AlarmStepStateDecrease,
			Value: 1,
		}
		steps.Add(st)
		for i := 0; i < 10; i++ {
			steps.Add(s1)
			steps.Add(s2)
		}
		cropNum := 30
		Convey("Crop does not work for the first 21 steps", func() {
			_, update := steps.Crop(&st, cropNum)
			So(update, ShouldBeFalse)
			Convey("Crop does work for the 41 steps", func() {
				for i := 0; i < 10; i++ {
					steps.Add(s1)
					steps.Add(s2)
				}

				steps, update := steps.Crop(&st, cropNum)

				So(steps[0].Type, ShouldEqual, st.Type)
				So(steps[1].Type, ShouldEqual, types.AlarmStepStateCounter)

				nstep := steps[1]

				So(update, ShouldBeTrue)
				So(len(steps), ShouldEqual, cropNum+2)
				// Ten steps were cropped from the 40 croppable steps
				testCounter := types.CropCounter{
					StateChanges:  10,
					Stateinc:      5,
					Statedec:      5,
					StateInfo:     0,
					StateMinor:    5,
					StateMajor:    5,
					StateCritical: 0,
				}
				So(nstep.StateCounter, ShouldResemble, testCounter)
				So(nstep.Author, ShouldEqual, "test")
				So(nstep.Message, ShouldEqual, "test")
				So(nstep.Value, ShouldEqual, st.Value)

			})
		})
	})
}

func TestAlarmStepCropOldData(t *testing.T) {
	Convey("Setup", t, func() {
		s, err := mongo.NewSession(mongo.Timeout)
		So(err, ShouldBeNil)

		c := s.DB(canopsis.DbName).C(alarm.AlarmCollectionName)
		_, err = c.RemoveAll(nil)
		So(err, ShouldBeNil)
		steps := make(types.AlarmSteps, 0)
		s1 := types.AlarmStep{
			Value:     1,
			Message:   "coucou",
			Author:    "coucou",
			Timestamp: types.CpsTime{Time: time.Now()},
		}
		stateCounter := types.CropCounter{}
		stateCounter.Stateinc = 1
		stateCounter.Statedec = 2

		s2 := types.AlarmStep{
			Type:         types.AlarmStepStateCounter,
			StateCounter: stateCounter,
		}

		So(steps.Add(s1), ShouldBeNil)
		So(steps.Add(s2), ShouldBeNil)

		al := types.Alarm{
			ID: "cropalarm",
			Value: types.AlarmValue{
				Steps: steps,
			},
		}

		Convey("Inserting the alarm is not a problem", func() {
			So(c.Insert(al), ShouldBeNil)

			Convey("Getting the alarm back isn't a problem either", func() {
				var al types.Alarm
				So(c.FindId("cropalarm").One(&al), ShouldBeNil)

				So(len(al.Value.Steps), ShouldEqual, 2)
				So(al.Value.Steps[0].Value, ShouldEqual, 1)
				So(al.Value.Steps[0].StateCounter, ShouldResemble, types.CropCounter{})
				So(al.Value.Steps[1].Value, ShouldEqual, 0)
				So(al.Value.Steps[1].Type, ShouldEqual, types.AlarmStepStateCounter)
				So(al.Value.Steps[1].StateCounter, ShouldNotBeNil)

				sc := al.Value.Steps[1].StateCounter

				So(sc.Statedec, ShouldEqual, 2)
				So(sc.Stateinc, ShouldEqual, 1)
				So(sc.StateInfo, ShouldEqual, 0)
				So(sc.StateMinor, ShouldEqual, 0)
				So(sc.StateMajor, ShouldEqual, 0)
				So(sc.StateCritical, ShouldEqual, 0)
			})
		})
	})
}
