package db_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestEventPatternToMongoQuery(t *testing.T) {
	regexCond, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "snmp")
	if err != nil {
		t.Fatalf("invalid regexp condition: %v", err)
	}

	f := func(pattern pattern.Event, expectedQuery primitive.M, expectedError error) {
		t.Helper()

		mongoQuery, err := db.EventPatternToMongoQuery(pattern, "event")
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error %v but got %v", expectedError, err)
		}
		diff := pretty.Compare(expectedQuery, mongoQuery)
		if diff != "" {
			t.Errorf("unexpected result %s", diff)
		}
	}

	// single condition
	f(pattern.Event{
		{
			{
				Field:     "event_type",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "check"),
			},
		},
	}, bson.M{"$or": []bson.M{
		{"$and": []bson.M{
			{"event.event_type": bson.M{"$eq": "check"}},
		}},
	}}, nil)

	// multiple conditions
	f(pattern.Event{
		{
			{
				Field:     "event_type",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "check"),
			},
			{
				Field:     "connector",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "snmp"),
			},
		},
	}, bson.M{"$or": []bson.M{
		{"$and": []bson.M{
			{"event.event_type": bson.M{"$eq": "check"}},
			{"event.connector": bson.M{"$eq": "snmp"}},
		}},
	}}, nil)

	// multiple groups
	f(pattern.Event{
		{
			{
				Field:     "event_type",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "check"),
			},
		},
		{
			{
				Field:     "connector",
				Condition: regexCond,
			},
		},
	}, bson.M{"$or": []bson.M{
		{"$and": []bson.M{
			{"event.event_type": bson.M{"$eq": "check"}},
		}},
		{"$and": []bson.M{
			{"event.connector": bson.M{"$regex": "snmp"}},
		}},
	}}, nil)

	// invalid condition type
	f(pattern.Event{
		{
			{
				Field:     "event_type",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			},
		},
	}, nil, pattern.ErrUnsupportedConditionType)

	// wrong condition value
	f(pattern.Event{
		{
			{
				Field:     "event_type",
				Condition: pattern.NewIntCondition(pattern.ConditionEqual, 1),
			},
		},
	}, nil, pattern.ErrWrongConditionValue)

	// extra infos condition
	f(pattern.Event{
		{
			{
				Field:     "extra.x_event_type",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "check"),
			},
			{
				Field:     "extra.x_state",
				FieldType: pattern.FieldTypeInt,
				Condition: pattern.NewIntCondition(pattern.ConditionLT, 3),
			},
			{
				Field:     "extra.x_has_perfdata",
				FieldType: pattern.FieldTypeBool,
				Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
			},
			{
				Field:     "extra.x_hosts",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"host1", "host2"}),
			},
			{
				Field:     "extra.x_flag_1",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			},
		},
	}, bson.M{"$or": []bson.M{
		{"$and": []bson.M{
			{"event.extra_infos.x_event_type": bson.M{"$eq": "check"}},
			{"event.extra_infos.x_state": bson.M{"$lt": 3}},
			{"event.extra_infos.x_has_perfdata": bson.M{"$eq": true}},
			{"event.extra_infos.x_hosts": bson.M{"$elemMatch": bson.M{"$in": []string{"host1", "host2"}}}},
			{"event.extra_infos.x_flag_1": bson.M{"$exists": true, "$ne": nil}},
		}},
	}}, nil)
}
