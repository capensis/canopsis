package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Rule struct {
	ID          string                    `bson:"_id"`
	Name        string                    `bson:"name"`
	Description string                    `bson:"description"`
	FreqLimit   int                       `bson:"freq_limit"`
	Duration    datetime.DurationWithUnit `bson:"duration"`
	Priority    int64                     `bson:"priority"`
	Author      string                    `bson:"author"`
	Created     datetime.CpsTime          `bson:"created"`
	Updated     datetime.CpsTime          `bson:"updated"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

// Matches returns true if alarm and entity match field patterns.
func (r *Rule) Matches(alarmWithEntity types.AlarmWithEntity) (bool, error) {
	return match.Match(&alarmWithEntity.Entity, &alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern)
}
