package match_test

import (
	"errors"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMatchAlarmPattern(t *testing.T) {
	dataSets := getMatchAlarmPatternDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := match.MatchAlarmPattern(data.pattern, &data.alarm)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func getMatchAlarmPatternDataSets() map[string]alarmDataSet {
	timeRelativeCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.DurationWithUnit{
		Value: 1,
		Unit:  datetime.DurationUnitMinute,
	})
	if err != nil {
		panic(err)
	}

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
						Field:     "v.unknown",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			alarm:       types.Alarm{},
			matchErr:    pattern.ErrUnsupportedField,
			matchResult: false,
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
		"given string info ne condition and not exist info should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "test name"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{},
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
			matchErr:    pattern.ErrWrongConditionValue,
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
		"given exist ack condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ack",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{},
				},
			},
			matchResult: true,
		},
		"given not exist ack condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ack",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{},
			},
			matchResult: false,
		},
		"given exist change_state condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.change_state",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					ChangeState: &types.AlarmStep{},
				},
			},
			matchResult: true,
		},
		"given not exist change_state condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.change_state",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					ChangeState: &types.AlarmStep{},
				},
			},
			matchResult: false,
		},
		"given not exist v.output condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.output",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{},
			},
			matchResult: true,
		},
		"given not exist v.activation_date condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.activation_date",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{},
			},
			matchResult: true,
		},
		"given time v.activation_date condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.activation_date",
						Condition: timeRelativeCond,
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					ActivationDate: &datetime.CpsTime{Time: time.Now().Add(-30 * time.Second)},
				},
			},
			matchResult: true,
		},
		"given match ticket's message condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.m",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						Message: "test",
					},
				},
			},
			matchResult: true,
		},
		"given match ticket's message condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.m",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						Message: "test 2",
					},
				},
			},
			matchResult: false,
		},
		"given match ticket's ticket condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						TicketInfo: types.TicketInfo{
							Ticket: "test",
						},
					},
				},
			},
			matchResult: true,
		},
		"given match ticket's ticket condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						TicketInfo: types.TicketInfo{
							Ticket: "test 2",
						},
					},
				},
			},
			matchResult: false,
		},
		"given match ticket data condition should match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket_data.data_1",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test_1"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						TicketInfo: types.TicketInfo{
							TicketData: map[string]string{
								"data_1": "test_1",
								"data_2": "test_2",
							},
						},
					},
				},
			},
			matchResult: true,
		},
		"given match ticket data condition should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket_data.data_2",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test_1"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						TicketInfo: types.TicketInfo{
							TicketData: map[string]string{
								"data_1": "test_1",
								"data_2": "test_2",
							},
						},
					},
				},
			},
			matchResult: false,
		},
		"given match ticket data condition without ticket data should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket_data.data_1",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						TicketInfo: types.TicketInfo{},
					},
				},
			},
			matchResult: false,
		},
		"given match ticket data condition without ticket should not match": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket_data.data_1",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{},
			},
			matchResult: false,
		},
		"given match ticket data condition with wrong condition should return err": {
			pattern: pattern.Alarm{
				{
					{
						Field:     "v.ticket.ticket_data.data_1",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 123),
					},
				},
			},
			alarm: types.Alarm{
				Value: types.AlarmValue{
					Ticket: &types.AlarmStep{
						TicketInfo: types.TicketInfo{
							TicketData: map[string]string{
								"data_1": "test_1",
								"data_2": "test_2",
							},
						},
					},
				},
			},
			matchErr:    pattern.ErrWrongConditionValue,
			matchResult: false,
		},
	}
}

type alarmDataSet struct {
	pattern     pattern.Alarm
	alarm       types.Alarm
	matchErr    error
	matchResult bool
}

func BenchmarkMatchAlarmPattern_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "v.display_name",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name 2"),
	}
	alarm := types.Alarm{
		Value: types.AlarmValue{
			DisplayName: "test name",
		},
	}

	benchmarkMatchAlarmPattern(b, cond, alarm)
}

func BenchmarkMatchAlarmPattern_Regexp(b *testing.B) {
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

	benchmarkMatchAlarmPattern(b, cond, alarm)
}

func BenchmarkMatchAlarmPattern_Infos_Equal(b *testing.B) {
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

	benchmarkMatchAlarmPattern(b, cond, alarm)
}

func BenchmarkMatchAlarmPattern_Infos_Regexp(b *testing.B) {
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

	benchmarkMatchAlarmPattern(b, cond, alarm)
}

func BenchmarkMatchAlarmPattern_UnmarshalBson_Equal(b *testing.B) {
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

	benchmarkMatchAlarmPatternUnmarshalBson(b, cond, []types.Alarm{alarm})
}

func BenchmarkMatchAlarmPattern_UnmarshalBson_Regexp(b *testing.B) {
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

	benchmarkMatchAlarmPatternUnmarshalBson(b, cond, []types.Alarm{alarm})
}

func BenchmarkMatchAlarmPattern_UnmarshalBson_Infos_Equal(b *testing.B) {
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

	benchmarkMatchAlarmPatternUnmarshalBson(b, cond, []types.Alarm{alarm})
}

func BenchmarkMatchAlarmPattern_UnmarshalBson_Infos_Regexp(b *testing.B) {
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

	benchmarkMatchAlarmPatternUnmarshalBson(b, cond, []types.Alarm{alarm})
}

func BenchmarkMatchAlarmPattern_ManyAlarms_UnmarshalBson_Regexp(b *testing.B) {
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

	benchmarkMatchAlarmPatternUnmarshalBson(b, cond, alarms)
}

func benchmarkMatchAlarmPattern(b *testing.B, fieldCond pattern.FieldCondition, alarm types.Alarm) {
	const size = 100
	p := make(pattern.Alarm, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := match.MatchAlarmPattern(p, &alarm)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func benchmarkMatchAlarmPatternUnmarshalBson(b *testing.B, fieldCond pattern.FieldCondition, alarms []types.Alarm) {
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
			_, err = match.MatchAlarmPattern(w.Pattern, &alarm)
			if err != nil {
				b.Fatalf("unexpected error %v", err)
			}
		}
	}
}
