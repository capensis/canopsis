package pattern_test

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"testing"
)

func TestMatchString(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          string
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test equal matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: "test",
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test equal not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: "test",
			},
			value:          "test2",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test not equal matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: "test",
			},
			value:          "test2",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test not equal not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: "test",
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test regexp matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "^test",
			},
			value:          "test",
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test regexp not matched",
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
			result, _, err := pattern.MatchString(dataset.conf, dataset.value)
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

func TestMatchInt(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          int
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test equal matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: 10,
			},
			value:          10,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test equal not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test not equal matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test not equal not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionNotEqual,
				Value: 10,
			},
			value:          10,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test gt matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionGT,
				Value: 10,
			},
			value:          11,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test gt not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionGT,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test lt matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionLT,
				Value: 10,
			},
			value:          9,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test lt not matched",
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
			result, err := pattern.MatchInt(dataset.conf, dataset.value)
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

func TestMatchBool(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          bool
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test equal matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: true,
			},
			value:          true,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test equal not matched",
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
			result, err := pattern.MatchBool(dataset.conf, dataset.value)
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

func TestMatchStringArray(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          []string
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test is empty matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: true,
			},
			value:          []string{},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is empty not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: true,
			},
			value:          []string{"test"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is not empty matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: false,
			},
			value:          []string{"test"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is not empty not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: false,
			},
			value:          []string{},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test has_not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasNot,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test has_not not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasNot,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test3", "test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test has_one_of matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasOneOf,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test3", "test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test has_one_of not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasOneOf,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test has_every matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasEvery,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test1", "test2", "test3", "test4"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test has_every not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionHasEvery,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test1", "test2"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is equal matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test1", "test2", "test3"},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is equal not matched bigger slice",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
				Value: []string{"test1", "test2", "test3"},
			},
			value:          []string{"test1", "test2", "test3", "test4"},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is equal not matched smaller slice",
			conf: pattern.Condition{
				Type:  pattern.ConditionEqual,
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
			result, err = pattern.MatchStringArray(dataset.conf, dataset.value)
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

type refStruct struct{}

func TestMatchRef(t *testing.T) {
	datasets := []struct {
		testName       string
		conf           pattern.Condition
		value          interface{}
		expectedErr    bool
		expectedResult bool
	}{
		{
			testName: "test is empty matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: true,
			},
			value:          nil,
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test is empty not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: true,
			},
			value:          &refStruct{},
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is not empty matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: false,
			},
			value:          nil,
			expectedErr:    false,
			expectedResult: false,
		},
		{
			testName: "test is not empty not matched",
			conf: pattern.Condition{
				Type:  pattern.ConditionIsEmpty,
				Value: false,
			},
			value:          &refStruct{},
			expectedErr:    false,
			expectedResult: true,
		},
		{
			testName: "test cond value is not a bool",
			conf: pattern.Condition{
				Type:  pattern.ConditionRegexp,
				Value: "test",
			},
			value:          &refStruct{},
			expectedErr:    true,
			expectedResult: false,
		},
		{
			testName: "test cond type is not supported",
			conf: pattern.Condition{
				Type:  "some",
				Value: 10,
			},
			value:          &refStruct{},
			expectedErr:    true,
			expectedResult: false,
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.testName, func(t *testing.T) {
			result, err := pattern.MatchRef(dataset.conf, dataset.value)
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
