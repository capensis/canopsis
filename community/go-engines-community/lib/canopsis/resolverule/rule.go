package resolverule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const DefaultRule = "default_rule"

type Rule struct {
	ID          string                    `bson:"_id,omitempty"`
	Name        string                    `bson:"name"`
	Description string                    `bson:"description"`
	Duration    datetime.DurationWithUnit `bson:"duration"`
	Priority    int64                     `bson:"priority"`
	Author      string                    `bson:"author"`
	Created     datetime.CpsTime          `bson:"created,omitempty"`
	Updated     datetime.CpsTime          `bson:"updated,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

// Matches returns true if alarm and entity match patterns.
func (r *Rule) Matches(alarmWithEntity types.AlarmWithEntity) (bool, error) {
	if r.ID == DefaultRule {
		return true, nil
	}

	return match.Match(&alarmWithEntity.Entity, &alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern)
}
