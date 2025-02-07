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
	f := func(p pattern.Alarm, expectedRes, expectedFields bson.M, expectedErr error) {
		t.Helper()
		query, err := db.AlarmPatternToMongoQuery(p, "alarm")
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected error %v but got %v", expectedErr, err)
		}

		if diff := pretty.Compare(query, expectedRes); diff != "" {
			t.Fatalf("unexpected result %s", diff)
		}

		fields := p.GetMongoFields("alarm")
		if diff := pretty.Compare(fields, expectedFields); diff != "" {
			t.Fatalf("unexpected result %s", diff)
		}
	}

	durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, datetime.DurationWithUnit{
		Value: 3,
		Unit:  "s",
	})
	if err != nil {
		panic(err)
	}
	from := time.Now().Add(-time.Hour).Unix()
	to := time.Now().Add(time.Hour).Unix()

	f(
		pattern.Alarm{
			{
				{
					Field:     "v.display_name",
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.display_name": bson.M{"$eq": "test name"}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.display_name",
					Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test name"),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.display_name": bson.M{"$ne": "test name"}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
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
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.display_name": bson.M{"$eq": "test name"}},
				{"alarm.v.output": bson.M{"$eq": "test output"}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
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
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.display_name": bson.M{"$eq": "test name"}},
			}},
			{"$and": []bson.M{
				{"alarm.v.output": bson.M{"$eq": "test output"}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.state.val",
					Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
				},
			},
		},
		nil,
		nil,
		pattern.ErrUnsupportedConditionType,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.state.val",
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
				},
			},
		},
		nil,
		nil,
		pattern.ErrWrongConditionValue,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.duration",
					Condition: durationCond,
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.duration": bson.M{"$gt": 3}},
			}},
		}},
		bson.M{
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
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeInt,
					Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{
					"$elemMatch": bson.M{"v.info_name": bson.M{"$eq": 3}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test"),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{
					"$elemMatch": bson.M{"v.info_name": bson.M{"$nin": bson.A{"test", nil}}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{"test"}),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{
					"$elemMatch": bson.M{"v.info_name": bson.M{
						"$nin": []string{"test"},
						"$ne":  nil,
					}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{
					"$elemMatch": bson.M{"v.info_name": bson.M{
						"$exists": true,
						"$nin":    bson.A{nil, ""},
					}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.infos_array": bson.M{"$in": bson.A{nil, bson.A{}}}},
					{"alarm.v.infos_array": bson.M{"$not": bson.M{"$elemMatch": bson.M{"v.info_name": bson.M{
						"$exists": true,
						"$nin":    bson.A{nil, ""},
					}}}}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeInt,
					Condition: pattern.NewIntCondition(pattern.ConditionNotEqual, 3),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{
					"$elemMatch": bson.M{"v.info_name": bson.M{"$nin": bson.A{3, nil}}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeStringArray,
					Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"test"}),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{
					"$elemMatch": bson.M{"v.info_name": bson.M{
						"$nin": []string{"test"},
						"$ne":  nil,
					}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeStringArray,
					Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.infos_array": bson.M{"$in": bson.A{nil, bson.A{}}}},
					{"alarm.v.infos_array": bson.M{"$not": bson.M{"$elemMatch": bson.M{"v.info_name": bson.M{
						"$exists": true,
						"$type":   "array",
						"$ne":     bson.A{},
					}}}}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					FieldType: pattern.FieldTypeStringArray,
					Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{"$elemMatch": bson.M{"v.info_name": bson.M{
					"$exists": true,
					"$type":   "array",
					"$ne":     bson.A{},
				}}}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.infos_array": bson.M{"$elemMatch": bson.M{"v.info_name": bson.M{
					"$exists": true,
					"$ne":     nil,
				}}}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.infos.info_name",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.infos_array": bson.M{"$in": bson.A{nil, bson.A{}}}},
					{"alarm.v.infos_array": bson.M{"$not": bson.M{"$elemMatch": bson.M{"v.info_name": bson.M{
						"$exists": true,
						"$ne":     nil,
					}}}}},
				}},
			}},
		}},
		bson.M{
			"alarm.v.infos_array": bson.M{"$objectToArray": "$alarm.v.infos"},
		},
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.ack",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.ack": bson.M{"$exists": true, "$ne": nil}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.ack",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.ack": bson.M{"$exists": false}},
					{"alarm.v.ack": bson.M{"$eq": nil}},
				}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.output",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.output": bson.M{
					"$exists": true,
					"$nin":    bson.A{nil, ""},
				}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.output",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.output": bson.M{"$exists": false}},
					{"alarm.v.output": bson.M{"$eq": nil}},
					{"alarm.v.output": bson.M{"$eq": ""}},
				}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.activation_date",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.activation_date": bson.M{"$exists": false}},
					{"alarm.v.activation_date": bson.M{"$eq": nil}},
				}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.activation_date",
					Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.activation_date": bson.M{
					"$gt": datetime.NewCpsTime(from),
					"$lt": datetime.NewCpsTime(to),
				}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
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
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.ticket.ticket": bson.M{"$eq": "test ticket"}},
				{"alarm.v.ticket.m": bson.M{"$eq": "test message"}},
				{"alarm.v.ticket.ticket_data.data_1": bson.M{"$eq": "test_1"}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.change_state",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"alarm.v.change_state": bson.M{"$exists": true, "$ne": nil}},
			}},
		}},
		nil,
		nil,
	)
	f(
		pattern.Alarm{
			{
				{
					Field:     "v.change_state",
					Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
				},
			},
		},
		bson.M{"$or": []bson.M{
			{"$and": []bson.M{
				{"$or": []bson.M{
					{"alarm.v.change_state": bson.M{"$exists": false}},
					{"alarm.v.change_state": bson.M{"$eq": nil}},
				}},
			}},
		}},
		nil,
		nil,
	)
}
