package pattern_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPbehaviorInfo_Match(t *testing.T) {
	dataSets := getPbehaviorInfoMatchDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := data.pattern.Match(data.pbehaviorInfo)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func TestPbehaviorInfo_ToMongoQuery(t *testing.T) {
	dataSets := getPbehaviorInfoMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := data.pattern.ToMongoQuery("alarm")
			if !errors.Is(err, data.mongoQueryErr) {
				t.Errorf("expected error %v but got %v", data.mongoQueryErr, err)
			}
			if diff := pretty.Compare(query, data.mongoQueryResult); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func getPbehaviorInfoMatchDataSets() map[string]PbehaviorInfoDataSet {
	return map[string]PbehaviorInfoDataSet{
		"given empty pattern should match": {
			pattern: pattern.PbehaviorInfo{},
			pbehaviorInfo: types.PbehaviorInfo{
				ID: "test id",
			},
			matchResult: true,
		},
		"given string field condition should match": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{
				ID: "test id",
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{
				ID: "test another id",
			},
			matchResult: false,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "created",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{},
			matchErr:      pattern.ErrUnsupportedField,
		},
		"given active canonical field condition and emtpty pbehavior infos should match": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "active"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{},
			matchResult:   true,
		},
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
			pattern: pattern.PbehaviorInfo{
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
	pattern          pattern.PbehaviorInfo
	pbehaviorInfo    types.PbehaviorInfo
	matchErr         error
	matchResult      bool
	mongoQueryErr    error
	mongoQueryResult bson.M
}
