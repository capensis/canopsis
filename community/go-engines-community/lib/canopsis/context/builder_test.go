package context_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func getEventResource() types.Event {
	e := types.Event{
		Connector:     "conn1",
		ConnectorName: "conn1",
		Component:     "comp1",
		Resource:      "res1",
		SourceType:    types.SourceTypeResource,
	}
	e.Format()

	return e
}

func getEventUpdatewatcher() types.Event {
	e := types.Event{
		Connector:     "conn1",
		ConnectorName: "conn1",
		Component:     "122",
		EventType:     types.EventTypeUpdateWatcher,
		Resource:      "res1",
		SourceType:    types.SourceTypeResource,
	}
	e.Format()

	return e
}

// entityPatternListWrapper is a type that wraps an EntityPatternList into a
// struct. It is required to test the unmarshalling of an array into an
// EntityPatternList because bson.Unmarshal does not work when called with an
// array.
type entityPatternListWrapper struct {
	PatternList pattern.EntityPatternList `bson:"list"`
}

func TestUpdateWatcherHandle(t *testing.T) {
	Convey("Setup", t, func() {
		var entities []types.Entity

		var aWatcher watcher.Watcher
		mapWatcher := bson.M{
			"_id":   "122",
			"name":  "watcher1",
			"infos": bson.M{},
			"entities": []bson.M{
				bson.M{
					"name": "mockEntity1",
				},
			},
			"state":           bson.M{},
			"output_template": "",
		}

		bsonWatcher, err := bson.Marshal(mapWatcher)
		So(err, ShouldBeNil)
		So(bson.Unmarshal(bsonWatcher, &aWatcher), ShouldBeNil)
		aWatcher.EnsureInitialized()

		var mockEntity1 types.Entity
		mapEntity1 := bson.M{
			"_id":   "123",
			"name":  "mockEntity1",
			"infos": bson.M{},
			"state": bson.M{},
		}

		bsonEntity1, err := bson.Marshal(mapEntity1)
		So(err, ShouldBeNil)
		So(bson.Unmarshal(bsonEntity1, &mockEntity1), ShouldBeNil)
		mockEntity1.EnsureInitialized()
		entities = append(entities, mockEntity1)

		var mockEntity2 types.Entity
		mapEntity2 := bson.M{
			"_id":   "124",
			"name":  "mockEntity2",
			"infos": bson.M{},
			"state": bson.M{},
		}

		bsonEntity2, err := bson.Marshal(mapEntity2)
		So(err, ShouldBeNil)
		So(bson.Unmarshal(bsonEntity2, &mockEntity2), ShouldBeNil)
		mockEntity2.EnsureInitialized()
		entities = append(entities, mockEntity2)

		tea := newTestEntityAdapter()
		tea.BulkInsert(mockEntity1)
		tea.BulkInsert(mockEntity2)
		twa := newTestWatcherAdapterWithEntities(entities)
		ec := context.NewBuilder(tea, twa, log.NewTestLogger())

		twa.Insert(aWatcher)

		Convey("Updating a watcher that is not in cache yet works", func() {
			e := getEventUpdatewatcher()
			eventEntity := ec.UpdateLinkedEntities(e)
			eventEntity = ec.UpdateWatchersLinks(e, eventEntity)
			updates := ec.Extract()

			So(eventEntity, ShouldNotBeNil)

			So(updates, ShouldContainKey, mockEntity1.ID)
			So(updates[mockEntity1.ID].Entity.Impacts, ShouldContain, aWatcher.ID)
			So(updates[mockEntity1.ID].NewImpacts, ShouldContain, aWatcher.ID)

			So(updates, ShouldNotContainKey, mockEntity2.ID)

			So(updates, ShouldContainKey, aWatcher.ID)
			So(updates[aWatcher.ID].Entity.Depends, ShouldContain, mockEntity1.ID)
			So(updates[aWatcher.ID].Entity.Depends, ShouldNotContain, mockEntity2.ID)
			So(updates[aWatcher.ID].NewDepends, ShouldContain, mockEntity1.ID)
			So(updates[aWatcher.ID].NewDepends, ShouldNotContain, mockEntity2.ID)

			Convey("Updating a watcher that already is in cache works", func() {
				var newEntitiesWrapper entityPatternListWrapper
				mapNewEntities := bson.M{
					"list": []bson.M{
						bson.M{
							"name": "mockEntity1",
						},
						bson.M{
							"name": "mockEntity2",
						},
					},
				}

				bsonNewEntities, err := bson.Marshal(mapNewEntities)
				So(err, ShouldBeNil)
				So(bson.Unmarshal(bsonNewEntities, &newEntitiesWrapper), ShouldBeNil)
				aWatcher.Entities = newEntitiesWrapper.PatternList
				twa.Insert(aWatcher)

				e := getEventUpdatewatcher()
				eventEntity := ec.UpdateLinkedEntities(e)
				eventEntity = ec.UpdateWatchersLinks(e, eventEntity)
				updates := ec.Extract()

				So(eventEntity, ShouldNotBeNil)

				So(updates, ShouldContainKey, mockEntity1.ID)
				So(updates[mockEntity1.ID].Entity.Impacts, ShouldContain, aWatcher.ID)
				So(updates[mockEntity1.ID].NewImpacts, ShouldContain, aWatcher.ID)

				So(updates, ShouldContainKey, mockEntity2.ID)
				So(updates[mockEntity2.ID].Entity.Impacts, ShouldContain, aWatcher.ID)
				So(updates[mockEntity2.ID].NewImpacts, ShouldContain, aWatcher.ID)

				So(updates[aWatcher.ID].Entity.Depends, ShouldContain, mockEntity1.ID)
				So(updates[aWatcher.ID].Entity.Depends, ShouldContain, mockEntity2.ID)
				So(updates[aWatcher.ID].NewDepends, ShouldContain, mockEntity1.ID)
				So(updates[aWatcher.ID].NewDepends, ShouldContain, mockEntity2.ID)
			})
		})
	})
}

func TestUpdateEntityWatchersHandle(t *testing.T) {
	Convey("Setup", t, func() {
		e := getEventResource()
		entityID := strings.Join([]string{e.Resource, e.Component}, "/")

		var aWatcher watcher.Watcher
		mapWatcher := bson.M{
			"_id":   "122",
			"name":  "watcher1",
			"infos": bson.M{},
			"entities": []bson.M{
				bson.M{
					"_id": entityID,
				},
			},
			"state":           bson.M{},
			"output_template": "",
		}

		bsonWatcher, err := bson.Marshal(mapWatcher)
		So(err, ShouldBeNil)
		So(bson.Unmarshal(bsonWatcher, &aWatcher), ShouldBeNil)
		aWatcher.EnsureInitialized()

		tea := newTestEntityAdapter()
		tea.BulkInsert(aWatcher.Entity)
		twa := newTestWatcherAdapter()
		ec := context.NewBuilder(tea, twa, log.NewTestLogger())

		twa.Insert(aWatcher)

		ec.LoadWatchers()
		eventEntity := ec.UpdateLinkedEntities(e)
		eventEntity = ec.UpdateWatchersLinks(e, eventEntity)
		updates := ec.Extract()

		So(eventEntity, ShouldNotBeNil)

		So(updates, ShouldContainKey, entityID)
		So(updates[entityID].Entity.Impacts, ShouldContain, aWatcher.ID)
		So(updates[entityID].NewImpacts, ShouldContain, aWatcher.ID)

		So(updates, ShouldContainKey, aWatcher.ID)
		So(updates[aWatcher.ID].Entity.Depends, ShouldContain, entityID)
		So(updates[aWatcher.ID].NewDepends, ShouldContain, entityID)
	})
}

func TestBuilder(t *testing.T) {
	Convey("Setup", t, func() {
		tea := newTestEntityAdapter()
		twa := newTestWatcherAdapter()
		ec := context.NewBuilder(tea, twa, log.NewTestLogger())
		e := getEventResource()
		eventEntity := ec.UpdateLinkedEntities(e)
		eventEntity = ec.UpdateWatchersLinks(e, eventEntity)

		So(eventEntity, ShouldNotBeNil)
		So(eventEntity.Type, ShouldEqual, e.SourceType)

		updates := ec.Extract()

		So(len(updates), ShouldEqual, 3)

		So(updates, ShouldContainKey, "conn1/conn1")
		So(updates, ShouldContainKey, "res1/comp1")
		So(updates, ShouldContainKey, "comp1")

		So(updates["conn1/conn1"].Entity.Type, ShouldEqual, types.EntityTypeConnector)
		So(updates["res1/comp1"].Entity.Type, ShouldEqual, types.EntityTypeResource)
		So(updates["comp1"].Entity.Type, ShouldEqual, types.EntityTypeComponent)

		So(updates["conn1/conn1"].Entity.Impacts, ShouldContain, "res1/comp1")
		So(updates["res1/comp1"].Entity.Impacts, ShouldContain, "comp1")
		So(updates["comp1"].Entity.Impacts, ShouldContain, "conn1/conn1")
		So(updates["conn1/conn1"].NewImpacts, ShouldContain, "res1/comp1")
		So(updates["res1/comp1"].NewImpacts, ShouldContain, "comp1")
		So(updates["comp1"].NewImpacts, ShouldContain, "conn1/conn1")

		So(updates["conn1/conn1"].Entity.Depends, ShouldContain, "comp1")
		So(updates["res1/comp1"].Entity.Depends, ShouldContain, "conn1/conn1")
		So(updates["comp1"].Entity.Depends, ShouldContain, "res1/comp1")
		So(updates["conn1/conn1"].NewDepends, ShouldContain, "comp1")
		So(updates["res1/comp1"].NewDepends, ShouldContain, "conn1/conn1")
		So(updates["comp1"].NewDepends, ShouldContain, "res1/comp1")
	})
}

func BenchmarkBuilder(b *testing.B) {
	tea := newTestEntityAdapter()
	twa := newTestWatcherAdapter()
	ec := context.NewBuilder(tea, twa, log.NewTestLogger())
	e := getEventResource()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ec.UpdateLinkedEntities(e)
	}
}
