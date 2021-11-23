package eventfilter_test

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

var eventCheck0 = types.Event{
	EventType:     types.EventTypeCheck,
	SourceType:    types.SourceTypeComponent,
	Connector:     "connector",
	ConnectorName: "connector-name",
	Component:     "component",
	State:         0,
	Debug:         true,
}
var eventCheck1 = types.Event{
	EventType:     types.EventTypeCheck,
	SourceType:    types.SourceTypeComponent,
	Component:     "component",
	Connector:     "connector",
	ConnectorName: "connector-name",
	State:         1,
	Debug:         true,
}
var eventCheck2 = types.Event{
	EventType:     types.EventTypeCheck,
	SourceType:    types.SourceTypeComponent,
	Component:     "component",
	Connector:     "connector",
	ConnectorName: "connector-name",
	State:         2,
	Debug:         true,
}
var eventCheck3 = types.Event{
	EventType:     types.EventTypeCheck,
	SourceType:    types.SourceTypeComponent,
	Connector:     "connector",
	ConnectorName: "connector-name",
	Component:     "component",
	State:         3,
	Debug:         true,
	Output:        "Warning: CPU Load is critical (90%)",
}

var untypedRule = bson.M{
	"_id":      "untyped",
	"patterns": []bson.M{},
	"priority": 10,
	"enabled":  true,
}
var invalidTypeRule = bson.M{
	"_id":      "invalid_type",
	"type":     "invalid_type",
	"patterns": []bson.M{},
	"priority": 10,
	"enabled":  true,
}
var mistypedTypeRule = bson.M{
	"_id":      "mistyped_type",
	"type":     12,
	"patterns": []bson.M{},
	"priority": 10,
	"enabled":  true,
}
var mistypedPatternRule = bson.M{
	"_id":  "mistyped_pattern",
	"type": "drop",
	"patterns": []bson.M{
		{"state": bson.M{">": "test"}},
	},
	"priority": 10,
	"enabled":  true,
}
var invalidRegexPatternRule = bson.M{
	"_id":  "invalid_regex_pattern",
	"type": "drop",
	"patterns": []bson.M{
		{"component": bson.M{
			"regex_match": "(.*",
		}},
	},
	"priority": 10,
	"enabled":  true,
}
var unexpectedFieldPatternRule = bson.M{
	"_id":  "unexpected_field_pattern",
	"type": "drop",
	"patterns": []bson.M{
		{"state": bson.M{
			">":                3,
			"unexpected_field": true,
		}},
	},
	"priority": 10,
	"enabled":  true,
}
var unexpectedFieldRule = bson.M{
	"_id":              "unexpected_field",
	"type":             "drop",
	"patterns":         []bson.M{},
	"priority":         10,
	"enabled":          true,
	"unexpected_field": true,
}
var unexpectedEnrichmentFieldRule = bson.M{
	"_id":        "unexpected_enrichment_field",
	"type":       "drop",
	"patterns":   []bson.M{},
	"priority":   10,
	"enabled":    true,
	"on_success": "drop",
}
var noActionsRule = bson.M{
	"_id":      "no_actions",
	"type":     "enrichment",
	"patterns": []bson.M{},
	"priority": 10,
	"enabled":  true,
}
var emptyActionsRule = bson.M{
	"_id":      "empty_actions",
	"type":     "enrichment",
	"patterns": []bson.M{},
	"actions":  []bson.M{},
	"priority": 10,
	"enabled":  true,
}
var invalidOutcomeRule = bson.M{
	"_id":  "invalid_outcome",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"type":  "set_field",
			"name":  "Component",
			"value": "component_name",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "invalid_outcome",
}
var mistypedOutcomeRule = bson.M{
	"_id":  "mistyped_outcome",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"type":  "set_field",
			"name":  "Component",
			"value": "component_name",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": 3,
}
var invalidActionRule = bson.M{
	"_id":  "invalid_action",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"type":  "invalid_action_type",
			"name":  "Output",
			"value": "modified output",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "drop",
}
var missingActionTypeRule = bson.M{
	"_id":  "missing_action_type",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"name":  "Output",
			"value": "modified output",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "drop",
}
var unexpectedActionFieldRule = bson.M{
	"_id":  "unexpected_action_field",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"type":             "set_field",
			"name":             "Output",
			"value":            "modified output",
			"unexpected_field": "test",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "drop",
}
var missingActionFieldRule = bson.M{
	"_id":  "missing_action_field",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"type": "set_field",
			"name": "Output",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "drop",
}
var dropRule = bson.M{
	"_id":  "valid_drop",
	"type": "drop",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"priority": 15,
	"enabled":  true,
}
var breakRule = bson.M{
	"_id":  "valid_break",
	"type": "break",
	"patterns": []bson.M{
		{"state": bson.M{">=": 3}},
	},
	"priority": 5,
}
var enrichmentRule = bson.M{
	"_id":  "valid_enrichment",
	"type": "enrichment",
	"actions": []bson.M{
		bson.M{
			"type":  "set_field",
			"name":  "Output",
			"value": "modified output",
		},
	},
	"priority":   10,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "drop",
}
var failingEnrichmentRule = bson.M{
	"_id":  "failing_enrichment",
	"type": "enrichment",
	"patterns": []bson.M{
		{"state": bson.M{">": 1}},
	},
	"actions": []bson.M{
		bson.M{
			"type":  "set_field",
			"name":  "Output",
			"value": 3,
		},
	},
	"priority":   20,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "drop",
}
var translationRule = bson.M{
	"_id":  "translation",
	"type": "enrichment",
	"patterns": []bson.M{
		{"output": bson.M{"regex_match": "Warning: CPU Load is critical \\((?P<value>.*)\\)"}},
	},
	"actions": []bson.M{
		{
			"type":  "set_field_from_template",
			"name":  "Output",
			"value": "Attention, la charge CPU est critique ({{.RegexMatch.Output.value}})",
		},
	},
	"priority":   100,
	"enabled":    true,
	"on_success": "pass",
	"on_failure": "pass",
}

func TestUnmarshal(t *testing.T) {
	Convey("Given a rule without a type", t, func() {
		bsonRule, err := bson.Marshal(untypedRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a rule with an invalid type", t, func() {
		bsonRule, err := bson.Marshal(invalidTypeRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a rule with a mistyped type", t, func() {
		bsonRule, err := bson.Marshal(mistypedTypeRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a drop rule with a mistyped pattern", t, func() {
		bsonRule, err := bson.Marshal(mistypedPatternRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a drop rule with a pattern with an invalid regex", t, func() {
		bsonRule, err := bson.Marshal(invalidRegexPatternRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a drop rule with a pattern with an unexpected field", t, func() {
		bsonRule, err := bson.Marshal(unexpectedFieldPatternRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a drop rule with an unexpected field", t, func() {
		bsonRule, err := bson.Marshal(unexpectedFieldRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given a drop rule with an on_success field", t, func() {
		bsonRule, err := bson.Marshal(unexpectedEnrichmentFieldRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule without actions", t, func() {
		bsonRule, err := bson.Marshal(noActionsRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with an empty actions field", t, func() {
		bsonRule, err := bson.Marshal(emptyActionsRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with an invalid outcome", t, func() {
		bsonRule, err := bson.Marshal(invalidOutcomeRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with a mistyped outcome", t, func() {
		bsonRule, err := bson.Marshal(mistypedOutcomeRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with an invalid action type", t, func() {
		bsonRule, err := bson.Marshal(invalidActionRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with an action without type", t, func() {
		bsonRule, err := bson.Marshal(missingActionTypeRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with an unexpected action field", t, func() {
		bsonRule, err := bson.Marshal(unexpectedActionFieldRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule with a missing action field", t, func() {
		bsonRule, err := bson.Marshal(missingActionFieldRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as invalid", func() {
				So(rule.Valid, ShouldBeFalse)
			})
		})
	})

}

func TestDropRule(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get()

	timeZone := mockTimezoneConfigProvider.Get()

	Convey("Given a valid drop rule", t, func() {
		bsonRule, err := bson.Marshal(dropRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as valid", func() {
				So(rule.Valid, ShouldBeTrue)
			})

			Convey("The values of type, priority should be set correctly", func() {
				So(rule.Type, ShouldEqual, eventfilter.RuleTypeDrop)
				So(rule.Priority, ShouldEqual, 15)
				So(rule.IsEnabled(), ShouldBeTrue)
			})

			Convey("The Pattern should match corresponding events", func() {
				So(rule.Patterns.Matches(eventCheck0), ShouldBeFalse)
				So(rule.Patterns.Matches(eventCheck1), ShouldBeFalse)
				So(rule.Patterns.Matches(eventCheck2), ShouldBeTrue)
				So(rule.Patterns.Matches(eventCheck3), ShouldBeTrue)
			})

			Convey("The rule's outcome should always be Drop", func() {
				report := eventfilter.Report{}
				_, outcome := rule.Apply(ctx, eventCheck0, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Drop)
				So(report.EntityUpdated, ShouldBeFalse)

				_, outcome = rule.Apply(ctx, eventCheck1, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Drop)
				So(report.EntityUpdated, ShouldBeFalse)

				_, outcome = rule.Apply(ctx, eventCheck2, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Drop)
				So(report.EntityUpdated, ShouldBeFalse)

				_, outcome = rule.Apply(ctx, eventCheck3, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Drop)
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})
}

func TestBreakRule(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get()
	timeZone := mockTimezoneConfigProvider.Get()

	Convey("Given a valid break rule", t, func() {
		bsonRule, err := bson.Marshal(breakRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as valid", func() {
				So(rule.Valid, ShouldBeTrue)
			})

			Convey("The values of type, priority should be set correctly", func() {
				So(rule.Type, ShouldEqual, eventfilter.RuleTypeBreak)
				So(rule.Priority, ShouldEqual, 5)
				So(rule.IsEnabled(), ShouldBeTrue)
			})

			Convey("The Pattern should match corresponding events", func() {
				So(rule.Patterns.Matches(eventCheck0), ShouldBeFalse)
				So(rule.Patterns.Matches(eventCheck1), ShouldBeFalse)
				So(rule.Patterns.Matches(eventCheck2), ShouldBeFalse)
				So(rule.Patterns.Matches(eventCheck3), ShouldBeTrue)
			})

			Convey("The rule's outcome should always be Break", func() {
				report := eventfilter.Report{}
				_, outcome := rule.Apply(ctx, eventCheck0, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Break)
				So(report.EntityUpdated, ShouldBeFalse)

				_, outcome = rule.Apply(ctx, eventCheck1, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Break)
				So(report.EntityUpdated, ShouldBeFalse)

				_, outcome = rule.Apply(ctx, eventCheck2, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Break)
				So(report.EntityUpdated, ShouldBeFalse)

				_, outcome = rule.Apply(ctx, eventCheck3, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(err, ShouldBeNil)
				So(outcome, ShouldEqual, eventfilter.Break)
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})
}

func TestEnrichmentRule(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get()
	timeZone := mockTimezoneConfigProvider.Get()

	Convey("Given an enrichment rule that changes the output of the events", t, func() {
		bsonRule, err := bson.Marshal(enrichmentRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as valid", func() {
				So(rule.Valid, ShouldBeTrue)
			})

			Convey("The values of type, priority should be set correctly", func() {
				So(rule.Type, ShouldEqual, eventfilter.RuleTypeEnrichment)
				So(rule.Priority, ShouldEqual, 10)
				So(rule.IsEnabled(), ShouldBeTrue)
				So(rule.OnSuccess, ShouldEqual, eventfilter.Pass)
				So(rule.OnFailure, ShouldEqual, eventfilter.Drop)
			})

			Convey("The rule's outcome should always be Break", func() {
				report := eventfilter.Report{}
				event, outcome := rule.Apply(ctx, eventCheck3, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(event.Output, ShouldEqual, "modified output")
				So(outcome, ShouldEqual, eventfilter.Pass)
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule that fails", t, func() {
		bsonRule, err := bson.Marshal(failingEnrichmentRule)
		So(err, ShouldBeNil)
		ctrl := gomock.NewController(t)
		mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
		mockTimezoneConfigProvider.EXPECT().Get()
		timeZone := mockTimezoneConfigProvider.Get()

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as valid", func() {
				So(rule.Valid, ShouldBeTrue)
			})

			Convey("The values of type, priority should be set correctly", func() {
				So(rule.Type, ShouldEqual, eventfilter.RuleTypeEnrichment)
				So(rule.Priority, ShouldEqual, 20)
				So(rule.IsEnabled(), ShouldBeTrue)
				So(rule.OnSuccess, ShouldEqual, eventfilter.Pass)
				So(rule.OnFailure, ShouldEqual, eventfilter.Drop)
			})

			Convey("The rule's outcome should be OnFailure", func() {
				report := eventfilter.Report{}
				_, outcome := rule.Apply(ctx, eventCheck3, pattern.NewEventRegexMatches(), &report, &timeZone, log.NewTestLogger())
				So(outcome, ShouldEqual, eventfilter.Drop)
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})

	Convey("Given an enrichment rule that translates the output of the events", t, func() {
		bsonRule, err := bson.Marshal(translationRule)
		So(err, ShouldBeNil)

		Convey("The rule should be decoded without errors", func() {
			var rule eventfilter.RuleUnpacker
			So(bson.Unmarshal(bsonRule, &rule), ShouldBeNil)

			Convey("The rule should be marked as valid", func() {
				So(rule.Valid, ShouldBeTrue)
			})

			Convey("The values of type, priority should be set correctly", func() {
				So(rule.Type, ShouldEqual, eventfilter.RuleTypeEnrichment)
				So(rule.Priority, ShouldEqual, 100)
				So(rule.IsEnabled(), ShouldBeTrue)
				So(rule.OnSuccess, ShouldEqual, eventfilter.Pass)
				So(rule.OnFailure, ShouldEqual, eventfilter.Pass)
			})

			Convey("The rule translates the output", func() {
				matches, match := rule.Patterns.GetRegexMatches(eventCheck3)
				So(match, ShouldBeTrue)

				report := eventfilter.Report{}
				event, outcome := rule.Apply(ctx, eventCheck3, matches, &report, &timeZone, log.NewTestLogger())
				So(event.Output, ShouldEqual, "Attention, la charge CPU est critique (90%)")
				So(outcome, ShouldEqual, eventfilter.Pass)
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})
}
