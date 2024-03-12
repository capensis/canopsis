package db_test

import (
	"errors"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAlarmPatternToMongoQuery(t *testing.T) {
	dataSets := getAlarmMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := db.AlarmPatternToMongoQuery(data.pattern, "alarm")
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

func getAlarmMongoQueryDataSets() map[string]alarmDataSet {
	durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, datetime.DurationWithUnit{
		Value: 3,
		Unit:  "s",
	})
	if err != nil {
		panic(err)
	}
	from := time.Now().Add(-time.Hour).Unix()
	to := time.Now().Add(time.Hour).Unix()

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
		"given one ne condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.display_name",
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test name"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.display_name": bson.M{"$ne": "test name"}},
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
		"given invalid condition type": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.state.val",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			mongoQueryErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.state.val",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
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
				"alarm.v.duration": bson.M{"$ifNull": bson.A{
					"$alarm.v.duration",
					bson.M{"$subtract": bson.A{
						bson.M{"$cond": bson.M{
							"if":   "$alarm.v.resolved",
							"then": "$alarm.v.resolved",
							"else": time.Now().Unix(),
						}},
						"$alarm.v.creation_date",
					}},
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
					{"alarm.v.infos_array.v.info_name": bson.M{"$eq": 3}},
				}},
			}},
		},
		"given infos string neq condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test"),
					},
				},
			},
			mongoQueryFields: bson.M{
				"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.infos_array.v.info_name": bson.M{
						"$nin": bson.A{"test", nil},
					}},
				}},
			}},
		},
		"given infos string nin condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{"test"}),
					},
				},
			},
			mongoQueryFields: bson.M{
				"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.infos_array.v.info_name": bson.M{
						"$nin": []string{"test"},
						"$ne":  nil,
					}},
				}},
			}},
		},
		"given infos int neq condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionNotEqual, 3),
					},
				},
			},
			mongoQueryFields: bson.M{
				"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.infos_array.v.info_name": bson.M{
						"$nin": bson.A{3, nil},
					}},
				}},
			}},
		},
		"given infos string array nin condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeStringArray,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"test"}),
					},
				},
			},
			mongoQueryFields: bson.M{
				"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.infos_array.v.info_name": bson.M{
						"$nin": []string{"test"},
						"$ne":  nil,
					}},
				}},
			}},
		},
		"given exist ref condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ack",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.ack": bson.M{"$exists": true, "$ne": nil}},
				}},
			}},
		},
		"given not exist ref condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ack",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$or": []bson.M{
						{"alarm.v.ack": bson.M{"$exists": false}},
						{"alarm.v.ack": bson.M{"$eq": nil}},
					}},
				}},
			}},
		},
		"given exist v.output condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.output",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.output": bson.M{
						"$exists": true,
						"$nin":    bson.A{nil, ""},
					}},
				}},
			}},
		},
		"given not exist v.output condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.output",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$or": []bson.M{
						{"alarm.v.output": bson.M{"$exists": false}},
						{"alarm.v.output": bson.M{"$eq": nil}},
						{"alarm.v.output": bson.M{"$eq": ""}},
					}},
				}},
			}},
		},
		"given not exist v.activation_date condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.activation_date",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$or": []bson.M{
						{"alarm.v.activation_date": bson.M{"$exists": false}},
						{"alarm.v.activation_date": bson.M{"$eq": nil}},
					}},
				}},
			}},
		},
		"given time v.activation_date condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.activation_date",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.activation_date": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
		},
		"given ticket conditions": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test ticket"),
					},
					{
						Field:     "v.ticket.m",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test message"),
					},
					{
						Field:     "v.ticket.ticket_data.data_1",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test_1"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.ticket.ticket": bson.M{"$eq": "test ticket"}},
					{"alarm.v.ticket.m": bson.M{"$eq": "test message"}},
					{"alarm.v.ticket.ticket_data.data_1": bson.M{"$eq": "test_1"}},
				}},
			}},
		},
		"given exist change_state condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.change_state",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm.v.change_state": bson.M{"$exists": true, "$ne": nil}},
				}},
			}},
		},
		"given not exist change_state condition": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.change_state",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$or": []bson.M{
						{"alarm.v.change_state": bson.M{"$exists": false}},
						{"alarm.v.change_state": bson.M{"$eq": nil}},
					}},
				}},
			}},
		},
	}
}

type alarmDataSet struct {
	pattern          pattern.Alarm
	mongoQueryErr    error
	mongoQueryResult bson.M
	mongoQueryFields bson.M
}
