package pattern_test

import (
	"encoding/json"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCondition_MatchString(t *testing.T) {
	regexpCond, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test")
	if err != nil {
		t.Errorf("unexpected err %v", err)
	}
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          string
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test equal is matched",
			conf:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test equal is not matched",
			conf:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          "test2",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test not equal is matched",
			conf:           pattern.NewStringCondition(pattern.ConditionNotEqual, "test"),
			value:          "test2",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test not equal is not matched",
			conf:           pattern.NewStringCondition(pattern.ConditionNotEqual, "test"),
			value:          "test",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test is one of is matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test2", "test"}),
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test is one of is not matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"test2", "test3"}),
			value:          "test",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test regexp is matched",
			conf:           regexpCond,
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test regexp is not matched",
			conf:           regexpCond,
			value:          "tesst",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test bad regexp err",
			conf:           pattern.NewStringCondition(pattern.ConditionRegexp, "["),
			value:          "test",
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a string",
			conf:           pattern.NewIntCondition(pattern.ConditionEqual, 123),
			value:          "test",
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			conf:           pattern.NewStringCondition("some", "test"),
			value:          "test",
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, _, err := dataset.conf.MatchString(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchString(t *testing.T) {
	dataSets := []pattern.Condition{
		{
			Type:  pattern.ConditionEqual,
			Value: "test",
		},
		{
			Type:  pattern.ConditionNotEqual,
			Value: "test",
		},
		{
			Type:  pattern.ConditionIsOneOf,
			Value: []string{"test"},
		},
		{
			Type:  pattern.ConditionIsNotOneOf,
			Value: []string{"test"},
		},
		{
			Type:  pattern.ConditionRegexp,
			Value: "^test.+$",
		},
	}

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, _, err = w.Cond.MatchString("test")
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, _, err = w.Cond.MatchString("test")
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_MatchInt(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          int64
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test equal is matched",
			conf:           pattern.NewIntCondition(pattern.ConditionEqual, 10),
			value:          10,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test equal is not matched",
			conf:           pattern.NewIntCondition(pattern.ConditionEqual, 10),
			value:          9,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test not equal is matched",
			conf:           pattern.NewIntCondition(pattern.ConditionNotEqual, 10),
			value:          9,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test not equal is not matched",
			conf:           pattern.NewIntCondition(pattern.ConditionNotEqual, 10),
			value:          10,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test gt is matched",
			conf:           pattern.NewIntCondition(pattern.ConditionGT, 10),
			value:          11,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test gt is not matched",
			conf:           pattern.NewIntCondition(pattern.ConditionGT, 10),
			value:          9,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test lt is matched",
			conf:           pattern.NewIntCondition(pattern.ConditionLT, 10),
			value:          9,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test lt is not matched",
			conf:           pattern.NewIntCondition(pattern.ConditionLT, 10),
			value:          11,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not an int",
			conf:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          10,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			conf:           pattern.NewIntCondition("some", 10),
			value:          10,
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err := dataset.conf.MatchInt(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchInt(t *testing.T) {
	dataSets := []pattern.Condition{
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

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
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

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
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
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          bool
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test equal is matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionEqual, true),
			value:          true,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test equal is not matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionEqual, true),
			value:          false,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			conf:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          false,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			conf:           pattern.NewBoolCondition("some", true),
			value:          false,
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err := dataset.conf.MatchBool(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchBool(t *testing.T) {
	dataSets := []pattern.Condition{
		{
			Type:  pattern.ConditionEqual,
			Value: true,
		},
		{
			Type:  pattern.ConditionEqual,
			Value: false,
		},
	}

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
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

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
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
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          []string
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test is empty is matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			value:          []string{},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test is empty is not matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			value:          []string{"test"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test is not empty is matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			value:          []string{"test"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is not matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			value:          []string{},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test has_is not matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"test1", "test2", "test3"}),
			value:          []string{"test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test has_not is not matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"test1", "test2", "test3"}),
			value:          []string{"test3", "test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test has_one_of is matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"test1", "test2", "test3"}),
			value:          []string{"test3", "test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test has_one_of is not matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"test1", "test2", "test3"}),
			value:          []string{"test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test has_every is matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"test1", "test2", "test3"}),
			value:          []string{"test1", "test2", "test3", "test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test has_every is not matched",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"test1", "test2", "test3"}),
			value:          []string{"test1", "test2"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			conf:           pattern.NewStringArrayCondition(pattern.ConditionIsEmpty, []string{"test"}),
			value:          []string{"test"},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a slice of strings",
			conf:           pattern.NewBoolCondition(pattern.ConditionHasEvery, true),
			value:          []string{"test"},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			conf:           pattern.NewStringArrayCondition("some", []string{"test"}),
			value:          []string{"test"},
			expectedErr:    true,
			expectedResult: false,
		},
	}

	var result bool
	var err error

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err = dataset.conf.MatchStringArray(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchStringArray(t *testing.T) {
	dataSets := []pattern.Condition{
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
			Value: []interface{}{"val1", "val2"},
		},
		{
			Type:  pattern.ConditionHasNot,
			Value: bson.A{"val1", "val2"},
		},
	}

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
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

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
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
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          interface{}
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test is empty is matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionExist, true),
			value:          nil,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test is empty is not matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionExist, true),
			value:          &struct{}{},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionExist, false),
			value:          nil,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test is not empty is not matched",
			conf:           pattern.NewBoolCondition(pattern.ConditionExist, false),
			value:          &struct{}{},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test cond value is not a bool",
			conf:           pattern.NewStringCondition(pattern.ConditionEqual, "test"),
			value:          &struct{}{},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test cond type is not supported",
			conf:           pattern.NewBoolCondition("some", true),
			value:          &struct{}{},
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err := dataset.conf.MatchRef(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchRef(t *testing.T) {
	dataSets := []pattern.Condition{
		{
			Type:  pattern.ConditionExist,
			Value: false,
		},
		{
			Type:  pattern.ConditionExist,
			Value: true,
		},
	}

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
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

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
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
	timeRelativeCond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, types.DurationWithUnit{
		Value: 100,
		Unit:  "s",
	})
	if err != nil {
		panic(err)
	}
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          time.Time
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test relative is matched",
			conf:           timeRelativeCond,
			value:          time.Now().Add(-50 * time.Second),
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test relative is not matched",
			conf:           timeRelativeCond,
			value:          time.Now().Add(-150 * time.Second),
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test absolute is matched",
			conf: pattern.NewTimeIntervalCondition(
				pattern.ConditionTimeAbsolute,
				time.Now().Add(-1000*time.Second).Unix(),
				time.Now().Add(1000*time.Second).Unix(),
			),
			value:          time.Now(),
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test absolute is not matched left",
			conf: pattern.NewTimeIntervalCondition(
				pattern.ConditionTimeAbsolute,
				time.Now().Add(-1000*time.Second).Unix(),
				time.Now().Add(1000*time.Second).Unix(),
			),
			value:          time.Now().Add(-2000 * time.Second),
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test absolute is not matched right",
			conf: pattern.NewTimeIntervalCondition(
				pattern.ConditionTimeAbsolute,
				time.Now().Add(-1000*time.Second).Unix(),
				time.Now().Add(1000*time.Second).Unix(),
			),
			value:          time.Now().Add(2000 * time.Second),
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test relative wrong condition value type",
			conf:           pattern.NewStringCondition(pattern.ConditionTimeRelative, "test"),
			value:          time.Now().Add(-50 * time.Second),
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test absolute condition value is wrong type",
			conf:           pattern.NewStringCondition(pattern.ConditionTimeAbsolute, "test"),
			value:          time.Now(),
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test unexpected type",
			conf: pattern.NewTimeIntervalCondition(
				"some",
				time.Now().Add(-1000*time.Second).Unix(),
				time.Now().Add(1000*time.Second).Unix(),
			),
			value:          time.Now(),
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err := dataset.conf.MatchTime(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchTime(t *testing.T) {
	dataSets := []pattern.Condition{
		{
			Type: pattern.ConditionTimeRelative,
			Value: types.DurationWithUnit{
				Value: 100,
				Unit:  "s",
			},
		},
		{
			Type: pattern.ConditionTimeAbsolute,
			Value: map[string]int64{
				"from": 10,
				"to":   10,
			},
		},
	}

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w := condWrapper{}
		err = json.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchTime(time.Now())
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		w = condWrapper{}
		err = bson.Unmarshal(b, &w)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

		_, err = w.Cond.MatchTime(time.Now())
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}
}

func TestCondition_MatchDuration(t *testing.T) {
	durationGtCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{Value: 5, Unit: "m"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	durationLtCond, err := pattern.NewDurationCondition(pattern.ConditionLT, types.DurationWithUnit{Value: 5, Unit: "m"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	durationSomeCond, err := pattern.NewDurationCondition("some", types.DurationWithUnit{Value: 5, Unit: "m"})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          int64
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName:       "test gt is matched",
			conf:           durationGtCond,
			value:          350,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test gt is not matched",
			conf:           durationGtCond,
			value:          250,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test lt is matched",
			conf:           durationLtCond,
			value:          250,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName:       "test lt is not matched",
			conf:           durationLtCond,
			value:          350,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName:       "test unexpected type",
			conf:           durationSomeCond,
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName:       "test value is not a map",
			conf:           pattern.NewStringCondition("some", "test"),
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err := dataset.conf.MatchDuration(dataset.value)
			if result != dataset.expectedResult {
				t.Errorf("expected %t got %t", dataset.expectedResult, result)
			}

			if dataset.expectedErr && err == nil {
				t.Error("an error is expected")
			}

			if !dataset.expectedErr && err != nil {
				t.Errorf("an error is not expected, but got %s", err.Error())
			}
		})
	}
}

func TestCondition_UnmarshalAndMatchDuration(t *testing.T) {
	dataSets := []pattern.Condition{
		{
			Type: pattern.ConditionGT,
			Value: map[string]interface{}{
				"value": 5,
				"unit":  "m",
			},
		},
		{
			Type: pattern.ConditionLT,
			Value: map[string]interface{}{
				"value": 5,
				"unit":  "m",
			},
		},
	}

	for _, dataSet := range dataSets {
		b, err := json.Marshal(condWrapper{Cond: dataSet})
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

		b, err = bson.Marshal(condWrapper{Cond: dataSet})
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

type condWrapper struct {
	Cond pattern.Condition `bson:"cond"`
}
