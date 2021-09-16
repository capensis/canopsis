package correlation

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type RulePatterns struct {
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns,omitempty" json:"entity_patterns,omitempty"`
}

func (rp *RulePatterns) IsMatched(event types.Event) bool {
	return rp.EntityPatterns.Matches(event.Entity)
}
