package eventfilter_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

var entityEnrichmentRule = bson.M{
	"_id":  "entity_enrichment",
	"type": "enrichment",
	"patterns": []bson.M{
		{"current_entity": nil},
	},
	"external_data": bson.M{
		"entity": bson.M{
			"type": "entity",
		},
	},
	"actions": []bson.M{
		bson.M{
			"type": "copy",
			"from": "ExternalData.entity",
			"to":   "Entity",
		},
	},
	"priority":   50,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "pass",
}

func testNewService(rules ...bson.M) eventfilter.Service {
	adapter := testNewAdapter(rules...)
	return eventfilter.NewService(adapter, log.NewTestLogger())
}

func TestProcessEvent(t *testing.T) {
	Convey("Given an event filter service with (in order of priority) a break rule, and enrichment rule and a drop rule", t, func() {
		service := testNewService(
			dropRule,
			breakRule,
			enrichmentRule,
		)

		entity := types.NewEntity("component", "name", types.EntityTypeComponent, nil, nil, nil)
		enrichmentCenter := MockEnrichmentCenter{
			Entity: &entity,
		}
		enrichFields := context.NewEnrichFields("", "")

		So(service.LoadDataSourceFactories(&enrichmentCenter, enrichFields, "."), ShouldBeNil)
		service.LoadRules()

		Convey("An event that is matched only by the first rule (break) is left unchanged", func() {
			event, report, err := service.ProcessEvent(eventCheck3)
			So(err, ShouldBeNil)
			So(event, ShouldResemble, eventCheck3)
			So(report.EntityUpdated, ShouldBeFalse)
		})

		Convey("An event that is matched only by the first two rules (break and enrichment) is enriched", func() {
			event, report, err := service.ProcessEvent(eventCheck0)
			So(err, ShouldBeNil)
			So(event.Output, ShouldEqual, "modified output")
			So(report.EntityUpdated, ShouldBeFalse)
		})

		Convey("An event that is matched by the three rules is dropped", func() {
			_, _, err := service.ProcessEvent(eventCheck2)
			So(err, ShouldNotBeNil)

			_, isDropError := err.(eventfilter.DropError)
			So(isDropError, ShouldBeTrue)
		})
	})
}

func TestEntityEnrichment(t *testing.T) {
	Convey("Given an event filter service with an entity enrichment rule", t, func() {
		service := testNewService(
			entityEnrichmentRule,
		)

		enrichmentCenterEntity := types.NewEntity("component", "name1", types.EntityTypeComponent, nil, nil, nil)
		eventEntity := types.NewEntity("component", "name2", types.EntityTypeComponent, nil, nil, nil)

		enrichmentCenter := MockEnrichmentCenter{
			Entity: &enrichmentCenterEntity,
		}
		enrichFields := context.NewEnrichFields("", "")

		So(service.LoadDataSourceFactories(&enrichmentCenter, enrichFields, "."), ShouldBeNil)
		service.LoadRules()

		Convey("An event with an entity is not modified", func() {
			event := types.Event{
				EventType:  types.EventTypeCheck,
				SourceType: types.SourceTypeComponent,
				Component:  "component",
				State:      3,
				Debug:      true,
				Entity:     &eventEntity,
			}
			event, report, err := service.ProcessEvent(event)
			So(err, ShouldBeNil)

			So(enrichmentCenter.HandleCalls, ShouldEqual, 0)
			So(event.Entity, ShouldNotBeNil)
			So(*event.Entity, ShouldResemble, eventEntity)
			So(report.EntityUpdated, ShouldBeFalse)
		})

		Convey("An event without an entity is enriched", func() {
			event := types.Event{
				EventType:  types.EventTypeCheck,
				SourceType: types.SourceTypeComponent,
				Component:  "component",
				State:      3,
				Debug:      true,
			}
			event, report, err := service.ProcessEvent(event)
			So(err, ShouldBeNil)

			So(enrichmentCenter.HandleCalls, ShouldEqual, 1)
			So(event.Entity, ShouldNotBeNil)
			So(*event.Entity, ShouldResemble, enrichmentCenterEntity)
			So(report.EntityUpdated, ShouldBeFalse)
		})
	})
}
