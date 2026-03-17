package db_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestEntityPatternToMongoQuery(t *testing.T) {
	dataSets := getEntityToMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := db.EntityPatternToMongoQuery(data.pattern, "entity")
			if !errors.Is(err, data.mongoQueryErr) {
				t.Errorf("expected error %v but got %v", data.mongoQueryErr, err)
			}
			if diff := pretty.Compare(query, data.mongoQueryResult); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func TestEntityPatternToSql(t *testing.T) {
	dataSets := getEntityToSqlDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			sql, err := db.EntityPatternToSql(data.pattern, "entity")
			if !errors.Is(err, data.sqlErr) {
				t.Errorf("expected error %v but got %v", data.sqlErr, err)
			}
			if sql != data.sqlResult {
				t.Errorf("expected\n%s\nbut got\n%s", data.sqlResult, sql)
			}
		})
	}
}

func getEntityToMongoQueryDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given one condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.name": bson.M{"$eq": "test name"}},
				}},
			}},
		},
		"given multiple conditions": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.name": bson.M{"$eq": "test name"}},
					{"entity.category": bson.M{"$eq": "test category"}},
				}},
			}},
		},
		"given multiple groups": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
				{
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.name": bson.M{"$eq": "test name"}},
				}},
				{"$and": []bson.M{
					{"entity.category": bson.M{"$eq": "test category"}},
				}},
			}},
		},
		"given invalid condition type": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			mongoQueryErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 2),
					},
				},
			},
			mongoQueryErr: pattern.ErrWrongConditionValue,
		},
		"given infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.infos.info_name.value": bson.M{"$eq": 3}},
				}},
			}},
		},
		"given component infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.component_infos.info_name.value": bson.M{"$eq": 3}},
				}},
			}},
		},
		"given infos timestamp condition with absolute time": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.infos.timestamp_info.value": bson.M{
						"$gt": datetime.NewCpsTime(1609459200),
						"$lt": datetime.NewCpsTime(1640995200),
					}},
				}},
			}},
		},
		"given component infos timestamp condition with absolute time": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.component_infos.timestamp_info.value": bson.M{
						"$gt": datetime.NewCpsTime(1609459200),
						"$lt": datetime.NewCpsTime(1640995200),
					}},
				}},
			}},
		},
		"given invalid timestamp condition type for infos": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "not a timestamp"),
					},
				},
			},
			mongoQueryErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid timestamp condition type for component infos": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "not a timestamp"),
					},
				},
			},
			mongoQueryErr: pattern.ErrUnsupportedConditionType,
		},
	}
}

func getEntityToSqlDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given one condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			sqlResult: "(entity.name = 'test name')",
		},
		"given multiple conditions": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			sqlResult: "(entity.name = 'test name' AND entity.category = 'test category')",
		},
		"given multiple groups": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
				{
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			sqlResult: "(entity.name = 'test name') OR (entity.category = 'test category')",
		},
		"given invalid condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionIsNotOneOf, "test name"),
					},
				},
			},
			sqlErr: pattern.ErrWrongConditionValue,
		},
		"given infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'info_name') = 'number' AND (CASE WHEN jsonb_typeof(infos->'info_name') = 'number' THEN (infos->'info_name')::numeric END) = 3))`,
		},
		"given component infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'info_name') = 'number' AND (CASE WHEN jsonb_typeof(component_infos->'info_name') = 'number' THEN (component_infos->'info_name')::numeric END) = 3))`,
		},
		"given infos timestamp condition with absolute time": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'timestamp_info') = 'number' AND (infos->'timestamp_info')::bigint > 1609459200 AND (infos->'timestamp_info')::bigint < 1640995200))`,
		},
		"given infos timestamp condition with relative time within condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'timestamp_info') = 'number' AND (infos->'timestamp_info')::bigint > EXTRACT(EPOCH FROM (NOW() - INTERVAL '3600 seconds'))::bigint))`,
		},
		"given infos timestamp condition with relative time older than condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(0, datetime.DurationUnitSecond), datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'timestamp_info') = 'number' AND (infos->'timestamp_info')::bigint < EXTRACT(EPOCH FROM (NOW() - INTERVAL '3600 seconds'))::bigint))`,
		},
		"given infos timestamp condition with relative time interval condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond), datetime.NewDurationWithUnit(1800, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'timestamp_info') = 'number' AND (infos->'timestamp_info')::bigint > EXTRACT(EPOCH FROM (NOW() - INTERVAL '3600 seconds'))::bigint AND (infos->'timestamp_info')::bigint < EXTRACT(EPOCH FROM (NOW() - INTERVAL '1800 seconds'))::bigint))`,
		},
		"given invalid timestamp condition type for infos sql": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "not a timestamp"),
					},
				},
			},
			sqlErr: pattern.ErrUnsupportedConditionType,
		},
		"given component_infos timestamp condition with absolute time": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'timestamp_info') = 'number' AND (component_infos->'timestamp_info')::bigint > 1609459200 AND (component_infos->'timestamp_info')::bigint < 1640995200))`,
		},
		"given component_infos timestamp condition with relative time within condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'timestamp_info') = 'number' AND (component_infos->'timestamp_info')::bigint > EXTRACT(EPOCH FROM (NOW() - INTERVAL '3600 seconds'))::bigint))`,
		},
		"given component_infos timestamp condition with relative time older than condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(0, datetime.DurationUnitSecond), datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'timestamp_info') = 'number' AND (component_infos->'timestamp_info')::bigint < EXTRACT(EPOCH FROM (NOW() - INTERVAL '3600 seconds'))::bigint))`,
		},
		"given component_infos timestamp condition with relative time interval condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond), datetime.NewDurationWithUnit(1800, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'timestamp_info') = 'number' AND (component_infos->'timestamp_info')::bigint > EXTRACT(EPOCH FROM (NOW() - INTERVAL '3600 seconds'))::bigint AND (component_infos->'timestamp_info')::bigint < EXTRACT(EPOCH FROM (NOW() - INTERVAL '1800 seconds'))::bigint))`,
		},
		"given invalid timestamp condition type for component_infos sql": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "not a timestamp"),
					},
				},
			},
			sqlErr: pattern.ErrUnsupportedConditionType,
		},
		"given one condition with string not equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test name"),
					},
				},
			},
			sqlResult: "((entity.name IS NULL OR entity.name != 'test name'))",
		},
		"given one condition with string is one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test name", "other name"}),
					},
				},
			},
			sqlResult: "(entity.name = ANY (ARRAY ['test name','other name']))",
		},
		"given one condition with string is not one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{"test name", "other name"}),
					},
				},
			},
			sqlResult: "((entity.name IS NULL OR NOT (entity.name = ANY (ARRAY ['test name','other name']))))",
		},
		"given one condition with string regexp": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test.*$")
							return condition
						}(),
					},
				},
			},
			sqlResult: "(entity.name ~ '^test.*$')",
		},
		"given one condition with string exist": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			sqlResult: "((entity.name IS NOT NULL AND entity.name != ''))",
		},
		"given one condition with int not equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "impact_level",
						Condition: pattern.NewIntCondition(pattern.ConditionNotEqual, 3),
					},
				},
			},
			sqlResult: "((entity.impact_level IS NULL OR entity.impact_level != 3))",
		},
		"given one condition with int greater than": {
			pattern: pattern.Entity{
				{
					{
						Field:     "impact_level",
						Condition: pattern.NewIntCondition(pattern.ConditionGT, 3),
					},
				},
			},
			sqlResult: "(entity.impact_level > 3)",
		},
		"given one condition with int less than": {
			pattern: pattern.Entity{
				{
					{
						Field:     "impact_level",
						Condition: pattern.NewIntCondition(pattern.ConditionLT, 3),
					},
				},
			},
			sqlResult: "(entity.impact_level < 3)",
		},
		"given infos string condition with not equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test name"),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'text_info') = 'string' AND infos->>'text_info' != 'test name'))`,
		},
		"given infos string condition with is one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test name", "other name"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'text_info') = 'string' AND infos->>'text_info' = ANY (ARRAY ['test name','other name'])))`,
		},
		"given infos string condition with is not one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{"test name", "other name"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'text_info') = 'string' AND NOT (infos->>'text_info' = ANY (ARRAY ['test name','other name']))))`,
		},
		"given infos string condition with regexp": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test.*$")
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'text_info') = 'string' AND infos->>'text_info' ~ '^test.*$'))`,
		},
		"given infos string condition with exist": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			sqlResult: `((infos ? 'text_info' AND jsonb_typeof(infos->'text_info') = 'string' AND infos->>'text_info' != ''))`,
		},
		"given infos int condition with not equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.number_info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionNotEqual, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'number_info') = 'number' AND (CASE WHEN jsonb_typeof(infos->'number_info') = 'number' THEN (infos->'number_info')::numeric END) != 3))`,
		},
		"given infos int condition with greater than": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.number_info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionGT, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'number_info') = 'number' AND (CASE WHEN jsonb_typeof(infos->'number_info') = 'number' THEN (infos->'number_info')::numeric END) > 3))`,
		},
		"given infos int condition with less than": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.number_info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionLT, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'number_info') = 'number' AND (CASE WHEN jsonb_typeof(infos->'number_info') = 'number' THEN (infos->'number_info')::numeric END) < 3))`,
		},
		"given infos bool condition with equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.bool_info",
						FieldType: pattern.FieldTypeBool,
						Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'bool_info') = 'boolean' AND (CASE WHEN jsonb_typeof(infos->'bool_info') = 'boolean' THEN (infos->'bool_info')::boolean END) = true))`,
		},
		"given infos string array condition with is empty": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'array_info') = 'array' AND jsonb_array_length((CASE WHEN jsonb_typeof(infos->'array_info') = 'array' THEN infos->'array_info' END)) = 0))`,
		},
		"given infos string array condition with has every": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"tag-1", "tag-2"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'array_info') = 'array' AND infos->'array_info' ?& ARRAY ['tag-1','tag-2']))`,
		},
		"given infos string array condition with has one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"tag-1", "tag-2"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'array_info') = 'array' AND infos->'array_info' ?| ARRAY ['tag-1','tag-2']))`,
		},
		"given infos string array condition with has not": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"tag-1", "tag-2"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'array_info') = 'array' AND NOT ((CASE WHEN jsonb_typeof(infos->'array_info') = 'array' THEN infos->'array_info' END) ?| ARRAY ['tag-1','tag-2'])))`,
		},
		"given infos ref condition with exist true": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.ref_info",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			sqlResult: "(infos ? 'ref_info')",
		},
		"given infos ref condition with exist false": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.ref_info",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			sqlResult: "(NOT (infos ? 'ref_info'))",
		},
		"given component_infos string condition with not equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test name"),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'text_info') = 'string' AND component_infos->>'text_info' != 'test name'))`,
		},
		"given component_infos string condition with is one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test name", "other name"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'text_info') = 'string' AND component_infos->>'text_info' = ANY (ARRAY ['test name','other name'])))`,
		},
		"given component_infos string condition with is not one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{"test name", "other name"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'text_info') = 'string' AND NOT (component_infos->>'text_info' = ANY (ARRAY ['test name','other name']))))`,
		},
		"given component_infos string condition with regexp": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test.*$")
							return condition
						}(),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'text_info') = 'string' AND component_infos->>'text_info' ~ '^test.*$'))`,
		},
		"given component_infos string condition with exist": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.text_info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			sqlResult: `((component_infos ? 'text_info' AND jsonb_typeof(component_infos->'text_info') = 'string' AND component_infos->>'text_info' != ''))`,
		},
		"given component_infos int condition with not equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.number_info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionNotEqual, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'number_info') = 'number' AND (CASE WHEN jsonb_typeof(component_infos->'number_info') = 'number' THEN (component_infos->'number_info')::numeric END) != 3))`,
		},
		"given component_infos int condition with greater than": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.number_info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionGT, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'number_info') = 'number' AND (CASE WHEN jsonb_typeof(component_infos->'number_info') = 'number' THEN (component_infos->'number_info')::numeric END) > 3))`,
		},
		"given component_infos int condition with less than": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.number_info",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionLT, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'number_info') = 'number' AND (CASE WHEN jsonb_typeof(component_infos->'number_info') = 'number' THEN (component_infos->'number_info')::numeric END) < 3))`,
		},
		"given component_infos bool condition with equal": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.bool_info",
						FieldType: pattern.FieldTypeBool,
						Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'bool_info') = 'boolean' AND (CASE WHEN jsonb_typeof(component_infos->'bool_info') = 'boolean' THEN (component_infos->'bool_info')::boolean END) = true))`,
		},
		"given component_infos string array condition with is empty": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'array_info') = 'array' AND jsonb_array_length((CASE WHEN jsonb_typeof(component_infos->'array_info') = 'array' THEN component_infos->'array_info' END)) = 0))`,
		},
		"given component_infos string array condition with has every": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"tag-1", "tag-2"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'array_info') = 'array' AND component_infos->'array_info' ?& ARRAY ['tag-1','tag-2']))`,
		},
		"given component_infos string array condition with has one of": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"tag-1", "tag-2"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'array_info') = 'array' AND component_infos->'array_info' ?| ARRAY ['tag-1','tag-2']))`,
		},
		"given component_infos string array condition with has not": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.array_info",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"tag-1", "tag-2"}),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'array_info') = 'array' AND NOT ((CASE WHEN jsonb_typeof(component_infos->'array_info') = 'array' THEN component_infos->'array_info' END) ?| ARRAY ['tag-1','tag-2'])))`,
		},
		"given component_infos ref condition with exist true": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.ref_info",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			sqlResult: "(component_infos ? 'ref_info')",
		},
		"given component_infos ref condition with exist false": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.ref_info",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			sqlResult: "(NOT (component_infos ? 'ref_info'))",
		},
	}
}

type entityDataSet struct {
	pattern          pattern.Entity
	mongoQueryErr    error
	mongoQueryResult bson.M
	sqlErr           error
	sqlResult        string
}
