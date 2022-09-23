package oldpattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// TimePattern is a type representing a pattern that can be applied to the value
// of a field of an event that contains a CpsTime.
type TimePattern struct {
	IntegerPattern
}

// Matches returns true if the value is matched by the pattern.
func (p TimePattern) Matches(value types.CpsTime) bool {
	return p.IntegerPattern.Matches(types.CpsNumber(value.Time.Unix()))
}

// TimeRefPattern is a type representing a pattern that can be applied to the
// value of a field of an event that contains a reference to a CpsTime.
type TimeRefPattern struct {
	IntegerRefPattern
}

// Matches returns true if the value is matched by the pattern.
func (p TimeRefPattern) Matches(value *types.CpsTime) bool {
	if value == nil {
		return p.IntegerRefPattern.Matches(nil)
	} else if p.EqualNil {
		return false
	} else {
		return p.IntegerPattern.Matches(types.CpsNumber(value.Time.Unix()))
	}
}
