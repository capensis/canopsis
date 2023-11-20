package db_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPbehaviorInfo_ToMongoQuery(t *testing.T) {
	dataSets := getPbehaviorInfoMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := db.PBehaviorInfoPatternToMongoQuery(data.pattern, "alarm")
			if !errors.Is(err, data.mongoQueryErr) {
				t.Errorf("expected error %v but got %v", data.mongoQueryErr, err)
			}
			if diff := pretty.Compare(query, data.mongoQueryResult); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func getPbehaviorInfoMongoQueryDataSets() map[string]PbehaviorInfoDataSet {
	return map[string]PbehaviorInfoDataSet{
		"given one condition": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.id": bson.M{"$eq": "test id"}},
				}},
			}},
		},
		"given multiple conditions": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
					{
						Field:     "pbehavior_info.type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test type"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.id": bson.M{"$eq": "test id"}},
					{"alarm.pbehavior_info.type": bson.M{"$eq": "test type"}},
				}},
			}},
		},
		"given multiple groups": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
				},
				{
					{
						Field:     "pbehavior_info.type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test type"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.id": bson.M{"$eq": "test id"}},
				}},
				{"$and": []bson.M{
					{"alarm.pbehavior_info.type": bson.M{"$eq": "test type"}},
				}},
			}},
		},
		"given invalid condition type": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionIsEmpty, "test id"),
					},
				},
			},
			mongoQueryErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 2),
					},
				},
			},
			mongoQueryErr: pattern.ErrWrongConditionValue,
		},
		"given equal to active canonical type condition": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "active"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$in": bson.A{nil, "active"}}},
				}},
			}},
		},
		"given not equal to active canonical type condition": {
			pattern: pattern.PBehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "active"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$nin": bson.A{nil, "active"}}},
				}},
			}},
		},
	}
}

type PbehaviorInfoDataSet struct {
	pattern          pattern.PBehaviorInfo
	mongoQueryErr    error
	mongoQueryResult bson.M
}
