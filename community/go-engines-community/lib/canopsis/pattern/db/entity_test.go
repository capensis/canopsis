package db_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
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
				t.Errorf("expected\n%s\nbut got\n%s", sql, data.sqlResult)
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
	}
}

type entityDataSet struct {
	pattern          pattern.Entity
	mongoQueryErr    error
	mongoQueryResult bson.M
	sqlErr           error
	sqlResult        string
}
