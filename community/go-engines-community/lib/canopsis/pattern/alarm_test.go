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
			fields := data.pattern.GetMongoFields("alarm")
			if diff := pretty.Compare(fields, data.mongoQueryFields); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func BenchmarkAlarm_Match_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "v.display_name",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name 2"),
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			DisplayName: "test name",
		},
	}

	benchmarkAlarmMatch(b, cond, alarm)
}

func BenchmarkAlarm_Match_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "v.display_name",
		Condition: regexpCondition,
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			DisplayName: "test name",
		},
	}

	benchmarkAlarmMatch(b, cond, alarm)
}

func BenchmarkAlarm_Match_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "v.infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test 2"),
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			Infos: map[string]map[string]interface{}{
				"rule1": {
					"test": "test",
				},
			},
		},
	}

	benchmarkAlarmMatch(b, cond, alarm)
}

func BenchmarkAlarm_Match_Infos_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "v.infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: regexpCondition,
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			Infos: map[string]map[string]interface{}{
				"rule1": {
					"test": "test",
				},
			},
		},
	}

	benchmarkAlarmMatch(b, cond, alarm)
}

func BenchmarkAlarm_UnmarshalBsonAndMatch_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "v.display_name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test name 2",
		},
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			DisplayName: "test name",
		},
	}

	benchmarkAlarmUnmarshalBsonAndMatch(b, cond, []types.Alarm{alarm})
}

func BenchmarkAlarm_UnmarshalBsonAndMatch_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "v.display_name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			DisplayName: "test name",
		},
	}

	benchmarkAlarmUnmarshalBsonAndMatch(b, cond, []types.Alarm{alarm})
}

func BenchmarkAlarm_UnmarshalBsonAndMatch_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "v.infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test 2",
		},
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			Infos: map[string]map[string]interface{}{
				"rule1": {
					"test": "test",
				},
			},
		},
	}

	benchmarkAlarmUnmarshalBsonAndMatch(b, cond, []types.Alarm{alarm})
}

func BenchmarkAlarm_UnmarshalBsonAndMatch_Infos_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "v.infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			Infos: map[string]map[string]interface{}{
				"rule1": {
					"test": "test",
				},
			},
		},
	}

	benchmarkAlarmUnmarshalBsonAndMatch(b, cond, []types.Alarm{alarm})
}

func BenchmarkAlarm_ManyAlarms_UnmarshalBsonAndMatch_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "v.display_name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	const size = 1000
	alarms := make([]types.Alarm, size)
	for i := 0; i < size; i++ {
		alarms[i] = types.Alarm{
			Value: types.AlarmValue{
				DisplayName: "test name",
			},
		}
	}

	benchmarkAlarmUnmarshalBsonAndMatch(b, cond, alarms)
}

func benchmarkAlarmMatch(b *testing.B, fieldCond pattern.FieldCondition, alarm types.Alarm) {
	const size = 100
	p := make(pattern.Alarm, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := p.Match(alarm)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func benchmarkAlarmUnmarshalBsonAndMatch(b *testing.B, fieldCond pattern.FieldCondition, alarms []types.Alarm) {
	const size = 100
	p := make(pattern.Alarm, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	type wrapper struct {
		Pattern pattern.Alarm `bson:"pattern"`
	}
	bytes, err := bson.Marshal(wrapper{Pattern: p})
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var w wrapper
		err := bson.Unmarshal(bytes, &w)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}

		for _, alarm := range alarms {
			_, err = w.Pattern.Match(alarm)
			if err != nil {
				b.Fatalf("unexpected error %v", err)
			}
		}
	}
}

func getAlarmMatchDataSets() map[string]alarmDataSet {
	return map[string]alarmDataSet{
		"given empty pattern should match": {
			pattern: pattern.Alarm{},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					DisplayName: "test name",
				},
			},
			matchResult: true,
		},
		"given string field condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.display_name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Field:     "v.display_name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Field:     "v.state.val",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Field:     "v.initial_output",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
						Field:     "v.infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
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
						Field:     "v.infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
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
						Field:     "v.infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
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
						Field:     "v.infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
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
	durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{
		Value: 3,
		Unit:  "s",
	})
	if err != nil {
		panic(err)
	}

	return map[string]alarmDataSet{
		"given one condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.display_name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.display_name": bson.M{"$eq": "test name"}},
				}},
			}},
		},
		"given multiple conditions": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.display_name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
					{
						Field:     "v.output",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test output"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.display_name": bson.M{"$eq": "test name"}},
					{"alarm.v.output": bson.M{"$eq": "test output"}},
				}},
			}},
		},
		"given multiple groups": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.display_name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
				{
					{
						Field:     "v.output",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test output"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.display_name": bson.M{"$eq": "test name"}},
				}},
				{"$and": []bson.M{
					{"alarm.v.output": bson.M{"$eq": "test output"}},
				}},
			}},
		},
		"given invalid condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.state.val",
						Condition: pattern.NewStringCondition(pattern.ConditionIsEmpty, "test name"),
					},
				},
			},
			mongoQueryErr: pattern.ErrWrongConditionValue,
		},
		"given duration condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.duration",
						Condition: durationCond,
					},
				},
			},
			mongoQueryFields: bson.M{
				"alarm.v.duration": bson.M{"$subtract": bson.A{
					bson.M{"$cond": bson.M{
						"if":   "$alarm.v.resolved",
						"then": "$alarm.v.resolved",
						"else": time.Now().Unix(),
					}},
					"$alarm.v.creation_date",
				}},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.duration": bson.M{"$gt": 3}},
				}},
			}},
		},
		"given infos condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			mongoQueryFields: bson.M{
				"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$and": []bson.M{
						{"alarm.v.infos_array.v.info_name": bson.M{"$type": bson.A{"long", "int", "decimal"}}},
						{"alarm.v.infos_array.v.info_name": bson.M{"$eq": 3}},
					}},
				}},
			}},
		},
	}
}

type alarmDataSet struct {
	pattern          pattern.Alarm
	alarm            types.Alarm
	matchErr         error
	matchResult      bool
	mongoQueryErr    error
	mongoQueryResult bson.M
	mongoQueryFields bson.M
}
