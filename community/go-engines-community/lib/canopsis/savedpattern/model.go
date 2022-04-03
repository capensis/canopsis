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
	ID               string                `bson:"_id"`
	Title            string                `bson:"title"`
	Type             string                `bson:"type"`
	IsCorporate      bool                  `bson:"is_corporate"`
	AlarmPattern     pattern.Alarm         `bson:"alarm_pattern,omitempty"`
	EntityPattern    pattern.Entity        `bson:"entity_pattern,omitempty"`
	PbehaviorPattern pattern.PbehaviorInfo `bson:"pbehavior_pattern,omitempty"`
	Author           string                `bson:"author"`
	Created          types.CpsTime         `bson:"created,omitempty"`
	Updated          types.CpsTime         `bson:"updated,omitempty"`
}

type AlarmPatternFields struct {
	AlarmPattern pattern.Alarm `bson:"alarm_pattern" json:"alarm_pattern,omitempty"`

	CorporateAlarmPattern      string `bson:"corporate_alarm_pattern" json:"corporate_alarm_pattern,omitempty"`
	CorporateAlarmPatternTitle string `bson:"corporate_alarm_pattern_title" json:"corporate_alarm_pattern_title,omitempty"`
}

type EntityPatternFields struct {
	EntityPattern pattern.Entity `bson:"entity_pattern" json:"entity_pattern,omitempty"`

	CorporateEntityPattern      string `bson:"corporate_entity_pattern" json:"corporate_entity_pattern,omitempty"`
	CorporateEntityPatternTitle string `bson:"corporate_entity_pattern_title" json:"corporate_entity_pattern_title,omitempty"`
}

type PbehaviorPatternFields struct {
	PbehaviorPattern pattern.PbehaviorInfo `bson:"pbehavior_pattern" json:"pbehavior_pattern,omitempty"`

	CorporatePbehaviorPattern      string `bson:"corporate_pbehavior_pattern" json:"corporate_pbehavior_pattern,omitempty"`
	CorporatePbehaviorPatternTitle string `bson:"corporate_pbehavior_pattern_title" json:"corporate_pbehavior_pattern_title,omitempty"`
}
