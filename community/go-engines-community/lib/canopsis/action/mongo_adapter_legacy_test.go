package action_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/action"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mkAdapter() (action.Adapter, mongo.Collection) {
	mgo, err := mongo.NewSession(mongo.Timeout)
	if err != nil {
		panic(err)
	}
	coll := action.DefaultCollection(mgo)
	_, err = coll.RemoveAll()
	if err != nil {
		panic(err)
	}
	adapter := action.NewAdapterLegacy(coll)
	return adapter, coll
}

func TestMongoAdapter_GetEnabled_ByPriority(t *testing.T) {
	adapter, collection := mkAdapter()

	Convey("get enabled action with priority should work fine", t, func() {
		_, err := collection.RemoveAll()
		So(err, ShouldBeNil)

		err = collection.Insert(action.Scenario{
			ID: "scenario_with_priority_10",
			Actions: []action.Action{
				{
					Type: "test",
					Parameters: map[string]interface{}{
						"author": "test",
						"output": "this is test",
					},
				},
			},
			Enabled:  true,
			Priority: 10,
		})
		So(err, ShouldBeNil)

		err = collection.Insert(action.Scenario{
			ID: "scenario_with_priority_2",
			Actions: []action.Action{
				{
					Type: "test",
					Parameters: map[string]interface{}{
						"author": "test",
						"output": "this is test",
					},
				},
			},
			Enabled:  true,
			Priority: 2,
		})
		So(err, ShouldBeNil)

		err = collection.Insert(action.Scenario{
			ID: "scenario_with_priority_1",
			Actions: []action.Action{
				{
					Type: "test",
					Parameters: map[string]interface{}{
						"author": "test",
						"output": "this is test",
					},
				},
			},
			Enabled:  true,
			Priority: 1,
		})
		So(err, ShouldBeNil)

		acts, err := adapter.GetEnabled()
		So(err, ShouldBeNil)
		So(acts, ShouldHaveLength, 3)
		So(acts[0].Priority, ShouldEqual, 1)
		So(acts[1].Priority, ShouldEqual, 2)
		So(acts[2].Priority, ShouldEqual, 10)

		err = collection.Insert(action.Scenario{
			ID: "scenario_with_priority_1_2",
			Actions: []action.Action{
				{
					Type: "test",
					Parameters: map[string]interface{}{
						"author": "test",
						"output": "this is test",
					},
				},
			},
			Enabled:  true,
			Priority: 1,
		})
		So(err, ShouldBeNil)
		acts, err = adapter.GetEnabled()
		So(err, ShouldBeNil)
		So(acts, ShouldHaveLength, 4)
		So(acts[0].Priority, ShouldEqual, 1)
		So(acts[1].Priority, ShouldEqual, 1)
		So(acts[2].Priority, ShouldEqual, 2)
		So(acts[3].Priority, ShouldEqual, 10)
		So(acts[0].ID, ShouldEqual, "scenario_with_priority_1")
		So(acts[1].ID, ShouldEqual, "scenario_with_priority_1_2")
	})
}
