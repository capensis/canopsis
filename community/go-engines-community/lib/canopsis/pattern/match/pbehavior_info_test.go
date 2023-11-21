package match_test

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func TestMatchPbehaviorInfoPattern(t *testing.T) {
	dataSets := getMatchPbehaviorInfoPatternDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := match.MatchPbehaviorInfoPattern(data.pattern, &data.pbehaviorInfo)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func getMatchPbehaviorInfoPatternDataSets() map[string]PbehaviorInfoDataSet {
	return map[string]PbehaviorInfoDataSet{
		"given empty pattern should match": {
			pattern: pattern.PbehaviorInfo{},
			pbehaviorInfo: types.PbehaviorInfo{
				ID: "test id",
			},
			matchResult: true,
		},
		"given string field condition should match": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{
				ID: "test id",
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.id",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test id"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{
				ID: "test another id",
			},
			matchResult: false,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "created",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{},
			matchErr:      pattern.ErrUnsupportedField,
		},
		"given active canonical field condition and emtpty pbehavior infos should match": {
			pattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, types.PbhCanonicalTypeActive),
					},
				},
			},
			pbehaviorInfo: types.PbehaviorInfo{},
			matchResult:   true,
		},
	}
}

type PbehaviorInfoDataSet struct {
	pattern       pattern.PbehaviorInfo
	pbehaviorInfo types.PbehaviorInfo
	matchErr      error
	matchResult   bool
}
