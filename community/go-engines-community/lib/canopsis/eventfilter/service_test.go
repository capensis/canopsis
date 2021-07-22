package eventfilter_test

import (
	"context"
	"sort"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	mockconfig "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mockeventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func testNewService(ctrl *gomock.Controller, data ...bson.M) eventfilter.Service {
	adapter := mockeventfilter.NewMockAdapter(ctrl)

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

	adapter.EXPECT().List(gomock.Any()).Return(rules, nil)
	mockTimezoneConfigProvider := mockconfig.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get()
	return eventfilter.NewService(adapter, mockTimezoneConfigProvider, log.NewTestLogger())
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

		So(service.LoadDataSourceFactories("."), ShouldBeNil)
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
