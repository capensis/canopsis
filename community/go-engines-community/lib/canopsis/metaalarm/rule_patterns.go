package metaalarm

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type RulePatterns struct {
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
}

func (rp *RulePatterns) IsMatched(event types.Event) bool {
	return rp.EntityPatterns.Matches(event.Entity)
}