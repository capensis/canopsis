package savedpattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	TypeAlarm     = "alarm"
	TypeEntity    = "entity"
	TypePbehavior = "pbehavior"
)

type SavedPattern struct {
	ID               string            `bson:"_id"`
	Title            string            `bson:"title"`
	Type             string            `bson:"type"`
	IsShared         bool              `bson:"is_shared"`
	AlarmPattern     pattern.Alarm     `bson:"alarm_pattern,omitempty"`
	EntityPattern    pattern.Entity    `bson:"entity_pattern,omitempty"`
	PbehaviorPattern pattern.Pbehavior `bson:"pbehavior_pattern,omitempty"`
	Author           string            `bson:"author"`
	Created          types.CpsTime     `bson:"created,omitempty"`
	Updated          types.CpsTime     `bson:"updated,omitempty"`
}
