package hook_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/hook"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

type eventPatternListWrapper struct {
	PatternList pattern.EventPatternList `bson:"list"`
}
type entityPatternListWrapper struct {
	PatternList pattern.EntityPatternList `bson:"list"`
}

func TestTriggeredByAlarmChangeType(t *testing.T) {
	Convey("Hook triggers", t, func() {
		hook := hook.Hook{}
		hook.Triggers = []string{"create", "declare_ticket"}

		createChange := types.AlarmChangeTypeCreate
		cancelChange := types.AlarmChangeTypeCancel

		So(hook.IsTriggeredByAlarmChangeType(createChange), ShouldBeTrue)
		So(hook.IsTriggeredByAlarmChangeType(cancelChange), ShouldBeFalse)
	})
}

func TestTriggeredByAlarmChange(t *testing.T) {
	Convey("Hook triggers", t, func() {
		hook := hook.Hook{}
		hook.Triggers = []string{"create", "declare_ticket"}

		createAlarmChange := types.AlarmChange{
			Type:                types.AlarmChangeTypeCreate,
			PreviousState:       types.CpsNumber(2),
			PreviousStateChange: types.CpsTime{Time: time.Now()},
		}
		cancelAlarmChange := types.AlarmChange{
			Type:                types.AlarmChangeTypeCancel,
			PreviousState:       types.CpsNumber(2),
			PreviousStateChange: types.CpsTime{Time: time.Now()},
		}

		So(hook.IsTriggeredByAlarmChange(createAlarmChange), ShouldBeTrue)
		So(hook.IsTriggeredByAlarmChange(cancelAlarmChange), ShouldBeFalse)
	})
}

func TestMatchesEventPatterns(t *testing.T) {
	Convey("Given a valid BSON pattern", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				bson.M{
					"state": bson.M{
						">": 0,
						"<": 3,
					},
					"component": "component",
					"resource":  bson.M{"regex_match": "service-(?P<id>\\d+)"},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w eventPatternListWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

		var p pattern.EventPatternList
		p = w.PatternList
		So(p.IsSet(), ShouldBeTrue)
		So(p.IsValid(), ShouldBeTrue)

		hook := hook.Hook{}
		hook.EventPatterns = p
		eventMatched := types.Event{
			Component: "component",
			Resource:  "service-69",
			State:     2,
		}

		eventNotMatched := types.Event{
			Component: "component",
			Resource:  "service-abd",
			State:     2,
		}

		So(hook.EventPatterns.Matches(eventMatched), ShouldBeTrue)
		So(hook.EventPatterns.Matches(eventNotMatched), ShouldBeFalse)

	})
}

func TestMatchesEntityPatterns(t *testing.T) {
	Convey("Given a valid BSON pattern", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				bson.M{
					"enabled": true,
					"type":    bson.M{"regex_match": "abc-.*-def"},
					"infos": bson.M{
						"output": bson.M{
							"value": "debian9",
						},
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternListWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

		var p pattern.EntityPatternList
		p = w.PatternList
		So(p.IsSet(), ShouldBeTrue)
		So(p.IsValid(), ShouldBeTrue)

		hook := hook.Hook{}
		hook.EntityPatterns = p

		entityMatched := types.Entity{
			Type:    "abc-ghi-def",
			Enabled: true,
			Infos: map[string]types.Info{
				"output": types.Info{
					Value: "debian9",
				},
			},
		}

		entityWrongOutput := types.Entity{
			Type:    "abc-ghi-def",
			Enabled: true,
			Infos: map[string]types.Info{
				"output": types.Info{
					Value: "ubuntu",
				},
			},
		}

		entityWrongType := types.Entity{
			Type:    "NO GOOD",
			Enabled: true,
			Infos: map[string]types.Info{
				"output": types.Info{
					Value: "debian9",
				},
			},
		}
		So(hook.EntityPatterns.Matches(&entityMatched), ShouldBeTrue)
		So(hook.EntityPatterns.Matches(&entityWrongOutput), ShouldBeFalse)
		So(hook.EntityPatterns.Matches(&entityWrongType), ShouldBeFalse)

	})
}
