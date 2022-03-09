package types_test

import (
	gocontext "context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
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
		err := steps.Add(st)
		So(err, ShouldBeNil)
		for i := 0; i < 10; i++ {
			err = steps.Add(s1)
			So(err, ShouldBeNil)
			err = steps.Add(s2)
			So(err, ShouldBeNil)
		}
		cropNum := 30
		Convey("Crop does not work for the first 21 steps", func() {
			_, update := steps.Crop(&st, cropNum)
			So(update, ShouldBeFalse)
			Convey("Crop does work for the 41 steps", func() {
				for i := 0; i < 10; i++ {
					err = steps.Add(s1)
					So(err, ShouldBeNil)
					err = steps.Add(s2)
					So(err, ShouldBeNil)
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
	ctx, cancel := gocontext.WithCancel(gocontext.Background())
	defer cancel()

	Convey("Setup", t, func() {
		s, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
		So(err, ShouldBeNil)

		c := s.Collection(mongo.AlarmMongoCollection)
		_, err = c.DeleteMany(ctx, bson.M{})
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
			_, err := c.InsertOne(ctx, al)
			So(err, ShouldBeNil)

			Convey("Getting the alarm back isn't a problem either", func() {
				var al types.Alarm
				So(c.FindOne(ctx, bson.M{"_id": "cropalarm"}).Decode(&al), ShouldBeNil)

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
