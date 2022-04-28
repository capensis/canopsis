package pattern_test

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type eventDataSet struct {
	pattern     pattern.Event
	event       types.Event
	matchErr    error
	matchResult bool
}

func TestEvent_Match(t *testing.T) {
	dataSets := getEventMatchDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, _, err := data.pattern.Match(data.event)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func BenchmarkEvent_Match_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "resource",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
	}
	event := types.Event{
		Resource: "test name",
	}

	benchmarkEventMatch(b, cond, event)
}

func BenchmarkEvent_Match_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "resource",
		Condition: regexpCondition,
	}
	event := types.Event{
		Resource: "test name",
	}

	benchmarkEventMatch(b, cond, event)
}

func BenchmarkEvent_Match_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "extra.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test 2"),
	}
	event := types.Event{
		ExtraInfos: map[string]interface{}{
			"test": "test",
		},
	}

	benchmarkEventMatch(b, cond, event)
}

func BenchmarkEvent_Match_Infos_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "extra.test",
		FieldType: pattern.FieldTypeString,
		Condition: regexpCondition,
	}
	event := types.Event{
		ExtraInfos: map[string]interface{}{
			"test": "test",
		},
	}

	benchmarkEventMatch(b, cond, event)
}

func BenchmarkEvent_UnmarshalBsonAndMatch_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "resource",
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	event := types.Event{
		Resource: "test name",
	}

	benchmarkEventUnmarshalBsonAndMatch(b, cond, event)
}

func BenchmarkEvent_UnmarshalBsonAndMatch_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "extra.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test 2",
		},
	}
	event := types.Event{
		ExtraInfos: map[string]interface{}{
			"test": "test",
		},
	}

	benchmarkEventUnmarshalBsonAndMatch(b, cond, event)
}

func BenchmarkEvent_UnmarshalBsonAndMatch_Infos_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "extra.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	event := types.Event{
		ExtraInfos: map[string]interface{}{
			"test": "test",
		},
	}

	benchmarkEventUnmarshalBsonAndMatch(b, cond, event)
}

func benchmarkEventMatch(b *testing.B, fieldCond pattern.FieldCondition, event types.Event) {
	const size = 100
	p := make(pattern.Event, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := p.Match(event)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func benchmarkEventUnmarshalBsonAndMatch(b *testing.B, fieldCond pattern.FieldCondition, event types.Event) {
	const size = 100
	p := make(pattern.Event, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	type wrapper struct {
		Pattern pattern.Event `bson:"pattern"`
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

		_, _, err = w.Pattern.Match(event)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func getEventMatchDataSets() map[string]eventDataSet {
	return map[string]eventDataSet{
		"given empty pattern should match": {
			pattern: pattern.Event{},
			event: types.Event{
				Resource: "test name",
			},
			matchResult: true,
		},
		"given string field condition should match": {
			pattern: pattern.Event{
				{
					{
						Field:     "resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event: types.Event{
				Resource: "test name",
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.Event{
				{
					{
						Field:     "resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event: types.Event{
				Resource: "test another name",
			},
			matchResult: false,
		},
		"given integer field condition and not integer field should return error": {
			pattern: pattern.Event{
				{
					{
						Field:     "resource",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 1),
					},
				},
			},
			event:    types.Event{},
			matchErr: pattern.ErrWrongConditionValue,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.Event{
				{
					{
						Field:     "created",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event:    types.Event{},
			matchErr: pattern.ErrUnsupportedField,
		},
		"given string extra field condition should match": {
			pattern: pattern.Event{
				{
					{
						Field:     "extra.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event: types.Event{
				ExtraInfos: map[string]interface{}{
					"info_name": "test name",
				},
			},
			matchResult: true,
		},
		"given string extra field condition should not match": {
			pattern: pattern.Event{
				{
					{
						Field:     "extra.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event: types.Event{
				ExtraInfos: map[string]interface{}{
					"info_name": "test another name",
				},
			},
			matchResult: false,
		},
		"given string extra field condition and not string extra field should not match": {
			pattern: pattern.Event{
				{
					{
						Field:     "extra.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event: types.Event{
				ExtraInfos: map[string]interface{}{
					"info_name": 2,
				},
			},
			matchResult: false,
		},
		"given extra field condition and unknown extra field should not match": {
			pattern: pattern.Event{
				{
					{
						Field:     "extra.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			event:       types.Event{},
			matchResult: false,
		},
	}
}
