package pattern_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCondition_MatchString(t *testing.T) {
	regexpCond, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test")
	if err != nil {
		t.Errorf("unexpected err %v", err)
	}
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          string
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test equal is matched",
			cond:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          "test",
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test equal is not matched",
			cond:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          "test2",
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test not equal is matched",
			cond:           pattern.NewStringCondition(pattern.ConditionNotEqual, "test"),
			value:          "test2",
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test not equal is not matched",
			cond:           pattern.NewStringCondition(pattern.ConditionNotEqual, "test"),
			value:          "test",
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test is one of is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test2", "test"}),
			value:          "test",
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is one of is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test2", "test3"}),
			value:          "test",
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test regexp is matched",
			cond:           regexpCond,
			value:          "test",
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test regexp is not matched",
			cond:           regexpCond,
			value:          "tesst",
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test bad regexp err",
			cond:           pattern.NewStringCondition(pattern.ConditionRegexp, "["),
			value:          "test",
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a string",
			cond:           pattern.NewIntCondition(pattern.ConditionEqual, 123),
			value:          "test",
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			cond:           pattern.NewStringCondition("some", "test"),
			value:          "test",
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
		{
			testName:       "test exist is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionExist, true),
			value:          "test",
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test not exist is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionExist, false),
			value:          "",
			expectedErr:    nil,
			expectedResult: true,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err := data.cond.MatchString(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}

			result, _, err = data.cond.MatchStringWithRegexpMatches(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchString(t *testing.T) {
	dataSet := map[string]struct {
		cond           pattern.Condition
		value          string
		expectedResult bool
	}{
		"given equal cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: "test",
			},
			value:          "test",
			expectedResult: true,
		},
		"given equal cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: "test",
			},
			value:          "test1",
			expectedResult: false,
		},
		"given not equal cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: "test",
			},
			value:          "test1",
			expectedResult: true,
		},
		"given not equal cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: "test",
			},
			value:          "test",
			expectedResult: false,
		},
		"given is one of cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionIsOneOf,
				Value: []string{"test"},
			},
			value:          "test",
			expectedResult: true,
		},
		"given is one of cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionIsOneOf,
				Value: []string{"test"},
			},
			value:          "test1",
			expectedResult: false,
		},
		"given is not one of cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionIsNotOneOf,
				Value: []string{"test"},
			},
			value:          "test1",
			expectedResult: true,
		},
		"given is not one of cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionIsNotOneOf,
				Value: []string{"test"},
			},
			value:          "test",
			expectedResult: false,
		},
		"given regexp cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "^test.+$",
			},
			value:          "test1",
			expectedResult: true,
		},
		"given regexp cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "^test.+$",
			},
			value:          "1test",
			expectedResult: false,
		},
		"given contain cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionContain,
				Value: "^test.+$",
			},
			value:          "test^test.+$test",
			expectedResult: true,
		},
		"given contain cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionContain,
				Value: "^test.+$",
			},
			value:          "test^test+$test",
			expectedResult: false,
		},
		"given not contain cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotContain,
				Value: "^test.+$",
			},
			value:          "test^test+$test",
			expectedResult: true,
		},
		"given not contain cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotContain,
				Value: "^test.+$",
			},
			value:          "test^test.+$test",
			expectedResult: false,
		},
		"given begin with cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionBeginWith,
				Value: "^test.+$",
			},
			value:          "^test.+$test",
			expectedResult: true,
		},
		"given begin with cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionBeginWith,
				Value: "^test.+$",
			},
			value:          "^test+$test",
			expectedResult: false,
		},
		"given not begin with cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotBeginWith,
				Value: "^test.+$",
			},
			value:          "^test+$test",
			expectedResult: true,
		},
		"given not begin with cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotBeginWith,
				Value: "^test.+$",
			},
			value:          "^test.+$test",
			expectedResult: false,
		},
		"given end with cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionEndWith,
				Value: "^test.+$",
			},
			value:          "test^test.+$",
			expectedResult: true,
		},
		"given end with cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionEndWith,
				Value: "^test.+$",
			},
			value:          "test^test+$",
			expectedResult: false,
		},
		"given not end with cond should match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotEndWith,
				Value: "^test.+$",
			},
			value:          "test^test+$",
			expectedResult: true,
		},
		"given not end with cond should not match": {
			cond: pattern.Condition{
				Type:  pattern.ConditionNotEndWith,
				Value: "^test.+$",
			},
			value:          "test^test.+$",
			expectedResult: false,
		},
	}

	for name, data := range dataSet {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(condWrapper{Cond: data.cond})
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}

			w := condWrapper{}
			err = json.Unmarshal(b, &w)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}

			result, err := w.Cond.MatchString(data.value)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			result, _, err = w.Cond.MatchStringWithRegexpMatches(data.value)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			b, err = bson.Marshal(condWrapper{Cond: data.cond})
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}

			w = condWrapper{}
			err = bson.Unmarshal(b, &w)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}

			result, err = w.Cond.MatchString(data.value)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			result, _, err = w.Cond.MatchStringWithRegexpMatches(data.value)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}
		})
	}
}

func TestCondition_MatchInt(t *testing.T) {
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          int64
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test equal is matched",
			cond:           pattern.NewIntCondition(pattern.ConditionEqual, 10),
			value:          10,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test equal is not matched",
			cond:           pattern.NewIntCondition(pattern.ConditionEqual, 10),
			value:          9,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test not equal is matched",
			cond:           pattern.NewIntCondition(pattern.ConditionNotEqual, 10),
			value:          9,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test not equal is not matched",
			cond:           pattern.NewIntCondition(pattern.ConditionNotEqual, 10),
			value:          10,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test gt is matched",
			cond:           pattern.NewIntCondition(pattern.ConditionGT, 10),
			value:          11,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test gt is not matched",
			cond:           pattern.NewIntCondition(pattern.ConditionGT, 10),
			value:          9,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test lt is matched",
			cond:           pattern.NewIntCondition(pattern.ConditionLT, 10),
			value:          9,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test lt is not matched",
			cond:           pattern.NewIntCondition(pattern.ConditionLT, 10),
			value:          11,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not an int",
			cond:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          10,
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			cond:           pattern.NewIntCondition("some", 10),
			value:          10,
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err := data.cond.MatchInt(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchInt(t *testing.T) {
	dataSet := []pattern.Condition{
		{
			Type:  pattern.ConditionEqual,
			Value: int32(20),
		},
		{
			Type:  pattern.ConditionNotEqual,
			Value: int64(20),
		},
		{
			Type:  pattern.ConditionGT,
			Value: int32(20),
		},
		{
			Type:  pattern.ConditionLT,
			Value: float64(20),
		},
	}

	for _, data := range dataSet {
		b, err := json.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchInt(100)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchInt(100)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_MatchBool(t *testing.T) {
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          bool
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test equal is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionEqual, true),
			value:          true,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test equal is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionEqual, true),
			value:          false,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			cond:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          false,
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			cond:           pattern.NewBoolCondition("some", true),
			value:          false,
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err := data.cond.MatchBool(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchBool(t *testing.T) {
	dataSet := []pattern.Condition{
		{
			Type:  pattern.ConditionEqual,
			Value: true,
		},
		{
			Type:  pattern.ConditionEqual,
			Value: false,
		},
	}

	for _, data := range dataSet {
		b, err := json.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchBool(true)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchBool(true)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_MatchStringArray(t *testing.T) {
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          []string
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test is empty is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			value:          []string{},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is empty is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			value:          []string{"test"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test is not empty is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			value:          []string{"test"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			value:          []string{},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"test1", "test2", "test3"}),
			value:          []string{"test4"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_not is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"test1", "test2", "test3"}),
			value:          []string{"test3", "test4"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_one_of is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"test1", "test2", "test3"}),
			value:          []string{"test3", "test4"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_one_of is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"test1", "test2", "test3"}),
			value:          []string{"test4"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_every is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"test1", "test2", "test3"}),
			value:          []string{"test1", "test2", "test3", "test4"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_every is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"test1", "test2", "test3"}),
			value:          []string{"test1", "test2"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionIsEmpty, []string{"test"}),
			value:          []string{"test"},
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a slice of strings",
			cond:           pattern.NewBoolCondition(pattern.ConditionHasEvery, true),
			value:          []string{"test"},
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			cond:           pattern.NewStringArrayCondition("some", []string{"test"}),
			value:          []string{"test"},
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
	}

	var result bool
	var err error

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err = data.cond.MatchStringArray(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchStringArray(t *testing.T) {
	dataSet := []pattern.Condition{
		{
			Type:  pattern.ConditionIsEmpty,
			Value: false,
		},
		{
			Type:  pattern.ConditionIsEmpty,
			Value: true,
		},
		{
			Type:  pattern.ConditionHasEvery,
			Value: []string{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasOneOf,
			Value: []any{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasNot,
			Value: bson.A{"val1", "val2"},
		},
	}

	for _, data := range dataSet {
		b, err := json.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchStringArray([]string{"test"})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchStringArray([]string{"test"})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_MatchRef(t *testing.T) {
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          any
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test is empty is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionExist, true),
			value:          nil,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test is empty is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionExist, true),
			value:          &struct{}{},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionExist, false),
			value:          nil,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionExist, false),
			value:          &struct{}{},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			cond:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          &struct{}{},
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			cond:           pattern.NewBoolCondition("some", true),
			value:          &struct{}{},
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err := data.cond.MatchRef(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchRef(t *testing.T) {
	dataSet := []pattern.Condition{
		{
			Type:  pattern.ConditionExist,
			Value: false,
		},
		{
			Type:  pattern.ConditionExist,
			Value: true,
		},
	}

	for _, data := range dataSet {
		b, err := json.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchRef(nil)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchRef(nil)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_MatchTime(t *testing.T) {
	f := func(cond pattern.Condition, value, testTime time.Time, expectedResult bool, expectedError error) {
		t.Helper()

		result, err := cond.MatchTime(value, testTime)
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error %v but got %v", expectedError, err)
		}

		if result != expectedResult {
			t.Errorf("expected %t but got %t", expectedResult, result)
		}
	}

	testDuration := datetime.NewDurationWithUnit(100, "s")
	timeRelativeCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, testDuration)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	timeOlderThanCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.DurationWithUnit{}, testDuration)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	timeBetweenCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(200, "s"), testDuration)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	testTime := time.Now()

	// test relative is matched
	f(timeRelativeCond, testTime.Add(-50*time.Second), testTime, true, nil)
	// test relative is not matched
	f(timeRelativeCond, testTime.Add(-150*time.Second), testTime, false, nil)
	// test relative is not matched to edge
	f(timeRelativeCond, testDuration.SubFrom(datetime.CpsTime{Time: testTime}).Time, testTime, false, nil)
	// test relative time older than is matched
	f(timeOlderThanCond, testTime.Add(-1*time.Hour), testTime, true, nil)
	// test relative time older than is matched to edge
	f(timeOlderThanCond, testDuration.SubFrom(datetime.CpsTime{Time: testTime.Add(-1 * time.Nanosecond)}).Time, testTime, true, nil)
	// test relative time older than is not matched
	f(timeOlderThanCond, testTime.Add(-50*time.Second), testTime, false, nil)
	// test relative time between is matched
	f(timeBetweenCond, testTime.Add(-150*time.Second), testTime, true, nil)
	// test relative time between is not matched
	f(timeBetweenCond, testTime.Add(-50*time.Second), testTime, false, nil)
	// test relative time between is not matched
	f(timeBetweenCond, testTime.Add(-250*time.Second), testTime, false, nil)
	// test absolute is matched
	f(pattern.NewTimeIntervalCondition(
		pattern.ConditionTimeAbsolute,
		testTime.Add(-1000*time.Second).Unix(),
		testTime.Add(1000*time.Second).Unix(),
	), testTime, testTime, true, nil)
	// test absolute is not matched left
	f(pattern.NewTimeIntervalCondition(
		pattern.ConditionTimeAbsolute,
		testTime.Add(-1000*time.Second).Unix(),
		testTime.Add(1000*time.Second).Unix(),
	), testTime.Add(-2000*time.Second), testTime, false, nil)
	// test absolute is not matched right
	f(pattern.NewTimeIntervalCondition(
		pattern.ConditionTimeAbsolute,
		testTime.Add(-1000*time.Second).Unix(),
		testTime.Add(1000*time.Second).Unix(),
	), testTime.Add(2000*time.Second), testTime, false, nil)
	// test relative wrong condition value type
	f(pattern.NewStringCondition(pattern.ConditionTimeRelative, "test"), testTime.Add(-50*time.Second), testTime, false, pattern.ErrWrongConditionValue)
	// test absolute condition value is wrong type
	f(pattern.NewStringCondition(pattern.ConditionTimeAbsolute, "test"), testTime, testTime, false, pattern.ErrWrongConditionValue)
	// test unexpected type
	f(pattern.NewTimeIntervalCondition(
		"some",
		testTime.Add(-1000*time.Second).Unix(),
		testTime.Add(1000*time.Second).Unix(),
	), testTime, testTime, false, pattern.ErrUnsupportedConditionType)
	// test relative time older than condition value is wrong type
	f(pattern.NewStringCondition(pattern.ConditionTimeRelative, "test"), testTime.Add(50*time.Second), testTime, false, pattern.ErrWrongConditionValue)
}

func TestCondition_UnmarshalAndMatchTime(t *testing.T) {
	f := func(cond pattern.Condition) {
		t.Helper()

		b, err := json.Marshal(condWrapper{Cond: cond})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		testTime := time.Now()
		_, err = w.Cond.MatchTime(testTime, testTime)
		if err != nil {
			t.Errorf("expected no error after JSON unmarshal but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: cond})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchTime(testTime, testTime)
		if err != nil {
			t.Errorf("expected no error after BSON unmarshal but got %v", err)
		}
	}

	// test relative wrong condition value type
	f(pattern.Condition{
		Type:  pattern.ConditionTimeRelative,
		Value: datetime.NewDurationWithUnit(100, "s"),
	})
	// test relative older than condition value
	f(pattern.Condition{
		Type: pattern.ConditionTimeRelative,
		Value: map[string]*datetime.DurationWithUnit{
			"to": {Value: 100, Unit: "s"},
		},
	})
	// test relative between durations
	f(pattern.Condition{
		Type: pattern.ConditionTimeRelative,
		Value: map[string]*datetime.DurationWithUnit{
			"from": {Value: 1, Unit: "h"},
			"to":   {Value: 100, Unit: "s"},
		},
	})
	// test absolute condition value
	f(pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 10, 10))
}

func TestCondition_MatchDuration(t *testing.T) {
	durationGtCond, err := pattern.NewDurationCondition(pattern.ConditionGT, datetime.DurationWithUnit{Value: 5, Unit: "m"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	durationLtCond, err := pattern.NewDurationCondition(pattern.ConditionLT, datetime.DurationWithUnit{Value: 5, Unit: "m"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	durationSomeCond, err := pattern.NewDurationCondition("some", datetime.DurationWithUnit{Value: 5, Unit: "m"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          int64
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test gt is matched",
			cond:           durationGtCond,
			value:          350,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test gt is not matched",
			cond:           durationGtCond,
			value:          250,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test lt is matched",
			cond:           durationLtCond,
			value:          250,
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test lt is not matched",
			cond:           durationLtCond,
			value:          350,
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test unexpected type",
			cond:           durationSomeCond,
			value:          250,
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
		{
			testName:       "test value is not a map",
			cond:           pattern.NewStringCondition("some", "test"),
			value:          250,
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err := data.cond.MatchDuration(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchDuration(t *testing.T) {
	dataSet := []pattern.Condition{
		{
			Type: pattern.ConditionGT,
			Value: map[string]any{
				"value": 5,
				"unit":  "m",
			},
		},
		{
			Type: pattern.ConditionLT,
			Value: map[string]any{
				"value": 5,
				"unit":  "m",
			},
		},
	}

	for _, data := range dataSet {
		b, err := json.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchDuration(10)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchDuration(10)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_TimeToMongoQuery(t *testing.T) {
	f := func(cond pattern.Condition, field string, testTime time.Time, expectedResult bson.M, expectedError error) {
		t.Helper()

		result, err := cond.TimeToMongoQuery(field, testTime)
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error %v but got %v", expectedError, err)
		}

		if diff := pretty.Compare(result, expectedResult); diff != "" {
			t.Errorf("unexpected result (-got +want): %s", diff)
		}
	}
	testDuration := datetime.DurationWithUnit{
		Value: 100,
		Unit:  "s",
	}
	timeRelativeCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, testDuration)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	timeOlderThanCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.DurationWithUnit{}, testDuration)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	timeBetweenCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(1, "h"), testDuration)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	timeBetweenValue, ok := timeBetweenCond.Value.(map[string]*datetime.DurationWithUnit)
	if !ok {
		t.Errorf("invalid type for timeBetweenCond.Value")
	}
	testTime := time.Now()
	timeBetweenResult := bson.M{
		"$gt": timeBetweenValue["from"].SubFrom(datetime.NewCpsTime(testTime.Unix())),
		"$lt": timeBetweenValue["to"].SubFrom(datetime.NewCpsTime(testTime.Unix())),
	}
	fromTime, toTime := time.Now().Add(-1000*time.Second).Unix(), time.Now().Add(1000*time.Second).Unix()
	timeAbsCond := pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, fromTime, toTime)
	// test time relative
	f(timeRelativeCond, "test", testTime, bson.M{"test": bson.M{"$gt": datetime.NewCpsTime(testTime.Add(time.Duration(-testDuration.Value) * time.Second).Unix())}}, nil)
	// test time older than
	f(timeOlderThanCond, "test", testTime, bson.M{"test": bson.M{"$lt": datetime.NewCpsTime(testTime.Add(time.Duration(-testDuration.Value) * time.Second).Unix())}}, nil)
	// test time between
	f(timeBetweenCond, "test", testTime, bson.M{"test": timeBetweenResult}, nil)
	// test time absolute
	f(timeAbsCond, "test", testTime, bson.M{"test": bson.M{"$gt": datetime.NewCpsTime(fromTime), "$lt": datetime.NewCpsTime(toTime)}}, nil)
	// unsupported condition type error
	f(pattern.Condition{Type: pattern.ConditionEqual, Value: "test"}, "test", testTime, nil, pattern.ErrUnsupportedConditionType)
	// ErrWrongConditionValue error
	f(pattern.Condition{Type: pattern.ConditionTimeRelative, Value: 10}, "test", testTime, nil, pattern.ErrWrongConditionValue)
	f(pattern.Condition{Type: pattern.ConditionTimeAbsolute, Value: 10}, "test", testTime, nil, pattern.ErrWrongConditionValue)
}

func TestCondition_MatchTags(t *testing.T) {
	dataSet := []struct {
		testName       string
		cond           pattern.Condition
		value          []string
		expectedErr    error
		expectedResult bool
	}{
		{
			testName:       "test is empty is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			value:          []string{},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is empty is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			value:          []string{"label-1: test-1"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test is not empty is matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			value:          []string{"label-1: test-1"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is not matched",
			cond:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			value:          []string{},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_not is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"}),
			value:          []string{"label-4: value-4"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_not is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"}),
			value:          []string{"label-3: value-3", "label-4: value-4"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_one_of is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"}),
			value:          []string{"label-3: value-3", "label-4: value-4"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_one_of is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"}),
			value:          []string{"label-4: value-4"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_every is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"}),
			value:          []string{"label-1: value-1", "label-2: value-2", "label-3: value-3", "label-4: value-4"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_every is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"}),
			value:          []string{"label-1: value-1", "label-2: value-2"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test has_labels is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasLabels, []string{"label-1", "label-4"}),
			value:          []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_labels is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasLabels, []string{"label-4"}),
			value:          []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test tag without value has_labels is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasLabels, []string{"label-2"}),
			value:          []string{"label-1: value-1", "label-2", "label-3: value-3"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_not_labels is matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNotLabels, []string{"label-4"}),
			value:          []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"},
			expectedErr:    nil,
			expectedResult: true,
		},
		{
			testName:       "test has_not_labels is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNotLabels, []string{"label-1", "label-4"}),
			value:          []string{"label-1: value-1", "label-2: value-2", "label-3: value-3"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test tag without value has_not_labels is not matched",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasNotLabels, []string{"label-2"}),
			value:          []string{"label-1: value-1", "label-2", "label-3: value-3"},
			expectedErr:    nil,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionIsEmpty, []string{"test"}),
			value:          []string{"test"},
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a slice of strings",
			cond:           pattern.NewBoolCondition(pattern.ConditionHasEvery, true),
			value:          []string{"test"},
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond value is empty slice of strings",
			cond:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{}),
			value:          []string{"test"},
			expectedErr:    pattern.ErrWrongConditionValue,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			cond:           pattern.NewStringArrayCondition("some", []string{"test"}),
			value:          []string{"test"},
			expectedErr:    pattern.ErrUnsupportedConditionType,
			expectedResult: false,
		},
	}

	var result bool
	var err error

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result, err = data.cond.MatchTags(data.value)
			if result != data.expectedResult {
				t.Errorf("expected %t but got %t", data.expectedResult, result)
			}

			if !errors.Is(err, data.expectedErr) {
				t.Errorf("expected error %v but got %v", data.expectedErr, err)
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchTags(t *testing.T) {
	dataSet := []pattern.Condition{
		{
			Type:  pattern.ConditionIsEmpty,
			Value: false,
		},
		{
			Type:  pattern.ConditionIsEmpty,
			Value: true,
		},
		{
			Type:  pattern.ConditionHasEvery,
			Value: []string{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasOneOf,
			Value: []any{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasNot,
			Value: bson.A{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasLabels,
			Value: []string{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasNotLabels,
			Value: []any{"val1", "val2"},
		},
	}

	for _, data := range dataSet {
		b, err := json.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchTags([]string{"test"})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: data})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchTags([]string{"test"})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

type condWrapper struct {
	Cond pattern.Condition `bson:"cond" json:"cond"`
}
