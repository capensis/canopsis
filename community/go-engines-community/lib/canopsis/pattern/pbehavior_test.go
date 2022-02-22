package pattern_test

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestPbehavior_Match(t *testing.T) {
	dataSets := getPbehaviorMatchDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := data.pattern.Match(data.pbehavior)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func TestPbehavior_ToMongoQuery(t *testing.T) {
	dataSets := getPbehaviorMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := data.pattern.ToMongoQuery("pbehavior")
			if !errors.Is(err, data.mongoQueryErr) {
				t.Errorf("expected error %v but got %v", data.mongoQueryErr, err)
			}
			if diff := pretty.Compare(query, data.mongoQueryResult); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func getPbehaviorMatchDataSets() map[string]PbehaviorDataSet {
	return map[string]PbehaviorDataSet{
		"given string field condition should match": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test id",
						},
					},
				},
			},
			pbehavior: pbehavior.PBehavior{
				ID: "test id",
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test id",
						},
					},
				},
			},
			pbehavior: pbehavior.PBehavior{
				ID: "test another id",
			},
			matchResult: false,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "created",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			pbehavior: pbehavior.PBehavior{},
			matchErr:  pattern.ErrUnsupportedField,
		},
	}
}

func getPbehaviorMongoQueryDataSets() map[string]PbehaviorDataSet {
	return map[string]PbehaviorDataSet{
		"given one condition": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test id",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"pbehavior._id": bson.M{"$eq": "test id"}},
					}},
				}}},
			},
		},
		"given multiple conditions": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test id",
						},
					},
					{
						Field: "type",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test type",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"pbehavior._id": bson.M{"$eq": "test id"}},
						{"pbehavior.type": bson.M{"$eq": "test type"}},
					}},
				}}},
			},
		},
		"given multiple groups": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test id",
						},
					},
				},
				{
					{
						Field: "type",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test type",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"pbehavior._id": bson.M{"$eq": "test id"}},
					}},
					{"$and": []bson.M{
						{"pbehavior.type": bson.M{"$eq": "test type"}},
					}},
				}}},
			},
		},
		"given invalid condition": {
			pattern: pattern.Pbehavior{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionIsEmpty,
							Value: "test id",
						},
					},
				},
			},
			mongoQueryErr: pattern.ErrWrongConditionValue,
		},
	}
}

type PbehaviorDataSet struct {
	pattern          pattern.Pbehavior
	pbehavior        pbehavior.PBehavior
	matchErr         error
	matchResult      bool
	mongoQueryErr    error
	mongoQueryResult []bson.M
}
