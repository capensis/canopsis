package eventfilter_test

import (
	"context"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"sort"
	"testing"

	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_context "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/context"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
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

func testNewService(ctrl *gomock.Controller, data ...bson.M) eventfilter.Service {
	adapter := mock_eventfilter.NewMockAdapter(ctrl)

	rules := make([]eventfilter.Rule, len(data))
	for i, v := range data {
		b, err := bson.Marshal(v)
		if err != nil {
			panic(err)
		}
		err = bson.Unmarshal(b, &rules[i])
		if err != nil {
			panic(err)
		}
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	dbClient := mock_mongo.NewMockDbClient(ctrl)

	adapter.EXPECT().List(gomock.Any()).Return(rules, nil)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get()
	return eventfilter.NewService(dbClient, adapter, mockTimezoneConfigProvider, log.NewTestLogger())
}

func TestProcessEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	Convey("Given an event filter service with (in order of priority) a break rule, and enrichment rule and a drop rule", t, func() {
		service := testNewService(
			ctrl,
			dropRule,
			breakRule,
			enrichmentRule,
		)

		enrichmentCenter := mock_context.NewMockEnrichmentCenter(ctrl)

		So(service.LoadDataSourceFactories(enrichmentCenter, "."), ShouldBeNil)
		err := service.LoadRules(ctx)
		So(err, ShouldBeNil)

		Convey("An event that is matched only by the first rule (break) is left unchanged", func() {
			event, report, err := service.ProcessEvent(ctx, eventCheck3)
			So(err, ShouldBeNil)
			So(event, ShouldResemble, eventCheck3)
			So(report.EntityUpdated, ShouldBeFalse)
		})

		Convey("An event that is matched only by the first two rules (break and enrichment) is enriched", func() {
			event, report, err := service.ProcessEvent(ctx, eventCheck0)
			So(err, ShouldBeNil)
			So(event.Output, ShouldEqual, "modified output")
			So(report.EntityUpdated, ShouldBeFalse)
		})

		Convey("An event that is matched by the three rules is dropped", func() {
			_, _, err := service.ProcessEvent(ctx, eventCheck2)
			So(err, ShouldNotBeNil)

			_, isDropError := err.(eventfilter.DropError)
			So(isDropError, ShouldBeTrue)
		})
	})
}

func TestEntityEnrichment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	Convey("Given an event filter service with an entity enrichment rule", t, func() {
		service := testNewService(
			ctrl,
			entityEnrichmentRule,
		)

		enrichmentCenterEntity := types.NewEntity("component", "name1", types.EntityTypeComponent, nil, nil, nil)
		eventEntity := types.NewEntity("component", "name2", types.EntityTypeComponent, nil, nil, nil)

		Convey("An event with an entity is not modified", func() {
			enrichmentCenter := mock_context.NewMockEnrichmentCenter(ctrl)

			So(service.LoadDataSourceFactories(enrichmentCenter, "."), ShouldBeNil)
			err := service.LoadRules(ctx)
			So(err, ShouldBeNil)

			event := types.Event{
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeComponent,
				Connector:     "connector",
				ConnectorName: "connector-name",
				Component:     "component",
				State:         3,
				Debug:         true,
				Entity:        &eventEntity,
			}
			event, report, err := service.ProcessEvent(ctx, event)
			So(err, ShouldBeNil)

			So(event.Entity, ShouldNotBeNil)
			So(*event.Entity, ShouldResemble, eventEntity)
			So(report.EntityUpdated, ShouldBeFalse)
		})

		Convey("An event without an entity is enriched", func() {
			enrichmentCenter := mock_context.NewMockEnrichmentCenter(ctrl)
			enrichmentCenter.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(&enrichmentCenterEntity, libcontext.UpdatedEntityServices{}, nil)

			So(service.LoadDataSourceFactories(enrichmentCenter, "."), ShouldBeNil)
			err := service.LoadRules(ctx)
			So(err, ShouldBeNil)

			event := types.Event{
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeComponent,
				Connector:     "connector",
				ConnectorName: "connector-name",
				Component:     "component",
				State:         3,
				Debug:         true,
			}
			event, report, err := service.ProcessEvent(ctx, event)
			So(err, ShouldBeNil)

			So(event.Entity, ShouldNotBeNil)
			So(*event.Entity, ShouldResemble, enrichmentCenterEntity)
			So(report.EntityUpdated, ShouldBeFalse)
		})
	})
}
