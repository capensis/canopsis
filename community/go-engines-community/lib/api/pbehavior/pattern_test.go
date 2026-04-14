package pbehavior

import (
	"errors"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSearchPatternToMongoQuery(t *testing.T) {
	now := time.Now()
	from := now.Add(-time.Hour).Unix()
	to := now.Add(time.Hour).Unix()

	timeRelativeCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.DurationWithUnit{
		Value: 3600,
		Unit:  datetime.DurationUnitSecond,
	})
	if err != nil {
		t.Fatalf("failed to create time relative condition: %v", err)
	}

	dataSets := map[string]struct {
		pattern     searchPattern
		expectedRes bson.M
		expectedErr error
	}{
		"given empty pattern": {
			pattern:     searchPattern{},
			expectedRes: nil,
			expectedErr: nil,
		},
		"given name field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name"}},
				}},
			}},
		},
		"given author field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "author",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-author"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"author": bson.M{"$eq": "test-author"}},
				}},
			}},
		},
		"given rrule field with exist true condition": {
			pattern: searchPattern{
				{
					{
						Field:     "rrule",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"rrule": bson.M{
						"$exists": true,
						"$nin":    bson.A{nil, ""},
					}},
				}},
			}},
		},
		"given rrule field with exist false condition": {
			pattern: searchPattern{
				{
					{
						Field:     "rrule",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$or": []bson.M{
						{"rrule": bson.M{"$exists": false}},
						{"rrule": bson.M{"$eq": nil}},
						{"rrule": bson.M{"$eq": ""}},
					}},
				}},
			}},
		},
		"given reason field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "reason",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-reason"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"reason": bson.M{"$eq": "test-reason"}},
				}},
			}},
		},
		"given type field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-type"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"type": bson.M{"$eq": "test-type"}},
				}},
			}},
		},
		"given enabled field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "enabled",
						Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"enabled": bson.M{"$eq": true}},
				}},
			}},
		},
		"given tstart field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "tstart",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"tstart": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
		},
		"given tstop field with relative time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "tstop",
						Condition: timeRelativeCond,
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"tstop": bson.M{
						"$gt": datetime.NewCpsTime(now.Add(-3600 * time.Second).Unix()),
					}},
				}},
			}},
		},
		"given rrule_end field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "rrule_end",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"rrule_end": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
		},
		"given created field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "created",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"created": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
		},
		"given updated field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "updated",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"updated": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
		},
		"given last_alarm_date field with absolute time condition": {
			pattern: searchPattern{
				{
					{
						Field:     "last_alarm_date",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, from, to),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"last_alarm_date": bson.M{
						"$gt": datetime.NewCpsTime(from),
						"$lt": datetime.NewCpsTime(to),
					}},
				}},
			}},
		},
		"given alarm_count field with equal condition": {
			pattern: searchPattern{
				{
					{
						Field:     "alarm_count",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 5),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"alarm_count": bson.M{"$eq": int64(5)}},
				}},
			}},
		},
		"given multiple conditions in one group": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
					{
						Field:     "enabled",
						Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
					},
					{
						Field:     "alarm_count",
						Condition: pattern.NewIntCondition(pattern.ConditionGT, 0),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name"}},
					{"enabled": bson.M{"$eq": true}},
					{"alarm_count": bson.M{"$gt": int64(0)}},
				}},
			}},
		},
		"given multiple groups": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name-1"),
					},
				},
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name-2"),
					},
				},
			},
			expectedRes: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name-1"}},
				}},
				{"$and": []bson.M{
					{"name": bson.M{"$eq": "test-name-2"}},
				}},
			}},
		},
		"given unsupported field": {
			pattern: searchPattern{
				{
					{
						Field:     "unsupported_field",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrUnsupportedField,
		},
		"given invalid condition type for string field": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value for string field": {
			pattern: searchPattern{
				{
					{
						Field:     "name",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 123),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrWrongConditionValue,
		},
		"given invalid condition type for bool field": {
			pattern: searchPattern{
				{
					{
						Field:     "enabled",
						Condition: pattern.NewBoolCondition(pattern.ConditionGT, true),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value for bool field": {
			pattern: searchPattern{
				{
					{
						Field:     "enabled",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 1),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrWrongConditionValue,
		},
		"given invalid condition type for time field": {
			pattern: searchPattern{
				{
					{
						Field:     "tstart",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "2024-01-01"),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition type for int field": {
			pattern: searchPattern{
				{
					{
						Field:     "alarm_count",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "5"),
					},
				},
			},
			expectedRes: nil,
			expectedErr: pattern.ErrWrongConditionValue,
		},
	}

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := data.pattern.ToMongoQuery()
			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
			if diff := pretty.Compare(query, data.expectedRes); diff != "" {
				t.Errorf("unexpected query result %s", diff)
			}
		})
	}
}
