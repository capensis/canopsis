package db_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPbehaviorInfoPatternToMongoQuery(t *testing.T) {
	dataSets := getPbehaviorInfoMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := db.PbehaviorInfoPatternToMongoQuery(data.pattern, "alarm")
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
			pattern: pattern.PbehaviorInfo{
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
			pattern: pattern.PbehaviorInfo{
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
			pattern: pattern.PbehaviorInfo{
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
			pattern: pattern.PbehaviorInfo{
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
			pattern: pattern.PbehaviorInfo{
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
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, types.PbhCanonicalTypeActive),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$in": bson.A{nil, types.PbhCanonicalTypeActive}}},
				}},
			}},
		},
		"given not equal to active canonical type condition": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, types.PbhCanonicalTypeActive),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$nin": bson.A{nil, types.PbhCanonicalTypeActive}}},
				}},
			}},
		},
		"given is_one_of without canonical type condition": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{pbehavior.TypePause, pbehavior.TypeMaintenance}),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$in": bson.A{pbehavior.TypePause, pbehavior.TypeMaintenance}}},
				}},
			}},
		},
		"given is_not_one_of without canonical type condition": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{pbehavior.TypePause, pbehavior.TypeMaintenance}),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$nin": bson.A{pbehavior.TypePause, pbehavior.TypeMaintenance}}},
				}},
			}},
		},
		"given is_one_of with active canonical type condition": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{types.PbhCanonicalTypeActive, pbehavior.TypeMaintenance}),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$in": bson.A{types.PbhCanonicalTypeActive, pbehavior.TypeMaintenance, nil}}},
				}},
			}},
		},
		"given is_not_one_of with active canonical type condition": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{types.PbhCanonicalTypeActive, pbehavior.TypeMaintenance}),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.pbehavior_info.canonical_type": bson.M{"$nin": bson.A{types.PbhCanonicalTypeActive, pbehavior.TypeMaintenance, nil}}},
				}},
			}},
		},
	}
}

type PbehaviorInfoDataSet struct {
	pattern          pattern.PbehaviorInfo
	mongoQueryErr    error
	mongoQueryResult bson.M
}
