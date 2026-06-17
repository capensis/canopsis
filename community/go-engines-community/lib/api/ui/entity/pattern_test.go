package entity

import (
	"errors"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSearchPatternToMongoQuery(t *testing.T) {
	now := time.Now()
	from := now.Add(-time.Hour).Unix()
	to := now.Add(time.Hour).Unix()

	dataSets := map[string]struct {
		pattern                searchPattern
		expectedRes            bson.M
		expectedLookups        map[string]struct{}
		expectedComputedFields map[string]struct{}
		expectedErr            error
	}{
		"given empty pattern": {
			pattern:                searchPattern{},
			expectedRes:            nil,
			expectedLookups:        nil,
			expectedComputedFields: nil,
			expectedErr:            nil,
		},
		"given _id field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "_id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-id"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"_id": bson.M{"$eq": "test-id"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
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
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given type field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "resource"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"type": bson.M{"$eq": "resource"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given category field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-category"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"category": bson.M{"$eq": "test-category"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given component field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "component",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-component"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"component": bson.M{"$eq": "test-component"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given connector field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "connector",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-connector"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"connector": bson.M{"$eq": "test-connector"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given import_source field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "import_source",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-import-source"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"import_source": bson.M{"$eq": "test-import-source"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given pbehavior_info.name field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "pbehavior_info.name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-pbh-name"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"pbehavior_info.name": bson.M{"$eq": "test-pbh-name"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given pbehavior_info.reason field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "pbehavior_info.reason",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-reason"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"pbehavior_info.reason": bson.M{"$eq": "test-reason"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given pbehavior_info.type field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "pbehavior_info.type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-type"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"pbehavior_info.type": bson.M{"$eq": "test-type"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given pbehavior_info.canonical_type field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"pbehavior_info.canonical_type": bson.M{"$eq": "pause"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given impact_level field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "impact_level",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"impact_level": bson.M{"$eq": int64(3)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given state field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "state",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"state": bson.M{"$eq": int64(3)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"alarm": {}},
			expectedComputedFields: map[string]struct{}{"state": {}},
		},
		"given status field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "status",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 1),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"status": bson.M{"$eq": int64(1)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"alarm": {}},
			expectedComputedFields: map[string]struct{}{"status": {}},
		},
		"given impact_state field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "impact_state",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 9),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"impact_state": bson.M{"$eq": int64(9)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"alarm": {}},
			expectedComputedFields: map[string]struct{}{"impact_state": {}},
		},
		"given ok_events field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "ok_events",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 10),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"ok_events": bson.M{"$eq": int64(10)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"event_stats": {}},
			expectedComputedFields: map[string]struct{}{"ok_events": {}},
		},
		"given ko_events field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "ko_events",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 5),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"ko_events": bson.M{"$eq": int64(5)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"event_stats": {}},
			expectedComputedFields: map[string]struct{}{"ko_events": {}},
		},
		"given idle_since field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "idle_since",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"idle_since": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given imported field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "imported",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"imported": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given last_pbehavior_date field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "last_pbehavior_date",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"last_pbehavior_date": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given last_event_date field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "last_event_date",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"last_event_date": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given alarm_last_update_date field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "alarm_last_update_date",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm_last_update_date": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"alarm": {}},
			expectedComputedFields: map[string]struct{}{"alarm_last_update_date": {}},
		},
		"given infos field with string condition": {
			pattern: searchPattern{
				{
					{
						Field:     "infos.test-info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-value"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"infos.test-info.value": bson.M{"$eq": "test-value"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given infos field with int condition": {
			pattern: searchPattern{
				{
					{
						Field:     "infos.test-info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 123),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"infos.test-info.value": bson.M{"$eq": int64(123)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given infos field with bool condition": {
			pattern: searchPattern{
				{
					{
						Field:     "infos.test-info",
						FieldType: pattern.FieldTypeBool,
						Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"infos.test-info.value": bson.M{"$eq": true}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given component_infos field with string condition": {
			pattern: searchPattern{
				{
					{
						Field:     "component_infos.test-info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-value"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"component_infos.test-info.value": bson.M{"$eq": "test-value"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given multiple conditions in one group": {
			pattern: searchPattern{
				{
					{
						Field:     "_id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-id"),
					},
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
					{
						Field:     "impact_level",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"_id": bson.M{"$eq": "test-id"}},
					{"name": bson.M{"$eq": "test-name"}},
					{"impact_level": bson.M{"$eq": int64(3)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given multiple groups": {
			pattern: searchPattern{
				{
					{
						Field:     "_id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-id-1"),
					},
				},
				{
					{
						Field:     "_id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-id-2"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"_id": bson.M{"$eq": "test-id-1"}},
				}},
				{"$and": []bson.M{
					{"_id": bson.M{"$eq": "test-id-2"}},
				}},
			}},
			expectedLookups:        map[string]struct{}{},
			expectedComputedFields: map[string]struct{}{},
		},
		"given multiple computed fields from different lookups": {
			pattern: searchPattern{
				{
					{
						Field:     "state",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
					{
						Field:     "ok_events",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 10),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"state": bson.M{"$eq": int64(3)}},
					{"ok_events": bson.M{"$eq": int64(10)}},
				}},
			}},
			expectedLookups:        map[string]struct{}{"alarm": {}, "event_stats": {}},
			expectedComputedFields: map[string]struct{}{"state": {}, "ok_events": {}},
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
			expectedRes:            nil,
			expectedLookups:        nil,
			expectedComputedFields: nil,
			expectedErr:            pattern.ErrUnsupportedField,
		},
		"given invalid condition type for string field": {
			pattern: searchPattern{
				{
					{
						Field:     "_id",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			expectedRes:            nil,
			expectedLookups:        nil,
			expectedComputedFields: nil,
			expectedErr:            pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value for string field": {
			pattern: searchPattern{
				{
					{
						Field:     "_id",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 123),
					},
				},
			},
			expectedRes:            nil,
			expectedLookups:        nil,
			expectedComputedFields: nil,
			expectedErr:            pattern.ErrWrongConditionValue,
		},
		"given invalid condition type for int field": {
			pattern: searchPattern{
				{
					{
						Field:     "impact_level",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			expectedRes:            nil,
			expectedLookups:        nil,
			expectedComputedFields: nil,
			expectedErr:            pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition type for time field": {
			pattern: searchPattern{
				{
					{
						Field:     "idle_since",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "2024-01-01"),
					},
				},
			},
			expectedRes:            nil,
			expectedLookups:        nil,
			expectedComputedFields: nil,
			expectedErr:            pattern.ErrUnsupportedConditionType,
		},
	}

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, lookups, computedFields, err := data.pattern.ToMongoQuery()
			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
			if diff := pretty.Compare(query, data.expectedRes); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
			if diff := pretty.Compare(lookups, data.expectedLookups); diff != "" {
				t.Errorf("unexpected lookups %s", diff)
			}
			if diff := pretty.Compare(computedFields, data.expectedComputedFields); diff != "" {
				t.Errorf("unexpected computed fields %s", diff)
			}
		})
	}
}
