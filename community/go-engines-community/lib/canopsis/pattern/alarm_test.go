package pattern_test

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestAlarm_Match(t *testing.T) {
	dataSets := getAlarmMatchDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := data.pattern.Match(data.alarm)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func TestAlarm_ToMongoQuery(t *testing.T) {
	dataSets := getAlarmMongoQueryDataSets()

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

func getAlarmMatchDataSets() map[string]alarmDataSet {
	return map[string]alarmDataSet{
		"given string field condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.display_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					DisplayName: "test name",
				},
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.display_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					DisplayName: "test another name",
				},
			},
			matchResult: false,
		},
		"given string field condition and not string field should return error": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.state.val",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{State: &types.AlarmStep{}},
			},
			matchErr: pattern.ErrWrongConditionValue,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.initial_output",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm:    types.Alarm{},
			matchErr: pattern.ErrUnsupportedField,
		},
		"given string info condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_name": "test name",
						},
					},
				},
			},
			matchResult: true,
		},
		"given string info condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_name": "test another name",
						},
					},
				},
			},
			matchResult: false,
		},
		"given string info condition and not string info should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_name": 2,
						},
					},
				},
			},
			matchResult: false,
		},
		"given string info condition and unknown info should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
			},
			alarm:       types.Alarm{},
			matchResult: false,
		},
		"given exist info condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: true,
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_name": 0,
						},
					},
				},
			},
			matchResult: true,
		},
		"given exist info condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: true,
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_another_name": 0,
						},
					},
				},
			},
			matchResult: false,
		},
		"given not exist info condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: false,
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_another_name": 0,
						},
					},
				},
			},
			matchResult: true,
		},
		"given not exist info condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.infos.info_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionExist,
							Value: false,
						},
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Infos: map[string]map[string]interface{}{
						"rule1": {
							"info_name": 0,
						},
					},
				},
			},
			matchResult: false,
		},
	}
}

func getAlarmMongoQueryDataSets() map[string]alarmDataSet {
	return map[string]alarmDataSet{
		"given one condition": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.display_name",
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
						{"alarm.v.display_name": bson.M{"$eq": "test name"}},
					}},
				}}},
			},
		},
		"given multiple conditions": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.display_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
					{
						Field: "v.output",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test output",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"alarm.v.display_name": bson.M{"$eq": "test name"}},
						{"alarm.v.output": bson.M{"$eq": "test output"}},
					}},
				}}},
			},
		},
		"given multiple groups": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.display_name",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test name",
						},
					},
				},
				{
					{
						Field: "v.output",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: "test output",
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"alarm.v.display_name": bson.M{"$eq": "test name"}},
					}},
					{"$and": []bson.M{
						{"alarm.v.output": bson.M{"$eq": "test output"}},
					}},
				}}},
			},
		},
		"given invalid condition": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.state.val",
						Condition: pattern.Condition{
							Type:  pattern.ConditionIsEmpty,
							Value: "test name",
						},
					},
				},
			},
			mongoQueryErr: pattern.ErrWrongConditionValue,
		},
		"given duration condition": {
			pattern: pattern.Alarm{
				{
					{
						Field: "v.duration",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: 3,
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$addFields": bson.M{
					"alarm.v.duration": bson.M{"$subtract": bson.A{
						bson.M{"$cond": bson.M{
							"if":   "$alarm.v.resolved",
							"then": "$alarm.v.resolved",
							"else": time.Now().Unix(),
						}},
						"$alarm.v.creation_date",
					}},
				}},
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"alarm.v.duration": bson.M{"$eq": 3}},
					}},
				}}},
			},
		},
		"given infos condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: 3,
						},
					},
				},
			},
			mongoQueryResult: []bson.M{
				{"$addFields": bson.M{
					"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
				}},
				{"$match": bson.M{"$or": []bson.M{
					{"$and": []bson.M{
						{"$and": []bson.M{
							{"alarm.v.infos_array.v.info_name": bson.M{"$type": bson.A{"long", "int", "decimal"}}},
							{"alarm.v.infos_array.v.info_name": bson.M{"$eq": 3}},
						}},
					}},
				}}},
			},
		},
	}
}

type alarmDataSet struct {
	pattern          pattern.Alarm
	alarm            types.Alarm
	matchErr         error
	matchResult      bool
	mongoQueryErr    error
	mongoQueryResult []bson.M
}
