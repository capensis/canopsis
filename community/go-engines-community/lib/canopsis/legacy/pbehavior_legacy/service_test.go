package pbehavior_legacy_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/legacy/pbehavior_legacy"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

func testNewService() (pbehavior_legacy.Service, error) {
	mongo, err := mongo.NewSession(mongo.Timeout)

	if err != nil {
		return nil, err
	}

	pbehaviorCollection := pbehavior_legacy.DefaultCollection(mongo)
	_, err = pbehaviorCollection.RemoveAll()
	if err != nil {
		return nil, err
	}
	pbehaviorAdapter := pbehavior_legacy.NewAdapter(pbehaviorCollection)
	ps := pbehavior_legacy.NewService(pbehaviorAdapter, log.NewTestLogger())

	return ps, nil
}

func getPBehavior() types.PBehaviorLegacy {
	start := types.CpsTime{Time: time.Now()}
	end := types.CpsTime{Time: time.Now().Add(time.Hour * 2)}

	p := types.PBehaviorLegacy{
		ID:            "fk_ra_hm",
		Author:        "Frank Klepacki",
		Connector:     "connector",
		ConnectorName: "connector_name",
		Enabled:       true,
		Filter:        `{"enabled": true}`,
		Name:          "Hell march",
		Reason:        "Red Alert",
		RRule:         "",
		Start:         &start,
		Stop:          &end,
		Type:          "Maintenance",
	}

	return p
}

func GetStep() types.AlarmStep {
	return types.AlarmStep{
		Type:      "",
		Timestamp: types.CpsTime{Time: time.Now()},
		Author:    "Kane",
		Message:   "He who controls the past, commands the future. He who commands the future, conquers the past.",
		Value:     types.AlarmStateMajor,
	}
}

func GetAlarmValue() types.AlarmValue {
	stateStep := GetStep()

	return types.AlarmValue{
		State:          &stateStep,
		Status:         nil,
		ACK:            nil,
		Canceled:       nil,
		Ticket:         nil,
		Resolved:       nil,
		Steps:          make(types.AlarmSteps, 0),
		Tags:           make([]string, 0),
		CreationDate:   types.CpsTime{Time: time.Now()},
		LastUpdateDate: types.CpsTime{Time: time.Now()},
		InitialOutput:  "",
		Snooze:         nil,
		HardLimit:      nil,
		Extra:          make(map[string]interface{}),
		Component:      "brotherhood",
		Connector:      "connector",
		ConnectorName:  "connector_name",
		Resource:       "",
	}
}

func getAlarm() types.Alarm {
	v := GetAlarmValue()

	a := types.Alarm{
		ID:       "c_and_c",
		Time:     types.CpsTime{Time: time.Now()},
		EntityID: "c_and_c",
		Value:    v,
	}

	return a
}

func TestPBehaviorCRUD(t *testing.T) {
	Convey("Given a pbehavior service", t, func() {
		ps, err := testNewService()
		So(err, ShouldBeNil)
		So(ps, ShouldNotBeNil)

		Convey("Then there is no pbehavior at first", func() {
			var pbs []types.PBehaviorLegacy
			pbs, err := ps.Get(bson.M{})
			So(err, ShouldBeNil)
			So(len(pbs), ShouldEqual, 0)
		})

		Convey("Then we can create a new one", func() {
			pb := getPBehavior()
			err := ps.Insert(pb)
			So(err, ShouldBeNil)

			pbs, err := ps.Get(bson.M{})
			So(err, ShouldBeNil)
			So(len(pbs), ShouldEqual, 1)
		})
	})
}

func TestAlarmHasPBehavior(t *testing.T) {

	Convey("Given a pbehavior service", t, func() {
		ps, err := testNewService()
		So(err, ShouldBeNil)
		So(ps, ShouldNotBeNil)
		alarm := getAlarm()

		Convey("Then a basic alarm should not be pbehaviored", func() {
			So(ps.AlarmHasPBehavior(alarm), ShouldBeFalse)
		})

		Convey("Then with a pbehavior on it, the alarm should be pbehaviored", func() {
			pbehavior := getPBehavior()
			pbehavior.Eids = []string{alarm.EntityID}
			// FixMe : implement ComputePbehavior() and remove the previous line
			err := ps.Insert(pbehavior)
			So(err, ShouldBeNil)
			So(ps.AlarmHasPBehavior(alarm), ShouldBeTrue)
		})

		Convey("Then with a pbehavior NOT on it, the alarm should NOT be pbehaviored", func() {
			pbehavior := getPBehavior()
			pbehavior.Eids = []string{}
			// FixMe : implement ComputePbehavior() and remove the previous line
			err := ps.Insert(pbehavior)
			So(err, ShouldBeNil)
			So(ps.AlarmHasPBehavior(alarm), ShouldBeFalse)
		})
	})

}
