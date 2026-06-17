package entityservice

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSearchPatternToMongoQuery(t *testing.T) {
	dataSets := map[string]struct {
		pattern     searchPattern
		expectedRes bson.M
		expectedErr error
	}{
		"given empty pattern": {
			pattern:     searchPattern{},
			expectedRes: nil,
			expectedErr: nil,
		},
		"given name field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name"}},
				}},
			}},
		},
		"given type field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "service"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"type": bson.M{"$eq": "service"}},
				}},
			}},
		},
		"given multiple conditions in one group": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
					{
						Field:     "type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "service"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name"}},
					{"type": bson.M{"$eq": "service"}},
				}},
			}},
		},
		"given multiple groups": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name-1"),
					},
				},
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name-2"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name-1"}},
				}},
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name-2"}},
				}},
			}},
		},
		"given unsupported field": {
			pattern: searchPattern{
				{
					{
						Field:     "unsupported_field",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrUnsupportedField,
		},
		"given invalid condition type for string field": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value for string field": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 123),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrWrongConditionValue,
		},
	}

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := data.pattern.ToMongoQuery()
			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
			if diff := pretty.Compare(query, data.expectedRes); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}
