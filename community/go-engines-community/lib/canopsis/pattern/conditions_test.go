package pattern_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

func TestCondition_MatchString(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          string
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test equal is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: "test",
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test equal is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: "test",
			},
			value:          "test2",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test not equal is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: "test",
			},
			value:          "test2",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test not equal is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: "test",
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is one of is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsOneOf,
				Value: []string{"test2", "test"},
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is one of is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsOneOf,
				Value: []string{"test2", "test3"},
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test regexp is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "^test",
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test regexp is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "^test",
			},
			value:          "tesst",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test bad regexp err",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "[",
			},
			value:          "test",
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond value is not a string",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: 123,
			},
			value:          "test",
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond type is not supported",
			conf: pattern.Condition{
				Type:  "some",
				Value: "test",
			},
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

func TestCondition_MatchInt(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          int
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test equal is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: 10,
			},
			value:          10,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test equal is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test not equal is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test not equal is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: 10,
			},
			value:          10,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test gt is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionGT,
				Value: 10,
			},
			value:          11,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test gt is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionGT,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test lt is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionLT,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test lt is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionLT,
				Value: 10,
			},
			value:          11,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test cond value is not an int",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "test",
			},
			value:          10,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond type is not supported",
			conf: pattern.Condition{
				Type:  "some",
				Value: 10,
			},
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

func TestCondition_MatchBool(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          bool
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test equal is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: true,
			},
			value:          true,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test equal is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: true,
			},
			value:          false,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test cond value is not a bool",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "test",
			},
			value:          false,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond type is not supported",
			conf: pattern.Condition{
				Type:  "some",
				Value: 10,
			},
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

func TestCondition_MatchStringArray(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          []string
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test is empty is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: true,
			},
			value:          []string{},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is empty is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: true,
			},
			value:          []string{"test"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is not empty is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: false,
			},
			value:          []string{"test"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is not empty is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: false,
			},
			value:          []string{},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test has_is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasNot,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test has_not is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasNot,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test3", "test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test has_one_of is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasOneOf,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test3", "test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test has_one_of is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasOneOf,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test has_every is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasEvery,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test1", "test2", "test3", "test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test has_every is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasEvery,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test1", "test2"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test cond value is not a bool",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: []string{"test"},
			},
			value:          []string{"test"},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond value is not a slice of strings",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasEvery,
				Value: true,
			},
			value:          []string{"test"},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond type is not supported",
			conf: pattern.Condition{
				Type:  "some",
				Value: []string{"test"},
			},
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

func TestCondition_MatchRef(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          interface{}
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test is empty is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionExist,
				Value: true,
			},
			value:          nil,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is empty is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionExist,
				Value: true,
			},
			value:          &struct{}{},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is not empty is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionExist,
				Value: false,
			},
			value:          nil,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is not empty is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionExist,
				Value: false,
			},
			value:          &struct{}{},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test cond value is not a bool",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "test",
			},
			value:          &struct{}{},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond type is not supported",
			conf: pattern.Condition{
				Type:  "some",
				Value: 10,
			},
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

func TestCondition_MatchTime(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          time.Time
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test relative is matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionTimeRelative,
				Value: 100,
			},
			value:          time.Now().Add(-50 * time.Second),
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test relative is not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionTimeRelative,
				Value: 100,
			},
			value:          time.Now().Add(-150 * time.Second),
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test absolute is matched",
			conf: pattern.Condition{
				Type: pattern.ConditionTimeAbsolute,
				Value: map[string]int64{
					"from": time.Now().Add(-1000 * time.Second).Unix(),
					"to":   time.Now().Add(1000 * time.Second).Unix(),
				},
			},
			value:          time.Now(),
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test absolute is not matched left",
			conf: pattern.Condition{
				Type: pattern.ConditionTimeAbsolute,
				Value: map[string]int64{
					"from": time.Now().Add(-1000 * time.Second).Unix(),
					"to":   time.Now().Add(1000 * time.Second).Unix(),
				},
			},
			value:          time.Now().Add(-2000 * time.Second),
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test absolute is not matched right",
			conf: pattern.Condition{
				Type: pattern.ConditionTimeAbsolute,
				Value: map[string]int64{
					"from": time.Now().Add(-1000 * time.Second).Unix(),
					"to":   time.Now().Add(1000 * time.Second).Unix(),
				},
			},
			value:          time.Now().Add(2000 * time.Second),
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test relative wrong condition value type",
			conf: pattern.Condition{
				Type:  pattern.ConditionTimeRelative,
				Value: "qwe",
			},
			value:          time.Now().Add(-50 * time.Second),
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test absolute condition value 'from' is absent",
			conf: pattern.Condition{
				Type: pattern.ConditionTimeAbsolute,
				Value: map[string]int64{
					"to": time.Now().Add(1000 * time.Second).Unix(),
				},
			},
			value:          time.Now(),
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test absolute condition value 'to' is absent",
			conf: pattern.Condition{
				Type: pattern.ConditionTimeAbsolute,
				Value: map[string]int64{
					"from": time.Now().Add(-1000 * time.Second).Unix(),
				},
			},
			value:          time.Now(),
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test absolute condition value is wrong type",
			conf: pattern.Condition{
				Type: pattern.ConditionTimeAbsolute,
				Value: map[string]string{
					"from": "abc",
				},
			},
			value:          time.Now(),
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test unexpected type",
			conf: pattern.Condition{
				Type: "some",
				Value: map[string]int64{
					"from": time.Now().Add(-1000 * time.Second).Unix(),
					"to":   time.Now().Add(1000 * time.Second).Unix(),
				},
			},
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

func TestCondition_MatchDuration(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          int64
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test gt is matched",
			conf: pattern.Condition{
				Type: pattern.ConditionGT,
				Value: map[string]interface{}{
					"unit":  "m",
					"value": 5,
				},
			},
			value:          350,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test gt is not matched",
			conf: pattern.Condition{
				Type: pattern.ConditionGT,
				Value: map[string]interface{}{
					"unit":  "m",
					"value": 5,
				},
			},
			value:          250,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test lt is matched",
			conf: pattern.Condition{
				Type: pattern.ConditionLT,
				Value: map[string]interface{}{
					"unit":  "m",
					"value": 5,
				},
			},
			value:          250,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test lt is not matched",
			conf: pattern.Condition{
				Type: pattern.ConditionLT,
				Value: map[string]interface{}{
					"unit":  "m",
					"value": 5,
				},
			},
			value:          350,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test unexpected type",
			conf: pattern.Condition{
				Type: "some",
				Value: map[string]interface{}{
					"unit":  "m",
					"value": 5,
				},
			},
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test value is not a map",
			conf: pattern.Condition{
				Type:  "some",
				Value: "qwe",
			},
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test unit is missed",
			conf: pattern.Condition{
				Type: "some",
				Value: map[string]interface{}{
					"value": 5,
				},
			},
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test value is missed",
			conf: pattern.Condition{
				Type: "some",
				Value: map[string]interface{}{
					"unit": "m",
				},
			},
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test unit is not a string",
			conf: pattern.Condition{
				Type: "some",
				Value: map[string]interface{}{
					"unit":  123,
					"value": 5,
				},
			},
			value:          250,
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test value is not an int",
			conf: pattern.Condition{
				Type: "some",
				Value: map[string]interface{}{
					"unit":  "m",
					"value": "qwe",
				},
			},
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
