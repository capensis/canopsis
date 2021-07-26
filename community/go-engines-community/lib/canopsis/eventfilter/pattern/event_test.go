package pattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	mgobson "github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type eventPatternWrapper struct {
	Pattern pattern.EventPattern `bson:"pattern"`
}

func TestPatternUnmarshalMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern", t, func() {
		mapPattern := mgobson.M{
			"state": mgobson.M{
				">": 0,
				"<": 3,
			},
			"component": "component",
			"resource":  mgobson.M{"regex_match": "service-(?P<id>\\d+)"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var p pattern.EventPattern
			So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

			Convey("The values of the pattern should be correct", func() {
				So(p.State.Gt.Set, ShouldBeTrue)
				So(p.State.Gt.Value, ShouldEqual, 0)
				So(p.State.Lt.Set, ShouldBeTrue)
				So(p.State.Lt.Value, ShouldEqual, 3)

				So(p.State.Equal.Set, ShouldBeFalse)
				So(p.State.Gte.Set, ShouldBeFalse)
				So(p.State.Lte.Set, ShouldBeFalse)

				So(p.Component.Equal.Set, ShouldBeTrue)
				So(p.Component.Equal.Value, ShouldEqual, "component")

				So(p.Component.RegexMatch.Set, ShouldBeFalse)

				So(p.Resource.RegexMatch.Set, ShouldBeTrue)
				So(p.Resource.RegexMatch.Value.String(), ShouldEqual, "service-(?P<id>\\d+)")

				So(p.Resource.Equal.Set, ShouldBeFalse)
			})
		})
	})

	Convey("Given a BSON pattern with an invalid regexp", t, func() {
		mapPattern := mgobson.M{
			"component": mgobson.M{
				"regex_match": "abc-(.*-def",
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mgobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected entity field", t, func() {
		mapPattern := mgobson.M{
			"current_entity": mgobson.M{
				"unexpected_field": "",
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mgobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected info field", t, func() {
		mapPattern := mgobson.M{
			"current_entity": mgobson.M{
				"infos": mgobson.M{
					"info_name": mgobson.M{
						"unexpected_field": "",
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mgobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a mistyped info field", t, func() {
		mapPattern := mgobson.M{
			"current_entity": mgobson.M{
				"infos": mgobson.M{
					"info_name": mgobson.M{
						"value": 3,
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mgobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})
}

func TestPatternUnmarshalMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern", t, func() {
		mapPattern := mongobson.M{
			"state": mongobson.M{
				">": 0,
				"<": 3,
			},
			"component": "component",
			"resource":  mongobson.M{"regex_match": "service-(?P<id>\\d+)"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var p pattern.EventPattern
			So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

			Convey("The values of the pattern should be correct", func() {
				So(p.State.Gt.Set, ShouldBeTrue)
				So(p.State.Gt.Value, ShouldEqual, 0)
				So(p.State.Lt.Set, ShouldBeTrue)
				So(p.State.Lt.Value, ShouldEqual, 3)

				So(p.State.Equal.Set, ShouldBeFalse)
				So(p.State.Gte.Set, ShouldBeFalse)
				So(p.State.Lte.Set, ShouldBeFalse)

				So(p.Component.Equal.Set, ShouldBeTrue)
				So(p.Component.Equal.Value, ShouldEqual, "component")

				So(p.Component.RegexMatch.Set, ShouldBeFalse)

				So(p.Resource.RegexMatch.Set, ShouldBeTrue)
				So(p.Resource.RegexMatch.Value.String(), ShouldEqual, "service-(?P<id>\\d+)")

				So(p.Resource.Equal.Set, ShouldBeFalse)
			})
		})
	})

	Convey("Given a BSON pattern with an invalid regexp", t, func() {
		mapPattern := mongobson.M{
			"component": mongobson.M{
				"regex_match": "abc-(.*-def",
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mongobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected entity field", t, func() {
		mapPattern := mongobson.M{
			"current_entity": mongobson.M{
				"unexpected_field": "",
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mongobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected info field", t, func() {
		mapPattern := mongobson.M{
			"current_entity": mongobson.M{
				"infos": mongobson.M{
					"info_name": mongobson.M{
						"unexpected_field": "",
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mongobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a mistyped info field", t, func() {
		mapPattern := mongobson.M{
			"current_entity": mongobson.M{
				"infos": mongobson.M{
					"info_name": mongobson.M{
						"value": 3,
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p pattern.EventPattern
			So(mongobson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})
}

func TestPatternMatchMgoDriver(t *testing.T) {
	Convey("Given an empty pattern", t, func() {
		mapPattern := mgobson.M{}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts all events", func() {
			event1 := types.Event{
				State:     2,
				Component: "test",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)

			event2 := types.Event{
				State:     3,
				Component: "component",
			}
			So(p.Matches(event2, &m), ShouldBeTrue)
		})
	})

	Convey("Given a pattern with an equality condition", t, func() {
		mapPattern := mgobson.M{
			"component": "component",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that match the condition", func() {
			event1 := types.Event{
				Component: "component",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that do not match the condition", func() {
			event1 := types.Event{
				Component: "not_component",
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern with an integer range", t, func() {
		mapPattern := mgobson.M{
			"state": mgobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that match the condition", func() {
			event1 := types.Event{
				State: 1,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)

			event2 := types.Event{
				State: 2,
			}
			So(p.Matches(event2, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that do not match the condition", func() {
			event1 := types.Event{
				State: 0,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)

			event2 := types.Event{
				State: 3,
			}
			So(p.Matches(event2, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an entity info", t, func() {
		mapPattern := mgobson.M{
			"current_entity": mgobson.M{
				"infos": mgobson.M{
					"info_name": mgobson.M{
						"value": "info_value",
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that have the right info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "info_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that have the wrong info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "wrong_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events where the info is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events where the entity is not defined", func() {
			event1 := types.Event{
				State:     2,
				Component: "test",
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks that the value of an entity info is empty", t, func() {
		mapPattern := mgobson.M{
			"current_entity": mgobson.M{
				"infos": mgobson.M{
					"info_name": mgobson.M{
						"value": "",
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that have the right info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that have the wrong info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "wrong_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events where the info is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks that the value of an entity info is null", t, func() {
		mapPattern := mgobson.M{
			"current_entity": mgobson.M{
				"infos": mgobson.M{
					"info_name": nil,
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events where the info is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events where the info is defined", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "wrong_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks that the entity is not defined", t, func() {
		mapPattern := mgobson.M{
			"current_entity": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that have no entity", func() {
			event1 := types.Event{
				State:     2,
				Component: "test",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events where the entity is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an extra info", t, func() {
		mapPattern := mgobson.M{
			"extra_info": "value",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = "value"
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the wrong value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = "wrong_value"
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an extra with has_every and has_not string array patterns", t, func() {
		mapPattern := mgobson.M{
			"extra_info": mgobson.M{
				"has_every": []string{"test1", "test2"},
				"has_not":   []string{"test3"},
			},
			"extra_info_2": mgobson.M{
				"is_empty": false,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		err = mgobson.Unmarshal(bsonPattern, &p)
		So(err, ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test4"}
			extraInfos["extra_info_2"] = []string{"test5"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the wrong value 1", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test3", "test4"}
			extraInfos["extra_info_2"] = []string{"test5"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern accepts events with the wrong value 2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test4"}
			extraInfos["extra_info_2"] = []string{"test5"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events with empty extra_info_2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test4"}
			extraInfos["extra_info_2"] = []string{}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without extra_info_2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an extra with has_one_of and has_not string array patterns", t, func() {
		mapPattern := mgobson.M{
			"extra_info": mgobson.M{
				"has_one_of": []string{"test1", "test2"},
				"has_not":    []string{"test3"},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		err = mgobson.Unmarshal(bsonPattern, &p)
		So(err, ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the wrong value 1", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test3", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern accepts events with the wrong value 2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern with multiple conditions", t, func() {
		mapPattern := mgobson.M{
			"state": mgobson.M{
				">": 0,
				"<": 3,
			},
			"component": "component",
			"resource":  mgobson.M{"regex_match": "service-(?P<id>\\d+)"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mgobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that match all the conditions", func() {
			event1 := types.Event{
				State:     1,
				Component: "component",
				Resource:  "service-12",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
			So(m.Resource["id"], ShouldEqual, "12")

			event2 := types.Event{
				State:     2,
				Component: "component",
				Resource:  "service-13",
			}
			So(p.Matches(event2, &m), ShouldBeTrue)
			So(m.Resource["id"], ShouldEqual, "13")
		})

		Convey("The pattern rejects events that do not match all the condition", func() {
			event1 := types.Event{
				State:     0,
				Component: "not_component",
				Resource:  "service-0",
			}
			So(p.Matches(event1, &m), ShouldBeFalse)

			event2 := types.Event{
				State:     1,
				Component: "not_component",
				Resource:  "service-0",
			}
			So(p.Matches(event2, &m), ShouldBeFalse)

			event3 := types.Event{
				State:     3,
				Component: "component",
				Resource:  "service-0",
			}
			So(p.Matches(event3, &m), ShouldBeFalse)

			event4 := types.Event{
				State:     2,
				Component: "component",
				Resource:  "service-z",
			}
			So(p.Matches(event4, &m), ShouldBeFalse)
		})
	})
}

func TestPatternMatchMongoDriver(t *testing.T) {
	Convey("Given an empty pattern", t, func() {
		mapPattern := mongobson.M{}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts all events", func() {
			event1 := types.Event{
				State:     2,
				Component: "test",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)

			event2 := types.Event{
				State:     3,
				Component: "component",
			}
			So(p.Matches(event2, &m), ShouldBeTrue)
		})
	})

	Convey("Given a pattern with an equality condition", t, func() {
		mapPattern := mongobson.M{
			"component": "component",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that match the condition", func() {
			event1 := types.Event{
				Component: "component",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that do not match the condition", func() {
			event1 := types.Event{
				Component: "not_component",
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern with an integer range", t, func() {
		mapPattern := mongobson.M{
			"state": mongobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that match the condition", func() {
			event1 := types.Event{
				State: 1,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)

			event2 := types.Event{
				State: 2,
			}
			So(p.Matches(event2, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that do not match the condition", func() {
			event1 := types.Event{
				State: 0,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)

			event2 := types.Event{
				State: 3,
			}
			So(p.Matches(event2, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an entity info", t, func() {
		mapPattern := mongobson.M{
			"current_entity": mongobson.M{
				"infos": mongobson.M{
					"info_name": mongobson.M{
						"value": "info_value",
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that have the right info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "info_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that have the wrong info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "wrong_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events where the info is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events where the entity is not defined", func() {
			event1 := types.Event{
				State:     2,
				Component: "test",
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks that the value of an entity info is empty", t, func() {
		mapPattern := mongobson.M{
			"current_entity": mongobson.M{
				"infos": mongobson.M{
					"info_name": mongobson.M{
						"value": "",
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that have the right info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events that have the wrong info", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "wrong_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events where the info is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks that the value of an entity info is null", t, func() {
		mapPattern := mongobson.M{
			"current_entity": mongobson.M{
				"infos": mongobson.M{
					"info_name": nil,
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events where the info is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events where the info is defined", func() {
			infos := make(map[string]types.Info)
			infos["info_name"] = types.NewInfo("info_name", "", "wrong_value")
			entity := types.NewEntity("r/c", "", "resource", infos, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks that the entity is not defined", t, func() {
		mapPattern := mongobson.M{
			"current_entity": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that have no entity", func() {
			event1 := types.Event{
				State:     2,
				Component: "test",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events where the entity is not defined", func() {
			entity := types.NewEntity("r/c", "", "resource", nil, nil, nil)
			event1 := types.Event{
				State:     2,
				Component: "test",
				Entity:    &entity,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an extra info", t, func() {
		mapPattern := mongobson.M{
			"extra_info": "value",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = "value"
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the wrong value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = "wrong_value"
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an extra with has_every and has_not string array patterns", t, func() {
		mapPattern := mongobson.M{
			"extra_info": mongobson.M{
				"has_every": []string{"test1", "test2"},
				"has_not":   []string{"test3"},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		err = mongobson.Unmarshal(bsonPattern, &p)
		So(err, ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the wrong value 1", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test3", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern accepts events with the wrong value 2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Given a pattern that checks the value of an extra with has_one_of and has_not string array patterns", t, func() {
		mapPattern := mongobson.M{
			"extra_info": mongobson.M{
				"has_one_of": []string{"test1", "test2"},
				"has_not":    []string{"test3"},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		err = mongobson.Unmarshal(bsonPattern, &p)
		So(err, ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the wrong value 1", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2", "test3", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern accepts events with the wrong value 2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test4"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Case with only has_one_of pattern", t, func() {
		mapPattern := mongobson.M{
			"extra_info": mongobson.M{
				"has_one_of": []string{"test1", "test2"},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		err = mongobson.Unmarshal(bsonPattern, &p)
		So(err, ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events with the right value 1", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the right value 2", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test2"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events with the right value 3", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test1", "test2"}
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
		})

		Convey("The pattern rejects events without value", func() {
			extraInfos := make(map[string]interface{})
			event1 := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})

		Convey("The pattern rejects events without infos", func() {
			event1 := types.Event{}
			So(p.Matches(event1, &m), ShouldBeFalse)
		})
	})

	Convey("Case when has_not is set, but event doesn't have an array", t, func() {
		mapPattern := mgobson.M{
			"extra_info": mgobson.M{
				"has_not": []string{"test1", "test2"},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		err = mgobson.Unmarshal(bsonPattern, &p)
		So(err, ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern should accept events with values not included in the has_not set", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test3"}
			event := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event, &m), ShouldBeTrue)
		})

		Convey("The pattern should reject events with values included in the has_not set", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{"test2"}
			event := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event, &m), ShouldBeFalse)
		})

		Convey("The pattern accepts events if the value is empty", func() {
			extraInfos := make(map[string]interface{})
			extraInfos["extra_info"] = []string{}
			event := types.Event{
				ExtraInfos: extraInfos,
			}
			So(p.Matches(event, &m), ShouldBeTrue)
		})

		Convey("The pattern accepts events if the value is null", func() {
			event := types.Event{}
			So(p.Matches(event, &m), ShouldBeTrue)
		})
	})

	Convey("Given a pattern with multiple conditions", t, func() {
		mapPattern := mongobson.M{
			"state": mongobson.M{
				">": 0,
				"<": 3,
			},
			"component": "component",
			"resource":  mongobson.M{"regex_match": "service-(?P<id>\\d+)"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var p pattern.EventPattern
		So(mongobson.Unmarshal(bsonPattern, &p), ShouldBeNil)

		m := pattern.NewEventRegexMatches()

		Convey("The pattern accepts events that match all the conditions", func() {
			event1 := types.Event{
				State:     1,
				Component: "component",
				Resource:  "service-12",
			}
			So(p.Matches(event1, &m), ShouldBeTrue)
			So(m.Resource["id"], ShouldEqual, "12")

			event2 := types.Event{
				State:     2,
				Component: "component",
				Resource:  "service-13",
			}
			So(p.Matches(event2, &m), ShouldBeTrue)
			So(m.Resource["id"], ShouldEqual, "13")
		})

		Convey("The pattern rejects events that do not match all the condition", func() {
			event1 := types.Event{
				State:     0,
				Component: "not_component",
				Resource:  "service-0",
			}
			So(p.Matches(event1, &m), ShouldBeFalse)

			event2 := types.Event{
				State:     1,
				Component: "not_component",
				Resource:  "service-0",
			}
			So(p.Matches(event2, &m), ShouldBeFalse)

			event3 := types.Event{
				State:     3,
				Component: "component",
				Resource:  "service-0",
			}
			So(p.Matches(event3, &m), ShouldBeFalse)

			event4 := types.Event{
				State:     2,
				Component: "component",
				Resource:  "service-z",
			}
			So(p.Matches(event4, &m), ShouldBeFalse)
		})
	})
}

func TestEventPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("service-(?P<id>\\d+)")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.EventPattern
	}{
		{
			TestName: "test for pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"state": mongobson.M{
						">": int64(0),
						"<": int64(3),
					},
					"component": "component",
					"resource":  mongobson.M{"regex_match": "service-(?P<id>\\d+)"},
					"current_entity": mongobson.M{
						"infos": mongobson.M{
							"info_name": mongobson.M{
								"value": "info_value",
							},
						},
					},
					"extra_info": mongobson.M{
						"has_every": mongobson.A{"test1", "test2"},
						"has_not":   mongobson.A{"test3"},
					},
					"extra_info_2": mongobson.M{
						"is_empty": false,
					},
				},
			},
			Pattern: pattern.EventPattern{
				State: pattern.IntegerPattern{
					IntegerConditions: pattern.IntegerConditions{
						Gt: utils.OptionalInt64{
							Set:   true,
							Value: 0,
						},
						Lt: utils.OptionalInt64{
							Set:   true,
							Value: 3,
						},
					},
				},
				Component: pattern.StringPattern{
					StringConditions: pattern.StringConditions{
						Equal: utils.OptionalString{
							Set:   true,
							Value: "component",
						},
					},
				},
				Resource: pattern.StringPattern{
					StringConditions: pattern.StringConditions{
						RegexMatch: utils.OptionalRegexp{
							Set:   true,
							Value: testRegexp,
						},
					},
				},
				Entity: pattern.EntityPattern{
					EntityFields: pattern.EntityFields{
						Infos: map[string]pattern.InfoPattern{
							"info_name": {
								InfoFields: pattern.InfoFields{
									Value: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set:   true,
												Value: "info_value",
											},
										},
									},
								},
							},
						},
					},
				},
				ExtraInfos: map[string]pattern.InterfacePattern{
					"extra_info": {
						StringArrayConditions: pattern.StringArrayConditions{
							HasEvery: utils.OptionalStringArray{
								Set:   true,
								Value: []string{"test1", "test2"},
							},
							HasNot: utils.OptionalStringArray{
								Set:   true,
								Value: []string{"test3"},
							},
						},
					},
					"extra_info_2": {
						StringArrayConditions: pattern.StringArrayConditions{
							IsEmpty: utils.OptionalBool{
								Set:   true,
								Value: false,
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.EventPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w eventPatternWrapper
			w.Pattern = dataset.Pattern

			resultBson, err := mongobson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled mongobson.M
			err = mongobson.Unmarshal(resultBson, &unmarshalled)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			if !reflect.DeepEqual(dataset.ExpectedUnmarshalled, unmarshalled) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled["pattern"], unmarshalled["pattern"])
			}
		})
	}
}

// eventPatternListWrapper is a type that wraps an EventPatternList into a
// struct. It is required to test the unmarshalling of an array into an
// EventPatternList because mgobson.Unmarshal does not work when called with an
// array.
type eventPatternListWrapper struct {
	PatternList pattern.EventPatternList `bson:"list"`
}

func TestValidEventPatternListMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"component": "component1",
				},
				mgobson.M{
					"component": "component2",
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding events", func() {
				event1 := types.Event{
					Component: "component1",
				}
				So(p.Matches(event1), ShouldBeTrue)

				event2 := types.Event{
					Component: "component2",
				}
				So(p.Matches(event2), ShouldBeTrue)

				event3 := types.Event{
					Component: "component3",
				}
				So(p.Matches(event3), ShouldBeFalse)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mgobson.M{}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding events", func() {
				event1 := types.Event{
					Component: "component1",
				}
				So(p.Matches(event1), ShouldBeTrue)

				event2 := types.Event{
					Component: "component2",
				}
				So(p.Matches(event2), ShouldBeTrue)

				event3 := types.Event{
					Component: "component3",
				}
				So(p.Matches(event3), ShouldBeTrue)
			})
		})
	})
}

func TestValidEventPatternListMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				mongobson.M{
					"component": "component1",
				},
				mongobson.M{
					"component": "component2",
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding events", func() {
				event1 := types.Event{
					Component: "component1",
				}
				So(p.Matches(event1), ShouldBeTrue)

				event2 := types.Event{
					Component: "component2",
				}
				So(p.Matches(event2), ShouldBeTrue)

				event3 := types.Event{
					Component: "component3",
				}
				So(p.Matches(event3), ShouldBeFalse)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mongobson.M{}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding events", func() {
				event1 := types.Event{
					Component: "component1",
				}
				So(p.Matches(event1), ShouldBeTrue)

				event2 := types.Event{
					Component: "component2",
				}
				So(p.Matches(event2), ShouldBeTrue)

				event3 := types.Event{
					Component: "component3",
				}
				So(p.Matches(event3), ShouldBeTrue)
			})
		})
	})
}

func TestInvalidEventPatternListMgoDriver(t *testing.T) {
	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"component": "component1",
				},
				mgobson.M{
					"component": mgobson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestInvalidEventPatternListMongoDriver(t *testing.T) {
	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				mongobson.M{
					"component": "component1",
				},
				mongobson.M{
					"component": mongobson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestEventPatternListMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("service-(?P<id>\\d+)")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.EventPatternList
	}{
		{
			TestName: "test for pattern list",
			ExpectedUnmarshalled: mongobson.M{
				"list": mongobson.A{
					mongobson.M{
						"state": mongobson.M{
							">": int64(0),
							"<": int64(3),
						},
						"component": "component",
						"resource":  mongobson.M{"regex_match": "service-(?P<id>\\d+)"},
						"current_entity": mongobson.M{
							"infos": mongobson.M{
								"info_name": mongobson.M{
									"value": "info_value",
								},
							},
						},
						"extra_info": mongobson.M{
							"has_every": mongobson.A{"test1", "test2"},
							"has_not":   mongobson.A{"test3"},
						},
					},
					mongobson.M{
						"state": mongobson.M{
							">": int64(0),
							"<": int64(3),
						},
						"component": "component",
						"resource":  mongobson.M{"regex_match": "service-(?P<id>\\d+)"},
						"current_entity": mongobson.M{
							"infos": mongobson.M{
								"info_name": mongobson.M{
									"value": "info_value",
								},
							},
						},
						"extra_info": mongobson.M{
							"has_every": mongobson.A{"test1", "test2"},
							"has_not":   mongobson.A{"test3"},
						},
						"extra_info_2": "test_info",
						"extra_info_3": mongobson.M{
							"is_empty": true,
						},
					},
				},
			},
			Pattern: pattern.EventPatternList{
				Set:   true,
				Valid: true,
				Patterns: []pattern.EventPattern{
					{
						State: pattern.IntegerPattern{
							IntegerConditions: pattern.IntegerConditions{
								Gt: utils.OptionalInt64{
									Set:   true,
									Value: 0,
								},
								Lt: utils.OptionalInt64{
									Set:   true,
									Value: 3,
								},
							},
						},
						Component: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								Equal: utils.OptionalString{
									Set:   true,
									Value: "component",
								},
							},
						},
						Resource: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								RegexMatch: utils.OptionalRegexp{
									Set:   true,
									Value: testRegexp,
								},
							},
						},
						Entity: pattern.EntityPattern{
							EntityFields: pattern.EntityFields{
								Infos: map[string]pattern.InfoPattern{
									"info_name": {
										InfoFields: pattern.InfoFields{
											Value: pattern.StringPattern{
												StringConditions: pattern.StringConditions{
													Equal: utils.OptionalString{
														Set:   true,
														Value: "info_value",
													},
												},
											},
										},
									},
								},
							},
						},
						ExtraInfos: map[string]pattern.InterfacePattern{
							"extra_info": {
								StringArrayConditions: pattern.StringArrayConditions{
									HasEvery: utils.OptionalStringArray{
										Set:   true,
										Value: []string{"test1", "test2"},
									},
									HasNot: utils.OptionalStringArray{
										Set:   true,
										Value: []string{"test3"},
									},
								},
							},
						},
					},
					{
						State: pattern.IntegerPattern{
							IntegerConditions: pattern.IntegerConditions{
								Gt: utils.OptionalInt64{
									Set:   true,
									Value: 0,
								},
								Lt: utils.OptionalInt64{
									Set:   true,
									Value: 3,
								},
							},
						},
						Component: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								Equal: utils.OptionalString{
									Set:   true,
									Value: "component",
								},
							},
						},
						Resource: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								RegexMatch: utils.OptionalRegexp{
									Set:   true,
									Value: testRegexp,
								},
							},
						},
						Entity: pattern.EntityPattern{
							EntityFields: pattern.EntityFields{
								Infos: map[string]pattern.InfoPattern{
									"info_name": {
										InfoFields: pattern.InfoFields{
											Value: pattern.StringPattern{
												StringConditions: pattern.StringConditions{
													Equal: utils.OptionalString{
														Set:   true,
														Value: "info_value",
													},
												},
											},
										},
									},
								},
							},
						},
						ExtraInfos: map[string]pattern.InterfacePattern{
							"extra_info": {
								StringArrayConditions: pattern.StringArrayConditions{
									HasEvery: utils.OptionalStringArray{
										Set:   true,
										Value: []string{"test1", "test2"},
									},
									HasNot: utils.OptionalStringArray{
										Set:   true,
										Value: []string{"test3"},
									},
								},
							},
							"extra_info_2": {
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
										Value: "test_info",
									},
								},
							},
							"extra_info_3": {
								StringArrayConditions: pattern.StringArrayConditions{
									IsEmpty: utils.OptionalBool{
										Set:   true,
										Value: true,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for nil",
			ExpectedUnmarshalled: mongobson.M{
				"list": nil,
			},
			Pattern: pattern.EventPatternList{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w eventPatternListWrapper
			w.PatternList = dataset.Pattern

			resultBson, err := mongobson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled mongobson.M
			err = mongobson.Unmarshal(resultBson, &unmarshalled)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			if !reflect.DeepEqual(dataset.ExpectedUnmarshalled, unmarshalled) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled["pattern"], unmarshalled["pattern"])
			}
		})
	}
}
