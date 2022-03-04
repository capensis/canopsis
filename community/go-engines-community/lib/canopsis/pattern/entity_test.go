package pattern_test

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestEntity_Match(t *testing.T) {
	dataSets := getEntityMatchDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := data.pattern.Match(data.entity)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func TestEntity_ToMongoQuery(t *testing.T) {
	dataSets := getEntityMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := data.pattern.ToMongoQuery("entity")
			if !errors.Is(err, data.mongoQueryErr) {
				t.Errorf("expected error %v but got %v", data.mongoQueryErr, err)
			}
			if diff := pretty.Compare(query, data.mongoQueryResult); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func getEntityMatchDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given empty pattern should match": {
			pattern: pattern.Entity{},
			entity: types.Entity{
				Name: "test name",
			},
			matchResult: true,
		},
		"given string field condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity: types.Entity{
				Name: "test name",
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity: types.Entity{
				Name: "test another name",
			},
			matchResult: false,
		},
		"given string field condition and not string field should return error": {
			pattern: pattern.Entity{
				{
					{
						Field: "impact_level",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity:   types.Entity{},
			matchErr: pattern.ErrWrongConditionValue,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.Entity{
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
			entity:   types.Entity{},
			matchErr: pattern.ErrUnsupportedField,
		},
		"given string info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given string info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test another name",
					},
				},
			},
			matchResult: false,
		},
		"given string info condition and not string info should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       2,
					},
				},
			},
			matchResult: false,
		},
		"given string info condition and unknown info should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			entity:      types.Entity{},
			matchResult: false,
		},
		"given exist info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field: "infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: true,
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given exist info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field: "infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: true,
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_another_name": {
						Name:        "info_another_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: false,
		},
		"given not exist info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field: "infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: false,
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_another_name": {
						Name:        "info_another_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given not exist info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field: "infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: false,
						},
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: false,
		},
	}
}

func getEntityMongoQueryDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given one condition": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"entity.name": bson.M{"$eq": "test name"}},
					}},
				}}},
			},
		},
		"given multiple conditions": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
					{
						Field: "category",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test category",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"entity.name": bson.M{"$eq": "test name"}},
						{"entity.category": bson.M{"$eq": "test category"}},
					}},
				}}},
			},
		},
		"given multiple groups": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
				{
					{
						Field: "category",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test category",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"entity.name": bson.M{"$eq": "test name"}},
					}},
					{"$and": []bson.M{
						{"entity.category": bson.M{"$eq": "test category"}},
					}},
				}}},
			},
		},
		"given invalid condition": {
			pattern: pattern.Entity{
				{
					{
						Field: "name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionIsEmpty,
							Value: "test name",
						},
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
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: 3,
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"$and": []bson.M{
							{"entity.infos.info_name.val": bson.M{"$type": bson.A{"long", "int", "decimal"}}},
							{"entity.infos.info_name.val": bson.M{"$eq": 3}},
						}},
					}},
				}}},
			},
		},
	}
}

type entityDataSet struct {
	pattern          pattern.Entity
	entity           types.Entity
	matchErr         error
	matchResult      bool
	mongoQueryErr    error
	mongoQueryResult []bson.M
}
